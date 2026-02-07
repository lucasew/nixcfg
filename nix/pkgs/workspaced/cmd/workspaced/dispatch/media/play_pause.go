package media

import (
	"workspaced/pkg/media"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "play-pause",
			Short: "Play or pause media",
			RunE: func(cmd *cobra.Command, args []string) error {
				return media.RunAction(cmd.Context(), "play-pause")
			},
		})
	})
}
