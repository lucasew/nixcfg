package tray

import (
	"context"
	"image"
)

// MenuItem represents an item in the tray menu.
type MenuItem struct {
	Label    string
	Callback func()
	Children []MenuItem
}

// State represents the desired state of the tray.
type State struct {
	Title string
	Icon  image.Image
	Menu  []MenuItem
}

// Driver is the interface for tray implementations.
// It follows a declarative pattern where the state is set and the driver updates the UI.
type Driver interface {
	Run(ctx context.Context) error
	SetState(s State)
	Close()
}
