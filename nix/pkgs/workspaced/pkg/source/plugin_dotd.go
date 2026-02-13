package source

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
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
	dotdGroups := make(map[string][]File) // Key: Final Target Path, Value: component files

	// 1. Group files by their intended concatenated target
	for _, f := range files {
		// A file is part of a DotD if its RelPath contains ".d.tmpl/"
		if strings.Contains(f.RelPath, ".d.tmpl/") || strings.Contains(f.RelPath, ".d.tmpl"+string(filepath.Separator)) {
			// Extract the target path by removing the ".d.tmpl" part and the filename
			parts := strings.Split(f.RelPath, ".d.tmpl")
			targetRelPath := parts[0] // e.g. ".bashrc" if RelPath was ".bashrc.d.tmpl/00-base.sh"
			finalTargetPath := filepath.Join(f.TargetBase, targetRelPath)
			dotdGroups[finalTargetPath] = append(dotdGroups[finalTargetPath], f)
		} else {
			result = append(result, f)
		}
	}

	// 2. Process each group
	for targetPath, groupFiles := range dotdGroups {
		// Sort by RelPath to ensure deterministic concatenation
		sort.Slice(groupFiles, func(i, j int) bool {
			return groupFiles[i].RelPath < groupFiles[j].RelPath
		})

		var concatenated strings.Builder
		for _, f := range groupFiles {
			content := f.Content
			// If content is empty but it's TypeStatic, read it (should have been expanded already if template)
			if len(content) == 0 && f.Type == TypeStatic {
				var err error
				content, err = os.ReadFile(filepath.Join(f.SourceBase, f.RelPath))
				if err != nil {
					return nil, fmt.Errorf("failed to read dotd component %s: %w", f.RelPath, err)
				}
			}
			concatenated.Write(content)
			if !strings.HasSuffix(string(content), "\n") {
				concatenated.WriteString("\n")
			}
		}

		// Use info from the first file for bases
		first := groupFiles[0]
		relPath, _ := filepath.Rel(first.TargetBase, targetPath)

		// Save concatenated result to temp dir
		tempFilePath := filepath.Join(p.tempDir, relPath)
		if err := os.MkdirAll(filepath.Dir(tempFilePath), 0755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(tempFilePath, []byte(concatenated.String()), 0644); err != nil {
			return nil, err
		}

		result = append(result, File{
			SourceName: "dotd-processor",
			RelPath:    relPath,
			SourceBase: p.tempDir,
			TargetBase: first.TargetBase,
			Content:    []byte(concatenated.String()),
			Mode:       0644,
			Type:       TypeDotD,
			Priority:   first.Priority,
		})
	}

	return result, nil
}
