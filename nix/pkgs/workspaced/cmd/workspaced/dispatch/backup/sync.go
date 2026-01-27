package backup

import (
	"workspaced/pkg/drivers/git"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "sync",
			Short: "QuickSync personal git repositories",
			RunE: func(c *cobra.Command, args []string) error {
				return git.QuickSync(c.Context())
			},
		})
	})
}
