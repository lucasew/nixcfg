package tray

import (
	"context"
	"workspaced/pkg/driver"
)

// GetDefault returns the appropriate tray driver for the current environment.
func GetDefault(ctx context.Context) (Driver, error) {
	return driver.Get[Driver](ctx)
}
