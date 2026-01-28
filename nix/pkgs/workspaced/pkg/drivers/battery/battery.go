package battery

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Status string

const (
	Charging    Status = "Charging"
	Discharging Status = "Discharging"
	Full        Status = "Full"
	Unknown     Status = "Unknown"
)

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
