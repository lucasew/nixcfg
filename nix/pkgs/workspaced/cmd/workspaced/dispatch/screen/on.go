package screen

import (
	"workspaced/pkg/drivers/screen"

	"github.com/spf13/cobra"
)

var onCmd = &cobra.Command{
	Use:   "on",
	Short: "Turn on the screen (DPMS)",
	RunE: func(c *cobra.Command, args []string) error {
		return screen.SetDPMS(c.Context(), true)
	},
}

func init() {
	Command.AddCommand(onCmd)
}
