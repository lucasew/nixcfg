package screen

import (
	"workspaced/pkg/drivers/screen"

	"github.com/spf13/cobra"
)

var offCmd = &cobra.Command{
	Use:   "off",
	Short: "Turn off the screen (DPMS)",
	RunE: func(c *cobra.Command, args []string) error {
		return screen.SetDPMS(c.Context(), false)
	},
}

func init() {
	Command.AddCommand(offCmd)
}
