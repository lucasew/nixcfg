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
	logger := common.GetLogger(ctx)
	actions := []Action{}
	desiredMap := make(map[string]DesiredState)

	for _, d := range desired {
		desiredMap[d.Target] = d
		current, managed := currentState.Files[d.Target]

		info, err := os.Lstat(d.Target)
		exists := err == nil

		if !exists {
			actions = append(actions, Action{Type: ActionCreate, Target: d.Target, Desired: d})
			continue
		}

		// Validate and determine if update is needed
		needsUpdate := false
		reason := ""

		if d.Mode == 0 { // Symlink desired
			if info.Mode()&os.ModeSymlink == 0 {
				needsUpdate = true
				reason = "target is not a symlink"
			} else {
				actualSource, err := os.Readlink(d.Target)
				if err != nil {
					needsUpdate = true
					reason = "cannot read symlink target"
				} else if actualSource != d.Source {
					needsUpdate = true
					reason = fmt.Sprintf("symlink points to wrong target (current: %s, desired: %s)", actualSource, d.Source)
				}
			}
		} else { // Regular file desired
			if !info.Mode().IsRegular() {
				needsUpdate = true
				reason = "target is not a regular file"
			} else {
				// Always validate content and permissions
				if info.Mode().Perm() != d.Mode.Perm() {
					needsUpdate = true
					reason = fmt.Sprintf("permissions mismatch (current: %o, desired: %o)", info.Mode().Perm(), d.Mode.Perm())
				} else {
					// Validate content
					srcData, err := os.ReadFile(d.Source)
					if err != nil {
						needsUpdate = true
						reason = "cannot read source file"
					} else {
						dstData, err := os.ReadFile(d.Target)
						if err != nil {
							needsUpdate = true
							reason = "cannot read target file"
						} else if string(srcData) != string(dstData) {
							needsUpdate = true
							reason = "content mismatch"
						}
					}
				}
			}
		}

		// Always update if validation failed, regardless of state tracking
		if needsUpdate {
			logger.Debug("validation failed", "target", d.Target, "reason", reason)
			actions = append(actions, Action{Type: ActionUpdate, Target: d.Target, Desired: d, Current: current})
		} else if !managed || current.Source != d.Source {
			// Not tracked in state or source changed - update to track it
			actions = append(actions, Action{Type: ActionUpdate, Target: d.Target, Desired: d, Current: current})
		} else {
			actions = append(actions, Action{Type: ActionNoop, Target: d.Target, Desired: d, Current: current})
		}
	}

	// Prune orphaned files
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
			logger.Info("pruning orphaned file", "target", action.Target)
			if _, err := os.Lstat(action.Target); err == nil {
				if err := os.Remove(action.Target); err != nil {
					return fmt.Errorf("failed to remove orphaned file %s: %w", action.Target, err)
				}
			}
			delete(state.Files, action.Target)

		case ActionCreate, ActionUpdate:
			if action.Type == ActionCreate {
				logger.Info("creating", "target", action.Target, "source", action.Desired.Source)
			} else {
				logger.Info("updating", "target", action.Target, "source", action.Desired.Source)
			}

			// Ensure parent directory exists
			if err := os.MkdirAll(filepath.Dir(action.Target), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory for %s: %w", action.Target, err)
			}

			// Always remove existing target to ensure clean state
			if info, err := os.Lstat(action.Target); err == nil {
				if action.Desired.Mode == 0 {
					// We want a symlink - always remove whatever exists
					logger.Debug("removing existing target for symlink", "target", action.Target, "type", info.Mode().Type())
					if err := os.Remove(action.Target); err != nil {
						return fmt.Errorf("failed to remove existing target %s: %w", action.Target, err)
					}
				} else {
					// We want a regular file
					if info.Mode().IsRegular() {
						// Existing file is regular - will overwrite
						logger.Debug("overwriting existing regular file", "target", action.Target)
					} else if info.Mode()&os.ModeSymlink != 0 {
						// Existing is symlink, remove it
						logger.Debug("removing existing symlink to replace with file", "target", action.Target)
						if err := os.Remove(action.Target); err != nil {
							return fmt.Errorf("failed to remove existing symlink %s: %w", action.Target, err)
						}
					} else {
						// Something else (directory, device, etc) - back it up
						bakPath := action.Target + ".bak.workspaced"
						logger.Warn("backing up non-regular file", "target", action.Target, "backup", bakPath, "type", info.Mode().Type())
						if err := os.Rename(action.Target, bakPath); err != nil {
							return fmt.Errorf("failed to backup %s: %w", action.Target, err)
						}
					}
				}
			}

			// Create the desired state
			if action.Desired.Mode == 0 {
				// Create symlink
				if err := os.Symlink(action.Desired.Source, action.Target); err != nil {
					return fmt.Errorf("failed to create symlink %s -> %s: %w", action.Target, action.Desired.Source, err)
				}
				logger.Debug("symlink created", "target", action.Target, "source", action.Desired.Source)
			} else {
				// Create/overwrite regular file
				data, err := os.ReadFile(action.Desired.Source)
				if err != nil {
					return fmt.Errorf("failed to read source file %s: %w", action.Desired.Source, err)
				}
				if err := os.WriteFile(action.Target, data, action.Desired.Mode); err != nil {
					return fmt.Errorf("failed to write target file %s: %w", action.Target, err)
				}
				logger.Debug("file written", "target", action.Target, "size", len(data), "mode", action.Desired.Mode)
			}

			// Track in state
			state.Files[action.Target] = ManagedInfo{Source: action.Desired.Source}
		}
	}
	return nil
}
