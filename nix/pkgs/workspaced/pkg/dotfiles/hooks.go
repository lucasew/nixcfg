package dotfiles

import (
	"context"
	"workspaced/pkg/deployer"
)

// Hook permite executar código antes/depois do deployment
type Hook interface {
	// Before é chamado antes de executar actions
	Before(ctx context.Context, actions []deployer.Action) error

	// After é chamado após executar actions (mesmo se houver erro)
	After(ctx context.Context, applied []deployer.Action, err error) error
}

// FuncHook implementa Hook usando funções
type FuncHook struct {
	BeforeFn func(ctx context.Context, actions []deployer.Action) error
	AfterFn  func(ctx context.Context, applied []deployer.Action, err error) error
}

func (h *FuncHook) Before(ctx context.Context, actions []deployer.Action) error {
	if h.BeforeFn != nil {
		return h.BeforeFn(ctx, actions)
	}
	return nil
}

func (h *FuncHook) After(ctx context.Context, applied []deployer.Action, err error) error {
	if h.AfterFn != nil {
		return h.AfterFn(ctx, applied, err)
	}
	return nil
}
