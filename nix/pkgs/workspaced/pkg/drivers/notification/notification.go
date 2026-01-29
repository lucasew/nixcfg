package notification

import (
	"context"
	"workspaced/pkg/common"
)

const (
	StatusNotificationID uint32 = 100
)

type Notification struct {
	ID       uint32
	Title    string
	Message  string
	Urgency  string // low, normal, critical
	Icon     string
	Progress float64 // 0.0-1.0, 0.0 means no progress
}

type Notifier interface {
	Notify(ctx context.Context, n *Notification) error
}

func (n *Notification) Notify(ctx context.Context) error {
	var notifier Notifier
	if common.IsBinaryAvailable(ctx, "termux-notification") {
		notifier = &TermuxNotifier{}
	} else {
		notifier = &NotifySendNotifier{}
	}
	return notifier.Notify(ctx, n)
}
