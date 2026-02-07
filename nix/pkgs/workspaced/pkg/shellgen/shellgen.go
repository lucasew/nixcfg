package shellgen

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

// Generator is a function that generates a snippet of shell initialization code.
// Generators are executed in parallel by Generate().
type Generator func() (string, error)

// rootCommand is the main CLI command, injected via SetRootCommand.
// It is required by generators that need to introspect the CLI (e.g., completion).
var rootCommand *cobra.Command

// SetRootCommand injects the root cobra command into the shellgen package.
// This allows the completion generator to introspect the command structure
// without creating a circular dependency or relying on global state initialized elsewhere.
func SetRootCommand(cmd *cobra.Command) {
	rootCommand = cmd
}

// generators maps a sortable key to a Generator function.
// The keys determine the order in which the generated snippets appear in the final output.
// Convention: "NN-name" where NN is a number (e.g., "00-colors", "10-completion").
var generators = map[string]Generator{
	"00-colors":     GenerateColors,
	"05-flags":      GenerateFlags,
	"06-daemon":     GenerateDaemon,
	"10-completion": GenerateCompletion,
	"15-mise":       GenerateMise,
	"20-history":    GenerateHistory,
}

// Generate executes all registered generators in parallel and returns the concatenated shell script.
//
// It performs the following steps:
// 1. Executes all generators concurrently to minimize latency.
// 2. Collects results and timing information (profiling enabled via WORKSPACED_PROFILE=1).
// 3. Sorts the outputs by their registry key to ensure deterministic order.
// 4. Concatenates the results into a single string.
//
// If any generator fails, an error is returned aggregating all failures.
func Generate() (string, error) {
	profile := os.Getenv("WORKSPACED_PROFILE") == "1"

	type result struct {
		key      string
		output   string
		err      error
		duration time.Duration
	}

	results := make(chan result, len(generators))
	var wg sync.WaitGroup

	// Execute all generators in parallel
	for key, gen := range generators {
		wg.Add(1)
		go func(k string, g Generator) {
			defer wg.Done()
			start := time.Now()
			output, err := g()
			results <- result{key: k, output: output, err: err, duration: time.Since(start)}
		}(key, gen)
	}

	// Wait and close
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	resultMap := make(map[string]string)
	timings := make(map[string]time.Duration)
	var errs []error
	for r := range results {
		if r.err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", r.key, r.err))
			continue
		}
		resultMap[r.key] = r.output
		timings[r.key] = r.duration
		if profile {
			fmt.Fprintf(os.Stderr, "    • %s: %v\n", r.key, r.duration)
		}
	}

	if len(errs) > 0 {
		return "", fmt.Errorf("generator errors: %v", errs)
	}

	// Build output in order (sorted by key)
	keys := make([]string, 0, len(resultMap))
	for k := range resultMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var output strings.Builder
	for _, key := range keys {
		output.WriteString(fmt.Sprintf("# Generated: %s\n", key))
		output.WriteString(resultMap[key])
		if !strings.HasSuffix(resultMap[key], "\n") {
			output.WriteString("\n")
		}
		output.WriteString("\n")
	}

	return output.String(), nil
}
