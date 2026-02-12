package wofi

import (
	"context"
	"fmt"
	"strings"
	"workspaced/pkg/driver"
	"workspaced/pkg/driver/menu"
	"workspaced/pkg/exec"
)

func init() {
	driver.Register[menu.Driver](&Provider{})
}

type Provider struct{}

func (p *Provider) Name() string { return "Wofi" }

func (p *Provider) CheckCompatibility(ctx context.Context) error {
	if exec.GetEnv(ctx, "WAYLAND_DISPLAY") == "" {
		return fmt.Errorf("%w: WAYLAND_DISPLAY not set", driver.ErrIncompatible)
	}
	if !exec.IsBinaryAvailable(ctx, "wofi") {
		return fmt.Errorf("%w: wofi not found", driver.ErrIncompatible)
	}
	return nil
}

func (p *Provider) New(ctx context.Context) (menu.Driver, error) {
	return &Driver{}, nil
}

type Driver struct{}

func (d *Driver) Choose(ctx context.Context, opts menu.Options) (*menu.Item, error) {
	var input strings.Builder
	for _, item := range opts.Items {
		label := item.Label
		if label == "" {
			label = item.Value
		}
		input.WriteString(label)
		input.WriteString("\n")
	}

	args := []string{"--dmenu", "-p", opts.Prompt}

	cmd := exec.RunCmd(ctx, "wofi", args...)
	cmd.Stdin = strings.NewReader(input.String())

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	selected := strings.TrimSpace(string(out))
	if selected == "" {
		return nil, nil
	}

	for _, item := range opts.Items {
		label := item.Label
		if label == "" {
			label = item.Value
		}
		if selected == label {
			return &item, nil
		}
	}

	return &menu.Item{Label: selected, Value: selected}, nil
}

func (d *Driver) RunApp(ctx context.Context) error {
	return exec.RunCmd(ctx, "wofi", "--show", "drun").Run()
}

func (d *Driver) SwitchWindow(ctx context.Context) error {
	return exec.RunCmd(ctx, "wofi", "--show", "drun").Run()
}
