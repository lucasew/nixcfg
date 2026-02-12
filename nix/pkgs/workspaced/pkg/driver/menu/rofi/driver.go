package rofi

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

func (p *Provider) Name() string { return "Rofi" }

func (p *Provider) CheckCompatibility(ctx context.Context) error {
	if exec.GetEnv(ctx, "DISPLAY") == "" && exec.GetEnv(ctx, "WAYLAND_DISPLAY") == "" {
		return fmt.Errorf("%w: neither DISPLAY nor WAYLAND_DISPLAY set", driver.ErrIncompatible)
	}
	if !exec.IsBinaryAvailable(ctx, "rofi") {
		return fmt.Errorf("%w: rofi not found", driver.ErrIncompatible)
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
		if item.Icon != "" {
			input.WriteString("\x00icon\x1f")
			input.WriteString(item.Icon)
		}
		input.WriteString("\n")
	}

	args := []string{"-dmenu", "-p", opts.Prompt}
	// Add some standard styling if desired, or keep it minimal
	args = append(args, "-show-icons")

	cmd := exec.RunCmd(ctx, "rofi", args...)
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

	// If not found in items but user typed something, return a virtual item
	return &menu.Item{Label: selected, Value: selected}, nil
}

func (d *Driver) RunApp(ctx context.Context) error {
	return exec.RunCmd(ctx, "rofi", "-show", "combi", "-combi-modi", "drun", "-show-icons").Run()
}

func (d *Driver) SwitchWindow(ctx context.Context) error {
	return exec.RunCmd(ctx, "rofi", "-show", "combi", "-combi-modi", "window", "-show-icons").Run()
}
