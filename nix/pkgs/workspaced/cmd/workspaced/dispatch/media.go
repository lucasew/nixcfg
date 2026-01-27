package dispatch

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var mediaCmd = &cobra.Command{
	Use:   "media",
	Short: "Control media playback",
}

func init() {
	Command.AddCommand(mediaCmd)

	cmds := []string{"play-pause", "next", "previous", "stop", "show"}
	for _, c := range cmds {
		cmdName := c
		subCmd := &cobra.Command{
			Use:   cmdName,
			Short: fmt.Sprintf("%s media", cmdName),
			RunE: func(c *cobra.Command, args []string) error {
				return runMediaAction(cmdName)
			},
		}
		mediaCmd.AddCommand(subCmd)
	}
}

func runMediaAction(action string) error {
	if action != "show" {
		if err := exec.Command("playerctl", action).Run(); err != nil {
			return fmt.Errorf("playerctl command failed: %w", err)
		}
	}

	if action != "play-pause" && action != "show" {
		time.Sleep(2 * time.Second)
	}

	format := "{{playerName}};{{mpris:artUrl}};{{status}};{{artist}};{{title}};{{position*100/mpris:length}};"
	out, err := exec.Command("playerctl", "metadata", "-f", format).Output()
	if err != nil {
		return fmt.Errorf("failed to get metadata: %w", err)
	}

	line := strings.TrimSpace(string(out))
	parts := strings.Split(line, ";")
	if len(parts) < 5 {
		fmt.Println("Metadata incomplete or empty")
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

	if err := exec.Command("notify-send", notifyArgs...).Run(); err != nil {
		return fmt.Errorf("notify-send failed: %w", err)
	}

	fmt.Printf("State: %s\n", state)
	return nil
}
