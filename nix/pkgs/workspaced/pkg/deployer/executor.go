package deployer

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"workspaced/pkg/logging"
	"workspaced/pkg/source"
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
				logger.Info("creating", "target", prettyPath(action.Target), "source", action.Desired.File.SourceInfo())
			} else {
				logger.Info("updating", "target", prettyPath(action.Target), "source", action.Desired.File.SourceInfo())
			}

			// Ensure parent directory exists
			if err := os.MkdirAll(filepath.Dir(action.Target), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory for %s: %w", action.Target, err)
			}

			// Always remove existing target to ensure clean state (especially if it's a symlink or dir)
			if _, err := os.Lstat(action.Target); err == nil {
				if err := os.RemoveAll(action.Target); err != nil {
					return fmt.Errorf("failed to remove existing target %s: %w", action.Target, err)
				}
			}

			// Handle symlinks
			if action.Desired.File.Type() == source.TypeSymlink {
				linkTarget, err := action.Desired.File.LinkTarget()
				if err != nil {
					return fmt.Errorf("failed to get link target for %s: %w", action.Desired.File.SourceInfo(), err)
				}
				if err := os.Symlink(linkTarget, action.Target); err != nil {
					return fmt.Errorf("failed to create symlink %s -> %s: %w", action.Target, linkTarget, err)
				}
				state.Files[action.Target] = ManagedInfo{SourceInfo: action.Desired.File.SourceInfo()}
				continue
			}

			// For regular files or templates
			// Open target for writing
			f, err := os.OpenFile(action.Target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, action.Desired.File.Mode())
			if err != nil {
				return fmt.Errorf("failed to open target file %s: %w", action.Target, err)
			}

			reader, err := action.Desired.File.Reader()
			if err != nil {
				f.Close()
				return fmt.Errorf("failed to get reader for %s: %w", action.Desired.File.SourceInfo(), err)
			}

			_, err = io.Copy(f, reader)
			reader.Close()
			f.Close()

			if err != nil {
				return fmt.Errorf("failed to write content to %s: %w", action.Target, err)
			}

			// Track in state
			state.Files[action.Target] = ManagedInfo{SourceInfo: action.Desired.File.SourceInfo()}
		}
	}

	return nil
}
