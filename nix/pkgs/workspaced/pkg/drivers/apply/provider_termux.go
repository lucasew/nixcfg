package apply

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"workspaced/pkg/common"
)

type TermuxProvider struct{}

func (p *TermuxProvider) Name() string {
	return "termux"
}

func (p *TermuxProvider) GetDesiredState(ctx context.Context) ([]DesiredState, error) {
	if !common.IsPhone() {
		return nil, nil
	}

	dotfiles, err := common.GetDotfilesRoot()
	if err != nil {
		return nil, err
	}

	shortcutsSrc := filepath.Join(dotfiles, "bin/_shortcuts/termux")
	if _, err := os.Stat(shortcutsSrc); os.IsNotExist(err) {
		return nil, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	genDir := filepath.Join(home, ".local/share/workspaced/generated/termux")
	if err := os.MkdirAll(genDir, 0755); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(shortcutsSrc)
	if err != nil {
		return nil, err
	}

	desired := []DesiredState{}
	logger := common.GetLogger(ctx)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		genPath := filepath.Join(genDir, name)

		logger.Info("generating termux shortcut", "name", name)
		content := fmt.Sprintf(`#!/data/data/com.termux/files/usr/bin/bash
TERMUX_EXEC_PATH="/data/data/com.termux/files/usr/lib/libtermux-exec.so"
if [ -f "$TERMUX_EXEC_PATH" ] && [[ "$LD_PRELOAD" != *"$TERMUX_EXEC_PATH"* ]]; then
	export LD_PRELOAD="$TERMUX_EXEC_PATH"
	exec "$0" "$@"
fi
export PATH="/data/data/com.termux/files/usr/bin"
. "%s/bin/source_me"
sd _shortcuts termux %s "$@"
`, dotfiles, name)

		if err := os.WriteFile(genPath, []byte(content), 0755); err != nil {
			return nil, err
		}

		desired = append(desired, DesiredState{
			Target: filepath.Join(home, ".shortcuts", name),
			Source: genPath,
			Mode:   0755,
		})
	}

	return desired, nil
}
