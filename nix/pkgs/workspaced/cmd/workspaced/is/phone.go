package is

import (
	"fmt"
	"workspaced/pkg/env"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "phone",
			Short: "Check if environment is a phone (Termux)",
			RunE: func(c *cobra.Command, args []string) error {
				if !env.IsPhone() {
					return fmt.Errorf("not phone")
				}
				return nil
			},
		})
	})
}
