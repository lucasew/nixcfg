package apply

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"workspaced/pkg/env"
)

type ProfileProvider struct{}

func (p *ProfileProvider) Name() string {
	return "profile"
}

func (p *ProfileProvider) GetDesiredState(ctx context.Context) ([]DesiredState, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dotfiles, err := env.GetDotfilesRoot()
	if err != nil {
		return nil, err
	}

	tmpDir := filepath.Join(home, ".config", "workspaced", "generated")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return nil, err
	}

	sourceMePath := filepath.Join(dotfiles, "bin", "source_me")
	profileContent := fmt.Sprintf("source %s\n", sourceMePath)

	tmpProfile := filepath.Join(tmpDir, "profile")
	if err := os.WriteFile(tmpProfile, []byte(profileContent), 0644); err != nil {
		return nil, err
	}

	return []DesiredState{
		{
			Target: filepath.Join(home, ".profile"),
			Source: tmpProfile,
			Mode:   0644,
		},
		{
			Target: filepath.Join(home, ".bashrc"),
			Source: tmpProfile,
			Mode:   0644,
		},
	}, nil
}
