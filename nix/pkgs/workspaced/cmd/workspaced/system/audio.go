package system

import (
	"workspaced/pkg/driver/audio"

	"github.com/spf13/cobra"
)

func newAudioCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "audio",
		Short: "Volume control",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "up",
		Short: "Increase volume",
		RunE: func(c *cobra.Command, args []string) error {
			return audio.IncreaseVolume(c.Context())
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "down",
		Short: "Decrease volume",
		RunE: func(c *cobra.Command, args []string) error {
			return audio.DecreaseVolume(c.Context())
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "mute",
		Short: "Toggle mute",
		RunE: func(c *cobra.Command, args []string) error {
			return audio.ToggleMute(c.Context())
		},
	})

	return cmd
}
