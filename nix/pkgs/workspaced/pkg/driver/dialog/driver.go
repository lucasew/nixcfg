package dialog

import (
	"context"
)

type Item struct {
	Label string
	Icon  string
	Value string
}

type ChooseOptions struct {
	Prompt string
	Items  []Item
}

// Chooser is for selecting from a list
type Chooser interface {
	Choose(ctx context.Context, opts ChooseOptions) (*Item, error)
}

// Prompter is for getting text input
type Prompter interface {
	Prompt(ctx context.Context, prompt string) (string, error)
}

// Confirmer is for binary yes/no questions
type Confirmer interface {
	Confirm(ctx context.Context, message string) (bool, error)
}

// Legacy compatibility
type Options = ChooseOptions
type Driver interface {
	Chooser
	RunApp(ctx context.Context) error
	SwitchWindow(ctx context.Context) error
}
