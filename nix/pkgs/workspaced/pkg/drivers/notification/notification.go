package notification

import (
	"context"
	"workspaced/pkg/common"
)

type Notification struct {
	ID      uint32
	Title   string
	Message string
	Urgency string // low, normal, critical
	Icon    string
	Hint    string // e.g. int:value:50 for progress bar
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
