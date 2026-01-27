package audio

import (
	"workspaced/pkg/drivers/audio"

	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "audio",
	Short: "Control audio volume",
}

func init() {
	actions := []struct {
		name  string
		short string
		arg   string
	}{
		{"up", "Increase volume", "+5%"},
		{"down", "Decrease volume", "-5%"},
		{"mute", "Toggle mute", "toggle"},
		{"show", "Show current volume", ""},
		{"status", "Show current volume (alias for show)", ""},
	}

	for _, a := range actions {
		action := a
		subCmd := &cobra.Command{
			Use:   action.name,
			Short: action.short,
			RunE: func(c *cobra.Command, args []string) error {
				if action.arg == "" {
					return audio.ShowStatus(c.Context())
				}
				return audio.SetVolume(c.Context(), action.arg)
			},
		}
		Command.AddCommand(subCmd)
	}
}
