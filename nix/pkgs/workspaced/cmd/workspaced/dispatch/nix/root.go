package nix

import (
	"workspaced/pkg/common/registry"

	"github.com/spf13/cobra"
)

var Registry registry.CommandRegistry

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nix",
		Short: "Nix operations",
	}
	return Registry.GetCommand(cmd)
}
