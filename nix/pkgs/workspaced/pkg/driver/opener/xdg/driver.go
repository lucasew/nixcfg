package xdg

import (
	"context"
	"fmt"
	"workspaced/pkg/driver"
	"workspaced/pkg/driver/opener"
	execdriver "workspaced/pkg/driver/exec"
	"workspaced/pkg/executil"
)

func init() {
	driver.Register[opener.Driver](&Provider{})
}

type Provider struct{}

func (p *Provider) ID() string   { return "opener_xdg" }
func (p *Provider) Name() string { return "xdg-open" }
func (p *Provider) DefaultWeight() int { return driver.DefaultWeight }

func (p *Provider) CheckCompatibility(ctx context.Context) error {
	if executil.GetEnv(ctx, "DISPLAY") == "" && executil.GetEnv(ctx, "WAYLAND_DISPLAY") == "" {
		return fmt.Errorf("%w: neither DISPLAY nor WAYLAND_DISPLAY set", driver.ErrIncompatible)
	}
	if !execdriver.IsBinaryAvailable(ctx, "xdg-open") {
		return fmt.Errorf("%w: xdg-open not found", driver.ErrIncompatible)
	}
	return nil
}

func (p *Provider) New(ctx context.Context) (opener.Driver, error) {
	return &Driver{}, nil
}

type Driver struct{}

func (d *Driver) Open(ctx context.Context, target string) error {
	return execdriver.MustRun(ctx, "xdg-open", target).Start()
}
