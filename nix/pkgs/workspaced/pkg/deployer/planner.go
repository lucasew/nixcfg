package deployer

import (
	"context"
	"fmt"
	"os"
	"workspaced/pkg/logging"
)

// Planner compara estado atual vs desejado e gera ações
type Planner struct{}

// NewPlanner cria um novo planner
func NewPlanner() *Planner {
	return &Planner{}
}

// Plan compara desired state com current state e retorna ações necessárias
func (p *Planner) Plan(ctx context.Context, desired []DesiredState, currentState *State) ([]Action, error) {
	logger := logging.GetLogger(ctx)
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
