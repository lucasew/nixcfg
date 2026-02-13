package template

import (
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Template management commands",
	}

	cmd.AddCommand(getMaterializeCommand())

	return cmd
}
