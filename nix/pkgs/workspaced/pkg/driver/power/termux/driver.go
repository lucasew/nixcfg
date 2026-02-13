package termux

import (
	"context"
	"fmt"
	"os"
	"workspaced/pkg/api"
	"workspaced/pkg/driver"
	"workspaced/pkg/driver/power"
	"workspaced/pkg/exec"
)

func init() {
	driver.Register[power.Driver](&Provider{})
}

type Provider struct{}

func (p *Provider) ID() string   { return "power_termux" }
func (p *Provider) Name() string { return "Termux" }
func (p *Provider) DefaultWeight() int { return driver.DefaultWeight }

func (p *Provider) CheckCompatibility(ctx context.Context) error {
	if os.Getenv("TERMUX_VERSION") == "" {
		return fmt.Errorf("%w: not running in Termux", driver.ErrIncompatible)
	}
	return nil
}

func (p *Provider) New(ctx context.Context) (power.Driver, error) {
	return &Driver{}, nil
}

type Driver struct{}

func (d *Driver) Lock(ctx context.Context) error {
	return fmt.Errorf("%w: screen lock not possible in Termux", api.ErrNotSupported)
}

func (d *Driver) Logout(ctx context.Context) error {
	return fmt.Errorf("%w: logout not possible in Termux", api.ErrNotSupported)
}

func (d *Driver) Suspend(ctx context.Context) error {
	return fmt.Errorf("%w: suspend not possible in Termux", api.ErrNotSupported)
}

func (d *Driver) Hibernate(ctx context.Context) error {
	return fmt.Errorf("%w: hibernate not possible in Termux", api.ErrNotSupported)
}

func (d *Driver) Reboot(ctx context.Context) error {
	// If rooted, might work.
	return exec.RunCmd(ctx, "reboot").Run()
}

func (d *Driver) Shutdown(ctx context.Context) error {
	// If rooted, might work.
	return exec.RunCmd(ctx, "shutdown", "-h", "now").Run()
}
