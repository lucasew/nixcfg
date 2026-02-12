package media

import (
	"context"
	"fmt"
	"time"
	"workspaced/pkg/driver"
	"workspaced/pkg/logging"
	"workspaced/pkg/notification"
)

func RunAction(ctx context.Context, action string) error {
	d, err := driver.Get[Driver](ctx)
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
	d, err := driver.Get[Driver](ctx)
	if err != nil {
		return err
	}

	meta, err := d.GetMetadata(ctx)
	if err != nil {
		return err
	}

	return Notify(ctx, meta)
}

func Notify(ctx context.Context, meta *Metadata) error {
	if meta == nil || meta.Title == "" {
		logging.GetLogger(ctx).Warn("no active player with title found")
		return nil
	}

	progress := 0.0
	if meta.Length > 0 {
		progress = float64(meta.Position) / float64(meta.Length)
	}

	iconPath := ""
	if meta.ArtUrl != "" {
		var err error
		iconPath, err = GetArtCachePath(ctx, meta.ArtUrl)
		if err != nil {
			logging.ReportError(ctx, err)
		}
	}

	title := meta.Title
	if title == "" {
		title = "Unknown Track"
	}
	message := meta.Artist
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
		"player", meta.Player,
		"title", title,
		"artist", message,
		"progress", progress,
		"icon", iconPath,
	)

	return notification.Notify(ctx, &n)
}

func Watch(ctx context.Context) {
	d, err := driver.Get[Driver](ctx)
	if err != nil {
		logging.GetLogger(ctx).Error("failed to get media driver for watch", "error", err)
		return
	}

	err = d.Watch(ctx, func(meta *Metadata) {
		if err := Notify(ctx, meta); err != nil {
			logging.ReportError(ctx, err)
		}
	})
	if err != nil {
		logging.GetLogger(ctx).Error("media watch failed", "error", err)
	}
}
