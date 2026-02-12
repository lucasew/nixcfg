package notify_send

import (
	"context"
	"fmt"
	"workspaced/pkg/driver"
	"workspaced/pkg/exec"
	"workspaced/pkg/notification"
)

func init() {
	driver.Register[notification.Driver](&Provider{})
}

type Provider struct{}

func (p *Provider) Name() string { return "notify-send" }

func (p *Provider) CheckCompatibility(ctx context.Context) error {
	if !exec.IsBinaryAvailable(ctx, "notify-send") {
		return fmt.Errorf("%w: notify-send not found", driver.ErrIncompatible)
	}
	return nil
}

func (p *Provider) New(ctx context.Context) (notification.Driver, error) {
	return &Driver{}, nil
}

type Driver struct{}

func (d *Driver) Notify(ctx context.Context, n *notification.Notification) error {
	args := []string{}
	if n.Urgency != "" {
		args = append(args, "-u", n.Urgency)
	}
	if n.Icon != "" {
		args = append(args, "-i", n.Icon)
	}
	if n.HasProgress {
		args = append(args, "-h", fmt.Sprintf("int:value:%d", int(n.Progress*100)))
	}
	if n.ID != 0 {
		args = append(args, "-r", fmt.Sprintf("%d", n.ID))
	}

	args = append(args, n.Title, n.Message)
	return exec.RunCmd(ctx, "notify-send", args...).Run()
}
