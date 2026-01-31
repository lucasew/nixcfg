package is

import (
	"fmt"
	"workspaced/pkg/common"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "in-store",
			Short: "Check if dotfiles are in nix store",
			RunE: func(c *cobra.Command, args []string) error {
				if !common.IsInStore() {
					return fmt.Errorf("not in store")
				}
				return nil
			},
		})
	})
}
