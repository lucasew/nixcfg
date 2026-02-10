package palette

import (
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"workspaced/pkg/palette/api"
	"workspaced/pkg/palette/genetic"
)

// GetDriver returns a palette extraction driver by name
func GetDriver(ctx context.Context, name string) (api.Driver, error) {
	switch name {
	case "genetic":
		return &genetic.Driver{}, nil
	default:
		return nil, fmt.Errorf("%w: %s", api.ErrDriverNotFound, name)
	}
}

// ExtractFromFile loads an image from a file and extracts a color palette
func ExtractFromFile(ctx context.Context, path string, driver string, opts api.Options) (*api.Palette, error) {
	// Load image from file
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %w", err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Get driver and extract
	d, err := GetDriver(ctx, driver)
	if err != nil {
		return nil, err
	}

	return d.Extract(ctx, img, opts)
}
