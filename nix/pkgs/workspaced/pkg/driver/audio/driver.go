package audio

import "context"

type Driver interface {
	SetVolume(ctx context.Context, arg string) error
}
