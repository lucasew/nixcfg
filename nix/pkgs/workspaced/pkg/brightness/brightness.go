package brightness

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"workspaced/pkg/brightness/api"
	"workspaced/pkg/driver"
	"workspaced/pkg/exec"
	"workspaced/pkg/logging"
	"workspaced/pkg/notification"
)

func SetBrightness(ctx context.Context, arg string) error {
	d, err := driver.Get[api.Driver](ctx)
	if err != nil {
		return err
	}
	if err := d.SetBrightness(ctx, arg); err != nil {
		return err
	}
	return ShowStatus(ctx)
}

func ShowStatus(ctx context.Context) error {
	out, err := exec.RunCmd(ctx, "brightnessctl", "-m").Output()
	if err != nil {
		return fmt.Errorf("failed to get brightness status: %w", err)
	}

	lines := strings.SplitSeq(strings.TrimSpace(string(out)), "\n")
	for line := range lines {
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

		n := notification.Notification{
			ID:          notification.StatusNotificationID,
			Title:       "Brightness",
			Message:     devname,
			Icon:        "display-brightness",
			Progress:    float64(l) / 100.0,
			HasProgress: true,
		}

		if err := notification.Notify(ctx, &n); err != nil {
			logging.ReportError(ctx, fmt.Errorf("failed to send brightness notification: %w", err))
		}
		logging.GetLogger(ctx).Info("brightness updated", "device", devname, "level", level)
	}
	return nil
}
