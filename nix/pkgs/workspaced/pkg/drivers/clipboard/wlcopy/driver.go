package wlcopy

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"io"
	"strings"
	"workspaced/pkg/exec"
)

type Driver struct{}

func (d *Driver) WriteImage(ctx context.Context, img image.Image) error {
	if !exec.IsBinaryAvailable(ctx, "wl-copy") {
		return fmt.Errorf("wl-copy not found")
	}
	pr, pw := io.Pipe()
	go func() {
		_ = png.Encode(pw, img)
		_ = pw.Close()
	}()

	cmd := exec.RunCmd(ctx, "wl-copy", "-t", "image/png")
	cmd.Stdin = pr
	return cmd.Run()
}

func (d *Driver) WriteText(ctx context.Context, text string) error {
	if !exec.IsBinaryAvailable(ctx, "wl-copy") {
		return fmt.Errorf("wl-copy not found")
	}
	cmd := exec.RunCmd(ctx, "wl-copy")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}
