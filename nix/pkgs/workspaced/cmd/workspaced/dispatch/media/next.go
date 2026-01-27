package media

import (
	"workspaced/pkg/drivers/media"

	"github.com/spf13/cobra"
)

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Next media",
	RunE: func(cmd *cobra.Command, args []string) error {
		return media.RunAction(cmd.Context(), "next")
	},
}

func init() {
	Command.AddCommand(nextCmd)
}
