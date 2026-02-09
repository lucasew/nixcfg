package tray

import (
	"fmt"
	"os"

	dbus_impl "workspaced/pkg/tray/dbus"
)

// GetDriver returns the appropriate tray driver for the current environment.
func GetDriver(id, title, icon string) (Driver, error) {
	// For now, check DBUS_SESSION_BUS_ADDRESS.
	// We can add more checks later (e.g., check for specific WM, or non-DBus systems if ever supported).
	if os.Getenv("DBUS_SESSION_BUS_ADDRESS") != "" {
		return dbus_impl.NewTray(id, title, icon), nil
	}

	return nil, fmt.Errorf("no suitable tray driver found (DBus not available)")
}
