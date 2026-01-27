package workspace

import (
	"workspaced/pkg/drivers/wm"

	"github.com/spf13/cobra"
)

var scratchpadCmd = &cobra.Command{
	Use:   "scratchpad",
	Short: "Toggle scratchpad visibility",
	RunE: func(c *cobra.Command, args []string) error {
		return wm.ToggleScratchpad(c.Context())
	},
}

func init() {
	Command.AddCommand(scratchpadCmd)
}
