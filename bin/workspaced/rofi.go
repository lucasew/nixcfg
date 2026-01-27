package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

var rofiCmd = &cobra.Command{
	Use:   "rofi",
	Short: "Rofi workspace switcher",
	Run: func(cmd *cobra.Command, args []string) {
		runOrRoute("rofi", args, func() (string, error) {
			return RunRofi(args, os.Environ())
		})
	},
}

type Config struct {
	Workspaces map[string]int `toml:"workspaces"`
}

func RunRofi(args []string, env []string) (string, error) {
	// Args parsing: check for --move
	move := false
	for _, arg := range args {
		if arg == "--move" {
			move = true
		}
	}

	// Load settings
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	settingsPath := filepath.Join(home, "settings.toml")

	mapping := map[string]int{
		"www":  1,
		"meet": 2,
	}

	if _, err := os.Stat(settingsPath); err == nil {
		var conf Config
		if _, err := toml.DecodeFile(settingsPath, &conf); err == nil {
			for k, v := range conf.Workspaces {
				mapping[k] = v
			}
		}
	}

	// Sort keys for deterministic menu
	var keys []string
	for k := range mapping {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	cmd := exec.Command("dmenu")
	cmd.Stdin = strings.NewReader(strings.Join(keys, "\n"))

	if len(env) > 0 {
		cmd.Env = env
	}

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("menu selection failed: %w", err)
	}

	selected := strings.TrimSpace(string(out))
	if selected == "" {
		return "No selection", nil
	}

	workspaceNum, ok := mapping[selected]
	if !ok {
		return "Invalid selection", nil
	}

	rpc := getRPC()

	rpcCmd := exec.Command(rpc)
	if len(env) > 0 {
		rpcCmd.Env = env
	}

	if move {
		rpcCmd.Args = append(rpcCmd.Args, "move", "container", "to", "workspace", "number", fmt.Sprintf("%d", workspaceNum))
	} else {
		rpcCmd.Args = append(rpcCmd.Args, "workspace", "number", fmt.Sprintf("%d", workspaceNum))
	}

	if err := rpcCmd.Run(); err != nil {
		return "", fmt.Errorf("swaymsg/i3-msg failed: %w", err)
	}

	return fmt.Sprintf("Switched to %s (%d)", selected, workspaceNum), nil
}
