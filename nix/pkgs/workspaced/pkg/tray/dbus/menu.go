package dbus

import (
	"context"
	"sync"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
	"workspaced/pkg/logging"
)

type DBusMenu struct {
	tray     *Tray
	revision uint32
	mu       sync.Mutex
}

func NewDBusMenu(t *Tray) *DBusMenu {
	return &DBusMenu{
		tray:     t,
		revision: 1,
	}
}

func (m *DBusMenu) Export(conn *dbus.Conn, path dbus.ObjectPath) error {
	// Export methods
	err := conn.Export(m, path, "com.canonical.dbusmenu")
	if err != nil {
		return err
	}

	// Export introspection
	n := &introspect.Node{
		Name: string(path),
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			prop.IntrospectData,
			{
				Name: "com.canonical.dbusmenu",
				Methods: []introspect.Method{
					{Name: "GetLayout", Args: []introspect.Arg{
						{Name: "parentId", Type: "i", Direction: "in"},
						{Name: "recursionDepth", Type: "i", Direction: "in"},
						{Name: "propertyNames", Type: "as", Direction: "in"},
						{Name: "revision", Type: "u", Direction: "out"},
						{Name: "layout", Type: "(ia{sv}av)", Direction: "out"},
					}},
					{Name: "GetGroupProperties", Args: []introspect.Arg{
						{Name: "ids", Type: "ai", Direction: "in"},
						{Name: "propertyNames", Type: "as", Direction: "in"},
						{Name: "updates", Type: "a(ia{sv})", Direction: "out"},
					}},
					{Name: "GetProperty", Args: []introspect.Arg{
						{Name: "id", Type: "i", Direction: "in"},
						{Name: "name", Type: "s", Direction: "in"},
						{Name: "value", Type: "v", Direction: "out"},
					}},
					{Name: "Event", Args: []introspect.Arg{
						{Name: "id", Type: "i", Direction: "in"},
						{Name: "eventId", Type: "s", Direction: "in"},
						{Name: "data", Type: "v", Direction: "in"},
						{Name: "timestamp", Type: "u", Direction: "in"},
					}},
					{Name: "AboutToShow", Args: []introspect.Arg{
						{Name: "id", Type: "i", Direction: "in"},
						{Name: "updatesNeeded", Type: "b", Direction: "out"},
					}},
				},
				Properties: []introspect.Property{
					{Name: "Version", Type: "u", Access: "read"},
					{Name: "TextDirection", Type: "s", Access: "read"},
					{Name: "Status", Type: "s", Access: "read"},
					{Name: "IconThemePath", Type: "as", Access: "read"},
				},
				Signals: []introspect.Signal{
					{Name: "ItemsPropertiesUpdated", Args: []introspect.Arg{
						{Name: "updates", Type: "a(ia{sv})"},
						{Name: "removedProps", Type: "a(ias)"},
					}},
					{Name: "LayoutUpdated", Args: []introspect.Arg{
						{Name: "revision", Type: "u"},
						{Name: "parent", Type: "i"},
					}},
					{Name: "ItemActivationRequested", Args: []introspect.Arg{
						{Name: "id", Type: "i"},
						{Name: "timestamp", Type: "u"},
					}},
				},
			},
		},
	}
	err = conn.Export(introspect.NewIntrospectable(n), path, "org.freedesktop.DBus.Introspectable")
	if err != nil {
		return err
	}

	// Export properties
	propsMap := map[string]interface{}{
		"Version":       uint32(3),
		"TextDirection": "ltr",
		"Status":        "normal",
		"IconThemePath": []string{},
	}

	convertedProps := make(map[string]*prop.Prop)
	for k, v := range propsMap {
		convertedProps[k] = &prop.Prop{
			Value:    v,
			Writable: false,
			Emit:     prop.EmitTrue,
		}
	}

	_, err = prop.Export(conn, path, prop.Map{
		"com.canonical.dbusmenu": convertedProps,
	})
	return err
}

// LayoutNode represents (ia{sv}av)
// i: id
// a{sv}: properties
// av: children (variants of LayoutNode)
type LayoutNode struct {
	ID         int32
	Properties map[string]dbus.Variant
	Children   []dbus.Variant
}

func (m *DBusMenu) GetLayout(parentId int32, recursionDepth int32, propertyNames []string) (uint32, LayoutNode, *dbus.Error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Acquire read lock on tray to access MenuItems safely
	m.tray.mu.RLock()
	defer m.tray.mu.RUnlock()

	root := LayoutNode{
		ID:         0,
		Properties: map[string]dbus.Variant{"children-display": dbus.MakeVariant("submenu")},
		Children:   []dbus.Variant{},
	}

	for _, item := range m.tray.MenuItems {
		child := LayoutNode{
			ID:         item.ID,
			Properties: map[string]dbus.Variant{"label": dbus.MakeVariant(item.Label)},
			Children:   []dbus.Variant{},
		}
		root.Children = append(root.Children, dbus.MakeVariant(child))
	}

	return m.revision, root, nil
}

func (m *DBusMenu) GetGroupProperties(ids []int32, propertyNames []string) ([]struct {
	ID         int32
	Properties map[string]dbus.Variant
}, *dbus.Error) {
	return []struct {
		ID         int32
		Properties map[string]dbus.Variant
	}{}, nil
}

func (m *DBusMenu) GetProperty(id int32, name string) (dbus.Variant, *dbus.Error) {
	return dbus.MakeVariant(""), nil
}

func (m *DBusMenu) Event(id int32, eventId string, data dbus.Variant, timestamp uint32) *dbus.Error {
	if eventId == "clicked" {
		m.tray.mu.RLock()
		defer m.tray.mu.RUnlock()
		for _, item := range m.tray.MenuItems {
			if item.ID == id {
				if item.Callback != nil {
					go item.Callback()
				}
				break
			}
		}
	}
	return nil
}

func (m *DBusMenu) AboutToShow(id int32) (bool, *dbus.Error) {
	return false, nil
}

func (m *DBusMenu) EmitLayoutUpdated() {
	m.mu.Lock()
	m.revision++
	rev := m.revision
	m.mu.Unlock()

	if m.tray.conn != nil {
		if err := m.tray.conn.Emit("/MenuBar", "com.canonical.dbusmenu.LayoutUpdated", rev, int32(0)); err != nil {
			logging.ReportError(context.Background(), err)
		}
	}
}
