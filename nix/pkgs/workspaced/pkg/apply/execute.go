package apply

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"workspaced/pkg/logging"
)

func (e *Engine) Execute(ctx context.Context, actions []Action, state *State) error {
	logger := logging.GetLogger(ctx)
	for _, action := range actions {
		switch action.Type {
		case ActionNoop:
			continue
		case ActionDelete:
			logger.Info("pruning orphaned file", "target", action.Target)
			if _, err := e.FS.Lstat(action.Target); err == nil {
				if err := e.FS.Remove(action.Target); err != nil {
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

			if err := e.applyAction(action); err != nil {
				return err
			}

			// Track in state
			state.Files[action.Target] = ManagedInfo{Source: action.Desired.Source}
		}
	}
	return nil
}

func (e *Engine) applyAction(action Action) error {
	// Ensure parent directory exists
	if err := e.FS.MkdirAll(filepath.Dir(action.Target), 0755); err != nil {
		return fmt.Errorf("failed to create parent directory for %s: %w", action.Target, err)
	}

	// Always remove existing target to ensure clean state
	// In a real refactor, we might want to check if removal is necessary (e.g. if overwriting file with file)
	// But to match previous logic safely, we remove if needed.
	if info, err := e.FS.Lstat(action.Target); err == nil {
		if action.Desired.Mode == 0 {
			// We want a symlink - always remove whatever exists
			if err := e.FS.Remove(action.Target); err != nil {
				return fmt.Errorf("failed to remove existing target %s: %w", action.Target, err)
			}
		} else {
			// We want a regular file
			if info.Mode().IsRegular() {
				// Existing file is regular - will overwrite
			} else if info.Mode()&os.ModeSymlink != 0 {
				// Existing is symlink, remove it
				if err := e.FS.Remove(action.Target); err != nil {
					return fmt.Errorf("failed to remove existing symlink %s: %w", action.Target, err)
				}
			} else {
				// Something else (directory, device, etc) - back it up
				bakPath := action.Target + ".bak.workspaced"
				if err := e.FS.Rename(action.Target, bakPath); err != nil {
					return fmt.Errorf("failed to backup %s: %w", action.Target, err)
				}
			}
		}
	}

	// Create the desired state
	if action.Desired.Mode == 0 {
		// Create symlink
		if err := e.FS.Symlink(action.Desired.Source, action.Target); err != nil {
			return fmt.Errorf("failed to create symlink %s -> %s: %w", action.Target, action.Desired.Source, err)
		}
	} else {
		// Create/overwrite regular file
		data, err := e.FS.ReadFile(action.Desired.Source)
		if err != nil {
			return fmt.Errorf("failed to read source file %s: %w", action.Desired.Source, err)
		}
		if err := e.FS.WriteFile(action.Target, data, action.Desired.Mode); err != nil {
			return fmt.Errorf("failed to write target file %s: %w", action.Target, err)
		}
	}
	return nil
}
