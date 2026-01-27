package workspace

import (
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "workspace",
	Short: "Workspace management commands",
}

func init() {
	Command.PersistentFlags().Bool("move", false, "Move container to workspace")
}
