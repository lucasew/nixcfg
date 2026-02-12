package brightness

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"workspaced/pkg/driver"
	"workspaced/pkg/driver/brightness"
	"workspaced/pkg/exec"
	"workspaced/pkg/logging"
	"workspaced/pkg/notification"
)

func SetBrightness(ctx context.Context, arg string) error {
	d, err := driver.Get[brightness.Driver](ctx)
	if err != nil {
		return err
	}
	if err := d.SetBrightness(ctx, arg); err != nil {
		return err
	}
	return ShowStatus(ctx)
}

func ShowStatus(ctx context.Context) error {
		d, err := driver.Get[brightness.Driver](ctx)
		status, err := d.Status(ctx)
		n := notification.Notification{
			ID:          notification.StatusNotificationID,
			Title:       "Brightness",
			Message:     status.Name,
			Icon:        "display-brightness",
			Progress:    float64(status.Brightness),
			HasProgress: true,
		}

		if err := notification.Notify(ctx, &n); err != nil {
			logging.ReportError(ctx, fmt.Errorf("failed to send brightness notification: %w", err))
		}
		logging.GetLogger(ctx).Info("brightness updated", "device", devname, "level", level)
	}
	return nil
}
