package pulse

import (
	"context"
	"fmt"
	"workspaced/pkg/driver"
	"workspaced/pkg/driver/audio"
	"workspaced/pkg/exec"
)

func init() {
	driver.Register[audio.Driver](&Provider{})
}

type Provider struct{}

func (p *Provider) Name() string { return "PulseAudio (pactl)" }

func (p *Provider) CheckCompatibility(ctx context.Context) error {
	if !exec.IsBinaryAvailable(ctx, "pactl") {
		return fmt.Errorf("%w: pactl not found", driver.ErrIncompatible)
	}
	// Check if PulseAudio server is reachable?
	// For now, binary existence is a good enough proxy for "can try"
	return nil
}

func (p *Provider) New(ctx context.Context) (audio.Driver, error) {
	return &Driver{}, nil
}

type Driver struct{}

func (d *Driver) SetVolume(ctx context.Context, arg string) error {
	sink := "@DEFAULT_SINK@"
	if err := exec.RunCmd(ctx, "pactl", "set-sink-volume", sink, arg).Run(); err != nil {
		return fmt.Errorf("failed to set volume: %w", err)
	}
	return nil
}
