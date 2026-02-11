package dbus

import (
	"context"
	"log/slog"
	"sync"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
	"workspaced/pkg/logging"
)

type DBusMenu struct {
	driver   *Driver
	revision uint32
	mu       sync.Mutex
}

func NewDBusMenu(d *Driver) *DBusMenu {
	return &DBusMenu{
		driver:   d,
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
					{Name: "EventGroup", Args: []introspect.Arg{
						{Name: "events", Type: "a(isvu)", Direction: "in"},
						{Name: "idErrors", Type: "ai", Direction: "out"},
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
						{Name: "updatedProps", Type: "a(ia{sv})"},
						{Name: "removedProps", Type: "a(ias)"},
					}},
					{Name: "LayoutUpdated", Args: []introspect.Arg{
						{Name: "revision", Type: "u"},
						{Name: "parent", Type: "i"},
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
	if err != nil {
		return err
	}

	// Signal that layout is ready
	m.EmitLayoutUpdated()
	return nil
}

// LayoutNode represents (ia{sv}av)
type LayoutNode struct {
	ID         int32
	Properties map[string]dbus.Variant
	Children   []dbus.Variant
}

func (m *DBusMenu) GetLayout(parentId int32, recursionDepth int32, propertyNames []string) (uint32, LayoutNode, *dbus.Error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.driver.mu.RLock()
	defer m.driver.mu.RUnlock()

	slog.Info("dbus menu GetLayout", "parentId", parentId, "recursionDepth", recursionDepth)

	if parentId != 0 {
		return m.revision, LayoutNode{ID: parentId}, nil
	}

	root := LayoutNode{
		ID:         0,
		Properties: map[string]dbus.Variant{"children-display": dbus.MakeVariant("submenu")},
		Children:   []dbus.Variant{},
	}

	for i, item := range m.driver.state.Menu {
		id := int32(i + 1)
		child := LayoutNode{
			ID: id,
			Properties: map[string]dbus.Variant{
				"label":   dbus.MakeVariant(item.Label),
				"enabled": dbus.MakeVariant(true),
				"visible": dbus.MakeVariant(true),
				"type":    dbus.MakeVariant("standard"),
			},
			Children: []dbus.Variant{},
		}
		root.Children = append(root.Children, dbus.MakeVariant(child))
	}

	return m.revision, root, nil
}

func (m *DBusMenu) GetGroupProperties(ids []int32, propertyNames []string) ([]struct {
	ID         int32
	Properties map[string]dbus.Variant
}, *dbus.Error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.driver.mu.RLock()
	defer m.driver.mu.RUnlock()

	slog.Info("dbus menu GetGroupProperties", "ids", ids)

	res := []struct {
		ID         int32
		Properties map[string]dbus.Variant
	}{}

	for _, id := range ids {
		if id == 0 {
			res = append(res, struct {
				ID         int32
				Properties map[string]dbus.Variant
			}{
				ID:         0,
				Properties: map[string]dbus.Variant{"children-display": dbus.MakeVariant("submenu")},
			})
			continue
		}

		idx := int(id) - 1
		if idx >= 0 && idx < len(m.driver.state.Menu) {
			item := m.driver.state.Menu[idx]
			res = append(res, struct {
				ID         int32
				Properties map[string]dbus.Variant
			}{
				ID: id,
				Properties: map[string]dbus.Variant{
					"label":   dbus.MakeVariant(item.Label),
					"enabled": dbus.MakeVariant(true),
					"visible": dbus.MakeVariant(true),
					"type":    dbus.MakeVariant("standard"),
				},
			})
		}
	}

	return res, nil
}

func (m *DBusMenu) GetProperty(id int32, name string) (dbus.Variant, *dbus.Error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.driver.mu.RLock()
	defer m.driver.mu.RUnlock()

	slog.Info("dbus menu GetProperty", "id", id, "name", name)

	if id == 0 {
		if name == "children-display" {
			return dbus.MakeVariant("submenu"), nil
		}
		return dbus.MakeVariant(""), nil
	}

	idx := int(id) - 1
	if idx >= 0 && idx < len(m.driver.state.Menu) {
		item := m.driver.state.Menu[idx]
		switch name {
		case "label":
			return dbus.MakeVariant(item.Label), nil
		case "enabled":
			return dbus.MakeVariant(true), nil
		case "visible":
			return dbus.MakeVariant(true), nil
		case "type":
			return dbus.MakeVariant("standard"), nil
		}
	}

	return dbus.MakeVariant(""), nil
}

func (m *DBusMenu) Event(id int32, eventId string, data dbus.Variant, timestamp uint32) *dbus.Error {
	slog.Info("!!! EVENT !!!", "id", id, "eventId", eventId, "data", data.Value(), "timestamp", timestamp)
	m.handleEvent(id, eventId, data, timestamp)
	return nil
}

type EventGroupItem struct {
	ID        int32
	EventID   string
	Data      dbus.Variant
	Timestamp uint32
}

func (m *DBusMenu) EventGroup(events []EventGroupItem) ([]int32, *dbus.Error) {
	slog.Info("!!! EVENT GROUP !!!", "len", len(events))
	for _, e := range events {
		m.handleEvent(e.ID, e.EventID, e.Data, e.Timestamp)
	}
	return []int32{}, nil
}

func (m *DBusMenu) handleEvent(id int32, eventId string, data dbus.Variant, timestamp uint32) {
	if eventId == "clicked" || eventId == "activate" {
		m.driver.mu.RLock()
		defer m.driver.mu.RUnlock()

		idx := int(id) - 1
		if idx >= 0 && idx < len(m.driver.state.Menu) {
			item := m.driver.state.Menu[idx]
			if item.Callback != nil {
				slog.Info("executing menu callback", "label", item.Label)
				go item.Callback()
			}
		}
	}
}

func (m *DBusMenu) AboutToShow(id int32) (bool, *dbus.Error) {
	slog.Info("dbus menu AboutToShow", "id", id)
	return false, nil
}

func (m *DBusMenu) EmitLayoutUpdated() {
	m.mu.Lock()
	m.revision++
	rev := m.revision
	m.mu.Unlock()

	if m.driver.conn != nil {
		if err := m.driver.conn.Emit("/MenuBar", "com.canonical.dbusmenu.LayoutUpdated", rev, int32(0)); err != nil {
			logging.ReportError(context.Background(), err)
		}
	}
}
