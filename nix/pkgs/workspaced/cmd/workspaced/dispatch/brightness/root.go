package brightness

import (
	"fmt"
	"strings"
	"workspaced/cmd/workspaced/dispatch/common"

	"github.com/spf13/cobra"
)

const brightnessNotificationID = "28419485"

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
		action := a // capture loop var
		subCmd := &cobra.Command{
			Use:   action.name,
			Short: action.short,
			RunE: func(c *cobra.Command, args []string) error {
				return runBrightnessAction(c, action.name, action.arg)
			},
		}
		Command.AddCommand(subCmd)
	}
}

func runBrightnessAction(c *cobra.Command, name, arg string) error {
	if arg != "" {
		if err := common.RunCmd(c, "brightnessctl", "s", arg).Run(); err != nil {
			return fmt.Errorf("failed to set brightness: %w", err)
		}
	}

	// Get status using machine-readable output: brightnessctl -m
	out, err := common.RunCmd(c, "brightnessctl", "-m").Output()
	if err != nil {
		return fmt.Errorf("failed to get brightness status: %w", err)
	}

	// Format: devname,devclass,current,level,max
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) < 5 {
			continue
		}
		devname := parts[0]
		level := parts[3] // current level in percent

		notifyArgs := []string{
			fmt.Sprintf("☀️ %s", devname),
			"-h", fmt.Sprintf("int:value:%s", strings.TrimSuffix(level, "%")),
			"-r", brightnessNotificationID,
		}

		if err := common.RunCmd(c, "notify-send", notifyArgs...).Run(); err != nil {
			return fmt.Errorf("failed to send notification: %w", err)
		}
		fmt.Printf("Brightness (%s): %s\n", devname, level)
	}

	return nil
}
