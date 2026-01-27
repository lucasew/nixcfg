package dispatch

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

type Config struct {
	Workspaces map[string]int `toml:"workspaces"`
}

func init() {
	Command.AddCommand(rofiCmd)
}

var rofiCmd = &cobra.Command{
	Use:   "rofi",
	Short: "Rofi workspace switcher",
	RunE: func(c *cobra.Command, args []string) error {
		// Args parsing: check for --move
		move := false
		for _, arg := range args {
			if arg == "--move" {
				move = true
			}
		}

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

		cmd := runCmd(c, "dmenu")
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

		ctx := c.Context()
		var env []string
		if ctx != nil {
			env, _ = ctx.Value("env").([]string)
		}
		rpc := getRPC(env)

		rpcCmd := runCmd(c, rpc)
		if move {
			rpcCmd.Args = append(rpcCmd.Args, "move", "container", "to", "workspace", "number", fmt.Sprintf("%d", workspaceNum))
		} else {
			rpcCmd.Args = append(rpcCmd.Args, "workspace", "number", fmt.Sprintf("%d", workspaceNum))
		}

		if err := rpcCmd.Run(); err != nil {
			return fmt.Errorf("swaymsg/i3-msg failed: %w", err)
		}

		fmt.Printf("Switched to %s (%d)\n", selected, workspaceNum)
		return nil
	},
}
