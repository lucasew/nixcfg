package screen

import (
	"workspaced/pkg/screen"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "lock",
			Short: "Lock the screen and turn it off",
			RunE: func(c *cobra.Command, args []string) error {
				return screen.Lock(c.Context())
			},
		})
	})
}
