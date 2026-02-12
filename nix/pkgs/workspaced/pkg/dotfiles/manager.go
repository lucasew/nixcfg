package dotfiles

import (
	"context"
	"fmt"
	"path/filepath"
	"workspaced/pkg/deployer"
	"workspaced/pkg/logging"
	"workspaced/pkg/source"
)

// Manager é a API principal para gerenciamento de dotfiles
type Manager struct {
	pipeline   *source.Pipeline
	stateStore deployer.StateStore
	planner    *deployer.Planner
	executor   *deployer.Executor
	hooks      []Hook
}

// Config configura o Manager
type Config struct {
	// Pipeline de plugins
	Pipeline *source.Pipeline

	// State store para persistência
	StateStore deployer.StateStore

	// Hooks opcionais
	Hooks []Hook
}

// NewManager cria um novo manager
func NewManager(cfg Config) (*Manager, error) {
	if cfg.Pipeline == nil {
		return nil, fmt.Errorf("pipeline is required")
	}

	if cfg.StateStore == nil {
		return nil, fmt.Errorf("state store is required")
	}

	return &Manager{
		pipeline:   cfg.Pipeline,
		stateStore: cfg.StateStore,
		planner:    deployer.NewPlanner(),
		executor:   deployer.NewExecutor(),
		hooks:      cfg.Hooks,
	}, nil
}

// ApplyOptions configura execução do Apply
type ApplyOptions struct {
	DryRun   bool // Se true, apenas mostra o que seria feito
	ShowDiff bool // Se true, mostra diff detalhado
}

// ApplyResult contém resultado do Apply
type ApplyResult struct {
	FilesCreated int
	FilesUpdated int
	FilesDeleted int
	FilesNoOp    int
	Actions      []deployer.Action
	Error        error
}

// Apply executa o ciclo completo de deployment
func (m *Manager) Apply(ctx context.Context, opts ApplyOptions) (*ApplyResult, error) {
	logger := logging.GetLogger(ctx)
	result := &ApplyResult{}

	// 1. Executar pipeline
	logger.Info("running pipeline", "plugins", len(m.pipeline.GetPlugins()))
	files, err := m.pipeline.Run(ctx, []source.File{})
	if err != nil {
		result.Error = err
		return result, fmt.Errorf("failed to run pipeline: %w", err)
	}

	logger.Info("pipeline completed", "files", len(files))

	// 2. Converter source.File para deployer.DesiredState
	// Calcula paths absolutos aqui (SourceBase+RelPath, TargetBase+RelPath)
	desired := make([]deployer.DesiredState, len(files))
	for i, f := range files {
		desired[i] = deployer.DesiredState{
			Source: filepath.Join(f.SourceBase, f.RelPath),
			Target: filepath.Join(f.TargetBase, f.RelPath),
			Mode:   f.Mode,
		}
	}

	// 3. Carregar estado atual
	logger.Info("loading state", "store", m.stateStore.Path())
	state, err := m.stateStore.Load()
	if err != nil {
		result.Error = err
		return result, fmt.Errorf("failed to load state: %w", err)
	}

	// 4. Planejar ações
	logger.Info("planning actions")
	actions, err := m.planner.Plan(ctx, desired, state)
	if err != nil {
		result.Error = err
		return result, fmt.Errorf("failed to plan: %w", err)
	}

	result.Actions = actions

	// Contar ações
	for _, a := range actions {
		switch a.Type {
		case deployer.ActionCreate:
			result.FilesCreated++
		case deployer.ActionUpdate:
			result.FilesUpdated++
		case deployer.ActionDelete:
			result.FilesDeleted++
		case deployer.ActionNoop:
			result.FilesNoOp++
		}
	}

	hasChanges := result.FilesCreated > 0 || result.FilesUpdated > 0 || result.FilesDeleted > 0

	if !hasChanges {
		logger.Info("no changes needed")
		return result, nil
	}

	logger.Info("changes planned",
		"create", result.FilesCreated,
		"update", result.FilesUpdated,
		"delete", result.FilesDeleted,
	)

	// Dry-run: para aqui
	if opts.DryRun {
		logger.Info("dry-run: skipping execution")
		return result, nil
	}

	// 5. Executar hooks Before
	for _, hook := range m.hooks {
		if err := hook.Before(ctx, actions); err != nil {
			result.Error = err
			return result, fmt.Errorf("hook before failed: %w", err)
		}
	}

	// 6. Executar ações
	logger.Info("executing actions")
	execErr := m.executor.Execute(ctx, actions, state)

	// 7. Executar hooks After (mesmo se houver erro)
	for _, hook := range m.hooks {
		if err := hook.After(ctx, actions, execErr); err != nil {
			logger.Error("hook after failed", "error", err)
			// Continua executando outros hooks
		}
	}

	if execErr != nil {
		result.Error = execErr
		return result, fmt.Errorf("failed to execute: %w", execErr)
	}

	// 8. Salvar estado
	logger.Info("saving state")
	if err := m.stateStore.Save(state); err != nil {
		result.Error = err
		return result, fmt.Errorf("failed to save state: %w", err)
	}

	logger.Info("apply completed successfully")
	return result, nil
}

// GetPipeline retorna pipeline configurado
func (m *Manager) GetPipeline() *source.Pipeline {
	return m.pipeline
}

// GetStateStore retorna state store configurado
func (m *Manager) GetStateStore() deployer.StateStore {
	return m.stateStore
}
