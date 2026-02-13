package history

import (
	"workspaced/cmd/workspaced/dispatch/history"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return history.GetCommand()
}
