package apply

import (
	"workspaced/pkg/common/registry"

	"github.com/spf13/cobra"
)

var Registry registry.CommandRegistry

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Declaratively apply system and user configurations",
	}
	return Registry.GetCommand(cmd)
}
