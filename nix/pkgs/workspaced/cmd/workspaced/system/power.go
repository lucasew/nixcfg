package system

import (
	"workspaced/pkg/driver/power"

	"github.com/spf13/cobra"
)

func newPowerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "power",
		Short: "Power and session management",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "lock",
		Short: "Lock the session",
		RunE: func(c *cobra.Command, args []string) error {
			return power.Lock(c.Context())
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "reboot",
		Short: "Reboot the system",
		RunE: func(c *cobra.Command, args []string) error {
			return power.Reboot(c.Context())
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "shutdown",
		Short: "Power off the system",
		RunE: func(c *cobra.Command, args []string) error {
			return power.Shutdown(c.Context())
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "suspend",
		Short: "Suspend the system",
		RunE: func(c *cobra.Command, args []string) error {
			return power.Suspend(c.Context())
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "wake <host>",
		Short: "Send Wake-on-LAN magic packet",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return power.Wake(c.Context(), args[0])
		},
	})

	return cmd
}
