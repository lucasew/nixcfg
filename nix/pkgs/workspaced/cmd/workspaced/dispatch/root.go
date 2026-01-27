package dispatch

import (
	"workspaced/cmd/workspaced/dispatch/audio"
	"workspaced/cmd/workspaced/dispatch/backup"
	"workspaced/cmd/workspaced/dispatch/brightness"
	"workspaced/cmd/workspaced/dispatch/is"
	"workspaced/cmd/workspaced/dispatch/media"
	"workspaced/cmd/workspaced/dispatch/menu"
	"workspaced/cmd/workspaced/dispatch/power"
	"workspaced/cmd/workspaced/dispatch/screen"
	"workspaced/cmd/workspaced/dispatch/screenshot"
	"workspaced/cmd/workspaced/dispatch/setup"
	"workspaced/cmd/workspaced/dispatch/wallpaper"
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
	Command.AddCommand(backup.Command)
	Command.AddCommand(brightness.Command)
	Command.AddCommand(is.Command)
	Command.AddCommand(media.Command)
	Command.AddCommand(menu.Command)
	Command.AddCommand(power.Command)
	Command.AddCommand(screen.Command)
	Command.AddCommand(screenshot.Command)
	Command.AddCommand(setup.Command)
	Command.AddCommand(wallpaper.Command)
	Command.AddCommand(workspace.Command)
}

func FindCommand(name string, args []string) (*cobra.Command, []string, error) {
	return Command.Find(append([]string{name}, args...))
}
