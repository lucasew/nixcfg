package battery

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Status represents the charging status of the battery.
type Status string

// Battery status constants.
const (
	Charging    Status = "Charging"
	Discharging Status = "Discharging"
	Full        Status = "Full"
	Unknown     Status = "Unknown"
)

// GetStatus retrieves the current battery status by reading from sysfs.
//
// It scans `/sys/class/power_supply/BAT*/status` to locate battery devices.
//
// Note: It currently picks the *first* matching battery found (e.g., BAT0).
// In multi-battery systems, this may not reflect the aggregate status.
func GetStatus(ctx context.Context) (Status, error) {
	matches, _ := filepath.Glob("/sys/class/power_supply/BAT*/status")
	if len(matches) == 0 {
		return Unknown, fmt.Errorf("no battery found")
	}

	// For now just get the first one
	data, err := os.ReadFile(matches[0])
	if err != nil {
		return Unknown, err
	}

	s := strings.TrimSpace(string(data))
	switch s {
	case "Charging":
		return Charging, nil
	case "Discharging":
		return Discharging, nil
	case "Full":
		return Full, nil
	default:
		return Status(s), nil
	}
}
