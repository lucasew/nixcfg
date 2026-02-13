package deployer

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"workspaced/pkg/logging"
	"workspaced/pkg/source"
)

// Planner compara estado atual vs desejado e gera ações
type Planner struct{}

// NewPlanner cria um novo planner
func NewPlanner() *Planner {
	return &Planner{}
}

func calculateHash(r io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// Plan compara desired state com current state e retorna ações necessárias
func (p *Planner) Plan(ctx context.Context, desired []DesiredState, currentState *State) ([]Action, error) {
	logger := logging.GetLogger(ctx)
	actions := []Action{}
	desiredMap := make(map[string]DesiredState)

	for _, d := range desired {
		target := d.Target()
		desiredMap[target] = d
		current, managed := currentState.Files[target]

		info, err := os.Lstat(target)
		exists := err == nil

		if !exists {
			actions = append(actions, Action{Type: ActionCreate, Target: target, Desired: d})
			continue
		}

		// Validate and determine if update is needed
		needsUpdate := false
		reason := ""

		// Check if type matches (File vs Symlink)
		desiredIsSymlink := d.File.Type() == source.TypeSymlink
		actualIsSymlink := info.Mode()&os.ModeSymlink != 0

		if desiredIsSymlink != actualIsSymlink {
			needsUpdate = true
			reason = "type mismatch (file vs symlink)"
		} else if desiredIsSymlink {
			// Compare symlink targets
			desiredTarget, err := d.File.LinkTarget()
			if err != nil {
				return nil, fmt.Errorf("failed to get desired link target for %s: %w", d.File.SourceInfo(), err)
			}
			actualTarget, err := os.Readlink(target)
			if err != nil {
				needsUpdate = true
				reason = "failed to read actual symlink"
			} else if desiredTarget != actualTarget {
				needsUpdate = true
				reason = fmt.Sprintf("link target mismatch: %s != %s", desiredTarget, actualTarget)
			}
		} else {
			// Regular file comparison
			if info.Mode().Perm() != d.File.Mode().Perm() {
				needsUpdate = true
				reason = "permissions mismatch"
			} else {
				// Compare content via hash
				reader, err := d.File.Reader()
				if err != nil {
					return nil, fmt.Errorf("failed to get reader for %s: %w", d.File.SourceInfo(), err)
				}
				desiredHash, err := calculateHash(reader)
				reader.Close()
				if err != nil {
					return nil, err
				}

				targetFile, err := os.Open(target)
				if err != nil {
					needsUpdate = true
					reason = "cannot open target file"
				} else {
					actualHash, err := calculateHash(targetFile)
					targetFile.Close()
					if err != nil {
						return nil, err
					}
					if desiredHash != actualHash {
						needsUpdate = true
						reason = "content mismatch"
					}
				}
			}
		}

		if needsUpdate {
			logger.Debug("validation failed", "target", target, "reason", reason)
			actions = append(actions, Action{Type: ActionUpdate, Target: target, Desired: d, Current: current})
		} else if !managed || current.SourceInfo != d.File.SourceInfo() {
			// Not tracked or source changed
			actions = append(actions, Action{Type: ActionUpdate, Target: target, Desired: d, Current: current})
		} else {
			actions = append(actions, Action{Type: ActionNoop, Target: target, Desired: d, Current: current})
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
