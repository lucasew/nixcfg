package brightness

import (
	"workspaced/pkg/brightness"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "down",
			Short: "Decrease brightness",
			RunE: func(c *cobra.Command, args []string) error {
				return brightness.SetBrightness(c.Context(), "5%-")
			},
		})
	})
}
