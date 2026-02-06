package clipboard

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"os"
	"workspaced/pkg/drivers/clipboard/api"
	"workspaced/pkg/drivers/clipboard/termux"
	"workspaced/pkg/drivers/clipboard/wlcopy"
	"workspaced/pkg/drivers/clipboard/xclip"
	"workspaced/pkg/exec"
)

func GetDriver(ctx context.Context) (api.Driver, error) {
	// Detect Termux first
	if os.Getenv("TERMUX_VERSION") != "" || exec.IsBinaryAvailable(ctx, "termux-clipboard-set") {
		return &termux.Driver{}, nil
	}

	rpc := exec.GetRPC(ctx)
	if rpc == "swaymsg" || rpc == "hyprctl" {
		if exec.IsBinaryAvailable(ctx, "wl-copy") {
			return &wlcopy.Driver{}, nil
		}
	}
	if exec.IsBinaryAvailable(ctx, "xclip") {
		return &xclip.Driver{}, nil
	}
	return nil, api.ErrDriverNotFound
}

// WriteImage encodes a stdlib image.Image to PNG and writes it to the clipboard.
func WriteImage(ctx context.Context, img image.Image) error {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return fmt.Errorf("failed to encode image to PNG: %w", err)
	}
	d, err := GetDriver(ctx)
	if err != nil {
		return err
	}
	// We need to decode back to image.Image or add WriteImageReader to interface?
	// Wait, the drivers now take image.Image.
	return d.WriteImage(ctx, img)
}

func WriteText(ctx context.Context, text string) error {
	d, err := GetDriver(ctx)
	if err != nil {
		return err
	}
	return d.WriteText(ctx, text)
}
