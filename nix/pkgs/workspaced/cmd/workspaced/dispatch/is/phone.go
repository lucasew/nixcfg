package is

import (
	"os"
	"workspaced/pkg/common"

	"github.com/spf13/cobra"
)

var phoneCmd = &cobra.Command{
	Use:   "phone",
	Short: "Check if environment is a phone (Termux)",
	Run: func(c *cobra.Command, args []string) {
		if !common.IsPhone() {
			os.Exit(1)
		}
	},
}

func init() {
	Command.AddCommand(phoneCmd)
}
