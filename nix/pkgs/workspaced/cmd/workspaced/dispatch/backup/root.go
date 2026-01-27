package backup

import (
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "backup",
	Short: "Data backup and synchronization",
}
