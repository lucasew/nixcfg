package backup

import (
	"workspaced/pkg/common/registry"

	"github.com/spf13/cobra"
)

var Registry registry.CommandRegistry

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Data backup and synchronization",
	}
	return Registry.GetCommand(cmd)
}
