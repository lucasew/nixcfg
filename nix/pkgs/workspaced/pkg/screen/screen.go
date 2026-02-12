package screen

import (
	"context"
	"workspaced/pkg/driver"
	"workspaced/pkg/exec"
	"workspaced/pkg/logging"
)

func Lock(ctx context.Context) error {
	logging.GetLogger(ctx).Info("locking session")
	return exec.RunCmd(ctx, "loginctl", "lock-session").Run()
}

func SetDPMS(ctx context.Context, on bool) error {
	d, err := driver.Get[Driver](ctx)
	if err != nil {
		return err
	}
	logging.GetLogger(ctx).Info("setting DPMS", "on", on)
	return d.SetDPMS(ctx, on)
}

func ToggleDPMS(ctx context.Context) error {
	d, err := driver.Get[Driver](ctx)
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
	d, err := driver.Get[Driver](ctx)
	if err != nil {
		return false, err
	}
	return d.IsDPMSOn(ctx)
}

func Reset(ctx context.Context) error {
	d, err := driver.Get[Driver](ctx)
	if err != nil {
		return err
	}
	logging.GetLogger(ctx).Info("resetting screen layout")
	return d.Reset(ctx)
}
