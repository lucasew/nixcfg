package api

import (
	"context"
	"workspaced/pkg/drivers/api"
)

// Re-export shared errors for local convenience if needed,
// or just use the shared ones.
var ErrDriverNotFound = api.ErrDriverNotFound

// Rect represents a geometry rectangle.
type Rect struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Workspace represents a workspace as returned by Sway/i3 IPC.
type Workspace struct {
	Name    string `json:"name"`
	Focused bool   `json:"focused"`
	Output  string `json:"output"`
}

// Output represents a display output as returned by Sway/i3 IPC.
type Output struct {
	Name             string `json:"name"`
	CurrentWorkspace string `json:"current_workspace"`
	Rect             Rect   `json:"rect"`
	Focused          bool   `json:"focused"`
}

// Node represents a node in the Sway/i3 tree.
type Node struct {
	Rect          Rect    `json:"rect"`
	Focused       bool    `json:"focused"`
	Nodes         []*Node `json:"nodes"`
	FloatingNodes []*Node `json:"floating_nodes"`
}

type Driver interface {
	SwitchToWorkspace(ctx context.Context, num int, move bool) error
	ToggleScratchpad(ctx context.Context) error
	GetFocusedOutput(ctx context.Context) (string, *Rect, error)
	GetFocusedWindowRect(ctx context.Context) (*Rect, error)
	GetOutputs(ctx context.Context) ([]Output, error)
	GetWorkspaces(ctx context.Context) ([]Workspace, error)
}
