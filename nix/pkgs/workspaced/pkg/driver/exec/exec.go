package exec

import (
	"context"
	"os/exec"
	"workspaced/pkg/driver"
)

// Driver provides platform-specific command execution.
type Driver interface {
	// Run creates an exec.Cmd configured for the platform.
	Run(ctx context.Context, name string, args ...string) *exec.Cmd

	// Which locates a command in PATH and returns its full path.
	Which(ctx context.Context, name string) (string, error)
}

// IsBinaryAvailable checks if a command exists in PATH using the selected driver.
func IsBinaryAvailable(ctx context.Context, name string) bool {
	d, err := driver.Get[Driver](ctx)
	if err != nil {
		return false
	}
	_, err = d.Which(ctx, name)
	return err == nil
}

// Run creates an exec.Cmd using the selected driver.
func Run(ctx context.Context, name string, args ...string) (*exec.Cmd, error) {
	d, err := driver.Get[Driver](ctx)
	if err != nil {
		return nil, err
	}
	return d.Run(ctx, name, args...), nil
}

// Which locates a command in PATH using the selected driver.
func Which(ctx context.Context, name string) (string, error) {
	d, err := driver.Get[Driver](ctx)
	if err != nil {
		return "", err
	}
	return d.Which(ctx, name)
}

// MustRun creates and returns an exec.Cmd using the selected driver.
// Panics if the driver cannot be loaded (should only happen during initialization).
// Use this for compatibility with code that expects *exec.Cmd directly.
func MustRun(ctx context.Context, name string, args ...string) *exec.Cmd {
	cmd, err := Run(ctx, name, args...)
	if err != nil {
		// Fallback to direct exec if driver fails
		return exec.CommandContext(ctx, name, args...)
	}
	return cmd
}
