package setup

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"workspaced/pkg/common"
)

// SetupTermuxShortcuts provisions Termux:Widget shortcuts by creating wrapper scripts in ~/.shortcuts.
//
// This function:
// 1. Verifies the environment is Termux (via IsPhone).
// 2. Scans for shortcut definitions in `$DOTFILES/bin/_shortcuts/termux`.
// 3. Creates executable wrappers in `~/.shortcuts` that invoke the corresponding dotfiles script via `source_me`.
//
// These shortcuts allow easy access to workspaced commands from the Android home screen
// using the Termux:Widget app.
func SetupTermuxShortcuts(ctx context.Context) error {
	if !common.IsPhone() {
		return fmt.Errorf("this command only works on phone (Termux)")
	}

	home, _ := os.UserHomeDir()
	shortcutDir := filepath.Join(home, ".shortcuts")
	_ = os.MkdirAll(shortcutDir, 0755)

	dotfiles, err := common.GetDotfilesRoot()
	if err != nil {
		return err
	}

	sourceDir := filepath.Join(dotfiles, "bin/_shortcuts/termux")
	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		return fmt.Errorf("failed to read shortcuts source: %w", err)
	}

	for _, entry := range entries {
		// Check for context cancellation
		if err := ctx.Err(); err != nil {
			return err
		}

		if !entry.IsDir() {
			name := entry.Name()
			content := fmt.Sprintf("#!/usr/bin/env bash\nexec %s/bin/source_me workspaced dispatch _shortcuts termux %s \"$@\"\n", dotfiles, name)

			destPath := filepath.Join(shortcutDir, name)
			if err := os.WriteFile(destPath, []byte(content), 0755); err != nil {
				return fmt.Errorf("failed to write shortcut %s: %w", name, err)
			}
			common.GetLogger(ctx).Info("created shortcut", "name", name)
		}
	}

	return nil
}
