package brightness

import (
	"workspaced/pkg/drivers/brightness"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"status"},
	Short:   "Show current brightness",
	RunE: func(c *cobra.Command, args []string) error {
		return brightness.ShowStatus(c.Context())
	},
}

func init() {
	Command.AddCommand(showCmd)
}
