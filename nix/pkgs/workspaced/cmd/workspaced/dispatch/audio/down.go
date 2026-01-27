package audio

import (
	"workspaced/pkg/drivers/audio"

	"github.com/spf13/cobra"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Decrease volume",
	RunE: func(c *cobra.Command, args []string) error {
		return audio.SetVolume(c.Context(), "-5%")
	},
}

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(downCmd)
	})
}
