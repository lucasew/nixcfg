package power

import (
	"workspaced/pkg/driver/power"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "shutdown",
			Short: "Power off the system",
			RunE: func(c *cobra.Command, args []string) error {
				return power.Shutdown(c.Context())
			},
		})
	})
}
