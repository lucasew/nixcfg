package dbus

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/godbus/dbus/v5"
	"workspaced/pkg/logging"
	"workspaced/pkg/tray"
)

func init() {
	tray.Register("dbus", func() tray.Driver {
		return NewDriver()
	})
}

type Driver struct {
	mu    sync.RWMutex
	state tray.State

	conn *dbus.Conn
	sni  *StatusNotifierItem
	menu *DBusMenu

	closeOnce sync.Once
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewDriver() *Driver {
	ctx, cancel := context.WithCancel(context.Background())
	return &Driver{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (d *Driver) SetState(s tray.State) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.state = s

	if d.menu != nil {
		d.menu.EmitLayoutUpdated()
	}
	// We could also emit NewIcon or NewTitle signals here if fully implementing SNI signals
}

func (d *Driver) Run(ctx context.Context) error {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return fmt.Errorf("failed to connect to session bus: %w", err)
	}
	d.conn = conn

	d.sni = NewStatusNotifierItem(d)
	d.menu = NewDBusMenu(d)

	// Export objects
	if err := d.sni.Export(d.conn, "/StatusNotifierItem"); err != nil {
		return fmt.Errorf("failed to export SNI: %w", err)
	}

	if err := d.menu.Export(d.conn, "/MenuBar"); err != nil {
		return fmt.Errorf("failed to export DBusMenu: %w", err)
	}

	// Request name
	serviceName := fmt.Sprintf("org.kde.StatusNotifierItem-%d-1", os.Getpid())
	reply, err := d.conn.RequestName(serviceName, dbus.NameFlagDoNotQueue)
	if err != nil {
		return fmt.Errorf("failed to request name: %w", err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		return fmt.Errorf("name already taken")
	}

	// Register with watcher
	watcher := d.conn.Object("org.kde.StatusNotifierWatcher", "/StatusNotifierWatcher")
	call := watcher.Call("org.kde.StatusNotifierWatcher.RegisterStatusNotifierItem", 0, serviceName)
	if call.Err != nil {
		logging.ReportError(ctx, fmt.Errorf("failed to register with watcher: %w", call.Err))
	}

	<-ctx.Done()
	return nil
}

func (d *Driver) Close() {
	d.closeOnce.Do(func() {
		d.cancel()
		if d.conn != nil {
			if err := d.conn.Close(); err != nil {
				logging.ReportError(context.Background(), err)
			}
		}
	})
}
