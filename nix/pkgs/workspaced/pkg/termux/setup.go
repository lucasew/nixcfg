package termux

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"workspaced/pkg/api"
	"workspaced/pkg/env"
	"workspaced/pkg/logging"
)

func SetupShortcuts(ctx context.Context) error {
	if !env.IsPhone() {
		return fmt.Errorf("%w: command only works on phone (Termux)", api.ErrNotSupported)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	shortcutDir := filepath.Join(home, ".shortcuts")
	if err := os.MkdirAll(shortcutDir, 0755); err != nil {
		return err
	}

	dotfiles, err := env.GetDotfilesRoot()
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
			content := fmt.Sprintf(`#!/data/data/com.termux/files/usr/bin/bash
TERMUX_EXEC_PATH="/data/data/com.termux/files/usr/lib/libtermux-exec.so"
if [ -f "$TERMUX_EXEC_PATH" ] && [[ "$LD_PRELOAD" != *"$TERMUX_EXEC_PATH"* ]]; then
	export LD_PRELOAD="$TERMUX_EXEC_PATH"
	exec "$0" "$@"
fi
exec %s/bin/source_me workspaced dispatch _shortcuts termux %s "$@"
`, dotfiles, name)

			destPath := filepath.Join(shortcutDir, name)
			if err := os.WriteFile(destPath, []byte(content), 0755); err != nil {
				return fmt.Errorf("failed to write shortcut %s: %w", name, err)
			}
			logging.GetLogger(ctx).Info("created shortcut", "name", name)
		}
	}

	return nil
}
