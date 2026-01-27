package dispatch

import (
	"workspaced/cmd/workspaced/dispatch/audio"
	"workspaced/cmd/workspaced/dispatch/brightness"
	"workspaced/cmd/workspaced/dispatch/media"
	"workspaced/cmd/workspaced/dispatch/workspace"

	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:              "dispatch",
	Short:            "Dispatch workspace commands",
	TraverseChildren: true,
}

func init() {
	Command.AddCommand(audio.Command)
	Command.AddCommand(brightness.Command)
	Command.AddCommand(media.Command)
	Command.AddCommand(workspace.Command)
}

func FindCommand(name string, args []string) (*cobra.Command, []string, error) {
	return Command.Find(append([]string{name}, args...))
}
