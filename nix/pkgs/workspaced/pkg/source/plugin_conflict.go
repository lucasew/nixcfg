package source

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"workspaced/pkg/logging"
)

// ConflictResolverPlugin resolve conflitos em arquivos com mesmo target
type ConflictResolverPlugin struct {
	strategy ConflictResolution
}

// NewConflictResolverPlugin cria plugin de resolução de conflitos
func NewConflictResolverPlugin(strategy ConflictResolution) *ConflictResolverPlugin {
	return &ConflictResolverPlugin{strategy: strategy}
}

func (p *ConflictResolverPlugin) Name() string {
	return fmt.Sprintf("conflict-resolver:%s", p.strategy)
}

func (p *ConflictResolverPlugin) Process(ctx context.Context, files []File) ([]File, error) {
	logger := logging.GetLogger(ctx)

	// Agrupar por TargetBase + RelPath (caminho completo de destino)
	byTarget := make(map[string][]File)
	for _, f := range files {
		key := filepath.Join(f.TargetBase, f.RelPath)
		byTarget[key] = append(byTarget[key], f)
	}

	result := []File{}
	conflicts := 0

	for targetKey, group := range byTarget {
		if len(group) == 1 {
			// Sem conflito
			result = append(result, group[0])
			continue
		}

		// Conflito detectado
		conflicts++
		logger.Debug("conflict detected", "target", targetKey, "count", len(group))

		// Resolver de acordo com estratégia
		resolved, err := p.resolveConflict(group[0].RelPath, group)
		if err != nil {
			return nil, err
		}

		if resolved != nil {
			result = append(result, *resolved)
		}
	}

	if conflicts > 0 {
		logger.Info("conflicts resolved", "count", conflicts, "strategy", p.strategy)
	}

	return result, nil
}

func (p *ConflictResolverPlugin) resolveConflict(relPath string, files []File) (*File, error) {
	switch p.strategy {
	case ResolveByPriority:
		// Ordenar por priority (maior primeiro)
		sort.Slice(files, func(i, j int) bool {
			return files[i].Priority > files[j].Priority
		})
		return &files[0], nil

	case ResolveByError:
		// Retornar erro detalhado
		msg := fmt.Sprintf("conflict at %s:\n", relPath)
		for _, f := range files {
			msg += fmt.Sprintf("  - %s (priority %d, type %s)\n", f.SourceName, f.Priority, f.Type)
		}
		return nil, fmt.Errorf("%s", msg)

	case ResolveBySkip:
		// Não aplica nenhum arquivo
		return nil, nil

	default:
		return nil, fmt.Errorf("unknown conflict resolution strategy: %v", p.strategy)
	}
}
