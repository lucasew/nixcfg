package screenshot

import (
	"workspaced/pkg/driver/screenshot"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "full",
			Short: "Capture full screen (all outputs)",
			RunE: func(c *cobra.Command, args []string) error {
				path, err := screenshot.Capture(c.Context(), screenshot.TargetAll)
				if err != nil {
					return err
				}
				c.Println(path)
				return nil
			},
		})
	})
}
