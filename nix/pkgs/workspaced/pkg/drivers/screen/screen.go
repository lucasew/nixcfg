package screen

import (
	"context"
	"fmt"
	"strings"
	"workspaced/pkg/common"
)

func Lock(ctx context.Context) error {
	common.GetLogger(ctx).Info("locking session")
	if err := common.RunCmd(ctx, "loginctl", "lock-session").Run(); err != nil {
		return err
	}
	return SetDPMS(ctx, false)
}

func SetDPMS(ctx context.Context, on bool) error {
	state := "off"
	if on {
		state = "on"
	}

	common.GetLogger(ctx).Info("setting DPMS", "state", state)
	rpc := common.GetRPC(ctx)
	if rpc == "swaymsg" {
		return common.RunCmd(ctx, "swaymsg", "output * dpms "+state).Run()
	}

	xsetArg := "off"
	if on {
		xsetArg = "on"
	}
	return common.RunCmd(ctx, "xset", "dpms", "force", xsetArg).Run()
}

func ToggleDPMS(ctx context.Context) error {
	isOn, err := IsDPMSOn(ctx)
	if err != nil {
		return err
	}
	return SetDPMS(ctx, !isOn)
}

func IsDPMSOn(ctx context.Context) (bool, error) {
	rpc := common.GetRPC(ctx)
	if rpc == "swaymsg" {
		out, err := common.RunCmd(ctx, "swaymsg", "-t", "get_outputs").Output()
		if err != nil {
			return false, err
		}
		return strings.Contains(string(out), `"dpms": true`), nil
	}

	out, err := common.RunCmd(ctx, "xset", "q").Output()
	if err != nil {
		return false, err
	}
	return strings.Contains(string(out), "Monitor is On"), nil
}

func Reset(ctx context.Context) error {
	if common.IsRiverwood() {
		// xrandr --output eDP-1 --mode 1366x768
		// xrandr --output HDMI-1 --mode 1366x768 --left-of eDP-1
		if err := common.RunCmd(ctx, "xrandr", "--output", "eDP-1", "--mode", "1366x768").Run(); err != nil {
			return err
		}
		return common.RunCmd(ctx, "xrandr", "--output", "HDMI-1", "--mode", "1366x768", "--left-of", "eDP-1").Run()
	}

	if common.IsWhiterun() {
		// xrandr --output HDMI-1 --mode 1368x768
		return common.RunCmd(ctx, "xrandr", "--output", "HDMI-1", "--mode", "1368x768").Run()
	}

	return fmt.Errorf("no reset logic for this host")
}
