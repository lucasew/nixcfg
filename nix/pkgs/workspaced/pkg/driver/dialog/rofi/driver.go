package rofi

import (
	"context"
	"fmt"
	"strings"
	"workspaced/pkg/driver"
	"workspaced/pkg/driver/dialog"
	"workspaced/pkg/exec"
)

func init() {
	driver.Register[dialog.Chooser](&ChooserProvider{})
	driver.Register[dialog.Driver](&FullDriverProvider{})
}

type ChooserProvider struct{}

func (p *ChooserProvider) ID() string                                      { return "rofi" }
func (p *ChooserProvider) Name() string                                    { return "Rofi" }
func (p *ChooserProvider) DefaultWeight() int                              { return driver.DefaultWeight }
func (p *ChooserProvider) CheckCompatibility(ctx context.Context) error    { return checkRofi(ctx) }
func (p *ChooserProvider) New(ctx context.Context) (dialog.Chooser, error) { return &Driver{}, nil }

type FullDriverProvider struct{}

func (p *FullDriverProvider) ID() string                                     { return "rofi" }
func (p *FullDriverProvider) Name() string                                   { return "Rofi" }
func (p *FullDriverProvider) DefaultWeight() int                             { return driver.DefaultWeight }
func (p *FullDriverProvider) CheckCompatibility(ctx context.Context) error   { return checkRofi(ctx) }
func (p *FullDriverProvider) New(ctx context.Context) (dialog.Driver, error) { return &Driver{}, nil }

func checkRofi(ctx context.Context) error {
	if exec.GetEnv(ctx, "DISPLAY") == "" && exec.GetEnv(ctx, "WAYLAND_DISPLAY") == "" {
		return fmt.Errorf("%w: neither DISPLAY nor WAYLAND_DISPLAY set", driver.ErrIncompatible)
	}
	if !exec.IsBinaryAvailable(ctx, "rofi") {
		return fmt.Errorf("%w: rofi not found", driver.ErrIncompatible)
	}
	return nil
}

type Driver struct{}

func (d *Driver) Choose(ctx context.Context, opts dialog.ChooseOptions) (*dialog.Item, error) {
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

	return &dialog.Item{Label: selected, Value: selected}, nil
}

func (d *Driver) RunApp(ctx context.Context) error {
	return exec.RunCmd(ctx, "rofi", "-show", "combi", "-combi-modi", "drun", "-show-icons").Run()
}

func (d *Driver) SwitchWindow(ctx context.Context) error {
	return exec.RunCmd(ctx, "rofi", "-show", "combi", "-combi-modi", "window", "-show-icons").Run()
}
