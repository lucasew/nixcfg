package doctor

import (
	"workspaced/pkg/driver"

	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "doctor",
	Short: "Check status of all registered drivers",
	Run: func(cmd *cobra.Command, args []string) {
		driver.Doctor(cmd.Context())
	},
}
