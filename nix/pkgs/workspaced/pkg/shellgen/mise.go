package shellgen

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// GenerateMise generates mise activation code
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
