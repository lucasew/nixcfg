package wallpaper

import (
	"workspaced/pkg/drivers/wallpaper"

	"github.com/spf13/cobra"
)

var videoCmd = &cobra.Command{
	Use:   "video <path>",
	Short: "Set an animated video as wallpaper",
	Args:  cobra.ExactArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return wallpaper.SetAnimated(c.Context(), args[0])
	},
}

func init() {
	Command.AddCommand(videoCmd)
}
