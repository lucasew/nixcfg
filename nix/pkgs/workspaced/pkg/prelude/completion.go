package prelude

import (
	"fmt"
	"strings"
)

// GenerateCompletion generates bash completion using cobra API directly (no exec)
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
