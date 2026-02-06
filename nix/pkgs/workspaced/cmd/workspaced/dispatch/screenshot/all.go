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
				path, err := screenshot.Capture(c.Context(), screenshot.All)
				if err != nil {
					return err
				}
				c.Println(path)
				return nil
			},
		})
	})
}
