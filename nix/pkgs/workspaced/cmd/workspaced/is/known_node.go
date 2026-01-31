package is

import (
	"fmt"
	"workspaced/pkg/common"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "known-node",
			Short: "Check if host is a known node",
			RunE: func(c *cobra.Command, args []string) error {
				logger := common.GetLogger(c.Context())
				if common.IsRiverwood() {
					logger.Info("riverwood")
					return nil
				}
				if common.IsWhiterun() {
					logger.Info("whiterun")
					return nil
				}
				if common.IsPhone() {
					logger.Info("phone")
					return nil
				}
				return fmt.Errorf("unknown node")
			},
		})
	})
}
