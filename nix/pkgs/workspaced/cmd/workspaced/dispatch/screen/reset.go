package screen

import (
	"workspaced/pkg/drivers/screen"

	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset screen resolution based on host",
	RunE: func(c *cobra.Command, args []string) error {
		return screen.Reset(c.Context())
	},
}

func init() {
	Command.AddCommand(resetCmd)
}
