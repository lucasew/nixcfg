package screen

import (
	"workspaced/pkg/driver/screen"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "toggle",
			Short: "Toggle screen state (DPMS)",
			RunE: func(c *cobra.Command, args []string) error {
				return screen.ToggleDPMS(c.Context())
			},
		})
	})
}
