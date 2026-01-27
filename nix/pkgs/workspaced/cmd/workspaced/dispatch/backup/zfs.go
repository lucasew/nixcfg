package backup

import (
	"workspaced/pkg/drivers/backup"

	"github.com/spf13/cobra"
)

var zfsCmd = &cobra.Command{
	Use:   "zfs",
	Short: "Replicate ZFS datasets",
	RunE: func(c *cobra.Command, args []string) error {
		return backup.ReplicateZFS(c.Context())
	},
}

func init() {
	Command.AddCommand(zfsCmd)
}
