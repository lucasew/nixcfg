package media

import (
	"workspaced/pkg/drivers/media"

	"github.com/spf13/cobra"
)

var playPauseCmd = &cobra.Command{
	Use:   "play-pause",
	Short: "Play or pause media",
	RunE: func(cmd *cobra.Command, args []string) error {
		return media.RunAction(cmd.Context(), "play-pause")
	},
}

func init() {
	Command.AddCommand(playPauseCmd)
}
