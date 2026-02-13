package brightness

import (
	"context"
	"workspaced/pkg/driver"
)

func IncreaseBrightness(ctx context.Context) error {
	d, err := driver.Get[Driver](ctx)
	if err != nil {
		return err
	}
	status, err := d.Status(ctx)
	if err != nil {
		return err
	}
	newLevel := status.Brightness + 0.05
	if newLevel > 1.0 {
		newLevel = 1.0
	}
	if err := d.SetBrightness(ctx, newLevel); err != nil {
		return err
	}
	return ShowStatus(ctx)
}

func DecreaseBrightness(ctx context.Context) error {
	d, err := driver.Get[Driver](ctx)
	if err != nil {
		return err
	}
	status, err := d.Status(ctx)
	if err != nil {
		return err
	}
	newLevel := status.Brightness - 0.05
	if newLevel < 0 {
		newLevel = 0
	}
	if err := d.SetBrightness(ctx, newLevel); err != nil {
		return err
	}
	return ShowStatus(ctx)
}
