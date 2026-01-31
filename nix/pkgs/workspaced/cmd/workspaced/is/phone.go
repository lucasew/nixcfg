package is

import (
	"fmt"
	"workspaced/pkg/common"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "phone",
			Short: "Check if environment is a phone (Termux)",
			RunE: func(c *cobra.Command, args []string) error {
				if !common.IsPhone() {
					return fmt.Errorf("not phone")
				}
				return nil
			},
		})
	})
}
