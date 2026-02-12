package brightness

import (
	"workspaced/pkg/driver"
	"workspaced/pkg/driver/brightness"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "down",
			Short: "Decrease brightness",
			RunE: func(c *cobra.Command, args []string) error {
				d, err := driver.Get[brightness.Driver](c.Context())
				if err != nil {
					return err
				}
				status, err := d.Status(c.Context())
				if err != nil {
					return err
				}
				err = d.SetBrightness(c.Context(), status.Brightness-0.05)
				if err != nil {
					return err
				}
				return brightness.ShowStatus(c.Context())

			},
		})
	})
}
