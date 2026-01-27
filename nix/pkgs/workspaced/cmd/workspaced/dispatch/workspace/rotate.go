package workspace

import (
	"workspaced/pkg/drivers/wm"

	"github.com/spf13/cobra"
)

var rotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "Rotate workspaces across outputs",
	RunE: func(c *cobra.Command, args []string) error {
		return wm.RotateWorkspaces(c.Context())
	},
}

func init() {
	Command.AddCommand(rotateCmd)
}
