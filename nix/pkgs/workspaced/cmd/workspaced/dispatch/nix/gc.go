package nix

import (
	"workspaced/pkg/drivers/nix"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		cmd := &cobra.Command{
			Use:   "gc-cleanup",
			Short: "Cleanup old Nix profiles by enqueuing rm commands",
			RunE: func(cmd *cobra.Command, args []string) error {
				return nix.CleanupProfiles(cmd.Context())
			},
		}
		parent.AddCommand(cmd)
	})
}
