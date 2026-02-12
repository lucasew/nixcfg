package menu

import (
	"workspaced/pkg/driver"
	"workspaced/pkg/driver/menu"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "window",
			Short: "Window switcher",
			RunE: func(c *cobra.Command, args []string) error {
				d, err := driver.Get[menu.Driver](c.Context())
				if err != nil {
					return err
				}
				return d.SwitchWindow(c.Context())
			},
		})
	})
}
