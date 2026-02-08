package xclip

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"io"
	"strings"
	dapi "workspaced/pkg/api"
	"workspaced/pkg/exec"
)

type Driver struct{}

func (d *Driver) WriteImage(ctx context.Context, img image.Image) error {
	if !exec.IsBinaryAvailable(ctx, "xclip") {
		return fmt.Errorf("%w: xclip", dapi.ErrBinaryNotFound)
	}
	pr, pw := io.Pipe()
	go func() {
		_ = png.Encode(pw, img)
		_ = pw.Close()
	}()

	cmd := exec.RunCmd(ctx, "xclip", "-selection", "clipboard", "-t", "image/png")
	cmd.Stdin = pr
	return cmd.Run()
}

func (d *Driver) WriteText(ctx context.Context, text string) error {
	if !exec.IsBinaryAvailable(ctx, "xclip") {
		return fmt.Errorf("%w: xclip", dapi.ErrBinaryNotFound)
	}
	cmd := exec.RunCmd(ctx, "xclip", "-selection", "clipboard")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}
