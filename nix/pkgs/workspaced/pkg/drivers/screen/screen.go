package screen

import (
	"context"
	"os"
	"strings"
	"workspaced/pkg/exec"
	"workspaced/pkg/logging"
	"workspaced/pkg/types"
)

func GetDriver(ctx context.Context) (Driver, error) {
	rpc := exec.GetRPC(ctx)
	if rpc == "swaymsg" {
		return &SwayDriver{}, nil
	}

	display := os.Getenv("DISPLAY")
	if env, ok := ctx.Value(types.EnvKey).([]string); ok {
		for _, e := range env {
			if strings.HasPrefix(e, "DISPLAY=") {
				display = strings.TrimPrefix(e, "DISPLAY=")
				break
			}
		}
	}

	if display != "" {
		return &X11Driver{}, nil
	}
	return nil, ErrDriverNotFound
}

func Lock(ctx context.Context) error {
	logging.GetLogger(ctx).Info("locking session")
	return exec.RunCmd(ctx, "loginctl", "lock-session").Run()
}

func SetDPMS(ctx context.Context, on bool) error {
	d, err := GetDriver(ctx)
	if err != nil {
		return err
	}
	logging.GetLogger(ctx).Info("setting DPMS", "on", on)
	return d.SetDPMS(ctx, on)
}

func ToggleDPMS(ctx context.Context) error {
	d, err := GetDriver(ctx)
	if err != nil {
		return err
	}
	isOn, err := d.IsDPMSOn(ctx)
	if err != nil {
		return err
	}
	return d.SetDPMS(ctx, !isOn)
}

func IsDPMSOn(ctx context.Context) (bool, error) {
	d, err := GetDriver(ctx)
	if err != nil {
		return false, err
	}
	return d.IsDPMSOn(ctx)
}

func Reset(ctx context.Context) error {
	d, err := GetDriver(ctx)
	if err != nil {
		return err
	}
	logging.GetLogger(ctx).Info("resetting screen layout")
	return d.Reset(ctx)
}
