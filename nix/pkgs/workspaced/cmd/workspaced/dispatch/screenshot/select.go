package screenshot

import (
	"errors"
	dapi "workspaced/pkg/api"
	"workspaced/pkg/screenshot"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "select",
			Short: "Capture selected area",
			RunE: func(c *cobra.Command, args []string) error {
				path, err := screenshot.Capture(c.Context(), screenshot.Selection)
				if err != nil {
					if errors.Is(err, dapi.ErrCanceled) {
						return nil
					}
					return err
				}
				if path != "" {
					c.Println(path)
				}
				return nil
			},
		})
	})
}
