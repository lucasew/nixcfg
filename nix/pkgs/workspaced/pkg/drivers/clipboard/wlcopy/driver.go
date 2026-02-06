package wlcopy

import (
	"context"
	"fmt"
	"io"
	"strings"
	"workspaced/pkg/exec"
)

type Driver struct{}

func (d *Driver) WriteImageReader(ctx context.Context, r io.Reader) error {
	if !exec.IsBinaryAvailable(ctx, "wl-copy") {
		return fmt.Errorf("wl-copy not found")
	}
	cmd := exec.RunCmd(ctx, "wl-copy", "-t", "image/png")
	cmd.Stdin = r
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
