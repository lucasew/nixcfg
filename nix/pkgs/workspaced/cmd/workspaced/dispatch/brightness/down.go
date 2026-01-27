package brightness

import (
	"workspaced/pkg/drivers/brightness"

	"github.com/spf13/cobra"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Decrease brightness",
	RunE: func(c *cobra.Command, args []string) error {
		return brightness.SetBrightness(c.Context(), "5%-")
	},
}

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(downCmd)
	})
}
