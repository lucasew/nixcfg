package is

import (
	"fmt"
	"workspaced/pkg/common"
	"workspaced/pkg/host"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "known-node",
			Short: "Check if host is a known node",
			RunE: func(c *cobra.Command, args []string) error {
				logger := common.GetLogger(c.Context())
				if host.IsRiverwood() {
					logger.Info("riverwood")
					return nil
				}
				if host.IsWhiterun() {
					logger.Info("whiterun")
					return nil
				}
				if host.IsPhone() {
					logger.Info("phone")
					return nil
				}
				return fmt.Errorf("unknown node")
			},
		})
	})
}
