package menu

import (
	"workspaced/pkg/common"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "launch",
			Short: "Application launcher",
			RunE: func(c *cobra.Command, args []string) error {
				return common.RunCmd(c.Context(), "rofi-launch").Run()
			},
		})
	})
}
