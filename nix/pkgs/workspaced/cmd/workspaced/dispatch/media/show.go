package media

import (
	"workspaced/pkg/drivers/media"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show media metadata",
	RunE: func(cmd *cobra.Command, args []string) error {
		return media.RunAction(cmd.Context(), "show")
	},
}

func init() {
	Command.AddCommand(showCmd)
}
