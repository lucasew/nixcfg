package opener

import (
	"context"
	"workspaced/pkg/driver"
)

func Open(ctx context.Context, target string) error {
	d, err := driver.Get[Driver](ctx)
	if err != nil {
		return err
	}
	return d.Open(ctx, target)
}
