package notification

import (
	"context"
	"fmt"
	"hash/fnv"
	"strconv"
	"workspaced/pkg/common"
)

type TermuxNotifier struct{}

func (t *TermuxNotifier) Notify(ctx context.Context, n *Notification) error {
	args := []string{}
	if n.ID == 0 {
		h := fnv.New32a()
		h.Write([]byte(n.Title))
		h.Write([]byte(n.Icon))
		n.ID = h.Sum32()
	}
	args = append(args, "-i", strconv.FormatUint(uint64(n.ID), 10))

	if n.Title != "" {
		args = append(args, "-t", n.Title)
	}

	switch n.Urgency {
	case "low":
		args = append(args, "--priority", "low")
	case "normal":
		args = append(args, "--priority", "normal")
	case "critical":
		args = append(args, "--priority", "high")
	}

	message := n.Message
	if n.Progress > 0 {
		if n.Progress < 1.0 {
			args = append(args, "--ongoing")
		}
		args = append(args, "--alert-once")
		width := 10
		percent := int(n.Progress * 100)
		completed := (percent * width) / 100
		bar := ""
		for i := 0; i < width; i++ {
			if i < completed {
				bar += "█"
			} else {
				bar += "░"
			}
		}
		message = fmt.Sprintf("%s\n%s %d%%", message, bar, percent)
	}

	if message != "" {
		args = append(args, "-c", message)
	}

	_, err := common.RunCmd(ctx, "termux-notification", args...).Output()
	if err != nil {
		return fmt.Errorf("failed to send termux notification: %w", err)
	}

	return nil
}
