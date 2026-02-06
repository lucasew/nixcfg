package clipboard

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"io"
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
	return nil, fmt.Errorf("no suitable clipboard driver found")
}

// WriteImage encodes a stdlib image.Image to PNG and writes it to the clipboard.
func WriteImage(ctx context.Context, img image.Image) error {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return fmt.Errorf("failed to encode image to PNG: %w", err)
	}
	return WriteImageReader(ctx, &buf)
}

// WriteImageReader writes raw image bytes (expected to be PNG) to the clipboard.
func WriteImageReader(ctx context.Context, r io.Reader) error {
	d, err := GetDriver(ctx)
	if err != nil {
		return err
	}
	return d.WriteImageReader(ctx, r)
}

func WriteText(ctx context.Context, text string) error {
	d, err := GetDriver(ctx)
	if err != nil {
		return err
	}
	return d.WriteText(ctx, text)
}
