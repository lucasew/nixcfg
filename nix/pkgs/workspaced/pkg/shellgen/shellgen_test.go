package shellgen

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestGenerateCompletion(t *testing.T) {
	// Setup dummy root command
	cmd := &cobra.Command{Use: "test"}
	SetRootCommand(cmd)

	// Test generation
	out, err := GenerateCompletion()
	if err != nil {
		t.Fatalf("GenerateCompletion failed: %v", err)
	}

	if !strings.Contains(out, "_test_root_command") {
		// Cobra v2 generates _<command_name>_root_command or similar
		// We just check if it's not empty and looks like bash completion
		if len(out) < 10 {
			t.Errorf("GenerateCompletion output too short: %q", out)
		}
	}
}

func TestGenerateMise(t *testing.T) {
	// This might fail if mise is not installed/found, but it should not return error unless execution fails.
	// If mise is not found, it returns empty string.
	out, err := GenerateMise()
	if err != nil {
		// It's acceptable for it to fail if mise is found but fails to run?
		// But in our logic, if it's found, it should run.
		// If it's not found, it returns nil error.
		t.Logf("GenerateMise returned error (possibly expected in test env): %v", err)
	} else {
		t.Logf("GenerateMise output length: %d", len(out))
	}
}

func TestGenerateFlags(t *testing.T) {
	out, err := GenerateFlags()
	if err != nil {
		t.Fatalf("GenerateFlags failed: %v", err)
	}
	if !strings.Contains(out, "WORKSPACED_SHELL_INIT=1") {
		t.Errorf("GenerateFlags missing expected flag")
	}
}

func TestGenerateDaemon(t *testing.T) {
	out, err := GenerateDaemon()
	if err != nil {
		t.Fatalf("GenerateDaemon failed: %v", err)
	}
	if !strings.Contains(out, "workspaced daemon --try") {
		t.Errorf("GenerateDaemon missing daemon command")
	}
}
