package source

import (
	"context"
	"fmt"
	"workspaced/pkg/logging"
)

// Plugin processa lista de arquivos e retorna nova lista
// Inspirado no padrão de plugins do Beancount
type Plugin interface {
	// Name retorna nome do plugin (para logging)
	Name() string

	// Process transforma lista de arquivos
	// Pode adicionar, remover, modificar arquivos
	Process(ctx context.Context, files []File) ([]File, error)
}

// Pipeline executa sequência de plugins
type Pipeline struct {
	plugins []Plugin
}

// NewPipeline cria pipeline com plugins
func NewPipeline(plugins ...Plugin) *Pipeline {
	return &Pipeline{plugins: plugins}
}

// AddPlugin adiciona plugin ao final do pipeline
func (p *Pipeline) AddPlugin(plugin Plugin) {
	p.plugins = append(p.plugins, plugin)
}

// Run executa pipeline completo
func (p *Pipeline) Run(ctx context.Context, initial []File) ([]File, error) {
	logger := logging.GetLogger(ctx)
	current := initial

	for i, plugin := range p.plugins {
		logger.Debug("running plugin", "index", i, "name", plugin.Name(), "input_count", len(current))

		result, err := plugin.Process(ctx, current)
		if err != nil {
			return nil, fmt.Errorf("plugin %s failed: %w", plugin.Name(), err)
		}

		logger.Debug("plugin completed", "name", plugin.Name(), "output_count", len(result))
		current = result
	}

	logger.Info("pipeline completed", "total_plugins", len(p.plugins), "final_count", len(current))
	return current, nil
}

// GetPlugins retorna lista de plugins configurados
func (p *Pipeline) GetPlugins() []Plugin {
	return p.plugins
}
