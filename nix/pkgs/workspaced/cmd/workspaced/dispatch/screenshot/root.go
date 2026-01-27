package screenshot

import (
	"workspaced/pkg/drivers/screenshot"

	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "screenshot",
	Short: "Screen capture management",
}

var fullCmd = &cobra.Command{
	Use:   "full",
	Short: "Capture full screen",
	RunE: func(c *cobra.Command, args []string) error {
		_, err := screenshot.Capture(c.Context(), false)
		return err
	},
}

var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "Capture selected area",
	RunE: func(c *cobra.Command, args []string) error {
		_, err := screenshot.Capture(c.Context(), true)
		return err
	},
}

func init() {
	Command.AddCommand(fullCmd)
	Command.AddCommand(selectCmd)
}
