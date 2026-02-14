package foot

import (
	"context"
	"fmt"
	"workspaced/pkg/driver"
	"workspaced/pkg/driver/terminal"
	execdriver "workspaced/pkg/driver/exec"
	"workspaced/pkg/executil"
)

func init() {
	driver.Register[terminal.Driver](&Provider{})
}

type Provider struct{}

func (p *Provider) ID() string   { return "terminal_foot" }
func (p *Provider) Name() string { return "Foot" }
func (p *Provider) DefaultWeight() int { return driver.DefaultWeight }

func (p *Provider) CheckCompatibility(ctx context.Context) error {
	if executil.GetEnv(ctx, "WAYLAND_DISPLAY") == "" {
		return fmt.Errorf("%w: foot requires WAYLAND_DISPLAY", driver.ErrIncompatible)
	}
	if !execdriver.IsBinaryAvailable(ctx, "foot") {
		return fmt.Errorf("%w: foot not found", driver.ErrIncompatible)
	}
	return nil
}

func (p *Provider) New(ctx context.Context) (terminal.Driver, error) {
	return &Driver{}, nil
}

type Driver struct{}

func (d *Driver) Open(ctx context.Context, opts terminal.Options) error {
	args := []string{}
	if opts.Title != "" {
		args = append(args, "-T", opts.Title)
	}
	if opts.Command != "" {
		args = append(args, opts.Command)
		args = append(args, opts.Args...)
	}

	cmd := execdriver.MustRun(ctx, "foot", args...)
	return cmd.Start()
}
