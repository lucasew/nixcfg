package brightness

import (
	"workspaced/pkg/drivers/brightness"

	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "brightness",
	Short: "Control screen brightness",
}

func init() {
	actions := []struct {
		name  string
		short string
		arg   string
	}{
		{"up", "Increase brightness", "+5%"},
		{"down", "Decrease brightness", "5%-"},
		{"show", "Show current brightness", ""},
		{"status", "Show current brightness (alias for show)", ""},
	}

	for _, a := range actions {
		action := a
		subCmd := &cobra.Command{
			Use:   action.name,
			Short: action.short,
			RunE: func(c *cobra.Command, args []string) error {
				if action.arg == "" {
					return brightness.ShowStatus(c.Context())
				}
				return brightness.SetBrightness(c.Context(), action.arg)
			},
		}
		Command.AddCommand(subCmd)
	}
}
