package terminal

import (
	"context"
	"workspaced/pkg/driver"
)

// Open opens the preferred terminal emulator.
func Open(ctx context.Context, opts Options) error {
	d, err := driver.Get[Driver](ctx)
	if err != nil {
		return err
	}
	return d.Open(ctx, opts)
}
