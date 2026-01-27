package is

import (
	"os"
	"workspaced/pkg/common"

	"github.com/spf13/cobra"
)

var inStoreCmd = &cobra.Command{
	Use:   "in-store",
	Short: "Check if dotfiles are in nix store",
	Run: func(c *cobra.Command, args []string) {
		if !common.IsInStore() {
			os.Exit(1)
		}
	},
}

func init() {
	Command.AddCommand(inStoreCmd)
}
