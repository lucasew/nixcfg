package screenshot

import (
	"workspaced/pkg/drivers/screenshot"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "all",
			Short: "Capture all outputs",
			RunE: func(c *cobra.Command, args []string) error {
				_, err := screenshot.Capture(c.Context(), screenshot.All)
				return err
			},
		})
	})
}
