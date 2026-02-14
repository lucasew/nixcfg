package termux

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"workspaced/pkg/api"
	"workspaced/pkg/driver"
	execdriver "workspaced/pkg/driver/exec"
	"workspaced/pkg/types"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "exec_termux"
}

func (p *Provider) Name() string {
	return "Termux"
}

func (p *Provider) DefaultWeight() int {
	return driver.DefaultWeight
}

func (p *Provider) CheckCompatibility(ctx context.Context) error {
	// Check if we're running in Termux
	if os.Getenv("TERMUX_VERSION") == "" {
		return fmt.Errorf("%w: not running in Termux", driver.ErrIncompatible)
	}
	return nil
}

func (p *Provider) New(ctx context.Context) (execdriver.Driver, error) {
	return &Driver{}, nil
}

type Driver struct{}

func (d *Driver) Run(ctx context.Context, name string, args ...string) *exec.Cmd {
	// Resolve the full path using custom Which to avoid SIGSYS on Android
	fullPath, err := d.Which(ctx, name)
	if err != nil {
		// If Which fails, fall back to the original name
		// This allows exec.CommandContext to handle the error properly
		fullPath = name
	}
	cmd := exec.CommandContext(ctx, fullPath, args...)
	return cmd
}

func (d *Driver) Which(ctx context.Context, name string) (string, error) {
	// Custom Which implementation to avoid SIGSYS errors on Android/Termux
	// Do not use os/exec.LookPath as it can trigger SIGSYS on Android with Go 1.24+

	if filepath.IsAbs(name) {
		if _, err := os.Stat(name); err == nil {
			slog.Debug("which", "binary", name, "result", name)
			return name, nil
		}
		slog.Debug("which", "binary", name, "result", api.ErrBinaryNotFound)
		return "", fmt.Errorf("%w: %s", api.ErrBinaryNotFound, name)
	}

	path := os.Getenv("PATH")
	if env, ok := ctx.Value(types.EnvKey).([]string); ok {
		for _, e := range env {
			if strings.HasPrefix(e, "PATH=") {
				path = strings.TrimPrefix(e, "PATH=")
				break
			}
		}
	}

	for _, dir := range filepath.SplitList(path) {
		fullPath := filepath.Join(dir, name)
		if info, err := os.Stat(fullPath); err == nil && !info.IsDir() {
			slog.Debug("which", "binary", name, "result", fullPath)
			return fullPath, nil
		}
	}
	slog.Debug("which", "binary", name, "result", api.ErrBinaryNotFound)
	return "", fmt.Errorf("%w: %s", api.ErrBinaryNotFound, name)
}

func init() {
	driver.Register[execdriver.Driver](&Provider{})
}
