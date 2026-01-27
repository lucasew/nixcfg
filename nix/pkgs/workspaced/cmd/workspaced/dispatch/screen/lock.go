package screen

import (
	"workspaced/pkg/drivers/screen"

	"github.com/spf13/cobra"
)

var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Lock the screen and turn it off",
	RunE: func(c *cobra.Command, args []string) error {
		return screen.Lock(c.Context())
	},
}

func init() {
	Command.AddCommand(lockCmd)
}
