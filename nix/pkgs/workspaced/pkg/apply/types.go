package apply

import (
	"context"
	"os"
)

type ActionType int

const (
	ActionCreate ActionType = iota
	ActionUpdate
	ActionDelete
	ActionNoop
)

func (a ActionType) String() string {
	switch a {
	case ActionCreate:
		return "CREATE"
	case ActionUpdate:
		return "UPDATE"
	case ActionDelete:
		return "DELETE"
	case ActionNoop:
		return "NOOP"
	}
	return "UNKNOWN"
}

type DesiredState struct {
	Target string      // Absolute path on the system
	Source string      // Absolute path to the source file (could be in repo or generated)
	Mode   os.FileMode // Type/Mode of the target. 0 means symlink.
}

type ManagedInfo struct {
	Source string `json:"source"`
}

type State struct {
	Files map[string]ManagedInfo `json:"files"` // Key: Target (Absolute path)
}

type Action struct {
	Type    ActionType
	Target  string
	Desired DesiredState
	Current ManagedInfo
}

type Provider interface {
	Name() string
	GetDesiredState(ctx context.Context) ([]DesiredState, error)
}
