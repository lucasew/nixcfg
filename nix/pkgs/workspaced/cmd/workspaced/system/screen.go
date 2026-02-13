package system

import (
	"workspaced/pkg/driver/screen"

	"github.com/spf13/cobra"
)

func newScreenCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "screen",
		Short: "Screen and display management",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "on",
		Short: "Turn screens on",
		RunE: func(c *cobra.Command, args []string) error {
			return screen.SetDPMS(c.Context(), true)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "off",
		Short: "Turn screens off",
		RunE: func(c *cobra.Command, args []string) error {
			return screen.SetDPMS(c.Context(), false)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "toggle",
		Short: "Toggle screen power",
		RunE: func(c *cobra.Command, args []string) error {
			return screen.ToggleDPMS(c.Context())
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "reset",
		Short: "Reset screen layout",
		RunE: func(c *cobra.Command, args []string) error {
			return screen.Reset(c.Context())
		},
	})

	return cmd
}
