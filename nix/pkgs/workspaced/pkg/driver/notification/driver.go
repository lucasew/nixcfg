package notification

import (
	"context"
)

// StatusNotificationID is the reserved ID for system status notifications (e.g. volume, brightness).
// Reusing this ID allows updating an existing notification instead of creating a new one.
const (
	StatusNotificationID   uint32 = 100
	NixBuildNotificationID uint32 = 101
	BackupNotificationID   uint32 = 102
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
