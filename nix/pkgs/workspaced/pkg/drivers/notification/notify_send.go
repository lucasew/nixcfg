package notification

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"workspaced/pkg/exec"
	"workspaced/pkg/logging"
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

	if n.HasProgress {
		args = append(args, "-h", fmt.Sprintf("int:value:%d", int(n.Progress*100)))
	}

	args = append(args, n.Title, n.Message)

	logging.GetLogger(ctx).Info("running notify-send", "args", strings.Join(args, " "))
	out, err := exec.RunCmd(ctx, "notify-send", args...).Output()
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
