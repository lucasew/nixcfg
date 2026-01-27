package screen

import (
	"workspaced/pkg/common/registry"

	"github.com/spf13/cobra"
)

var Registry registry.CommandRegistry

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "screen",
		Short: "Screen and power management",
	}
	return Registry.GetCommand(cmd)
}
