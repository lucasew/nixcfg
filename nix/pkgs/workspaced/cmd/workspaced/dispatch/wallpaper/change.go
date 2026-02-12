package wallpaper

import (
	"workspaced/pkg/driver/wallpaper"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "change [path]",
			Short: "Change wallpaper to a random image or specific path",
			RunE: func(c *cobra.Command, args []string) error {
				path := ""
				if len(args) > 0 {
					path = args[0]
				}
				return wallpaper.SetStatic(c.Context(), path)
			},
		})
	})
}
