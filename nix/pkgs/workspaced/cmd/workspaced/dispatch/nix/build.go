package nix

import (
	"workspaced/pkg/drivers/nix"

	"github.com/spf13/cobra"
)

// Build sub-command
func init() {
	Registry.Register(func(parent *cobra.Command) {
		var noCache bool
		cmd := &cobra.Command{
			Use:   "build <ref>",
			Short: "Build a Nix flake reference with RAM caching",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				ref := args[0]
				path, err := nix.Build(cmd.Context(), ref, !noCache)
				if err != nil {
					return err
				}
				cmd.Println(path)
				return nil
			},
		}
		cmd.Flags().BoolVar(&noCache, "no-cache", false, "Disable RAM cache")
		parent.AddCommand(cmd)
	})
}
