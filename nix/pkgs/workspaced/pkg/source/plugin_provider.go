package source

import (
	"context"
	"fmt"
)

// ProviderPlugin adapta um Provider para Plugin
// Permite usar providers legacy no sistema de pipeline
type ProviderPlugin struct {
	provider Provider
	priority int
}

// NewProviderPlugin cria plugin a partir de provider legacy
func NewProviderPlugin(provider Provider, priority int) *ProviderPlugin {
	return &ProviderPlugin{
		provider: provider,
		priority: priority,
	}
}

func (p *ProviderPlugin) Name() string {
	return fmt.Sprintf("provider:%s", p.provider.Name())
}

func (p *ProviderPlugin) Process(ctx context.Context, files []File) ([]File, error) {
	desired, err := p.provider.GetDesiredState(ctx)
	if err != nil {
		return nil, fmt.Errorf("provider %s failed: %w", p.provider.Name(), err)
	}

	// Converter DesiredState para source.File
	newFiles := make([]File, len(desired))
	for i, d := range desired {
		// Providers legacy sempre retornam BufferFiles ou StaticFiles construídos a partir de DesiredState legacy.
		// No novo modelo, DesiredState já contém um File interface.
		newFiles[i] = d.File
	}

	// Append aos arquivos existentes
	return append(files, newFiles...), nil
}
