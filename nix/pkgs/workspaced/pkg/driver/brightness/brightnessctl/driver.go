package brightnessctl

import (
	"context"
	"fmt"
	"workspaced/pkg/brightness/api"
	"workspaced/pkg/driver"
	"workspaced/pkg/exec"
)

func init() {
	driver.Register[api.Driver](&Provider{})
}

type Provider struct{}

func (p *Provider) Name() string { return "brightnessctl" }

func (p *Provider) CheckCompatibility(ctx context.Context) error {
	if !exec.IsBinaryAvailable(ctx, "brightnessctl") {
		return fmt.Errorf("%w: brightnessctl not found", driver.ErrIncompatible)
	}
	return nil
}

func (p *Provider) New(ctx context.Context) (api.Driver, error) {
	return &Driver{}, nil
}

type Driver struct{}

func (d *Driver) SetBrightness(ctx context.Context, arg string) error {
	if arg != "" {
		if err := exec.RunCmd(ctx, "brightnessctl", "s", arg).Run(); err != nil {
			return fmt.Errorf("failed to set brightness: %w", err)
		}
	}
	return nil
}
