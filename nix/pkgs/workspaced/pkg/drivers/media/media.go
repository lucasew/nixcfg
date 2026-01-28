package media

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
	"workspaced/pkg/common"
	"workspaced/pkg/drivers/notification"
)

var n = &notification.Notification{}

func RunAction(ctx context.Context, action string) error {
	if action != "show" {
		if err := common.RunCmd(ctx, "playerctl", action).Run(); err != nil {
			return fmt.Errorf("playerctl command failed: %w", err)
		}
	}

	if action != "play-pause" && action != "show" {
		time.Sleep(2 * time.Second)
	}

	return ShowStatus(ctx)
}

func ShowStatus(ctx context.Context) error {
	format := "{{playerName}};{{mpris:artUrl}};{{status}};{{artist}};{{title}};{{position*100/mpris:length}};"
	out, err := common.RunCmd(ctx, "playerctl", "metadata", "-f", format).Output()
	if err != nil {
		return fmt.Errorf("failed to get metadata: %w", err)
	}

	line := strings.TrimSpace(string(out))
	parts := strings.Split(line, ";")
	if len(parts) < 5 {
		return nil
	}

	player := parts[0]
	icon := parts[1]
	state := parts[2]
	artist := parts[3]
	title := parts[4]

	emoji := "❔"
	switch state {
	case "Playing":
		emoji = "▶️"
	case "Paused":
		emoji = "⏸️"
	case "Stopped":
		emoji = "⏹️"
	}

	n.Title = fmt.Sprintf("%s %s", emoji, player)
	n.Message = fmt.Sprintf("%s - %s", artist, title)
	if l, err := strconv.Atoi(parts[5]); err == nil {
		n.Progress = l
	}
	n.Icon = icon

	common.GetLogger(ctx).Info("media status", "player", player, "state", state, "artist", artist, "title", title)

	return n.Notify(ctx)
}
