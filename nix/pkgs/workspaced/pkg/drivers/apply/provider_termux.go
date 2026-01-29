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
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		genPath := filepath.Join(genDir, name)

		content := fmt.Sprintf(`#!/data/data/com.termux/files/usr/bin/bash
export LD_PRELOAD="/data/data/com.termux/files/usr/lib/libtermux-exec.so"
export PATH="/data/data/com.termux/files/usr/bin"
exec bash %s/bin/source_me sd _shortcuts termux %s "$@"
`, dotfiles, name)

		if err := os.WriteFile(genPath, []byte(content), 0755); err != nil {
			return nil, err
		}

		desired = append(desired, DesiredState{
			Target: filepath.Join(home, ".shortcuts", name),
			Source: genPath,
		})
	}

	return desired, nil
}
