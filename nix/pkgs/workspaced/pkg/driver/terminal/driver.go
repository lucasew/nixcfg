package terminal

import (
	"context"
)

type Options struct {
	Title   string
	Command string
	Args    []string
}

type Driver interface {
	Open(ctx context.Context, opts Options) error
}
