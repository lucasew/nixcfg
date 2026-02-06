package menu

import (
	"workspaced/pkg/registry"

	"github.com/spf13/cobra"
)

var Registry registry.CommandRegistry

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "menu",
		Short: "Interactive menus",
	}
	return Registry.GetCommand(cmd)
}
