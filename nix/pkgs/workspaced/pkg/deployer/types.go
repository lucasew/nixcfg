package deployer

import (
	"path/filepath"
	"workspaced/pkg/source"
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
		return "+"
	case ActionUpdate:
		return "*"
	case ActionDelete:
		return "-"
	case ActionNoop:
		return " "
	}
	return "?"
}

// DesiredState alias para source.DesiredState
type DesiredState = source.DesiredState

// Helper para pegar o target path de um DesiredState
func GetTarget(d DesiredState) string {
	return filepath.Join(d.File.TargetBase(), d.File.RelPath())
}

// ManagedInfo contém informações sobre arquivo gerenciado
type ManagedInfo struct {
	SourceInfo string `json:"source_info"`
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

// Provider alias para source.Provider
type Provider = source.Provider
