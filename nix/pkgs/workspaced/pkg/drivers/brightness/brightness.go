package brightness

import (
	"context"
	"fmt"
	"strings"
	"workspaced/pkg/common"
)

const NotificationID = "28419485"

func SetBrightness(ctx context.Context, arg string) error {
	if arg != "" {
		if err := common.RunCmd(ctx, "brightnessctl", "s", arg).Run(); err != nil {
			return fmt.Errorf("failed to set brightness: %w", err)
		}
	}
	return ShowStatus(ctx)
}

func ShowStatus(ctx context.Context) error {
	out, err := common.RunCmd(ctx, "brightnessctl", "-m").Output()
	if err != nil {
		return fmt.Errorf("failed to get brightness status: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) < 5 {
			continue
		}
		devname := parts[0]
		level := parts[3]

		notifyArgs := []string{
			fmt.Sprintf("☀️ %s", devname),
			"-h", fmt.Sprintf("int:value:%s", strings.TrimSuffix(level, "%")),
			"-r", NotificationID,
		}
		common.RunCmd(ctx, "notify-send", notifyArgs...).Run()
	}
	return nil
}
