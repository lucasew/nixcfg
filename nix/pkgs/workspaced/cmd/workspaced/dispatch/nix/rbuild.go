package nix

import (
	"workspaced/pkg/drivers/nix"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		var target string
		var copyBack bool
		var useNom bool

		cmd := &cobra.Command{
			Use:   "rbuild <ref>",
			Short: "Performs a remote build of a Nix flake reference",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				ctx := cmd.Context()
				ref := args[0]

				resultPath, err := nix.RemoteBuild(ctx, ref, target, copyBack)
				if err != nil {
					return err
				}

				cmd.Println(resultPath)
				return nil
			},
		}

		cmd.Flags().StringVarP(&target, "target", "t", "", "Remote host to build on (default: whiterun)")
		cmd.Flags().BoolVar(&copyBack, "copy-back", true, "Copy result back to local store")
		cmd.Flags().BoolVar(&useNom, "nom", true, "Use nix-output-monitor (nom)")

		parent.AddCommand(cmd)
	})
}
