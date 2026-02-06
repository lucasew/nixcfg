package is

import (
	"fmt"
	"workspaced/pkg/env"
	"workspaced/pkg/logging"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "known-node",
			Short: "Check if host is a known node",
			RunE: func(c *cobra.Command, args []string) error {
				logger := logging.GetLogger(c.Context())
				if env.IsRiverwood() {
					logger.Info("riverwood")
					return nil
				}
				if env.IsWhiterun() {
					logger.Info("whiterun")
					return nil
				}
				if env.IsPhone() {
					logger.Info("phone")
					return nil
				}
				return fmt.Errorf("unknown node")
			},
		})
	})
}
