package source

import (
	"context"
	"fmt"
	"workspaced/pkg/deployer"
)

// ProviderPlugin adapta um deployer.Provider para Plugin
// Permite usar providers legacy no sistema de pipeline
type ProviderPlugin struct {
	provider deployer.Provider
	priority int
}

// NewProviderPlugin cria plugin a partir de provider legacy
func NewProviderPlugin(provider deployer.Provider, priority int) *ProviderPlugin {
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

	// Converter deployer.DesiredState para source.File
	// Providers legados retornam paths absolutos, então precisamos extrair RelPath
	newFiles := make([]File, len(desired))
	for i, d := range desired {
		fileType := TypeStatic
		if d.Mode == 0 {
			fileType = TypeSymlink
		}

		// Para providers legacy, Source e Target são absolutos
		// Usamos Target como base + basename como RelPath (simplificação)
		// Isso funciona para providers como DconfProvider que geram marker files
		relPath := filepath.Base(d.Target)
		targetBase := filepath.Dir(d.Target)
		sourceBase := filepath.Dir(d.Source)

		newFiles[i] = File{
			SourceName: p.provider.Name(),
			RelPath:    relPath,
			SourceBase: sourceBase,
			TargetBase: targetBase,
			Type:       fileType,
			Mode:       d.Mode,
			Priority:   p.priority,
		}
	}

	// Append aos arquivos existentes
	return append(files, newFiles...), nil
}
