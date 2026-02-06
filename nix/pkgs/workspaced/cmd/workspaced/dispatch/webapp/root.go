package webapp

import (
	"workspaced/pkg/registry"

	"github.com/spf13/cobra"
)

var Registry registry.CommandRegistry

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "webapp",
		Short: "Manage webapps",
	}
	return Registry.GetCommand(cmd)
}
