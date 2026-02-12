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
	"workspaced/pkg/driver"
	mdrv "workspaced/pkg/driver/media"
	"workspaced/pkg/logging"
	"workspaced/pkg/notification"

	"github.com/godbus/dbus/v5"
)

func getArtCachePath(ctx context.Context, url string) (string, error) {
	if after, ok := strings.CutPrefix(url, "file://"); ok {
		return after, nil
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
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logging.ReportError(ctx, err)
		}
	}()

	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := out.Close(); err != nil {
			logging.ReportError(ctx, err)
		}
	}()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", err
	}

	return path, nil
}

func RunAction(ctx context.Context, action string) error {
	d, err := driver.Get[mdrv.Driver](ctx)
	if err != nil {
		return err
	}

	switch action {
	case "next":
		err = d.Next(ctx)
	case "previous":
		err = d.Previous(ctx)
	case "play-pause":
		err = d.PlayPause(ctx)
	case "stop":
		err = d.Stop(ctx)
	case "show":
		// just show
	default:
		return fmt.Errorf("unknown action: %s", action)
	}

	if err != nil {
		return err
	}

	// Small delay to let the player update metadata
	if action == "next" || action == "previous" || action == "play-pause" {
		time.Sleep(200 * time.Millisecond)
	}

	return ShowStatus(ctx)
}

func ShowStatus(ctx context.Context) error {
	d, err := driver.Get[mdrv.Driver](ctx)
	if err != nil {
		return err
	}

	best, err := d.GetMetadata(ctx)
	if err != nil {
		return err
	}

	if best == nil || best.Title == "" {
		logging.GetLogger(ctx).Warn("no active player with title found")
		return nil
	}

	progress := 0.0
	if best.Length > 0 {
		progress = float64(best.Position) / float64(best.Length)
	}

	iconPath := ""
	if best.ArtUrl != "" {
		var err error
		iconPath, err = getArtCachePath(ctx, best.ArtUrl)
		if err != nil {
			logging.ReportError(ctx, err)
		}
	}

	title := best.Title
	if title == "" {
		title = "Unknown Track"
	}
	message := best.Artist
	if message == "" {
		message = "Unknown Artist"
	}

	n := notification.Notification{
		ID:          notification.StatusNotificationID,
		Title:       title,
		Message:     message,
		Icon:        iconPath,
		Progress:    progress,
		HasProgress: true,
	}

	logging.GetLogger(ctx).Info("sending media notification",
		"player", best.Player,
		"title", title,
		"artist", message,
		"progress", progress,
		"icon", iconPath,
	)

	return notification.Notify(ctx, &n)
}

func Watch(ctx context.Context) {
	logger := logging.GetLogger(ctx)
	conn, err := dbus.SessionBus()
	if err != nil {
		logger.Error("failed to connect to session bus", "error", err)
		return
	}

	rule := "type='signal',interface='org.freedesktop.DBus.Properties',member='PropertiesChanged',path='/org/mpris/MediaPlayer2'"
	if err := conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, rule).Err; err != nil {
		logger.Error("failed to add dbus match", "error", err)
		return
	}

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
				if err := ShowStatus(ctx); err != nil {
					logging.ReportError(ctx, err)
				}
			}
		}
	}
}
