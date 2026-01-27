package workspace

import (
	"workspaced/pkg/common/registry"

	"github.com/spf13/cobra"
)

var Registry registry.CommandRegistry

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workspace",
		Short: "Workspace management commands",
	}
	cmd.PersistentFlags().Bool("move", false, "Move container to workspace")
	return Registry.GetCommand(cmd)
}
