package tray

import (
	"context"
)

// Driver is the interface for tray implementations.
type Driver interface {
	Run(ctx context.Context) error
	AddMenuItem(label string, callback func())
	Close()
}
