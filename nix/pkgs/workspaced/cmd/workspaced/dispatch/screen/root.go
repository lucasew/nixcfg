package screen

import (
	"workspaced/pkg/drivers/screen"

	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "screen",
	Short: "Screen and power management",
}

var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Lock the screen and turn it off",
	RunE: func(c *cobra.Command, args []string) error {
		return screen.Lock(c.Context())
	},
}

var offCmd = &cobra.Command{
	Use:   "off",
	Short: "Turn off the screen (DPMS)",
	RunE: func(c *cobra.Command, args []string) error {
		return screen.SetDPMS(c.Context(), false)
	},
}

var onCmd = &cobra.Command{
	Use:   "on",
	Short: "Turn on the screen (DPMS)",
	RunE: func(c *cobra.Command, args []string) error {
		return screen.SetDPMS(c.Context(), true)
	},
}

var toggleCmd = &cobra.Command{
	Use:   "toggle",
	Short: "Toggle screen state (DPMS)",
	RunE: func(c *cobra.Command, args []string) error {
		return screen.ToggleDPMS(c.Context())
	},
}

func init() {
	Command.AddCommand(lockCmd)
	Command.AddCommand(offCmd)
	Command.AddCommand(onCmd)
	Command.AddCommand(toggleCmd)
}
