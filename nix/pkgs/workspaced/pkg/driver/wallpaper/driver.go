package wallpaper

import (
	"context"
)

type Driver interface {
	SetStatic(ctx context.Context, path string) error
}
