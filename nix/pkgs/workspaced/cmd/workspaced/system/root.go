package system

import (
	"workspaced/cmd/workspaced/system/screenshot"
	"workspaced/cmd/workspaced/system/workspace"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "system",
		Short: "System and hardware management commands",
	}

	cmd.AddCommand(newAudioCommand())
	cmd.AddCommand(newBrightnessCommand())
	cmd.AddCommand(newPowerCommand())
	cmd.AddCommand(newScreenCommand())
	cmd.AddCommand(screenshot.GetCommand())
	cmd.AddCommand(workspace.GetCommand())

	return cmd
}
