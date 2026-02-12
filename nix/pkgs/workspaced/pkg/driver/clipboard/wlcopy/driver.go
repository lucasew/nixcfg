package wlcopy

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"io"
	"strings"
	dapi "workspaced/pkg/api"
	"workspaced/pkg/clipboard/api"
	"workspaced/pkg/driver"
	"workspaced/pkg/exec"
)

func init() {
	driver.Register[api.Driver](&Provider{})
}

type Provider struct{}

func (p *Provider) Name() string { return "Wayland (wl-copy)" }

func (p *Provider) CheckCompatibility(ctx context.Context) error {
	// Logic from original GetDriver
	rpc := exec.GetRPC(ctx)
	if rpc != "swaymsg" && rpc != "hyprctl" {
		return fmt.Errorf("not running inside sway or hyprland (rpc=%s)", rpc)
	}
	if !exec.IsBinaryAvailable(ctx, "wl-copy") {
		return fmt.Errorf("wl-copy binary not found")
	}
	return nil
}

func (p *Provider) New(ctx context.Context) (api.Driver, error) {
	return &Driver{}, nil
}

type Driver struct{}

func (d *Driver) WriteImage(ctx context.Context, img image.Image) error {
	if !exec.IsBinaryAvailable(ctx, "wl-copy") {
		return fmt.Errorf("%w: wl-copy", dapi.ErrBinaryNotFound)
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
		return fmt.Errorf("%w: wl-copy", dapi.ErrBinaryNotFound)
	}
	cmd := exec.RunCmd(ctx, "wl-copy")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}
