package termux

import (
	"context"
	"fmt"
	"os"
	"workspaced/pkg/driver"
	"workspaced/pkg/driver/opener"
	execdriver "workspaced/pkg/driver/exec"
)

func init() {
	driver.Register[opener.Driver](&Provider{})
}

type Provider struct{}

func (p *Provider) ID() string   { return "opener_termux" }
func (p *Provider) Name() string { return "termux-open" }
func (p *Provider) DefaultWeight() int { return driver.DefaultWeight }

func (p *Provider) CheckCompatibility(ctx context.Context) error {
	if os.Getenv("TERMUX_VERSION") == "" {
		return fmt.Errorf("%w: not running in Termux", driver.ErrIncompatible)
	}
	if !execdriver.IsBinaryAvailable(ctx, "termux-open") {
		return fmt.Errorf("%w: termux-open not found", driver.ErrIncompatible)
	}
	return nil
}

func (p *Provider) New(ctx context.Context) (opener.Driver, error) {
	return &Driver{}, nil
}

type Driver struct{}

func (d *Driver) Open(ctx context.Context, target string) error {
	return execdriver.MustRun(ctx, "termux-open", target).Start()
}
