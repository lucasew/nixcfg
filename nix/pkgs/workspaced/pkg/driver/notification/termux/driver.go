package termux

import (
	"context"
	"fmt"
	"workspaced/pkg/driver"
	"workspaced/pkg/driver/notification"
	execdriver "workspaced/pkg/driver/exec"
)

func init() {
	driver.Register[notification.Driver](&Provider{})
}

type Provider struct{}

func (p *Provider) ID() string   { return "notification_termux" }
func (p *Provider) Name() string { return "Termux" }
func (p *Provider) DefaultWeight() int { return driver.DefaultWeight }

func (p *Provider) CheckCompatibility(ctx context.Context) error {
	if !execdriver.IsBinaryAvailable(ctx, "termux-notification") {
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
	return execdriver.MustRun(ctx, "termux-notification", args...).Run()
}
