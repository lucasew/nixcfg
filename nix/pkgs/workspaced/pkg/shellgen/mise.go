package shellgen

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// GenerateMise generates the shell activation code for the `mise` tool.
//
// Unlike other generators that produce static content, this function invokes
// `mise activate bash` to dynamically generate the environment setup script.
// It assumes `mise` is installed at `~/.local/bin/mise` (standard install location).
// If `mise` is not found, it returns an empty string to avoid breaking shell init.
func GenerateMise() (string, error) {
	misePath := filepath.Join(os.Getenv("HOME"), ".local", "bin", "mise")

	// Check if mise exists
	if _, err := os.Stat(misePath); os.IsNotExist(err) {
		return "", nil // Skip if mise not installed
	}

	// Execute mise activate bash
	cmd := exec.Command(misePath, "activate", "bash")
	cmd.Env = os.Environ()

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to generate mise activation: %w", err)
	}

	return string(output), nil
}
