package workspace

import (
	"workspaced/pkg/drivers/wm"

	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "workspace",
	Short: "Workspace management commands",
}

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Go to the next available workspace",
	RunE: func(c *cobra.Command, args []string) error {
		move, _ := c.Flags().GetBool("move")
		return wm.NextWorkspace(c.Context(), move)
	},
}

var rotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "Rotate workspaces across outputs",
	RunE: func(c *cobra.Command, args []string) error {
		return wm.RotateWorkspaces(c.Context())
	},
}

var scratchpadCmd = &cobra.Command{
	Use:   "scratchpad",
	Short: "Toggle scratchpad visibility",
	RunE: func(c *cobra.Command, args []string) error {
		return wm.ToggleScratchpad(c.Context())
	},
}

func init() {
	Command.PersistentFlags().Bool("move", false, "Move container to workspace")
	Command.AddCommand(nextCmd)
	Command.AddCommand(rotateCmd)
	Command.AddCommand(scratchpadCmd)
}
