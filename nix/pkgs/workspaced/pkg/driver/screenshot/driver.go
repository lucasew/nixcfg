package screenshot

import (
	"context"
	"image"
	api "workspaced/pkg/driver/wm"
)

type TargetType int

const (
	TargetAll TargetType = iota
	TargetOutput
	TargetWindow
	TargetSelection
)

type Driver interface {
	Capture(ctx context.Context, rect *api.Rect) (image.Image, error)
	SelectArea(ctx context.Context) (*api.Rect, error)
}
