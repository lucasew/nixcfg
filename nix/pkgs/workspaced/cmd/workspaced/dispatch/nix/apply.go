package nix

import (
	"workspaced/pkg/drivers/nix"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		cmd := &cobra.Command{
			Use:   "apply [action]",
			Short: "Apply NixOS configuration (nixos-rebuild)",
			Args:  cobra.MaximumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				action := "switch"
				if len(args) > 0 {
					action = args[0]
				}
				flake, _ := cmd.Flags().GetString("flake")
				return nix.Rebuild(cmd.Context(), action, flake)
			},
		}
		cmd.Flags().StringP("flake", "f", "", "Flake reference to use")
		parent.AddCommand(cmd)
	})
}
