package kitty

import (
	"context"
	"fmt"
	"workspaced/pkg/driver"
	"workspaced/pkg/driver/terminal"
	execdriver "workspaced/pkg/driver/exec"
)

func init() {
	driver.Register[terminal.Driver](&Provider{})
}

type Provider struct{}

func (p *Provider) ID() string   { return "terminal_kitty" }
func (p *Provider) Name() string { return "Kitty" }
func (p *Provider) DefaultWeight() int { return driver.DefaultWeight }

func (p *Provider) CheckCompatibility(ctx context.Context) error {
	if !execdriver.IsBinaryAvailable(ctx, "kitty") {
		return fmt.Errorf("%w: kitty not found", driver.ErrIncompatible)
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
		args = append(args, "--title", opts.Title)
	}
	if opts.Command != "" {
		args = append(args, opts.Command)
		args = append(args, opts.Args...)
	}

	cmd := execdriver.MustRun(ctx, "kitty", args...)
	return cmd.Start()
}
