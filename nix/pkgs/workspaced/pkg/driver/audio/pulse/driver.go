package pulse

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"workspaced/pkg/driver"
	"workspaced/pkg/driver/audio"
	"workspaced/pkg/exec"
)

var sink = "@DEFAULT_SINK@"

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

func (d *Driver) SetVolume(ctx context.Context, level float64) error {

	if err := exec.RunCmd(ctx, "pactl", "set-sink-volume", sink, fmt.Sprintf("%d%%", int(level*100))).Run(); err != nil {
		return fmt.Errorf("failed to set volume: %w", err)
	}
	return nil
}

func (d *Driver) GetVolume(ctx context.Context) (float64, error) {
	volumeOut, err := exec.RunCmd(ctx, "pactl", "get-sink-volume", sink).Output()
	if err != nil {
		return 0, err
	}
	volumeStr := strings.TrimSpace(string(volumeOut))
	volumeStr = strings.TrimPrefix(volumeStr, "Volume: ")
	volumeStr = strings.TrimSuffix(volumeStr, "%")
	volume, err := strconv.Atoi(volumeStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse volume: %w", err)
	}
	return float64(volume) / 100, nil
}

func (d *Driver) GetMute(ctx context.Context) (bool, error) {
	muteOut, err := exec.RunCmd(ctx, "pactl", "get-sink-mute", sink).Output()
	if err != nil {
		return false, err
	}
	return strings.Contains(string(muteOut), "yes"), nil
}

func (d *Driver) SinkName(ctx context.Context) (string, error) {
	nameOut, err := exec.RunCmd(ctx, "pactl", "get-default-sink").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(nameOut)), nil
}

func (d *Driver) ToggleMute(ctx context.Context) error {
	mute, err := d.GetMute(ctx)
	if err != nil {
		return err
	}
	if mute {
		return exec.RunCmd(ctx, "pactl", "set-sink-mute", sink, "no").Run()
	}
	return exec.RunCmd(ctx, "pactl", "set-sink-mute", sink, "yes").Run()
}
