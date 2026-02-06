package exec

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"workspaced/pkg/host"
	"workspaced/pkg/types"
)

var (
	ErrCommandNotFound = fmt.Errorf("command not found")
)

// EssentialPaths defines the list of directories that must be present in the PATH
// for the application to function correctly on NixOS.
// These typically include locations for wrapped binaries and current system software.
var EssentialPaths = []string{"/run/wrappers/bin", "/run/current-system/sw/bin"}

func init() {
	if home, err := os.UserHomeDir(); err == nil {
		EssentialPaths = append(EssentialPaths, filepath.Join(home, ".nix-profile/bin"))
	}
	if root, err := host.GetDotfilesRoot(); err == nil && root != "" {
		EssentialPaths = append(EssentialPaths, filepath.Join(root, "bin/shim"))
	}
	if dataDir, err := host.GetUserDataDir(); err == nil && dataDir != "" {
		EssentialPaths = append(EssentialPaths, filepath.Join(dataDir, "shim/global"))
	}
	newPath := strings.Split(os.Getenv("PATH"), ":")

	for _, path := range EssentialPaths {
		if !slices.Contains(newPath, path) {
			newPath = append([]string{path}, newPath...)
		}
	}
	if err := os.Setenv("PATH", strings.Join(newPath, ":")); err != nil {
		panic(err)
	}
}

// RunCmd creates an exec.Cmd with environment variables injected from the context.
// It ensures that the PATH includes EssentialPaths.
// It uses the custom Which implementation to avoid SIGSYS errors on Android/Termux.
func RunCmd(ctx context.Context, name string, args ...string) *exec.Cmd {
	// Resolve the full path using our custom Which to avoid SIGSYS on Android
	fullPath, err := Which(ctx, name)
	if err != nil {
		// If Which fails, fall back to the original name
		// This allows exec.CommandContext to handle the error properly
		fullPath = name
	}
	cmd := exec.CommandContext(ctx, fullPath, args...)
	return cmd
}

// InheritContextWriters configures the command's Stdout and Stderr to write to the writers
// stored in the context, allowing output capture or redirection.
func InheritContextWriters(ctx context.Context, cmd *exec.Cmd) {
	if stdout, ok := ctx.Value(types.StdoutKey).(io.Writer); ok {
		cmd.Stdout = stdout
	}
	if stderr, ok := ctx.Value(types.StderrKey).(io.Writer); ok {
		cmd.Stderr = stderr
	}
}

// Which locates a command in the PATH without using os/exec.LookPath
// to avoid SIGSYS errors on Android/Go 1.24.
func Which(ctx context.Context, name string) (string, error) {
	if filepath.IsAbs(name) {
		if _, err := os.Stat(name); err == nil {
			return name, nil
		}
		return "", fmt.Errorf("file not found: %s", name)
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
			return fullPath, nil
		}
	}
	return "", fmt.Errorf("%w: %s", ErrCommandNotFound, name)
}

// IsBinaryAvailable checks if a command exists in the PATH using the internal Which implementation.
func IsBinaryAvailable(ctx context.Context, name string) bool {
	_, err := Which(ctx, name)
	return err == nil
}
