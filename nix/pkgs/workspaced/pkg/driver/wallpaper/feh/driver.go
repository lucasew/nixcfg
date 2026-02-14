package feh

import (
	"context"
	"fmt"
	"os"
	"workspaced/pkg/driver"
	"workspaced/pkg/driver/wallpaper"
	execdriver "workspaced/pkg/driver/exec"
)

func init() {
	driver.Register[wallpaper.Driver](&Provider{})
}

type Provider struct{}

func (p *Provider) ID() string   { return "x11_feh" }
func (p *Provider) Name() string { return "X11 (feh)" }
func (p *Provider) DefaultWeight() int { return driver.DefaultWeight }

func (p *Provider) CheckCompatibility(ctx context.Context) error {
	if os.Getenv("DISPLAY") == "" {
		return fmt.Errorf("%w: DISPLAY not set", driver.ErrIncompatible)
	}
	if _, err := execdriver.Which(ctx, "feh"); err != nil {
		return fmt.Errorf("%w: feh not found", driver.ErrIncompatible)
	}
	return nil
}

func (p *Provider) New(ctx context.Context) (wallpaper.Driver, error) {
	return &Driver{}, nil
}

type Driver struct{}

func (d *Driver) SetStatic(ctx context.Context, path string) error {
	feh, err := execdriver.Which(ctx, "feh")
	if err != nil {
		return err
	}
	cmd := execdriver.MustRun(ctx, "systemd-run", "--user", "-u", "wallpaper-change", "--collect", feh, "--bg-fill", path)
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("can't run feh in systemd unit: %w", err)
	}
	return nil
}
