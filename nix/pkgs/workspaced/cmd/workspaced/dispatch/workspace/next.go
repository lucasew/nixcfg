package workspace

import (
	"workspaced/pkg/drivers/wm"

	"github.com/spf13/cobra"
)

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Go to the next available workspace",
	RunE: func(c *cobra.Command, args []string) error {
		move, _ := c.Flags().GetBool("move")
		return wm.NextWorkspace(c.Context(), move)
	},
}

func init() {
	Command.AddCommand(nextCmd)
}
