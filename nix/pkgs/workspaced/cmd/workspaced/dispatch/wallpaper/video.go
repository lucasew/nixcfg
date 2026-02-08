package wallpaper

import (
	"workspaced/pkg/wallpaper"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "video <path>",
			Short: "Set an animated video as wallpaper",
			Args:  cobra.ExactArgs(1),
			RunE: func(c *cobra.Command, args []string) error {
				return wallpaper.SetAnimated(c.Context(), args[0])
			},
		})
	})
}
