package is

import (
	"fmt"
	"os"
	"workspaced/pkg/common"

	"github.com/spf13/cobra"
)

var knownNodeCmd = &cobra.Command{
	Use:   "known-node",
	Short: "Check if host is a known node",
	Run: func(c *cobra.Command, args []string) {
		if common.IsRiverwood() {
			fmt.Println("riverwood")
			return
		}
		if common.IsWhiterun() {
			fmt.Println("whiterun")
			return
		}
		if common.IsPhone() {
			fmt.Println("phone")
			return
		}
		os.Exit(1)
	},
}

func init() {
	Command.AddCommand(knownNodeCmd)
}
