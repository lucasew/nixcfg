package notification

import (
	"workspaced/pkg/notification/api"
)

// StatusNotificationID is the reserved ID for system status notifications (e.g. volume, brightness).
// Reusing this ID allows updating an existing notification instead of creating a new one.
const StatusNotificationID = api.StatusNotificationID

// Notification represents a system notification.
type Notification = api.Notification

// Notifier is the interface for sending notifications.
type Notifier = api.Driver
