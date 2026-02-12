package linux

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"workspaced/pkg/driver"
	"workspaced/pkg/driver/battery"
)

func init() {
	driver.Register[battery.Driver](&Provider{})
}

type Provider struct{}

// CheckCompatibility implements [driver.DriverProvider].
func (p *Provider) CheckCompatibility(ctx context.Context) error {
	matches, _ := filepath.Glob("/sys/class/power_supply/BAT*/status")
	if len(matches) == 0 {
		return fmt.Errorf("%w: /sys/class/power_supply/BAT*/status", driver.ErrIncompatible)
	}
	return nil
}

// Name implements [driver.DriverProvider].
func (p *Provider) Name() string {
	return "linux"
}

// New implements [driver.DriverProvider].
func (p *Provider) New(ctx context.Context) (battery.Driver, error) {
	return &Driver{}, nil
}

type Driver struct{}

// GetStatus retrieves the current battery status.
// It scans /sys/class/power_supply/BAT*/status to find the first available battery
// and reads its status.
func (d *Driver) BatteryStatus(ctx context.Context) (battery.Status, error) {
	matches, _ := filepath.Glob("/sys/class/power_supply/BAT*/status")
	if len(matches) == 0 {
		return battery.Unknown, fmt.Errorf("no battery found")
	}

	// For now just get the first one
	data, err := os.ReadFile(matches[0])
	if err != nil {
		return battery.Unknown, err
	}

	s := strings.TrimSpace(string(data))
	switch s {
	case "Charging":
		return battery.Charging, nil
	case "Discharging":
		return battery.Discharging, nil
	case "Full":
		return battery.Full, nil
	default:
		return battery.Status(s), nil
	}
}
