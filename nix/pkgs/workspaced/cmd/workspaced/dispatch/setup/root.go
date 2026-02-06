package setup

import (
	"workspaced/pkg/registry"

	"github.com/spf13/cobra"
)

var Registry registry.CommandRegistry

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Environment setup commands",
	}
	return Registry.GetCommand(cmd)
}
