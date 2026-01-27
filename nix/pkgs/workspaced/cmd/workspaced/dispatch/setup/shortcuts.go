package setup

import (
	"workspaced/pkg/drivers/setup"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "shortcuts",
			Short: "Setup Termux shortcuts",
			RunE: func(c *cobra.Command, args []string) error {
				return setup.SetupTermuxShortcuts(c.Context())
			},
		})
	})
}
