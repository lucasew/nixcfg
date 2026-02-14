package executil

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
	"workspaced/pkg/types"
)

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

// GetEnv retrieves an environment variable from the context or the system.
func GetEnv(ctx context.Context, key string) string {
	if env, ok := ctx.Value(types.EnvKey).([]string); ok {
		for _, e := range env {
			if strings.HasPrefix(e, key+"=") {
				return e[len(key)+1:]
			}
		}
	}
	return os.Getenv(key)
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

// GetBinaryMtime returns the modification time of the current executable.
func GetBinaryMtime() (time.Time, error) {
	exePath, err := os.Executable()
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get executable path: %w", err)
	}

	info, err := os.Stat(exePath)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to stat executable: %w", err)
	}

	return info.ModTime(), nil
}
