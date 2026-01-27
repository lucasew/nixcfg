package wallpaper

import (
	"workspaced/pkg/drivers/wallpaper"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "apod",
			Short: "Fetch NASA Astronomy Picture of the Day and set as wallpaper",
			RunE: func(c *cobra.Command, args []string) error {
				return wallpaper.SetAPOD(c.Context())
			},
		})
	})
}
