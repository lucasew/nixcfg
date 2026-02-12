package x11

import (
	"context"
	"fmt"
	"os"
	"strings"
	"workspaced/pkg/api"
	"workspaced/pkg/driver"
	"workspaced/pkg/env"
	"workspaced/pkg/exec"
	"workspaced/pkg/screen"
	"workspaced/pkg/types"
)

func init() {
	driver.Register[screen.Driver](&Provider{})
}

type Provider struct{}

func (p *Provider) Name() string { return "X11" }

func (p *Provider) CheckCompatibility(ctx context.Context) error {
	display := os.Getenv("DISPLAY")
	if env, ok := ctx.Value(types.EnvKey).([]string); ok {
		for _, e := range env {
			if strings.HasPrefix(e, "DISPLAY=") {
				display = strings.TrimPrefix(e, "DISPLAY=")
				break
			}
		}
	}

	if display == "" {
		return fmt.Errorf("%w: DISPLAY not set", driver.ErrIncompatible)
	}
	return nil
}

func (p *Provider) New(ctx context.Context) (screen.Driver, error) {
	return &Driver{}, nil
}

type Driver struct{}

func (d *Driver) SetDPMS(ctx context.Context, on bool) error {
	xsetArg := "off"
	if on {
		xsetArg = "on"
	}
	return exec.RunCmd(ctx, "xset", "dpms", "force", xsetArg).Run()
}

func (d *Driver) IsDPMSOn(ctx context.Context) (bool, error) {
	out, err := exec.RunCmd(ctx, "xset", "q").Output()
	if err != nil {
		return false, err
	}
	return strings.Contains(string(out), "Monitor is On"), nil
}

func (d *Driver) Reset(ctx context.Context) error {
	if env.IsRiverwood() {
		// Ensure eDP-1 is primary and on the left, HDMI-A-1 on the right
		return exec.RunCmd(ctx, "xrandr",
			"--output", "eDP-1", "--auto", "--primary", "--pos", "0x0",
			"--output", "HDMI-A-1", "--auto", "--right-of", "eDP-1",
		).Run()
	}
	if env.IsWhiterun() {
		return exec.RunCmd(ctx, "xrandr", "--output", "HDMI-1", "--auto").Run()
	}
	return api.ErrNotImplemented
}
