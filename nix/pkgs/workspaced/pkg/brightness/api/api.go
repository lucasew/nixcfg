package api

import (
	"context"
)

type Driver interface {
	SetBrightness(ctx context.Context, arg string) error
}
