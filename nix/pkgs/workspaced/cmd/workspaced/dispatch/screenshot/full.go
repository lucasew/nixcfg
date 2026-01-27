package screenshot

import (
	"workspaced/pkg/drivers/screenshot"

	"github.com/spf13/cobra"
)

var fullCmd = &cobra.Command{
	Use:   "full",
	Short: "Capture full screen",
	RunE: func(c *cobra.Command, args []string) error {
		_, err := screenshot.Capture(c.Context(), false)
		return err
	},
}

func init() {
	Command.AddCommand(fullCmd)
}
