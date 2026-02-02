package prelude

import (
	"fmt"
	"os"
	"os/exec"
)

// GenerateCompletion generates bash completion by calling workspaced
func GenerateCompletion() (string, error) {
	// Execute workspaced completion bash to get the completion code
	// This is only done once when generating the cache
	cmd := exec.Command("workspaced", "completion", "bash")
	cmd.Env = os.Environ()

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to generate completion: %w", err)
	}

	return string(output), nil
}
