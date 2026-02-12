package audio

import "context"

type Driver interface {
	SetVolume(ctx context.Context, volume float64) error
	GetVolume(ctx context.Context) (float64, error)
	ToggleMute(ctx context.Context) error
	GetMute(ctx context.Context) (bool, error)
	SinkName(ctx context.Context) (string, error)
}
