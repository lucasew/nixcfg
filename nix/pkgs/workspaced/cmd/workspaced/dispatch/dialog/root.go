package dialog

import (
	"workspaced/pkg/registry"

	"github.com/spf13/cobra"
)

var Registry registry.CommandRegistry

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dialog",
		Short: "Interactive dialogs",
	}
	return Registry.GetCommand(cmd)
}
