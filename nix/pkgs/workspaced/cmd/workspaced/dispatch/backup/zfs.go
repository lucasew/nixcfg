package backup

import (
	"workspaced/pkg/drivers/backup"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "zfs",
			Short: "Replicate ZFS datasets",
			RunE: func(c *cobra.Command, args []string) error {
				return backup.ReplicateZFS(c.Context())
			},
		})
	})
}
