package sudo

import (
	"workspaced/pkg/drivers/sudo"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		var slug string
		cmd := &cobra.Command{
			Use:                "add <command> [args...]",
			Short:              "Manually add a command to the queue",
			Args:               cobra.MinimumNArgs(1),
			DisableFlagParsing: false,
			RunE: func(cmd *cobra.Command, args []string) error {
				sc := &sudo.SudoCommand{
					Slug:    slug,
					Command: args[0],
					Args:    args[1:],
				}
				return sudo.Enqueue(cmd.Context(), sc)
			},
		}
		cmd.Flags().SetInterspersed(false)
		cmd.Flags().StringVarP(&slug, "slug", "s", "", "Slug for the command")
		parent.AddCommand(cmd)
	})
}
