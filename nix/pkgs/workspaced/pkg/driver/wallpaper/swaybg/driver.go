package swaybg

import (
	"context"
	"fmt"
	"workspaced/pkg/driver"
	"workspaced/pkg/driver/wallpaper"
	"workspaced/pkg/exec"
)

func init() {
	driver.Register[wallpaper.Driver](&Provider{})
}

type Provider struct{}

func (p *Provider) ID() string   { return "wayland_swaybg" }
func (p *Provider) Name() string { return "Wayland (swaybg)" }
func (p *Provider) DefaultWeight() int { return driver.DefaultWeight }

func (p *Provider) CheckCompatibility(ctx context.Context) error {
	if exec.GetEnv(ctx, "WAYLAND_DISPLAY") == "" {
		return fmt.Errorf("%w: WAYLAND_DISPLAY not set", driver.ErrIncompatible)
	}
	if _, err := exec.Which(ctx, "swaybg"); err != nil {
		return fmt.Errorf("%w: swaybg not found", driver.ErrIncompatible)
	}
	return nil
}

func (p *Provider) New(ctx context.Context) (wallpaper.Driver, error) {
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
