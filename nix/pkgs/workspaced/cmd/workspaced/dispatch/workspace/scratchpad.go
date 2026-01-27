package workspace

import (
	"workspaced/pkg/drivers/wm"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "scratchpad",
			Short: "Toggle scratchpad visibility",
			RunE: func(c *cobra.Command, args []string) error {
				return wm.ToggleScratchpad(c.Context())
			},
		})
	})
}
