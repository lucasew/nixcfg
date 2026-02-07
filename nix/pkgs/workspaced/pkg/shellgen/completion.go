package shellgen

import (
	"fmt"
	"strings"
)

// GenerateCompletion generates Bash completion scripts for the `workspaced` CLI.
//
// Instead of spawning a subprocess to run `workspaced completion bash`, this function
// uses the internal Cobra API (`GenBashCompletionV2`) directly on the `rootCommand`.
// This avoids the overhead of a fork/exec cycle, speeding up shell startup.
//
// Precondition: `SetRootCommand` must have been called with the active Cobra root command.
func GenerateCompletion() (string, error) {
	if rootCommand == nil {
		return "", fmt.Errorf("root command not set, call SetRootCommand first")
	}

	// Generate completion directly using cobra API (much faster than exec)
	var buf strings.Builder
	if err := rootCommand.GenBashCompletionV2(&buf, true); err != nil {
		return "", fmt.Errorf("failed to generate bash completion: %w", err)
	}

	return buf.String(), nil
}
