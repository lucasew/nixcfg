package clipboard

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"workspaced/pkg/clipboard/api"
	"workspaced/pkg/driver"
)

// WriteImage encodes a stdlib image.Image to PNG and writes it to the clipboard.
func WriteImage(ctx context.Context, img image.Image) error {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return fmt.Errorf("failed to encode image to PNG: %w", err)
	}
	d, err := driver.Get[api.Driver](ctx)
	if err != nil {
		return err
	}
	// We need to decode back to image.Image or add WriteImageReader to interface?
	// Wait, the drivers now take image.Image.
	return d.WriteImage(ctx, img)
}

func WriteText(ctx context.Context, text string) error {
	d, err := driver.Get[api.Driver](ctx)
	if err != nil {
		return err
	}
	return d.WriteText(ctx, text)
}
