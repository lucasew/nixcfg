package power

import (
	"workspaced/pkg/drivers/power"

	"github.com/spf13/cobra"
)

var wakeCmd = &cobra.Command{
	Use:   "wake <host>",
	Short: "Send Wake-on-LAN magic packet",
	Args:  cobra.ExactArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return power.Wake(c.Context(), args[0])
	},
}

func init() {
	Command.AddCommand(wakeCmd)
}
