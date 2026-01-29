package audio

import (
	"workspaced/pkg/drivers/audio"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:     "show",
			Aliases: []string{"status"},
			Short:   "Show current volume",
			RunE: func(c *cobra.Command, args []string) error {
				return audio.ShowStatus(c.Context())
			},
		})
	})
}
