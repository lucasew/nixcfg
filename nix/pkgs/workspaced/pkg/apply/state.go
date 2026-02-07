package apply

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func (e *Engine) GetStatePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(home, ".config/workspaced/state.json")
	if err := e.FS.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", err
	}
	return path, nil
}

func (e *Engine) LoadState() (*State, error) {
	path, err := e.GetStatePath()
	if err != nil {
		return nil, err
	}
	state := &State{Files: make(map[string]ManagedInfo)}

	if _, err := e.FS.Stat(path); err == nil {
		data, err := e.FS.ReadFile(path)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, state); err != nil {
			return nil, err
		}
	}
	if state.Files == nil {
		state.Files = make(map[string]ManagedInfo)
	}
	return state, nil
}

func (e *Engine) SaveState(state *State) error {
	path, err := e.GetStatePath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return e.FS.WriteFile(path, data, 0644)
}
