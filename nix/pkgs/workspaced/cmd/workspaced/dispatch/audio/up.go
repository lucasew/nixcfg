package audio

import (
	"workspaced/pkg/drivers/audio"

	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Increase volume",
	RunE: func(c *cobra.Command, args []string) error {
		return audio.SetVolume(c.Context(), "+5%")
	},
}

func init() {
	Command.AddCommand(upCmd)
}
