package brightness

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"workspaced/pkg/common"
	"workspaced/pkg/drivers/notification"
)

var n = &notification.Notification{}

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

		n.Title = fmt.Sprintf("☀️ %s", devname)
		levelVal := strings.TrimSuffix(level, "%")
		if l, err := strconv.Atoi(levelVal); err == nil {
			n.Progress = float64(l) / 100.0
		}
		n.Notify(ctx)
		common.GetLogger(ctx).Info("brightness updated", "device", devname, "level", level)
	}
	return nil
}
