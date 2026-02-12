package clipboard

import (
	"context"
	"image"
	"workspaced/pkg/api"
)

var ErrDriverNotFound = api.ErrDriverNotFound

type Driver interface {
	WriteImage(ctx context.Context, img image.Image) error
	WriteText(ctx context.Context, text string) error
}
