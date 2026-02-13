package screenshot

import (
	"workspaced/pkg/registry"

	"github.com/spf13/cobra"
)

var Registry registry.CommandRegistry

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "screenshot",
		Short: "Screen capture management",
	}
	return Registry.GetCommand(cmd)
}
