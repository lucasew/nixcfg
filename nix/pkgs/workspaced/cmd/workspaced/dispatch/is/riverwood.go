package is

import (
	"os"
	"workspaced/pkg/common"

	"github.com/spf13/cobra"
)

var riverwoodCmd = &cobra.Command{
	Use:   "riverwood",
	Short: "Check if host is riverwood",
	Run: func(c *cobra.Command, args []string) {
		if !common.IsRiverwood() {
			os.Exit(1)
		}
	},
}

func init() {
	Command.AddCommand(riverwoodCmd)
}
