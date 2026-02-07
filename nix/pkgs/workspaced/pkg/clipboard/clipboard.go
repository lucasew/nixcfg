package clipboard

import (
	"context"
	"fmt"
	"image"
	"os"
	"workspaced/pkg/clipboard/api"
	"workspaced/pkg/clipboard/termux"
	"workspaced/pkg/clipboard/wlcopy"
	"workspaced/pkg/clipboard/xclip"
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
	d, err := GetDriver(ctx)
	if err != nil {
		return err
	}
	return d.WriteImage(ctx, img)
}

func WriteText(ctx context.Context, text string) error {
	d, err := GetDriver(ctx)
	if err != nil {
		return err
	}
	return d.WriteText(ctx, text)
}
