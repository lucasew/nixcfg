package power

import (
	"workspaced/pkg/common/registry"

	"github.com/spf13/cobra"
)

var Registry registry.CommandRegistry

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "power",
		Short: "Power management commands",
	}
	return Registry.GetCommand(cmd)
}
