package brightness

import (
	"workspaced/pkg/driver/brightness"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:     "show",
			Aliases: []string{"status"},
			Short:   "Show current brightness",
			RunE: func(c *cobra.Command, args []string) error {
				return brightness.ShowStatus(c.Context())
			},
		})
	})
}
