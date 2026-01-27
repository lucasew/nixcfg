package screen

import (
	"context"
	"strings"
	"workspaced/pkg/common"
)

func Lock(ctx context.Context) error {
	return common.RunCmd(ctx, "loginctl", "lock-session").Run()
}

func SetDPMS(ctx context.Context, on bool) error {
	state := "off"
	if on {
		state = "on"
	}

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
