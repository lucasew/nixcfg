package media

import (
	"workspaced/pkg/drivers/media"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop media",
	RunE: func(cmd *cobra.Command, args []string) error {
		return media.RunAction(cmd.Context(), "stop")
	},
}

func init() {
	Command.AddCommand(stopCmd)
}
