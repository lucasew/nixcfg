package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"workspaced/cmd/workspaced/dispatch/common"
	"workspaced/cmd/workspaced/dispatch/types"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

type Config struct {
	Workspaces map[string]int `toml:"workspaces"`
}

var Command = &cobra.Command{
	Use:   "workspace",
	Short: "Workspace management commands",
}

var workspaceSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Workspace switcher using a menu (dmenu/rofi)",
	RunE: func(c *cobra.Command, args []string) error {
		move, _ := c.Flags().GetBool("move")

		// Inlined logic
		home, err := os.UserHomeDir()
		if err != nil {
			return err
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

		var keys []string
		for k := range mapping {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		cmd := common.RunCmd(c, "dmenu")
		cmd.Stdin = strings.NewReader(strings.Join(keys, "\n"))

		out, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("menu selection failed: %w", err)
		}

		selected := strings.TrimSpace(string(out))
		if selected == "" {
			return nil
		}

		workspaceNum, ok := mapping[selected]
		if !ok {
			return fmt.Errorf("invalid selection")
		}

		return switchToWorkspace(c, workspaceNum, move)
	},
}

var workspaceNextCmd = &cobra.Command{
	Use:   "next",
	Short: "Go to the next available workspace",
	RunE: func(c *cobra.Command, args []string) error {
		move, _ := c.Flags().GetBool("move")

		runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
		if runtimeDir == "" {
			runtimeDir = filepath.Join(os.TempDir(), fmt.Sprintf("workspaced-%d", os.Getuid()))
		}
		workspacedDir := filepath.Join(runtimeDir, "workspaced")
		if err := os.MkdirAll(workspacedDir, 0700); err != nil {
			return fmt.Errorf("failed to create runtime dir: %w", err)
		}

		wsFile := filepath.Join(workspacedDir, "last_ws")
		lastWS := 10
		if data, err := os.ReadFile(wsFile); err == nil {
			if val, err := strconv.Atoi(strings.TrimSpace(string(data))); err == nil {
				lastWS = val
			}
		}

		nextWS := lastWS + 1
		if err := os.WriteFile(wsFile, []byte(strconv.Itoa(nextWS)), 0600); err != nil {
			return fmt.Errorf("failed to save last workspace: %w", err)
		}

		return switchToWorkspace(c, nextWS, move)
	},
}

func switchToWorkspace(c *cobra.Command, num int, move bool) error {
	ctx := c.Context()
	var env []string
	if ctx != nil {
		if val, ok := ctx.Value(types.EnvKey).([]string); ok {
			env = val
		}
	}
	rpc := common.GetRPC(env)

	rpcCmd := common.RunCmd(c, rpc)
	if move {
		rpcCmd.Args = append(rpcCmd.Args, "move", "container", "to", "workspace", "number", strconv.Itoa(num))
	} else {
		rpcCmd.Args = append(rpcCmd.Args, "workspace", "number", strconv.Itoa(num))
	}

	if err := rpcCmd.Run(); err != nil {
		return fmt.Errorf("swaymsg/i3-msg failed: %w", err)
	}

	fmt.Printf("Switched to workspace %d (move=%v)\n", num, move)
	return nil
}

func init() {
	Command.PersistentFlags().Bool("move", false, "Move container to workspace")
	Command.AddCommand(workspaceSearchCmd)
	Command.AddCommand(workspaceNextCmd)
}
