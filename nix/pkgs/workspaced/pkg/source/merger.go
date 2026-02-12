package source

import (
	"context"
	"fmt"
	"sort"
)

// Merger junta múltiplas sources e resolve conflitos
type Merger struct {
	sources  []Source
	strategy ConflictResolution
}

// NewMerger cria um merger com as sources fornecidas
func NewMerger(sources []Source, strategy ConflictResolution) *Merger {
	return &Merger{
		sources:  sources,
		strategy: strategy,
	}
}

// MergeResult contém resultado do merge
type MergeResult struct {
	Files     []File     // Arquivos finais (após resolver conflitos)
	Conflicts []Conflict // Conflitos detectados
	Resolved  int        // Número de conflitos resolvidos
	Skipped   int        // Número de arquivos ignorados
}

// Merge escaneia todas as sources e resolve conflitos
func (m *Merger) Merge(ctx context.Context) (*MergeResult, error) {
	// 1. Escanear todas as sources
	allFiles := []File{}
	for _, src := range m.sources {
		files, err := src.Scan(ctx)
		if err != nil {
			return nil, fmt.Errorf("source %s failed to scan: %w", src.Name(), err)
		}
		allFiles = append(allFiles, files...)
	}

	// 2. Agrupar por target path
	filesByTarget := make(map[string][]File)
	for _, f := range allFiles {
		filesByTarget[f.TargetPath] = append(filesByTarget[f.TargetPath], f)
	}

	// 3. Detectar e resolver conflitos
	result := &MergeResult{
		Files:     []File{},
		Conflicts: []Conflict{},
	}

	for targetPath, files := range filesByTarget {
		if len(files) == 1 {
			// Sem conflito
			result.Files = append(result.Files, files[0])
			continue
		}

		// Conflito detectado
		// Ordenar por priority (maior primeiro)
		sort.Slice(files, func(i, j int) bool {
			return files[i].Priority > files[j].Priority
		})

		conflict := Conflict{
			TargetPath: targetPath,
			Files:      files,
		}

		// Resolver de acordo com estratégia
		resolved, err := m.resolveConflict(conflict)
		if err != nil {
			return nil, err
		}

		result.Conflicts = append(result.Conflicts, conflict)

		if resolved != nil {
			result.Files = append(result.Files, *resolved)
			result.Resolved++
		} else {
			result.Skipped++
		}
	}

	return result, nil
}

// resolveConflict aplica estratégia de resolução
func (m *Merger) resolveConflict(conflict Conflict) (*File, error) {
	switch m.strategy {
	case ResolveByPriority:
		// Retorna o arquivo com maior priority (primeiro da lista)
		return &conflict.Files[0], nil

	case ResolveByError:
		// Monta mensagem de erro detalhada
		msg := fmt.Sprintf("conflict at %s:\n", conflict.TargetPath)
		for _, f := range conflict.Files {
			msg += fmt.Sprintf("  - %s (priority %d)\n", f.SourceName, f.Priority)
		}
		return nil, fmt.Errorf(msg)

	case ResolveBySkip:
		// Não aplica nenhum arquivo
		return nil, nil

	default:
		return nil, fmt.Errorf("unknown conflict resolution strategy: %v", m.strategy)
	}
}
