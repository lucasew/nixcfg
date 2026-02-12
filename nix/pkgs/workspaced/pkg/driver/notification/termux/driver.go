package termux

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

func (p *Provider) Name() string { return "Termux" }

func (p *Provider) CheckCompatibility(ctx context.Context) error {
	if !exec.IsBinaryAvailable(ctx, "termux-notification") {
		return fmt.Errorf("%w: termux-notification not found", driver.ErrIncompatible)
	}
	return nil
}

func (p *Provider) New(ctx context.Context) (notification.Driver, error) {
	return &Driver{}, nil
}

type Driver struct{}

func (d *Driver) Notify(ctx context.Context, n *notification.Notification) error {
	args := []string{
		"--title", n.Title,
		"--content", n.Message,
	}
	if n.ID != 0 {
		args = append(args, "--id", fmt.Sprintf("%d", n.ID))
	}
	return exec.RunCmd(ctx, "termux-notification", args...).Run()
}
