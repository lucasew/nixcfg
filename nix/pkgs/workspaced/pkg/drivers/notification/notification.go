package notification

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"workspaced/pkg/common"
)

type Notification struct {
	ID      uint32
	Title   string
	Message string
	Urgency string // low, normal, critical
	Icon    string
	Hint    string // e.g. int:value:50 for progress bar
}

func (n *Notification) Notify(ctx context.Context) error {
	args := []string{}
	if n.ID > 0 {
		args = append(args, "-r", strconv.FormatUint(uint64(n.ID), 10))
	} else {
		args = append(args, "-p")
	}

	if n.Urgency != "" {
		args = append(args, "-u", n.Urgency)
	}

	if n.Icon != "" {
		args = append(args, "-i", n.Icon)
	}

	if n.Hint != "" {
		args = append(args, "-h", n.Hint)
	}

	args = append(args, n.Title, n.Message)

	out, err := common.RunCmd(ctx, "notify-send", args...).Output()
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	if n.ID == 0 {
		idStr := strings.TrimSpace(string(out))
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err == nil {
			n.ID = uint32(id)
		}
	}

	return nil
}
