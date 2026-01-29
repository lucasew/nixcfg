package apply

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"workspaced/pkg/common"
)

func GetStatePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(home, ".config/workspaced/state.json")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", err
	}
	return path, nil
}

func LoadState() (*State, error) {
	path, err := GetStatePath()
	if err != nil {
		return nil, err
	}
	state := &State{Files: make(map[string]ManagedInfo)}
	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
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

func SaveState(state *State) error {
	path, err := GetStatePath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func Plan(ctx context.Context, desired []DesiredState, currentState *State) ([]Action, error) {
	actions := []Action{}
	desiredMap := make(map[string]DesiredState)

	for _, d := range desired {
		desiredMap[d.Target] = d
		current, managed := currentState.Files[d.Target]

		actualSource, err := os.Readlink(d.Target)
		exists := true
		if err != nil {
			exists = false
		}

		if !exists {
			actions = append(actions, Action{Type: ActionCreate, Target: d.Target, Desired: d})
		} else if actualSource != d.Source {
			actions = append(actions, Action{Type: ActionUpdate, Target: d.Target, Desired: d, Current: current})
		} else if !managed || current.Source != d.Source {
			// Even if link is correct, update state if not managed
			actions = append(actions, Action{Type: ActionUpdate, Target: d.Target, Desired: d, Current: current})
		} else {
			actions = append(actions, Action{Type: ActionNoop, Target: d.Target, Desired: d, Current: current})
		}
	}

	for target, current := range currentState.Files {
		if _, ok := desiredMap[target]; !ok {
			actions = append(actions, Action{Type: ActionDelete, Target: target, Current: current})
		}
	}

	return actions, nil
}

func Execute(ctx context.Context, actions []Action, state *State) error {
	logger := common.GetLogger(ctx)
	for _, action := range actions {
		switch action.Type {
		case ActionNoop:
			continue
		case ActionDelete:
			logger.Info("pruning orphaned link", "target", action.Target)
			// Only delete if it's a symlink pointing where we thought it was
			if actual, err := os.Readlink(action.Target); err == nil && actual == action.Current.Source {
				if err := os.Remove(action.Target); err != nil {
					return fmt.Errorf("failed to remove orphaned link %s: %w", action.Target, err)
				}
			}
			delete(state.Files, action.Target)

		case ActionCreate, ActionUpdate:
			logger.Info("applying link", "target", action.Target, "source", action.Desired.Source)

			// Ensure parent directory
			if err := os.MkdirAll(filepath.Dir(action.Target), 0755); err != nil {
				return err
			}

			// Conflict handling
			if info, err := os.Lstat(action.Target); err == nil {
				if info.Mode()&os.ModeSymlink != 0 {
					// It's a symlink, just remove it
					_ = os.Remove(action.Target)
				} else {
					// It's a real file/dir, backup it
					bakPath := action.Target + ".bak.workspaced"
					logger.Warn("conflict detected, backing up real file", "target", action.Target, "backup", bakPath)
					if err := os.Rename(action.Target, bakPath); err != nil {
						return fmt.Errorf("failed to backup conflict %s: %w", action.Target, err)
					}
				}
			}

			if err := os.Symlink(action.Desired.Source, action.Target); err != nil {
				return fmt.Errorf("failed to create symlink %s -> %s: %w", action.Target, action.Desired.Source, err)
			}
			state.Files[action.Target] = ManagedInfo{Source: action.Desired.Source}
		}
	}
	return nil
}
