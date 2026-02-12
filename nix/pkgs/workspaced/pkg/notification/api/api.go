package api

import (
	"context"
	"workspaced/pkg/driver"
)

// StatusNotificationID is the reserved ID for system status notifications (e.g. volume, brightness).
// Reusing this ID allows updating an existing notification instead of creating a new one.
const (
	StatusNotificationID uint32 = 100
)

// Notification represents a system notification.
type Notification struct {
	ID          uint32
	Title       string
	Message     string
	Urgency     string // low, normal, critical
	Icon        string
	Progress    float64 // 0.0-1.0
	HasProgress bool
}

type Driver interface {
	Notify(ctx context.Context, n *Notification) error
}

func (n *Notification) Notify(ctx context.Context) error {
	d, err := driver.Get[Driver](ctx)
	if err != nil {
		return err
	}
	return d.Notify(ctx, n)
}
