package open

import (
	"workspaced/pkg/driver/opener"

	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open [target]",
		Short: "Open a file or URL using the preferred opener",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return opener.Open(c.Context(), args[0])
		},
	}

	return cmd
}
