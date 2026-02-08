package power

import (
	"workspaced/pkg/power"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "wake <host>",
			Short: "Send Wake-on-LAN magic packet",
			Args:  cobra.ExactArgs(1),
			RunE: func(c *cobra.Command, args []string) error {
				return power.Wake(c.Context(), args[0])
			},
		})
	})
}
