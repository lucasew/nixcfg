package is

import (
	"os"
	"workspaced/pkg/common"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "phone",
			Short: "Check if environment is a phone (Termux)",
			Run: func(c *cobra.Command, args []string) {
				if !common.IsPhone() {
					os.Exit(1)
				}
			},
		})
	})
}
