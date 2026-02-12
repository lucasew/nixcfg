package screenshot

import (
	"workspaced/pkg/driver/screenshot"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "window",
			Short: "Capture current window",
			RunE: func(c *cobra.Command, args []string) error {
				path, err := screenshot.Capture(c.Context(), screenshot.TargetWindow)
				if err != nil {
					return err
				}
				c.Println(path)
				return nil
			},
		})
	})
}
