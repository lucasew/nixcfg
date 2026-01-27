package is

import (
	"os"
	"workspaced/pkg/common"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "riverwood",
			Short: "Check if host is riverwood",
			Run: func(c *cobra.Command, args []string) {
				if !common.IsRiverwood() {
					os.Exit(1)
				}
			},
		})
	})
}
