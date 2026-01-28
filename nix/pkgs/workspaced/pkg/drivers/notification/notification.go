package notification

import (
	"context"
	"workspaced/pkg/common"
)

type Notification struct {
	ID       uint32
	Title    string
	Message  string
	Urgency  string // low, normal, critical
	Icon     string
	Progress int // 0-100, 0 means no progress
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
