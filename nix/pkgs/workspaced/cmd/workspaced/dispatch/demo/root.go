package demo

import (
	"workspaced/pkg/common/registry"

	"github.com/spf13/cobra"
)

var Registry registry.CommandRegistry

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "demo",
		Short: "Demo commands",
	}
	return Registry.GetCommand(cmd)
}
