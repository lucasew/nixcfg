package audio

import (
	"fmt"
	"strings"
	"workspaced/cmd/workspaced/dispatch/common"

	"github.com/spf13/cobra"
)

const audioNotificationID = "25548177"

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
		action := a // capture loop var
		subCmd := &cobra.Command{
			Use:   action.name,
			Short: action.short,
			RunE: func(c *cobra.Command, args []string) error {
				return runAudioAction(c, action.name, action.arg)
			},
		}
		Command.AddCommand(subCmd)
	}
}

func runAudioAction(c *cobra.Command, name, arg string) error {
	sink := "@DEFAULT_SINK@"

	if arg != "" {
		if err := common.RunCmd(c, "pactl", "set-sink-volume", sink, arg).Run(); err != nil {
			return fmt.Errorf("failed to set volume: %w", err)
		}
	}

	// Get current level
	out, err := common.RunCmd(c, "pactl", "get-sink-volume", sink).Output()
	if err != nil {
		return fmt.Errorf("failed to get volume: %w", err)
	}

	// Parse level
	level := "0"
	parts := strings.Fields(string(out))
	for _, p := range parts {
		if strings.Contains(p, "%") {
			level = strings.Trim(p, "%")
			break
		}
	}

	emoji := "🔊"
	if level == "0" {
		emoji = "🔇"
	}

	// Get default sink name for notification
	sinkNameOut, _ := common.RunCmd(c, "pactl", "get-default-sink").Output()
	sinkName := strings.TrimSpace(string(sinkNameOut))

	notifyArgs := []string{
		fmt.Sprintf("%s Volume", emoji),
		sinkName,
		"-h", fmt.Sprintf("int:value:%s", level),
		"-i", "audio",
		"-r", audioNotificationID,
	}

	if err := common.RunCmd(c, "notify-send", notifyArgs...).Run(); err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	fmt.Printf("Volume: %s%%\n", level)
	return nil
}
