package termux

import (
	"context"
	"fmt"
	"image"
	"os"
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

func (p *Provider) Name() string { return "Termux" }

func (p *Provider) CheckCompatibility(ctx context.Context) error {
	if os.Getenv("TERMUX_VERSION") == "" && !exec.IsBinaryAvailable(ctx, "termux-clipboard-set") {
		return fmt.Errorf("%w: termux not detected", driver.ErrIncompatible)
	}
	return nil
}

func (p *Provider) New(ctx context.Context) (api.Driver, error) {
	return &Driver{}, nil
}

type Driver struct{}

func (d *Driver) WriteImage(ctx context.Context, img image.Image) error {
	return fmt.Errorf("%w: writing images to clipboard is not supported on Termux", dapi.ErrNotSupported)
}

func (d *Driver) WriteText(ctx context.Context, text string) error {
	if !exec.IsBinaryAvailable(ctx, "termux-clipboard-set") {
		return fmt.Errorf("%w: termux-clipboard-set (install termux-api)", dapi.ErrBinaryNotFound)
	}
	cmd := exec.RunCmd(ctx, "termux-clipboard-set")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}
