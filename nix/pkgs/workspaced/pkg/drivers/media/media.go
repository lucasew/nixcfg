package media

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"workspaced/pkg/common"
	"workspaced/pkg/drivers/notification"

	"github.com/godbus/dbus/v5"
)

func getArtCachePath(url string) (string, error) {
	if strings.HasPrefix(url, "file://") {
		return strings.TrimPrefix(url, "file://"), nil
	}
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return url, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	cacheDir := filepath.Join(home, ".cache/workspaced/media_art")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", err
	}

	hash := fmt.Sprintf("%x", md5.Sum([]byte(url)))
	path := filepath.Join(cacheDir, hash)

	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", err
	}

	return path, nil
}

func RunAction(ctx context.Context, action string) error {
	if action != "show" {
		if err := common.RunCmd(ctx, "playerctl", action).Run(); err != nil {
			return fmt.Errorf("playerctl command failed: %w", err)
		}
	}

	// Small delay to let the player update metadata
	if action == "next" || action == "previous" || action == "play-pause" {
		time.Sleep(200 * time.Millisecond)
	}

	return ShowStatus(ctx)
}

func ShowStatus(ctx context.Context) error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return err
	}

	var names []string
	err = conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&names)
	if err != nil {
		return err
	}

	var mprisPlayers []string
	for _, name := range names {
		if strings.HasPrefix(name, "org.mpris.MediaPlayer2.") {
			mprisPlayers = append(mprisPlayers, name)
		}
	}

	if len(mprisPlayers) == 0 {
		return nil
	}

	// Use the first available player for now
	player := mprisPlayers[0]
	obj := conn.Object(player, "/org/mpris/MediaPlayer2")

	metadata, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.Metadata")
	if err != nil {
		return err
	}

	m := metadata.Value().(map[string]dbus.Variant)

	title := ""
	if v, ok := m["xesam:title"]; ok {
		title = v.Value().(string)
	}

	artist := ""
	if v, ok := m["xesam:artist"]; ok {
		artists := v.Value().([]string)
		artist = strings.Join(artists, ", ")
	}

	artUrl := ""
	if v, ok := m["mpris:artUrl"]; ok {
		artUrl = v.Value().(string)
	}

	length := int64(0)
	if v, ok := m["mpris:length"]; ok {
		length = v.Value().(int64)
	}

	position := int64(0)
	posVar, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.Position")
	if err == nil {
		position = posVar.Value().(int64)
	}

	statusVar, _ := obj.GetProperty("org.mpris.MediaPlayer2.Player.PlaybackStatus")
	status := statusVar.Value().(string)

	progress := 0.0
	if length > 0 {
		progress = float64(position) / float64(length)
	}

	iconPath := ""
	if artUrl != "" {
		iconPath, _ = getArtCachePath(artUrl)
	}

	n := &notification.Notification{
		ID:       notification.StatusNotificationID,
		Title:    title,
		Message:  artist,
		Icon:     iconPath,
		Progress: progress,
	}

	common.GetLogger(ctx).Info("media status", "player", player, "title", title, "artist", artist, "status", status)

	return n.Notify(ctx)
}

func Watch(ctx context.Context) {
	logger := common.GetLogger(ctx)
	conn, err := dbus.SessionBus()
	if err != nil {
		logger.Error("failed to connect to session bus", "error", err)
		return
	}

	rule := "type='signal',interface='org.freedesktop.DBus.Properties',member='PropertiesChanged',path='/org/mpris/MediaPlayer2'"
	conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, rule)

	c := make(chan *dbus.Signal, 10)
	conn.Signal(c)

	for {
		select {
		case <-ctx.Done():
			return
		case signal := <-c:
			if len(signal.Body) < 2 {
				continue
			}
			if signal.Body[0].(string) != "org.mpris.MediaPlayer2.Player" {
				continue
			}
			changed := signal.Body[1].(map[string]dbus.Variant)
			if _, ok := changed["Metadata"]; ok {
				logger.Debug("metadata changed, updating status")
				_ = ShowStatus(ctx)
			}
		}
	}
}
