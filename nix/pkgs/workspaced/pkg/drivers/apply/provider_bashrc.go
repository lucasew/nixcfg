package apply

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"workspaced/pkg/common"
)

type BashrcProvider struct{}

func (p *BashrcProvider) Name() string {
	return "bashrc"
}

func (p *BashrcProvider) GetDesiredState(ctx context.Context) ([]DesiredState, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dotfiles, err := common.GetDotfilesRoot()
	if err != nil {
		return nil, err
	}

	// Create a temporary file with the bashrc content
	tmpDir := filepath.Join(home, ".config", "workspaced", "generated")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return nil, err
	}

	sourceMePath := filepath.Join(dotfiles, "bin", "source_me")
	bashrcContent := fmt.Sprintf("source %s\n", sourceMePath)

	tmpBashrc := filepath.Join(tmpDir, "bashrc")
	if err := os.WriteFile(tmpBashrc, []byte(bashrcContent), 0644); err != nil {
		return nil, err
	}

	return []DesiredState{
		{
			Target: filepath.Join(home, ".bashrc"),
			Source: tmpBashrc,
			Mode:   0644,
		},
	}, nil
}
