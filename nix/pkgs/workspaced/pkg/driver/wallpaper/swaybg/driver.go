package swaybg

import (
	"context"
	"fmt"
	"workspaced/pkg/driver"
	"workspaced/pkg/exec"
	"workspaced/pkg/wallpaper/api"
)

func init() {
	driver.Register[api.Driver](&Provider{})
}

type Provider struct{}

func (p *Provider) Name() string { return "Wayland (swaybg)" }

func (p *Provider) CheckCompatibility(ctx context.Context) error {
	rpc := exec.GetRPC(ctx)
	if rpc != "swaymsg" {
		return fmt.Errorf("%w: current session is '%s', expected 'swaymsg'", driver.ErrIncompatible, rpc)
	}
	if _, err := exec.Which(ctx, "swaybg"); err != nil {
		return fmt.Errorf("%w: swaybg not found", driver.ErrIncompatible)
	}
	return nil
}

func (p *Provider) New(ctx context.Context) (api.Driver, error) {
	return &Driver{}, nil
}

type Driver struct{}

func (d *Driver) SetStatic(ctx context.Context, path string) error {
	swaybg, err := exec.Which(ctx, "swaybg")
	if err != nil {
		return err
	}

	if err = exec.RunCmd(ctx, "systemd-run", "--user", "-u", "wallpaper-change", "--collect", swaybg, "-i", path).Run(); err != nil {
		return fmt.Errorf("can't run swaybg in systemd unit: %w", err)
	}
	return nil
}
