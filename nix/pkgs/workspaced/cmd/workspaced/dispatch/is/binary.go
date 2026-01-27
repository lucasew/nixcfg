package is

import (
	"os"
	"workspaced/pkg/common"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "binary <name>",
			Short: "Check if binary is available",
			Args:  cobra.ExactArgs(1),
			Run: func(c *cobra.Command, args []string) {
				if !common.IsBinaryAvailable(c.Context(), args[0]) {
					os.Exit(1)
				}
			},
		})
	})
}
