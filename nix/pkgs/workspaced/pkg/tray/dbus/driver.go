package dbus

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/godbus/dbus/v5"
	"workspaced/pkg/logging"
)

// MenuItem represents a single item in the tray menu.
type MenuItem struct {
	ID       int32
	Label    string
	Callback func()
	Children []*MenuItem
}

type Tray struct {
	mu        sync.RWMutex
	ID        string
	Title     string
	Icon      string
	MenuItems []*MenuItem
	nextID    int32

	conn *dbus.Conn
	sni  *StatusNotifierItem
	menu *DBusMenu

	closeOnce sync.Once
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewTray(id, title, icon string) *Tray {
	ctx, cancel := context.WithCancel(context.Background())
	return &Tray{
		ID:     id,
		Title:  title,
		Icon:   icon,
		nextID: 1,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (t *Tray) AddMenuItem(label string, callback func()) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.MenuItems = append(t.MenuItems, &MenuItem{
		ID:       t.nextID,
		Label:    label,
		Callback: callback,
	})
	t.nextID++
	if t.menu != nil {
		// Signal layout update if running
		t.menu.EmitLayoutUpdated()
	}
}

func (t *Tray) Run(ctx context.Context) error {
	// Use ConnectSessionBus to get a private connection that we can close safely
	// without affecting other parts of the application sharing the default bus.
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return fmt.Errorf("failed to connect to session bus: %w", err)
	}
	t.conn = conn

	t.sni = NewStatusNotifierItem(t)
	t.menu = NewDBusMenu(t)

	// Export objects
	if err := t.sni.Export(t.conn, "/StatusNotifierItem"); err != nil {
		return fmt.Errorf("failed to export SNI: %w", err)
	}

	if err := t.menu.Export(t.conn, "/MenuBar"); err != nil {
		return fmt.Errorf("failed to export DBusMenu: %w", err)
	}

	// Request name
	serviceName := fmt.Sprintf("org.kde.StatusNotifierItem-%d-1", os.Getpid())
	reply, err := t.conn.RequestName(serviceName, dbus.NameFlagDoNotQueue)
	if err != nil {
		return fmt.Errorf("failed to request name: %w", err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		return fmt.Errorf("name already taken")
	}

	// Register with watcher
	watcher := t.conn.Object("org.kde.StatusNotifierWatcher", "/StatusNotifierWatcher")
	call := watcher.Call("org.kde.StatusNotifierWatcher.RegisterStatusNotifierItem", 0, serviceName)
	if call.Err != nil {
		logging.ReportError(ctx, fmt.Errorf("failed to register with watcher: %w", call.Err))
	}

	<-ctx.Done()
	return nil
}

func (t *Tray) Close() {
	t.closeOnce.Do(func() {
		t.cancel()
		if t.conn != nil {
			if err := t.conn.Close(); err != nil {
				logging.ReportError(context.Background(), err)
			}
		}
	})
}
