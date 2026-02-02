package shell

import (
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shell",
		Short: "Shell integration commands",
	}

	cmd.AddCommand(getInitCommand())

	return cmd
}
