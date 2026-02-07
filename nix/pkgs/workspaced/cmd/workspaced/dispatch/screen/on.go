package screen

import (
	"workspaced/pkg/screen"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "on",
			Short: "Turn on the screen (DPMS)",
			RunE: func(c *cobra.Command, args []string) error {
				return screen.SetDPMS(c.Context(), true)
			},
		})
	})
}
