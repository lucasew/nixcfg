package maim

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/png"
	"strconv"
	"strings"
	"workspaced/pkg/driver"
	"workspaced/pkg/driver/screenshot"
	api "workspaced/pkg/driver/wm"
	"workspaced/pkg/exec"
)

func init() {
	driver.Register[screenshot.Driver](&Provider{})
}

type Provider struct{}

func (p *Provider) Name() string { return "Maim (X11)" }

func (p *Provider) CheckCompatibility(ctx context.Context) error {
	if exec.GetEnv(ctx, "DISPLAY") == "" {
		return fmt.Errorf("%w: DISPLAY not set", driver.ErrIncompatible)
	}
	if !exec.IsBinaryAvailable(ctx, "maim") {
		return fmt.Errorf("%w: maim not found", driver.ErrIncompatible)
	}
	return nil
}

func (p *Provider) New(ctx context.Context) (screenshot.Driver, error) {
	return &Driver{}, nil
}

type Driver struct{}

func (d *Driver) SelectArea(ctx context.Context) (*api.Rect, error) {
	// maim uses slop for selection.
	// maim -g $(slop) ... is common.
	// We can run slop directly to get the geometry.
	if !exec.IsBinaryAvailable(ctx, "slop") {
		return nil, fmt.Errorf("slop not found for selection")
	}
	out, err := exec.RunCmd(ctx, "slop", "-f", "%x %y %w %h").Output()
	if err != nil {
		return nil, err // likely canceled
	}
	raw := strings.TrimSpace(string(out))
	parts := strings.Fields(raw)
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid slop output: %q", raw)
	}
	x, _ := strconv.Atoi(parts[0])
	y, _ := strconv.Atoi(parts[1])
	w, _ := strconv.Atoi(parts[2])
	h, _ := strconv.Atoi(parts[3])
	return &api.Rect{X: x, Y: y, Width: w, Height: h}, nil
}

func (d *Driver) Capture(ctx context.Context, rect *api.Rect) (image.Image, error) {
	args := []string{}
	if rect != nil {
		// maim geometry format: WxH+X+Y
		args = append(args, "-g", fmt.Sprintf("%dx%d+%d+%d", rect.Width, rect.Height, rect.X, rect.Y))
	}

	cmd := exec.RunCmd(ctx, "maim", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("maim failed: %w", err)
	}

	img, _, err := image.Decode(bytes.NewReader(out))
	if err != nil {
		return nil, fmt.Errorf("failed to decode maim output: %w", err)
	}

	return img, nil
}
