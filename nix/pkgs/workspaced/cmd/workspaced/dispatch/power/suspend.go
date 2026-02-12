package power

import (
	"workspaced/pkg/driver/power"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "suspend",
			Short: "Suspend the system",
			RunE: func(c *cobra.Command, args []string) error {
				return power.Suspend(c.Context())
			},
		})
	})
}
