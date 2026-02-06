package xclip

import (
	"context"
	"fmt"
	"io"
	"strings"
	"workspaced/pkg/exec"
)

type Driver struct{}

func (d *Driver) WriteImageReader(ctx context.Context, r io.Reader) error {
	if !exec.IsBinaryAvailable(ctx, "xclip") {
		return fmt.Errorf("xclip not found")
	}
	cmd := exec.RunCmd(ctx, "xclip", "-selection", "clipboard", "-t", "image/png")
	cmd.Stdin = r
	return cmd.Run()
}

func (d *Driver) WriteText(ctx context.Context, text string) error {
	if !exec.IsBinaryAvailable(ctx, "xclip") {
		return fmt.Errorf("xclip not found")
	}
	cmd := exec.RunCmd(ctx, "xclip", "-selection", "clipboard")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}
