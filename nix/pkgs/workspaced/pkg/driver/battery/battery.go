package battery

import (
	"context"
)

type Driver interface {
	BatteryStatus(ctx context.Context) (Status, error)
}

// Status represents the charging status of the battery.
type Status string

// Battery status constants.
const (
	Charging    Status = "Charging"
	Discharging Status = "Discharging"
	Full        Status = "Full"
	Unknown     Status = "Unknown"
)
