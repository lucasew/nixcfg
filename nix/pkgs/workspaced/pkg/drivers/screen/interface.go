package screen

import (
	"context"
	"workspaced/pkg/drivers/api"
)

var (
	ErrDriverNotFound = api.ErrDriverNotFound
)

type Driver interface {
	SetDPMS(ctx context.Context, on bool) error
	IsDPMSOn(ctx context.Context) (bool, error)
	Reset(ctx context.Context) error
}
