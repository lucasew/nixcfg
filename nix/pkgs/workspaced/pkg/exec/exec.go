package exec

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"workspaced/pkg/api"
	"workspaced/pkg/types"

	_ "workspaced/pkg/env" // Ensure PATH is set up
)

var (
	ErrCommandNotFound = fmt.Errorf("command not found")
)

// RunCmd creates an exec.Cmd with environment variables injected from the context.
// It ensures that the PATH includes EssentialPaths (via initialized environment).
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

// GetRPC determines the appropriate IPC command for the current window manager.
// It checks for HYPRLAND_INSTANCE_SIGNATURE for Hyprland,
// and WAYLAND_DISPLAY to decide between "swaymsg" (Wayland) and "i3-msg" (X11).
func GetRPC(ctx context.Context) string {
	if env, ok := ctx.Value(types.EnvKey).([]string); ok {
		for _, e := range env {
			if strings.HasPrefix(e, "HYPRLAND_INSTANCE_SIGNATURE=") {
				return "hyprctl"
			}
			if strings.HasPrefix(e, "WAYLAND_DISPLAY=") {
				return "swaymsg"
			}
		}
	}
	if os.Getenv("HYPRLAND_INSTANCE_SIGNATURE") != "" {
		return "hyprctl"
	}
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		return "swaymsg"
	}
	return "i3-msg"
}

// Which locates a command in the PATH without using os/exec.LookPath
// to avoid SIGSYS errors on Android/Go 1.24.
func Which(ctx context.Context, name string) (string, error) {
	if filepath.IsAbs(name) {
		if _, err := os.Stat(name); err == nil {
			return name, nil
		}
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
			return fullPath, nil
		}
	}
	return "", fmt.Errorf("%w: %s", api.ErrBinaryNotFound, name)
}

// IsBinaryAvailable checks if a command exists in the PATH using the internal Which implementation.
func IsBinaryAvailable(ctx context.Context, name string) bool {
	_, err := Which(ctx, name)
	return err == nil
}

// GetBinaryHash returns the SHA256 hash of the current executable
func GetBinaryHash() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	file, err := os.Open(exePath)
	if err != nil {
		return "", fmt.Errorf("failed to open executable: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to hash executable: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
