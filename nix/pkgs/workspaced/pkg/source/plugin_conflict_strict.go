package source

import (
	"context"
	"fmt"
	"path/filepath"
	"workspaced/pkg/logging"
)

// StrictConflictResolverPlugin garante exclusividade total de caminhos
type StrictConflictResolverPlugin struct{}

func NewStrictConflictResolverPlugin() *StrictConflictResolverPlugin {
	return &StrictConflictResolverPlugin{}
}

func (p *StrictConflictResolverPlugin) Name() string {
	return "strict-conflict-resolver"
}

func (p *StrictConflictResolverPlugin) Process(ctx context.Context, files []File) ([]File, error) {
	logger := logging.GetLogger(ctx)
	logger.Debug("running strict conflict resolution")

	ownedPaths := make(map[string]File) // FinalPath -> File info

	for _, f := range files {
		finalPath := filepath.Join(f.TargetBase, f.RelPath)

		if existing, ok := ownedPaths[finalPath]; ok {
			return nil, fmt.Errorf("strict conflict detected for path %q: provided by source %q and %q (zero-substitution policy)",
				finalPath, existing.SourceName, f.SourceName)
		}

		ownedPaths[finalPath] = f
	}

	return files, nil
}
