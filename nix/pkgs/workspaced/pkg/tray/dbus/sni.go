package dbus

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
)

type StatusNotifierItem struct {
	driver *Driver
	props  *prop.Properties
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

func NewStatusNotifierItem(d *Driver) *StatusNotifierItem {
	return &StatusNotifierItem{driver: d}
}

func imageToSNIPixmap(img image.Image) ([]SNIPixmap, error) {
	if img == nil {
		return []SNIPixmap{}, nil
	}

	width := int32(img.Bounds().Dx())
	height := int32(img.Bounds().Dy())
	data := new(bytes.Buffer)

	// SNI expects ARGB32 in network byte order
	for y := 0; y < int(height); y++ {
		for x := 0; x < int(width); x++ {
			r, g, b, a := img.At(x, y).RGBA()
			// RGBA() returns 0-65535, scale down to 0-255
			// Also it's pre-multiplied alpha, but usually ARGB32 expects standard.
			// Let's assume standard simple conversion for now.

			// binary.Write to bytes.Buffer strictly returns nil error unless OOM, but we must check.
			var err error
			err = binary.Write(data, binary.BigEndian, uint8(a>>8))
			if err != nil { return nil, fmt.Errorf("failed to write pixmap data: %w", err) }
			err = binary.Write(data, binary.BigEndian, uint8(r>>8))
			if err != nil { return nil, fmt.Errorf("failed to write pixmap data: %w", err) }
			err = binary.Write(data, binary.BigEndian, uint8(g>>8))
			if err != nil { return nil, fmt.Errorf("failed to write pixmap data: %w", err) }
			err = binary.Write(data, binary.BigEndian, uint8(b>>8))
			if err != nil { return nil, fmt.Errorf("failed to write pixmap data: %w", err) }
		}
	}

	return []SNIPixmap{{
		Width:  width,
		Height: height,
		Data:   data.Bytes(),
	}}, nil
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

	pixmaps, err := imageToSNIPixmap(s.driver.state.Icon)
	// If image conversion fails, we proceed with empty pixmaps but maybe log?
	// For now, if we can't convert, we just don't show it.
	if err != nil {
		pixmaps = []SNIPixmap{}
	}

	// Export properties
	propsMap := map[string]interface{}{
		"Category":            "ApplicationStatus",
		"Id":                  "workspaced", // Should come from config?
		"Title":               s.driver.state.Title,
		"Status":              "Active",
		"WindowId":            int32(0),
		"IconThemePath":       "",
		"Menu":                dbus.ObjectPath("/MenuBar"),
		"ItemIsMenu":          true,
		"IconName":            "", // We use Pixmap mainly if Image provided
		"IconPixmap":          pixmaps,
		"OverlayIconName":     "",
		"OverlayIconPixmap":   []SNIPixmap{},
		"AttentionIconName":   "",
		"AttentionIconPixmap": []SNIPixmap{},
		"ToolTip":             SNIToolTip{Icon: "", Image: nil, Title: s.driver.state.Title, Description: ""},
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
