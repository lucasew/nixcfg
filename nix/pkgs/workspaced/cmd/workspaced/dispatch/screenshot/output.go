package screenshot

import (
	"workspaced/pkg/screenshot"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "output",
			Short: "Capture current output (monitor)",
			RunE: func(c *cobra.Command, args []string) error {
				path, err := screenshot.Capture(c.Context(), screenshot.Output)
				if err != nil {
					return err
				}
				c.Println(path)
				return nil
			},
		})
	})
}
