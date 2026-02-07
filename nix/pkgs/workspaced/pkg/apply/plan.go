package apply

import (
	"context"
	"fmt"
	"os"
	"workspaced/pkg/logging"
)

func (e *Engine) Plan(ctx context.Context, desired []DesiredState, currentState *State) ([]Action, error) {
	logger := logging.GetLogger(ctx)
	actions := []Action{}
	desiredMap := make(map[string]DesiredState)

	for _, d := range desired {
		desiredMap[d.Target] = d
		current, managed := currentState.Files[d.Target]

		info, err := e.FS.Lstat(d.Target)
		exists := err == nil

		if !exists {
			actions = append(actions, Action{Type: ActionCreate, Target: d.Target, Desired: d})
			continue
		}

		needsUpdate, reason := e.validate(d, info)

		if needsUpdate {
			logger.Debug("validation failed", "target", d.Target, "reason", reason)
			actions = append(actions, Action{Type: ActionUpdate, Target: d.Target, Desired: d, Current: current})
		} else if !managed || current.Source != d.Source {
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

func (e *Engine) validate(d DesiredState, info os.FileInfo) (bool, string) {
	if d.Mode == 0 { // Symlink desired
		return e.validateSymlink(d, info)
	}
	return e.validateFile(d, info)
}

func (e *Engine) validateSymlink(d DesiredState, info os.FileInfo) (bool, string) {
	if info.Mode()&os.ModeSymlink == 0 {
		return true, "target is not a symlink"
	}
	actualSource, err := e.FS.Readlink(d.Target)
	if err != nil {
		return true, "cannot read symlink target"
	}
	if actualSource != d.Source {
		return true, fmt.Sprintf("symlink points to wrong target (current: %s, desired: %s)", actualSource, d.Source)
	}
	return false, ""
}

func (e *Engine) validateFile(d DesiredState, info os.FileInfo) (bool, string) {
	if !info.Mode().IsRegular() {
		return true, "target is not a regular file"
	}
	// Always validate content and permissions
	if info.Mode().Perm() != d.Mode.Perm() {
		return true, fmt.Sprintf("permissions mismatch (current: %o, desired: %o)", info.Mode().Perm(), d.Mode.Perm())
	}
	// Validate content
	srcData, err := e.FS.ReadFile(d.Source)
	if err != nil {
		return true, "cannot read source file"
	}
	dstData, err := e.FS.ReadFile(d.Target)
	if err != nil {
		return true, "cannot read target file"
	}
	if string(srcData) != string(dstData) {
		return true, "content mismatch"
	}
	return false, ""
}
