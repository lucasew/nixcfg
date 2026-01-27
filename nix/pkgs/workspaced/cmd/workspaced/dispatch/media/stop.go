package media

import (
	"workspaced/pkg/drivers/media"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "stop",
			Short: "Stop media",
			RunE: func(cmd *cobra.Command, args []string) error {
				return media.RunAction(cmd.Context(), "stop")
			},
		})
	})
}
