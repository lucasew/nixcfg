package xclip

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"io"
	"strings"
	dapi "workspaced/pkg/api"
	"workspaced/pkg/driver"
	"workspaced/pkg/driver/clipboard"
	"workspaced/pkg/exec"
)

func init() {
	driver.Register[clipboard.Driver](&Provider{})
}

type Provider struct{}

func (p *Provider) ID() string   { return "clipboard_xclip" }
func (p *Provider) Name() string { return "X11 (xclip)" }
func (p *Provider) DefaultWeight() int { return driver.DefaultWeight }

func (p *Provider) CheckCompatibility(ctx context.Context) error {
	if !exec.IsBinaryAvailable(ctx, "xclip") {
		return fmt.Errorf("xclip binary not found")
	}
	// Fallback driver, usually always valid if binary exists
	return nil
}

func (p *Provider) New(ctx context.Context) (clipboard.Driver, error) {
	return &Driver{}, nil
}

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
