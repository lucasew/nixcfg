package dbus

import (
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
)

type StatusNotifierItem struct {
	tray  *Tray
	props *prop.Properties
}

type SNIPixmap struct {
	Width  int32
	Height int32
	Data   []byte
}

type SNIToolTip struct {
	Icon        string
	Image       []SNIPixmap
	Title       string
	Description string
}

func NewStatusNotifierItem(t *Tray) *StatusNotifierItem {
	return &StatusNotifierItem{tray: t}
}

func (s *StatusNotifierItem) Export(conn *dbus.Conn, path dbus.ObjectPath) error {
	// Export methods
	err := conn.Export(s, path, "org.kde.StatusNotifierItem")
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
				Name: "org.kde.StatusNotifierItem",
				Methods: []introspect.Method{
					{Name: "ContextMenu", Args: []introspect.Arg{{Name: "x", Type: "i", Direction: "in"}, {Name: "y", Type: "i", Direction: "in"}}},
					{Name: "Activate", Args: []introspect.Arg{{Name: "x", Type: "i", Direction: "in"}, {Name: "y", Type: "i", Direction: "in"}}},
					{Name: "SecondaryActivate", Args: []introspect.Arg{{Name: "x", Type: "i", Direction: "in"}, {Name: "y", Type: "i", Direction: "in"}}},
					{Name: "Scroll", Args: []introspect.Arg{{Name: "delta", Type: "i", Direction: "in"}, {Name: "orientation", Type: "s", Direction: "in"}}},
				},
				Properties: []introspect.Property{
					{Name: "Category", Type: "s", Access: "read"},
					{Name: "Id", Type: "s", Access: "read"},
					{Name: "Title", Type: "s", Access: "read"},
					{Name: "Status", Type: "s", Access: "read"},
					{Name: "WindowId", Type: "i", Access: "read"},
					{Name: "IconThemePath", Type: "s", Access: "read"},
					{Name: "Menu", Type: "o", Access: "read"},
					{Name: "ItemIsMenu", Type: "b", Access: "read"},
					{Name: "IconName", Type: "s", Access: "read"},
					{Name: "IconPixmap", Type: "a(iiay)", Access: "read"},
					{Name: "OverlayIconName", Type: "s", Access: "read"},
					{Name: "OverlayIconPixmap", Type: "a(iiay)", Access: "read"},
					{Name: "AttentionIconName", Type: "s", Access: "read"},
					{Name: "AttentionIconPixmap", Type: "a(iiay)", Access: "read"},
					{Name: "ToolTip", Type: "(sa(iiay)ss)", Access: "read"},
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
		"Category":            "ApplicationStatus",
		"Id":                  s.tray.ID,
		"Title":               s.tray.Title,
		"Status":              "Active",
		"WindowId":            int32(0),
		"IconThemePath":       "",
		"Menu":                dbus.ObjectPath("/MenuBar"),
		"ItemIsMenu":          true,
		"IconName":            s.tray.Icon,
		"IconPixmap":          []SNIPixmap{},
		"OverlayIconName":     "",
		"OverlayIconPixmap":   []SNIPixmap{},
		"AttentionIconName":   "",
		"AttentionIconPixmap": []SNIPixmap{},
		"ToolTip":             SNIToolTip{Icon: "", Image: nil, Title: "", Description: ""},
	}

	convertedProps := make(map[string]*prop.Prop)
	for k, v := range propsMap {
		convertedProps[k] = &prop.Prop{
			Value:    v,
			Writable: false,
			Emit:     prop.EmitTrue,
		}
	}

	p, err := prop.Export(conn, path, prop.Map{
		"org.kde.StatusNotifierItem": convertedProps,
	})
	if err != nil {
		return err
	}
	s.props = p
	return nil
}

func (s *StatusNotifierItem) ContextMenu(x, y int32) *dbus.Error {
	return nil
}

func (s *StatusNotifierItem) Activate(x, y int32) *dbus.Error {
	return nil
}

func (s *StatusNotifierItem) SecondaryActivate(x, y int32) *dbus.Error {
	return nil
}

func (s *StatusNotifierItem) Scroll(delta int32, orientation string) *dbus.Error {
	return nil
}
