package source

import (
	"context"
	"workspaced/pkg/deployer"
)

// ProviderSource adapta um deployer.Provider para source.Source
// Permite usar providers legacy com o novo sistema
type ProviderSource struct {
	provider deployer.Provider
	priority int
}

// NewProviderSource cria uma source a partir de um provider
func NewProviderSource(provider deployer.Provider, priority int) *ProviderSource {
	return &ProviderSource{
		provider: provider,
		priority: priority,
	}
}

func (s *ProviderSource) Name() string {
	return s.provider.Name()
}

func (s *ProviderSource) Priority() int {
	return s.priority
}

func (s *ProviderSource) Scan(ctx context.Context) ([]File, error) {
	desired, err := s.provider.GetDesiredState(ctx)
	if err != nil {
		return nil, err
	}

	files := make([]File, len(desired))
	for i, d := range desired {
		// Determinar tipo baseado no Mode
		fileType := TypeStatic
		if d.Mode == 0 {
			fileType = TypeSymlink
		}

		files[i] = File{
			SourceName: s.provider.Name(),
			SourcePath: d.Source,
			TargetPath: d.Target,
			Type:       fileType,
			Mode:       d.Mode,
			Priority:   s.priority,
		}
	}

	return files, nil
}
