package notification

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"workspaced/pkg/common"
)

type NotifySendNotifier struct{}

func (s *NotifySendNotifier) Notify(ctx context.Context, n *Notification) error {
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

	if n.Progress > 0 {
		args = append(args, "-h", fmt.Sprintf("int:value:%d", n.Progress))
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
