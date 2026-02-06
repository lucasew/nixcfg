package screenshot

import (
	"workspaced/pkg/drivers/screenshot"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "output",
			Short: "Capture current output (monitor)",
			RunE: func(c *cobra.Command, args []string) error {
				_, err := screenshot.Capture(c.Context(), screenshot.Output)
				return err
			},
		})
	})
}
