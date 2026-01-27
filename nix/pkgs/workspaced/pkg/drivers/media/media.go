package media

import (
	"context"
	"fmt"
	"strings"
	"time"
	"workspaced/pkg/common"
)

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

	notifyArgs := []string{
		fmt.Sprintf("%s %s", emoji, player),
		fmt.Sprintf("%s - %s", artist, title),
		"-h", fmt.Sprintf("int:value:%s", parts[5]),
		"-i", icon,
		"-r", "28693965",
	}

	return common.RunCmd(ctx, "notify-send", notifyArgs...).Run()
}
