package shellgen

import (
	"fmt"
	"os"
	"os/exec"
)

// GenerateMise generates the shell activation code for the `mise` tool.
//
// Unlike other generators that produce static content, this function invokes
// `mise activate bash` to dynamically generate the environment setup script.
// It first attempts to find `mise` in the PATH using `exec.LookPath`.
// If `mise` is not found, it returns an empty string to avoid breaking shell init.
func GenerateMise() (string, error) {
	// Dynamically find mise in PATH
	misePath, err := exec.LookPath("mise")
	if err != nil {
		return "", nil // Skip if mise not found in PATH
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
