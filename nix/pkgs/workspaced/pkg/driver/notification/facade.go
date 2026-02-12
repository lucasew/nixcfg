package notification

import (
	"context"
	"workspaced/pkg/driver"
)

func Notify(ctx context.Context, n *Notification) error {
	d, err := driver.Get[Driver](ctx)
	if err != nil {
		return err
	}
	return d.Notify(ctx, n)
}
