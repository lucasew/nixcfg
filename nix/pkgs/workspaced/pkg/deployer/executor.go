package deployer

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"workspaced/pkg/logging"
)

// prettyPath converte caminho absoluto para relativo ao $HOME
func prettyPath(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	if strings.HasPrefix(path, home+"/") {
		return "~/" + strings.TrimPrefix(path, home+"/")
	}

	return path
}

// Executor executa ações de deployment
type Executor struct{}

// NewExecutor cria um novo executor
func NewExecutor() *Executor {
	return &Executor{}
}

// Execute aplica lista de ações e atualiza estado
func (e *Executor) Execute(ctx context.Context, actions []Action, state *State) error {
	logger := logging.GetLogger(ctx)

	for _, action := range actions {
		switch action.Type {
		case ActionNoop:
			continue

		case ActionDelete:
			logger.Info("pruning orphaned file", "target", prettyPath(action.Target))
			if _, err := os.Lstat(action.Target); err == nil {
				if err := os.Remove(action.Target); err != nil {
					return fmt.Errorf("failed to remove orphaned file %s: %w", action.Target, err)
				}
			}
			delete(state.Files, action.Target)

		case ActionCreate, ActionUpdate:
			if action.Type == ActionCreate {
				logger.Info("creating", "target", prettyPath(action.Target), "source", prettyPath(action.Desired.Source))
			} else {
				logger.Info("updating", "target", prettyPath(action.Target), "source", prettyPath(action.Desired.Source))
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
