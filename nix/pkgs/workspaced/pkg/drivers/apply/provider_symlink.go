package apply

import (
	"context"
	"os"
	"path/filepath"
	"workspaced/pkg/host"
)

type SymlinkProvider struct{}

func (p *SymlinkProvider) Name() string {
	return "symlink"
}

func (p *SymlinkProvider) GetDesiredState(ctx context.Context) ([]DesiredState, error) {
	dotfiles, err := host.GetDotfilesRoot()
	if err != nil {
		return nil, err
	}
	configDir := filepath.Join(dotfiles, "config")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return nil, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	desired := []DesiredState{}
	err = filepath.Walk(configDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(configDir, path)
		if err != nil {
			return err
		}

		desired = append(desired, DesiredState{
			Target: filepath.Join(home, rel),
			Source: path,
		})
		return nil
	})

	if err != nil {
		return nil, err
	}

	return desired, nil
}
