package audio

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"workspaced/pkg/drivers/notification"
	"workspaced/pkg/exec"
	"workspaced/pkg/logging"
)

// SetVolume sets the volume of the default sink using pactl.
// The arg parameter accepts format like "50%", "+10%", "-10%".
// It also triggers a status notification after the change.
func SetVolume(ctx context.Context, arg string) error {
	sink := "@DEFAULT_SINK@"
	if err := exec.RunCmd(ctx, "pactl", "set-sink-volume", sink, arg).Run(); err != nil {
		return fmt.Errorf("failed to set volume: %w", err)
	}
	return ShowStatus(ctx)
}

// ShowStatus retrieves the current volume and mute status of the default sink
// and displays a notification with a progress bar.
// It parses the output of `pactl get-sink-volume` and `pactl get-sink-mute`.
func ShowStatus(ctx context.Context) error {
	sink := "@DEFAULT_SINK@"
	out, err := exec.RunCmd(ctx, "pactl", "get-sink-volume", sink).Output()
	if err != nil {
		return fmt.Errorf("failed to get volume: %w", err)
	}

	level := 0
	parts := strings.Fields(string(out))
	for _, p := range parts {
		if strings.Contains(p, "%") {
			l, err := strconv.Atoi(strings.Trim(p, "%"))
			if err == nil {
				level = l
				break
			}
		}
	}

	muteOut, _ := exec.RunCmd(ctx, "pactl", "get-sink-mute", sink).Output()
	isMuted := strings.Contains(string(muteOut), "yes")

	icon := "audio-volume-high"
	if isMuted || level == 0 {
		icon = "audio-volume-muted"
	} else if level < 33 {
		icon = "audio-volume-low"
	} else if level < 66 {
		icon = "audio-volume-medium"
	}

	sinkNameOut, _ := exec.RunCmd(ctx, "pactl", "get-default-sink").Output()
	sinkName := strings.TrimSpace(string(sinkNameOut))

	n := &notification.Notification{
		ID:          notification.StatusNotificationID,
		Title:       "Volume",
		Message:     sinkName,
		Icon:        icon,
		Progress:    float64(level) / 100.0,
		HasProgress: true,
	}

	logging.GetLogger(ctx).Info("volume updated", "level", level, "sink", sinkName, "muted", isMuted)

	return n.Notify(ctx)
}
