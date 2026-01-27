package is

import (
	"os"
	"workspaced/pkg/common"

	"github.com/spf13/cobra"
)

var whiterunCmd = &cobra.Command{
	Use:   "whiterun",
	Short: "Check if host is whiterun",
	Run: func(c *cobra.Command, args []string) {
		if !common.IsWhiterun() {
			os.Exit(1)
		}
	},
}

func init() {
	Command.AddCommand(whiterunCmd)
}
