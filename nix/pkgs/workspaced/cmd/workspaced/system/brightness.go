package system

import (
	"workspaced/pkg/driver/brightness"

	"github.com/spf13/cobra"
)

func newBrightnessCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "brightness",
		Short: "Brightness control",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "up",
		Short: "Increase brightness",
		RunE: func(c *cobra.Command, args []string) error {
			return brightness.IncreaseBrightness(c.Context())
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "down",
		Short: "Decrease brightness",
		RunE: func(c *cobra.Command, args []string) error {
			return brightness.DecreaseBrightness(c.Context())
		},
	})

	return cmd
}
