package deployer

import (
	"context"
	"os"
)

// ActionType representa tipo de ação no deployment
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

// DesiredState representa estado desejado de um arquivo
type DesiredState struct {
	Target string      // Absolute path on the system
	Source string      // Absolute path to the source file
	Mode   os.FileMode // Type/Mode of the target. 0 means symlink.
}

// ManagedInfo contém informações sobre arquivo gerenciado
type ManagedInfo struct {
	Source string `json:"source"`
}

// State representa estado atual do sistema
type State struct {
	Files map[string]ManagedInfo `json:"files"` // Key: Target (Absolute path)
}

// Action representa ação a ser executada
type Action struct {
	Type    ActionType
	Target  string
	Desired DesiredState
	Current ManagedInfo
}

// Provider gera estados desejados (interface legacy, mantida para compatibilidade)
type Provider interface {
	Name() string
	GetDesiredState(ctx context.Context) ([]DesiredState, error)
}
