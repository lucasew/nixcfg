package api

import (
	"context"
	"image"
)

type Driver interface {
	WriteImage(ctx context.Context, img image.Image) error
	WriteText(ctx context.Context, text string) error
}
