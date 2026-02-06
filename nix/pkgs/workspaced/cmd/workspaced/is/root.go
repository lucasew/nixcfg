package is

import (
	"workspaced/pkg/registry"

	"github.com/spf13/cobra"
)

var Registry registry.CommandRegistry

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "is",
		Short: "Environment detection commands",
	}
	return Registry.GetCommand(cmd)
}
