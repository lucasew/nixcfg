package brightness

import "context"

type Device struct {
	Name       string
	Brightness float32
}

type Driver interface {
	SetBrightness(ctx context.Context, arg string) error
	Status(ctx context.Context) (*Device, error)
}
