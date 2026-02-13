package launch

import (
	"workspaced/pkg/driver/terminal"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "launch",
		Short: "Launch applications",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "terminal",
		Short: "Launch the preferred terminal",
		RunE: func(c *cobra.Command, args []string) error {
			opts := terminal.Options{
				Title: "Terminal",
			}
			if len(args) > 0 {
				opts.Command = args[0]
				opts.Args = args[1:]
			}
			return terminal.Open(c.Context(), opts)
		},
	})

	cmd.AddCommand(newWebappCommand())

	return cmd
}
