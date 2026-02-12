package menu

import (
	"context"
)

type Item struct {
	Label string
	Icon  string
	Value string
}

type Options struct {
	Prompt string
	Items  []Item
}

type Driver interface {
	Choose(ctx context.Context, opts Options) (*Item, error)
	RunApp(ctx context.Context) error
	SwitchWindow(ctx context.Context) error
}
