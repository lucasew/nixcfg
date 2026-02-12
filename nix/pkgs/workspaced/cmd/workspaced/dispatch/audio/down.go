package audio

import (
	"workspaced/pkg/driver"
	"workspaced/pkg/driver/audio"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "down",
			Short: "Decrease volume",
			RunE: func(c *cobra.Command, args []string) error {
				d, err := driver.Get[audio.Driver](c.Context())
				if err != nil {
					return err
				}
				volume, err := d.GetVolume(c.Context())
				if err != nil {
					return err
				}
				return d.SetVolume(c.Context(), volume-0.05)
			},
		})
	})
}
