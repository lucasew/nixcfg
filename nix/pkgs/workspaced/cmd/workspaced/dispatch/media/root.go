package media

import (
	"workspaced/pkg/drivers/media"

	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "media",
	Short: "Control media playback",
}

func init() {
	cmds := []string{"play-pause", "next", "previous", "stop", "show"}
	for _, c := range cmds {
		cmdName := c
		subCmd := &cobra.Command{
			Use:   cmdName,
			Short: cmdName + " media",
			RunE: func(cmd *cobra.Command, args []string) error {
				return media.RunAction(cmd.Context(), cmdName)
			},
		}
		Command.AddCommand(subCmd)
	}
}
