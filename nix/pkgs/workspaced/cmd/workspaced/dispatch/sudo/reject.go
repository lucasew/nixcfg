package sudo

import (
	"workspaced/pkg/sudo"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		cmd := &cobra.Command{
			Use:   "reject <slug>",
			Short: "Reject a pending command",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return sudo.Remove(args[0])
			},
		}
		parent.AddCommand(cmd)
	})
}
