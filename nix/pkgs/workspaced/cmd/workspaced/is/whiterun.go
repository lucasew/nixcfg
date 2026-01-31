package is

import (
	"fmt"
	"workspaced/pkg/common"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "whiterun",
			Short: "Check if host is whiterun",
			RunE: func(c *cobra.Command, args []string) error {
				if !common.IsWhiterun() {
					return fmt.Errorf("not whiterun")
				}
				return nil
			},
		})
	})
}
