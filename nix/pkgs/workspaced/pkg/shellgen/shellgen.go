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

// Generator is a function that generates shell code
type Generator func() (string, error)

// rootCommand is set by SetRootCommand and used by generators that need it
var rootCommand *cobra.Command

// SetRootCommand sets the root command for generators that need it (e.g., completion)
func SetRootCommand(cmd *cobra.Command) {
	rootCommand = cmd
}

// generators maps order/name to generator functions
var generators = map[string]Generator{
	"05-flags":      GenerateFlags,
	"06-daemon":     GenerateDaemon,
	"10-completion": GenerateCompletion,
	"15-mise":       GenerateMise,
	"20-history":    GenerateHistory,
}

// Generate executes all generators in parallel and returns ordered output
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
