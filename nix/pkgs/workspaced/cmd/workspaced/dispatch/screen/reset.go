package screen

import (
	"workspaced/pkg/screen"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "reset",
			Short: "Reset screen resolution based on host",
			RunE: func(c *cobra.Command, args []string) error {
				return screen.Reset(c.Context())
			},
		})
	})
}
