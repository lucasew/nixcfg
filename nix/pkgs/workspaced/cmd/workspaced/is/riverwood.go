package is

import (
	"fmt"
	"workspaced/pkg/env"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "riverwood",
			Short: "Check if host is riverwood",
			RunE: func(c *cobra.Command, args []string) error {
				if !env.IsRiverwood() {
					return fmt.Errorf("not riverwood")
				}
				return nil
			},
		})
	})
}
