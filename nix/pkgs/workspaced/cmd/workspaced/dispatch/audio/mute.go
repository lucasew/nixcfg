package audio

import (
	"workspaced/pkg/drivers/audio"

	"github.com/spf13/cobra"
)

var muteCmd = &cobra.Command{
	Use:   "mute",
	Short: "Toggle mute",
	RunE: func(c *cobra.Command, args []string) error {
		return audio.SetVolume(c.Context(), "toggle")
	},
}

func init() {
	Command.AddCommand(muteCmd)
}
