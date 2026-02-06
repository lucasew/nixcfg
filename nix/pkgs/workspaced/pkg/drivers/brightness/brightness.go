package brightness

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"workspaced/pkg/drivers/notification"
	"workspaced/pkg/exec"
	"workspaced/pkg/logging"
)

func SetBrightness(ctx context.Context, arg string) error {
	if arg != "" {
		if err := exec.RunCmd(ctx, "brightnessctl", "s", arg).Run(); err != nil {
			return fmt.Errorf("failed to set brightness: %w", err)
		}
	}
	return ShowStatus(ctx)
}

func ShowStatus(ctx context.Context) error {
	out, err := exec.RunCmd(ctx, "brightnessctl", "-m").Output()
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

		levelVal := strings.TrimSuffix(level, "%")
		l, err := strconv.Atoi(levelVal)
		if err != nil {
			continue
		}

		n := &notification.Notification{
			ID:          notification.StatusNotificationID,
			Title:       "Brightness",
			Message:     devname,
			Icon:        "display-brightness",
			Progress:    float64(l) / 100.0,
			HasProgress: true,
		}

		_ = n.Notify(ctx)
		logging.GetLogger(ctx).Info("brightness updated", "device", devname, "level", level)
	}
	return nil
}
