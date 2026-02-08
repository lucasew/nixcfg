package audio

import (
	"workspaced/pkg/audio"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "mute",
			Short: "Toggle mute",
			RunE: func(c *cobra.Command, args []string) error {
				return audio.SetVolume(c.Context(), "toggle")
			},
		})
	})
}
