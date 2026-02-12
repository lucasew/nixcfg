package source

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"workspaced/pkg/template"
)

// DotDProcessorPlugin processa diretórios .d.tmpl (concatenação)
type DotDProcessorPlugin struct {
	engine  *template.Engine
	data    interface{}
	tempDir string
}

// NewDotDProcessorPlugin cria plugin de processamento .d.tmpl
func NewDotDProcessorPlugin(engine *template.Engine, data interface{}, tempDir string) (*DotDProcessorPlugin, error) {
	// Expand tempDir
	expanded := os.ExpandEnv(tempDir)
	if strings.HasPrefix(expanded, "~/") {
		home, _ := os.UserHomeDir()
		expanded = filepath.Join(home, expanded[2:])
	}

	// Criar diretório
	if err := os.MkdirAll(expanded, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	return &DotDProcessorPlugin{
		engine:  engine,
		data:    data,
		tempDir: expanded,
	}, nil
}

func (p *DotDProcessorPlugin) Name() string {
	return "dotd-processor"
}

func (p *DotDProcessorPlugin) Process(ctx context.Context, files []File) ([]File, error) {
	result := []File{}
	processed := make(map[string]bool)
	dotdDirs := make(map[string][]File) // Key: RelPath do dir .d.tmpl, Value: arquivos dentro

	// 1. Identificar diretórios .d.tmpl e agrupar arquivos
	for _, f := range files {
		// Verificar se RelPath está dentro de .d.tmpl/
		if strings.Contains(f.RelPath, ".d.tmpl/") || strings.Contains(f.RelPath, ".d.tmpl"+string(filepath.Separator)) {
			// Extrair caminho do diretório .d.tmpl relativo
			parts := strings.Split(f.RelPath, ".d.tmpl")
			if len(parts) >= 2 {
				dotdDirRel := parts[0] + ".d.tmpl"
				dotdDirs[dotdDirRel] = append(dotdDirs[dotdDirRel], f)
			}
		} else {
			// Arquivo normal, não está em .d.tmpl
			result = append(result, f)
		}
	}

	// 2. Processar cada diretório .d.tmpl
	for dotdDirRel, dirFiles := range dotdDirs {
		if processed[dotdDirRel] {
			continue
		}

		// Pegar info do primeiro arquivo para inferir bases
		var sourceName string
		var sourceBase string
		var targetBase string
		var priority int
		if len(dirFiles) > 0 {
			sourceName = dirFiles[0].SourceName
			sourceBase = dirFiles[0].SourceBase
			targetBase = dirFiles[0].TargetBase
			priority = dirFiles[0].Priority
		}

		// Caminho absoluto do diretório .d.tmpl no source
		dotdDirAbs := filepath.Join(sourceBase, dotdDirRel)

		// Concatenar arquivos do diretório
		concatenated, err := p.engine.ProcessDotD(ctx, dotdDirAbs, p.data)
		if err != nil {
			return nil, fmt.Errorf("failed to process .d.tmpl directory %s: %w", dotdDirRel, err)
		}

		// Validar que não tem markers multi-file
		if _, isMulti := template.ParseMultiFile(concatenated); isMulti {
			return nil, fmt.Errorf(
				".d.tmpl directory %s contains multi-file markers - this is not supported (use multi-file OR .d.tmpl, not both)",
				dotdDirRel,
			)
		}

		// RelPath do arquivo concatenado (remove .d.tmpl)
		relDir := filepath.Dir(dotdDirRel)
		targetName := strings.TrimSuffix(filepath.Base(dotdDirRel), ".d.tmpl")
		relPath := filepath.Join(relDir, targetName)

		// Escrever concatenado em temp
		tempFilePath := filepath.Join(p.tempDir, relPath)
		if err := os.MkdirAll(filepath.Dir(tempFilePath), 0755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(tempFilePath, concatenated, 0644); err != nil {
			return nil, err
		}

		result = append(result, File{
			SourceName: sourceName,
			RelPath:    relPath,
			SourceBase: p.tempDir,
			TargetBase: targetBase,
			Content:    concatenated,
			Mode:       0644,
			Type:       TypeDotD,
			Priority:   priority,
		})

		processed[dotdDirRel] = true
	}

	return result, nil
}
