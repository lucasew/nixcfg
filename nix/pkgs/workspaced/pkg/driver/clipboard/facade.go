package clipboard

import (
	"context"
	"image"
	"workspaced/pkg/driver"
)

// WriteImage writes a stdlib image.Image to the clipboard using the available driver.
func WriteImage(ctx context.Context, img image.Image) error {
	d, err := driver.Get[Driver](ctx)
	if err != nil {
		return err
	}
	return d.WriteImage(ctx, img)
}

// WriteText writes text to the clipboard using the available driver.
func WriteText(ctx context.Context, text string) error {
	d, err := driver.Get[Driver](ctx)
	if err != nil {
		return err
	}
	return d.WriteText(ctx, text)
}
