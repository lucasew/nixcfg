package notification

import (
	"context"
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"
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

	if n.Message != "" {
		args = append(args, "-c", n.Message)
	}

	switch n.Urgency {
	case "low":
		args = append(args, "--priority", "low")
	case "normal":
		args = append(args, "--priority", "normal")
	case "critical":
		args = append(args, "--priority", "high")
	}

	if n.Hint != "" && strings.HasPrefix(n.Hint, "int:value:") {
		valStr := strings.TrimPrefix(n.Hint, "int:value:")
		if _, err := strconv.Atoi(valStr); err == nil {
			args = append(args, "--progress", valStr)
		}
	}

	_, err := common.RunCmd(ctx, "termux-notification", args...).Output()
	if err != nil {
		return fmt.Errorf("failed to send termux notification: %w", err)
	}

	return nil
}
