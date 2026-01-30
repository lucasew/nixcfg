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
	defer func() { _ = resp.Body.Close() }()

	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = out.Close() }()

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

	type playerInfo struct {
		name     string
		status   string
		title    string
		artist   string
		artUrl   string
		length   int64
		position int64
	}

	var players []playerInfo
	for _, p := range mprisPlayers {
		obj := conn.Object(p, "/org/mpris/MediaPlayer2")

		statusVar, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.PlaybackStatus")
		if err != nil {
			common.GetLogger(ctx).Debug("failed to get status for player", "player", p, "error", err)
			continue
		}
		status := statusVar.Value().(string)

		metadataVar, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.Metadata")
		if err != nil {
			common.GetLogger(ctx).Debug("failed to get metadata for player", "player", p, "error", err)
			continue
		}
		m := metadataVar.Value().(map[string]dbus.Variant)

		title := ""
		if v, ok := m["xesam:title"]; ok {
			title = v.Value().(string)
		}

		artist := ""
		if v, ok := m["xesam:artist"]; ok {
			switch val := v.Value().(type) {
			case []string:
				artist = strings.Join(val, ", ")
			case []interface{}:
				var artists []string
				for _, a := range val {
					artists = append(artists, a.(string))
				}
				artist = strings.Join(artists, ", ")
			case string:
				artist = val
			}
		}

		artUrl := ""
		if v, ok := m["mpris:artUrl"]; ok {
			artUrl = v.Value().(string)
		}

		length := int64(0)
		if v, ok := m["mpris:length"]; ok {
			switch val := v.Value().(type) {
			case int64:
				length = val
			case uint64:
				length = int64(val)
			}
		}

		position := int64(0)
		posVar, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.Position")
		if err == nil {
			switch val := posVar.Value().(type) {
			case int64:
				position = val
			case uint64:
				position = int64(val)
			}
		}

		common.GetLogger(ctx).Debug("found player", "name", p, "title", title, "artist", artist, "status", status)

		players = append(players, playerInfo{
			name:     p,
			status:   status,
			title:    title,
			artist:   artist,
			artUrl:   artUrl,
			length:   length,
			position: position,
		})
	}

	if len(players) == 0 {
		return nil
	}

	// Select best player: Playing > Paused > Stopped, and non-empty title
	var best *playerInfo
	for i := range players {
		p := &players[i]
		if p.title == "" {
			continue
		}
		if best == nil {
			best = p
			continue
		}
		// Priority: Playing > Paused > Stopped
		statusPriority := map[string]int{"Playing": 3, "Paused": 2, "Stopped": 1}
		if statusPriority[p.status] > statusPriority[best.status] {
			best = p
		}
	}

	if best == nil || best.title == "" {
		common.GetLogger(ctx).Warn("no active player with title found")
		return nil
	}

	progress := 0.0
	if best.length > 0 {
		progress = float64(best.position) / float64(best.length)
	}

	iconPath := ""
	if best.artUrl != "" {
		iconPath, _ = getArtCachePath(best.artUrl)
	}

	title := best.title
	if title == "" {
		title = "Unknown Track"
	}
	message := best.artist
	if message == "" {
		message = "Unknown Artist"
	}

	n := &notification.Notification{
		ID:          notification.StatusNotificationID,
		Title:       title,
		Message:     message,
		Icon:        iconPath,
		Progress:    progress,
		HasProgress: true,
	}

	common.GetLogger(ctx).Info("sending media notification",
		"player", best.name,
		"title", title,
		"artist", message,
		"progress", progress,
		"icon", iconPath,
	)

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
