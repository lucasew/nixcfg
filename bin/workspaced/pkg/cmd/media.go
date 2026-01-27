package cmd

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func RunMedia(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("no command provided")
	}
	cmd := args[0]

	// Execute playerctl command (unless it's 'show', which implies just showing status)
	if cmd != "show" {
		if err := exec.Command("playerctl", cmd).Run(); err != nil {
			return "", fmt.Errorf("playerctl command failed: %w", err)
		}
	}

	// Wait logic from original script
	if cmd != "play-pause" && cmd != "show" {
		time.Sleep(2 * time.Second)
	}

	// Metadata format
	format := "{{playerName}};{{mpris:artUrl}};{{status}};{{artist}};{{title}};{{position*100/mpris:length}};"
	out, err := exec.Command("playerctl", "metadata", "-f", format).Output()
	if err != nil {
		return "", fmt.Errorf("failed to get metadata: %w", err)
	}

	line := strings.TrimSpace(string(out))
	parts := strings.Split(line, ";")
	if len(parts) < 5 {
		return "Metadata incomplete or empty", nil
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

	if err := exec.Command("notify-send", notifyArgs...).Run(); err != nil {
		return "", fmt.Errorf("notify-send failed: %w", err)
	}

	return fmt.Sprintf("State: %s\n", state), nil
}
