package is

import (
	"os"
	"workspaced/pkg/common"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "known-node",
			Short: "Check if host is a known node",
			Run: func(c *cobra.Command, args []string) {
				logger := common.GetLogger(c.Context())
				if common.IsRiverwood() {
					logger.Info("riverwood")
					return
				}
				if common.IsWhiterun() {
					logger.Info("whiterun")
					return
				}
				if common.IsPhone() {
					logger.Info("phone")
					return
				}
				os.Exit(1)
			},
		})
	})
}
