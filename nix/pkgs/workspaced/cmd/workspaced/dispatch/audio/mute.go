package audio

import (
	"workspaced/pkg/driver"
	"workspaced/pkg/driver/audio"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "mute",
			Short: "Toggle mute",
			RunE: func(c *cobra.Command, args []string) error {
				d, err := driver.Get[audio.Driver](c.Context())
				if err != nil {
					return err
				}
				err = d.ToggleMute(c.Context())
				if err != nil {
					return err
				}
				return audio.ShowStatus(c.Context())
			},
		})
	})
}
