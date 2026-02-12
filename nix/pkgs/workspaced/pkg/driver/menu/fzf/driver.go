package fzf

import (
	"context"
	"fmt"
	"os"
	"strings"
	"workspaced/pkg/driver"
	"workspaced/pkg/driver/menu"
	"workspaced/pkg/exec"
)

func init() {
	driver.Register[menu.Driver](&Provider{})
}

type Provider struct{}

func (p *Provider) Name() string { return "FZF" }

func (p *Provider) CheckCompatibility(ctx context.Context) error {
	// Se tiver display gráfico, não usar FZF por padrão (dar preferência ao rofi/wofi)
	if exec.GetEnv(ctx, "DISPLAY") != "" || exec.GetEnv(ctx, "WAYLAND_DISPLAY") != "" {
		return fmt.Errorf("%w: graphical display available, skipping FZF", driver.ErrIncompatible)
	}
	if !exec.IsBinaryAvailable(ctx, "fzf") {
		return fmt.Errorf("%w: fzf not found", driver.ErrIncompatible)
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

	args := []string{"--prompt", opts.Prompt + "> "}

	cmd := exec.RunCmd(ctx, "fzf", args...)
	cmd.Stdin = strings.NewReader(input.String())
	cmd.Stderr = os.Stderr // fzf renders on stderr

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
	return fmt.Errorf("RunApp not implemented for FZF")
}

func (d *Driver) SwitchWindow(ctx context.Context) error {
	return fmt.Errorf("SwitchWindow not implemented for FZF")
}
