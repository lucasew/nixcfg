package audio

import (
	"context"
	"fmt"
	"strings"
	"workspaced/pkg/common"
	"workspaced/pkg/drivers/notification"
)

var n = &notification.Notification{
	Icon: "audio",
}

func SetVolume(ctx context.Context, arg string) error {
	sink := "@DEFAULT_SINK@"
	if err := common.RunCmd(ctx, "pactl", "set-sink-volume", sink, arg).Run(); err != nil {
		return fmt.Errorf("failed to set volume: %w", err)
	}
	return ShowStatus(ctx)
}

func ShowStatus(ctx context.Context) error {
	sink := "@DEFAULT_SINK@"
	out, err := common.RunCmd(ctx, "pactl", "get-sink-volume", sink).Output()
	if err != nil {
		return fmt.Errorf("failed to get volume: %w", err)
	}

	level := "0"
	parts := strings.Fields(string(out))
	for _, p := range parts {
		if strings.Contains(p, "%") {
			level = strings.Trim(p, "%")
			break
		}
	}

	emoji := "🔊"
	if level == "0" {
		emoji = "🔇"
	}

	sinkNameOut, _ := common.RunCmd(ctx, "pactl", "get-default-sink").Output()
	sinkName := strings.TrimSpace(string(sinkNameOut))

	n.Title = fmt.Sprintf("%s Volume", emoji)
	n.Message = sinkName
	n.Hint = fmt.Sprintf("int:value:%s", level)

	return n.Notify(ctx)
}
