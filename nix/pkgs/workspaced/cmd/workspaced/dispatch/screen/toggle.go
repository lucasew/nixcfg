package screen

import (
	"workspaced/pkg/drivers/screen"

	"github.com/spf13/cobra"
)

var toggleCmd = &cobra.Command{
	Use:   "toggle",
	Short: "Toggle screen state (DPMS)",
	RunE: func(c *cobra.Command, args []string) error {
		return screen.ToggleDPMS(c.Context())
	},
}

func init() {
	Command.AddCommand(toggleCmd)
}
