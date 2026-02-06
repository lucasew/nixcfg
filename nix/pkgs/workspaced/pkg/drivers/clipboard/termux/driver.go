package termux

import (
	"context"
	"fmt"
	"image"
	"strings"
	"workspaced/pkg/exec"
)

type Driver struct{}

func (d *Driver) WriteImage(ctx context.Context, img image.Image) error {
	return fmt.Errorf("writing images to clipboard is not supported on Termux")
}

func (d *Driver) WriteText(ctx context.Context, text string) error {
	if !exec.IsBinaryAvailable(ctx, "termux-clipboard-set") {
		return fmt.Errorf("termux-clipboard-set not found (install termux-api)")
	}
	cmd := exec.RunCmd(ctx, "termux-clipboard-set")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}
