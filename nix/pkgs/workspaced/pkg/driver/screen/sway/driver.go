package sway

import (
	"context"
	"fmt"
	"strings"
	"workspaced/pkg/api"
	"workspaced/pkg/driver"
	"workspaced/pkg/env"
	"workspaced/pkg/exec"
	"workspaced/pkg/driver/screen"
)

func init() {
	driver.Register[screen.Driver](&Provider{})
}

type Provider struct{}

func (p *Provider) Name() string { return "Wayland (sway)" }

func (p *Provider) CheckCompatibility(ctx context.Context) error {
	if exec.GetEnv(ctx, "WAYLAND_DISPLAY") == "" {
		return fmt.Errorf("%w: WAYLAND_DISPLAY not set", driver.ErrIncompatible)
	}
	if !exec.IsBinaryAvailable(ctx, "swaymsg") {
		return fmt.Errorf("%w: swaymsg not found", driver.ErrIncompatible)
	}
	return nil
}

func (p *Provider) New(ctx context.Context) (screen.Driver, error) {
	return &Driver{}, nil
}

type Driver struct{}

func (d *Driver) SetDPMS(ctx context.Context, on bool) error {
	state := "off"
	if on {
		state = "on"
	}
	return exec.RunCmd(ctx, "swaymsg", "output * dpms "+state).Run()
}

func (d *Driver) IsDPMSOn(ctx context.Context) (bool, error) {
	out, err := exec.RunCmd(ctx, "swaymsg", "-t", "get_outputs").Output()
	if err != nil {
		return false, err
	}
	return strings.Contains(string(out), `"dpms": true`), nil
}

func (d *Driver) Reset(ctx context.Context) error {
	if env.IsRiverwood() {
		// eDP-1 (notebook) on the LEFT (0,0), HDMI-A-1 on the RIGHT
		if err := exec.RunCmd(ctx, "swaymsg", "output", "eDP-1", "mode", "1366x768", "pos", "0", "0").Run(); err != nil {
			return err
		}
		return exec.RunCmd(ctx, "swaymsg", "output", "HDMI-A-1", "mode", "1366x768", "pos", "1366", "0").Run()
	}
	if env.IsWhiterun() {
		return exec.RunCmd(ctx, "swaymsg", "output", "HDMI-A-1", "mode", "1368x768").Run()
	}
	return api.ErrNotImplemented
}
