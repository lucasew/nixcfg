package dbus

import (
	"context"
	"fmt"
	"sync"
	"workspaced/pkg/driver"
	"workspaced/pkg/driver/notification"

	"github.com/godbus/dbus/v5"
)

func init() {
	driver.Register[notification.Driver](&Provider{})
}

type Provider struct{}

func (p *Provider) Name() string { return "DBus" }

func (p *Provider) CheckCompatibility(ctx context.Context) error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return fmt.Errorf("%w: failed to connect to session bus: %v", driver.ErrIncompatible, err)
	}
	// Check if the notifications service is available
	var names []string
	err = conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&names)
	if err != nil {
		return fmt.Errorf("%w: failed to list dbus names: %v", driver.ErrIncompatible, err)
	}

	found := false
	for _, name := range names {
		if name == "org.freedesktop.Notifications" {
			found = true
			break
		}
	}

	if !found {
		// Try to see if it can be started
		err = conn.BusObject().Call("org.freedesktop.DBus.StartServiceByName", 0, "org.freedesktop.Notifications", uint32(0)).Err
		if err != nil {
			return fmt.Errorf("%w: notifications service not found and could not be started: %v", driver.ErrIncompatible, err)
		}
	}

	return nil
}

func (p *Provider) New(ctx context.Context) (notification.Driver, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}
	return &Driver{
		conn: conn,
		ids:  make(map[uint32]uint32),
	}, nil
}

type Driver struct {
	conn *dbus.Conn
	mu   sync.Mutex
	ids  map[uint32]uint32 // requested ID -> actual server ID
}

func (d *Driver) Notify(ctx context.Context, n *notification.Notification) error {
	obj := d.conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")

	d.mu.Lock()
	replacesID := n.ID
	if actualID, ok := d.ids[n.ID]; ok && n.ID != 0 {
		replacesID = actualID
	}
	d.mu.Unlock()

	hints := make(map[string]dbus.Variant)

	// Urgency
	var urgency byte = 1 // normal
	switch n.Urgency {
	case "low":
		urgency = 0
	case "critical":
		urgency = 2
	}
	hints["urgency"] = dbus.MakeVariant(urgency)

	// Progress
	if n.HasProgress {
		hints["value"] = dbus.MakeVariant(int32(n.Progress * 100))
	}

	// Actions (empty for now)
	actions := []string{}

	// Expire timeout (-1 for default)
	expireTimeout := int32(-1)

	var serverID uint32
	err := obj.Call("org.freedesktop.Notifications.Notify", 0,
		"workspaced",  // app_name
		replacesID,    // replaces_id
		n.Icon,        // app_icon
		n.Title,       // summary
		n.Message,     // body
		actions,       // actions
		hints,         // hints
		expireTimeout, // expire_timeout
	).Store(&serverID)

	if err != nil {
		return err
	}

	if n.ID != 0 {
		d.mu.Lock()
		d.ids[n.ID] = serverID
		d.mu.Unlock()
	}

	return nil
}
