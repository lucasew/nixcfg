package setup

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"workspaced/pkg/common"
	"workspaced/pkg/host"
)

func SetupTermuxShortcuts(ctx context.Context) error {
	if !host.IsPhone() {
		return fmt.Errorf("this command only works on phone (Termux)")
	}

	home, _ := os.UserHomeDir()
	shortcutDir := filepath.Join(home, ".shortcuts")
	_ = os.MkdirAll(shortcutDir, 0755)

	dotfiles, err := host.GetDotfilesRoot()
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
			common.GetLogger(ctx).Info("created shortcut", "name", name)
		}
	}

	return nil
}
