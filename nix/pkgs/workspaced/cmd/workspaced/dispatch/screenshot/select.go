package screenshot

import (
	"workspaced/pkg/drivers/screenshot"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "select",
			Short: "Capture selected area",
			RunE: func(c *cobra.Command, args []string) error {
				_, err := screenshot.Capture(c.Context(), screenshot.Selection)
				return err
			},
		})
	})
}
