package notification

import (
	"context"
	"workspaced/pkg/common"
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

// Notifier is the interface for sending notifications.
type Notifier interface {
	Notify(ctx context.Context, n *Notification) error
}

// Notify sends the notification using the best available backend.
// It checks for 'termux-notification' to decide between the Termux backend
// and the standard desktop 'notify-send' backend.
func (n *Notification) Notify(ctx context.Context) error {
	var notifier Notifier
	if common.IsBinaryAvailable(ctx, "termux-notification") {
		notifier = &TermuxNotifier{}
	} else {
		notifier = &NotifySendNotifier{}
	}
	return notifier.Notify(ctx, n)
}
