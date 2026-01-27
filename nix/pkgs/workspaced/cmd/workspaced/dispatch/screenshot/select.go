package screenshot

import (
	"workspaced/pkg/drivers/screenshot"

	"github.com/spf13/cobra"
)

var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "Capture selected area",
	RunE: func(c *cobra.Command, args []string) error {
		_, err := screenshot.Capture(c.Context(), true)
		return err
	},
}

func init() {
	Command.AddCommand(selectCmd)
}
