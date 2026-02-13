package workspace

import (
	"workspaced/pkg/driver/wm"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "scratchpad",
			Short: "Toggle scratchpad visibility with status info",
			RunE: func(c *cobra.Command, args []string) error {
				return wm.ToggleScratchpadWithInfo(c.Context())
			},
		})
	})
}
