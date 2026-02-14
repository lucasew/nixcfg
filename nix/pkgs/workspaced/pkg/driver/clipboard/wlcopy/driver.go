package wlcopy

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
	execdriver "workspaced/pkg/driver/exec"
	"workspaced/pkg/executil"
)

func init() {
	driver.Register[clipboard.Driver](&Provider{})
}

type Provider struct{}

func (p *Provider) ID() string   { return "clipboard_wlcopy" }
func (p *Provider) Name() string { return "Wayland (wl-copy)" }
func (p *Provider) DefaultWeight() int { return driver.DefaultWeight }

func (p *Provider) CheckCompatibility(ctx context.Context) error {
	if executil.GetEnv(ctx, "WAYLAND_DISPLAY") == "" {
		return fmt.Errorf("%w: WAYLAND_DISPLAY not set", driver.ErrIncompatible)
	}
	if !execdriver.IsBinaryAvailable(ctx, "wl-copy") {
		return fmt.Errorf("%w: wl-copy not found", driver.ErrIncompatible)
	}
	return nil
}

func (p *Provider) New(ctx context.Context) (clipboard.Driver, error) {
	return &Driver{}, nil
}

type Driver struct{}

func (d *Driver) WriteImage(ctx context.Context, img image.Image) error {
	if !execdriver.IsBinaryAvailable(ctx, "wl-copy") {
		return fmt.Errorf("%w: wl-copy", dapi.ErrBinaryNotFound)
	}
	pr, pw := io.Pipe()
	go func() {
		_ = png.Encode(pw, img)
		_ = pw.Close()
	}()

	cmd := execdriver.MustRun(ctx, "wl-copy", "-t", "image/png")
	cmd.Stdin = pr
	return cmd.Run()
}

func (d *Driver) WriteText(ctx context.Context, text string) error {
	if !execdriver.IsBinaryAvailable(ctx, "wl-copy") {
		return fmt.Errorf("%w: wl-copy", dapi.ErrBinaryNotFound)
	}
	cmd := execdriver.MustRun(ctx, "wl-copy")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}
