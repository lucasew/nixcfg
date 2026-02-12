package brightness

import (
	"context"
	"workspaced/pkg/driver"

	"workspaced/pkg/notification"
)

func ShowStatus(ctx context.Context) error {
	d, err := driver.Get[Driver](ctx)
	if err != nil {
		return err
	}

	status, err := d.Status(ctx)
	if err != nil {
		return err
	}

	n := notification.Notification{
		ID:          notification.StatusNotificationID,
		Title:       "Brightness",
		Message:     status.Name,
		Icon:        "display-brightness",
		Progress:    float64(status.Brightness),
		HasProgress: true,
	}
	return notification.Notify(ctx, &n)

}
