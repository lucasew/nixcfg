package is

import (
	"fmt"
	"workspaced/pkg/host"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "whiterun",
			Short: "Check if host is whiterun",
			RunE: func(c *cobra.Command, args []string) error {
				if !host.IsWhiterun() {
					return fmt.Errorf("not whiterun")
				}
				return nil
			},
		})
	})
}
