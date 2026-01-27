package menu

import (
	"fmt"
	"sort"
	"strings"
	"workspaced/pkg/common"
	"workspaced/pkg/drivers/wm"

	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "menu",
	Short: "Interactive menus",
}

var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Workspace switcher",
	RunE: func(c *cobra.Command, args []string) error {
		move, _ := c.Flags().GetBool("move")
		config, err := common.LoadConfig()
		if err != nil {
			return err
		}

		var keys []string
		for k := range config.Workspaces {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		cmd := common.RunCmd(c.Context(), "dmenu")
		cmd.Stdin = strings.NewReader(strings.Join(keys, "\n"))

		out, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("menu selection failed: %w", err)
		}

		selected := strings.TrimSpace(string(out))
		if selected == "" {
			return nil
		}

		workspaceNum, ok := config.Workspaces[selected]
		if !ok {
			return fmt.Errorf("invalid selection")
		}

		return wm.SwitchToWorkspace(c.Context(), workspaceNum, move)
	},
}

var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Application launcher",
	RunE: func(c *cobra.Command, args []string) error {
		return common.RunCmd(c.Context(), "rofi-launch").Run()
	},
}

var windowCmd = &cobra.Command{
	Use:   "window",
	Short: "Window switcher",
	RunE: func(c *cobra.Command, args []string) error {
		return common.RunCmd(c.Context(), "rofi-window").Run()
	},
}

func init() {
	workspaceCmd.Flags().Bool("move", false, "Move container to workspace")
	Command.AddCommand(workspaceCmd)
	Command.AddCommand(launchCmd)
	Command.AddCommand(windowCmd)
}
