package audio

import (
	"context"
	"workspaced/pkg/driver"
	"workspaced/pkg/logging"
	"workspaced/pkg/notification"
)

// ShowStatus retrieves the current volume and mute status of the default sink
// and displays a notification with a progress bar.
// It parses the output of `pactl get-sink-volume` and `pactl get-sink-mute`.
func ShowStatus(ctx context.Context) error {
	d, err := driver.Get[Driver](ctx)
	if err != nil {
		return err
	}
	level, err := d.GetVolume(ctx)
	if err != nil {
		return err
	}
	isMuted, err := d.GetMute(ctx)
	if err != nil {
		return err
	}
	sinkName, err := d.SinkName(ctx)
	if err != nil {
		return err
	}

	icon := "audio-volume-high"
	if isMuted || level == 0 {
		icon = "audio-volume-muted"
	} else if level < .33 {
		icon = "audio-volume-low"
	} else if level < .66 {
		icon = "audio-volume-medium"
	}

	logging.GetLogger(ctx).Info("volume updated", "level", level, "sink", sinkName, "muted", isMuted)

	n := notification.Notification{
		ID:          notification.StatusNotificationID,
		Title:       "Volume",
		Message:     sinkName,
		Icon:        icon,
		Progress:    float64(level) / 100.0,
		HasProgress: true,
	}
	return notification.Notify(ctx, &n)
}
