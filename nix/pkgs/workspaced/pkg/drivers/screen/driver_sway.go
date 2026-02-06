package screen

import (
	"context"
	"fmt"
	"strings"
	"workspaced/pkg/exec"
	"workspaced/pkg/host"
)

type SwayDriver struct{}

func (d *SwayDriver) SetDPMS(ctx context.Context, on bool) error {
	state := "off"
	if on {
		state = "on"
	}
	return exec.RunCmd(ctx, "swaymsg", "output * dpms "+state).Run()
}

func (d *SwayDriver) IsDPMSOn(ctx context.Context) (bool, error) {
	out, err := exec.RunCmd(ctx, "swaymsg", "-t", "get_outputs").Output()
	if err != nil {
		return false, err
	}
	return strings.Contains(string(out), `"dpms": true`), nil
}

func (d *SwayDriver) Reset(ctx context.Context) error {
	if host.IsRiverwood() {
		// eDP-1 (notebook) on the LEFT (0,0), HDMI-A-1 on the RIGHT
		if err := exec.RunCmd(ctx, "swaymsg", "output", "eDP-1", "mode", "1366x768", "pos", "0", "0").Run(); err != nil {
			return err
		}
		return exec.RunCmd(ctx, "swaymsg", "output", "HDMI-A-1", "mode", "1366x768", "pos", "1366", "0").Run()
	}
	if host.IsWhiterun() {
		return exec.RunCmd(ctx, "swaymsg", "output", "HDMI-A-1", "mode", "1368x768").Run()
	}
	return fmt.Errorf("no sway reset logic for this host")
}
