package media

import (
	"workspaced/pkg/drivers/media"

	"github.com/spf13/cobra"
)

var previousCmd = &cobra.Command{
	Use:   "previous",
	Short: "Previous media",
	RunE: func(cmd *cobra.Command, args []string) error {
		return media.RunAction(cmd.Context(), "previous")
	},
}

func init() {
	Command.AddCommand(previousCmd)
}
