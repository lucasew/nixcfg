package screen

import (
	"context"
	"fmt"
	"strings"
	"workspaced/pkg/env"
	"workspaced/pkg/exec"
)

type X11Driver struct{}

func (d *X11Driver) SetDPMS(ctx context.Context, on bool) error {
	xsetArg := "off"
	if on {
		xsetArg = "on"
	}
	return exec.RunCmd(ctx, "xset", "dpms", "force", xsetArg).Run()
}

func (d *X11Driver) IsDPMSOn(ctx context.Context) (bool, error) {
	out, err := exec.RunCmd(ctx, "xset", "q").Output()
	if err != nil {
		return false, err
	}
	return strings.Contains(string(out), "Monitor is On"), nil
}

func (d *X11Driver) Reset(ctx context.Context) error {
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
	return fmt.Errorf("no x11 reset logic for this host")
}
