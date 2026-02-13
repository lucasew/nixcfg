package systemd

import (
	"context"
	"fmt"
	"workspaced/pkg/driver"
	"workspaced/pkg/driver/power"
	"workspaced/pkg/exec"
)

func init() {
	driver.Register[power.Driver](&Provider{})
}

type Provider struct{}

func (p *Provider) ID() string   { return "power_systemd" }
func (p *Provider) Name() string { return "Systemd" }
func (p *Provider) DefaultWeight() int { return driver.DefaultWeight }

func (p *Provider) CheckCompatibility(ctx context.Context) error {
	if !exec.IsBinaryAvailable(ctx, "loginctl") {
		return fmt.Errorf("%w: loginctl not found", driver.ErrIncompatible)
	}
	return nil
}

func (p *Provider) New(ctx context.Context) (power.Driver, error) {
	return &Driver{}, nil
}

type Driver struct{}

func (d *Driver) Lock(ctx context.Context) error {
	return exec.RunCmd(ctx, "loginctl", "lock-session").Run()
}

func (d *Driver) Logout(ctx context.Context) error {
	return exec.RunCmd(ctx, "loginctl", "terminate-session", "self").Run()
}

func (d *Driver) Suspend(ctx context.Context) error {
	return exec.RunCmd(ctx, "systemctl", "suspend").Run()
}

func (d *Driver) Hibernate(ctx context.Context) error {
	return exec.RunCmd(ctx, "systemctl", "hibernate").Run()
}

func (d *Driver) Reboot(ctx context.Context) error {
	return exec.RunCmd(ctx, "systemctl", "reboot").Run()
}

func (d *Driver) Shutdown(ctx context.Context) error {
	return exec.RunCmd(ctx, "systemctl", "poweroff").Run()
}
