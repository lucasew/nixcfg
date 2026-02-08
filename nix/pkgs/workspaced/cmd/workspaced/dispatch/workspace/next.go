package workspace

import (
	"workspaced/pkg/wm"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "next",
			Short: "Go to the next available workspace",
			RunE: func(c *cobra.Command, args []string) error {
				move, _ := c.Flags().GetBool("move")
				return wm.NextWorkspace(c.Context(), move)
			},
		})
	})
}
