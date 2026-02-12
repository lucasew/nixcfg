package brightness

import (
	"context"
)

type Device struct {
	Name       string
	Brightness float64
}

type Driver interface {
	SetBrightness(ctx context.Context, brightness float64) error
	Status(ctx context.Context) (*Device, error)
}
