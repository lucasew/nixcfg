package audio

import (
	"workspaced/pkg/drivers/audio"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "down",
			Short: "Decrease volume",
			RunE: func(c *cobra.Command, args []string) error {
				return audio.SetVolume(c.Context(), "-5%")
			},
		})
	})
}
