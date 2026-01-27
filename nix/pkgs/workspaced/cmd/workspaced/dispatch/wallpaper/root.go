package wallpaper

import (
	"workspaced/pkg/drivers/wallpaper"

	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "wallpaper",
	Short: "Wallpaper management",
}

var changeCmd = &cobra.Command{
	Use:   "change [path]",
	Short: "Change wallpaper to a random image or specific path",
	RunE: func(c *cobra.Command, args []string) error {
		path := ""
		if len(args) > 0 {
			path = args[0]
		}
		return wallpaper.SetStatic(c.Context(), path)
	},
}

var videoCmd = &cobra.Command{
	Use:   "video <path>",
	Short: "Set an animated video as wallpaper",
	Args:  cobra.ExactArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return wallpaper.SetAnimated(c.Context(), args[0])
	},
}

var apodCmd = &cobra.Command{
	Use:   "apod",
	Short: "Fetch NASA Astronomy Picture of the Day and set as wallpaper",
	RunE: func(c *cobra.Command, args []string) error {
		return wallpaper.SetAPOD(c.Context())
	},
}

func init() {
	Command.AddCommand(changeCmd)
	Command.AddCommand(videoCmd)
	Command.AddCommand(apodCmd)
}
