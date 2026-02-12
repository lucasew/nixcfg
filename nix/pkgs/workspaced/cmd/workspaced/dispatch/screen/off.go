package screen

import (
	"workspaced/pkg/driver/screen"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "off",
			Short: "Turn off the screen (DPMS)",
			RunE: func(c *cobra.Command, args []string) error {
				return screen.SetDPMS(c.Context(), false)
			},
		})
	})
}
