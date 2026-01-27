package backup

import (
	"workspaced/pkg/drivers/backup"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "run",
			Short: "Run full backup (rsync)",
			RunE: func(c *cobra.Command, args []string) error {
				return backup.RunFullBackup(c.Context())
			},
		})
	})
}
