package tray

import (
	"context"
	"fmt"
	"image"
	"os"
	"sync"
)

// MenuItem represents an item in the tray menu.
type MenuItem struct {
	Label    string
	Callback func()
	Children []MenuItem
}

// State represents the desired state of the tray.
type State struct {
	Title string
	Icon  image.Image
	Menu  []MenuItem
}

// Driver is the interface for tray implementations.
// It follows a declarative pattern where the state is set and the driver updates the UI.
type Driver interface {
	Run(ctx context.Context) error
	SetState(s State)
	Close()
}

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]func() Driver)
)

// Register makes a tray driver available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, driver func() Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("tray: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("tray: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

// Get returns the driver with the given name.
func Get(name string) (Driver, error) {
	driversMu.RLock()
	factory, ok := drivers[name]
	driversMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("tray: unknown driver %q (forgotten import?)", name)
	}
	return factory(), nil
}

// GetDefault returns the appropriate tray driver for the current environment.
func GetDefault() (Driver, error) {
	// For now, prioritize DBus if available
	if os.Getenv("DBUS_SESSION_BUS_ADDRESS") != "" {
		if d, err := Get("dbus"); err == nil {
			return d, nil
		}
	}

	// Fallback or explicit order
	// if d, err := Get("other"); err == nil { return d, nil }

	return nil, fmt.Errorf("no suitable tray driver found")
}
