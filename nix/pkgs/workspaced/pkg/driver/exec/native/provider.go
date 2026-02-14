package native

import (
	"context"
	"fmt"
	"os/exec"
	"workspaced/pkg/api"
	"workspaced/pkg/driver"
	execdriver "workspaced/pkg/driver/exec"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "exec_native"
}

func (p *Provider) Name() string {
	return "Native"
}

func (p *Provider) DefaultWeight() int {
	return driver.DefaultWeight
}

func (p *Provider) CheckCompatibility(ctx context.Context) error {
	// Always compatible on non-Termux systems
	return nil
}

func (p *Provider) New(ctx context.Context) (execdriver.Driver, error) {
	return &Driver{}, nil
}

type Driver struct{}

func (d *Driver) Run(ctx context.Context, name string, args ...string) *exec.Cmd {
	// Use standard exec.CommandContext - works fine on normal Linux/macOS
	return exec.CommandContext(ctx, name, args...)
}

func (d *Driver) Which(ctx context.Context, name string) (string, error) {
	// Use standard LookPath - works fine on normal systems
	path, err := exec.LookPath(name)
	if err != nil {
		return "", fmt.Errorf("%w: %s", api.ErrBinaryNotFound, name)
	}
	return path, nil
}

func init() {
	driver.Register[execdriver.Driver](&Provider{})
}
