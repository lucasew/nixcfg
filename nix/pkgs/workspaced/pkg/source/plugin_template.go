package source

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"workspaced/pkg/template"
)

// TemplateExpanderPlugin renderiza templates e expande multi-file
type TemplateExpanderPlugin struct {
	engine  *template.Engine
	data    interface{}
	tempDir string
}

// NewTemplateExpanderPlugin cria plugin de expansão de templates
func NewTemplateExpanderPlugin(engine *template.Engine, data interface{}, tempDir string) (*TemplateExpanderPlugin, error) {
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

	return &TemplateExpanderPlugin{
		engine:  engine,
		data:    data,
		tempDir: expanded,
	}, nil
}

func (p *TemplateExpanderPlugin) Name() string {
	return "template-expander"
}

func (p *TemplateExpanderPlugin) Process(ctx context.Context, files []File) ([]File, error) {
	result := []File{}

	for _, f := range files {
		// Detectar se é template
		filename := filepath.Base(f.RelPath)
		parts := strings.Split(filename, ".")
		isTemplate := (len(parts) >= 2 && parts[len(parts)-1] == "tmpl") ||
			(len(parts) >= 3 && parts[len(parts)-2] == "tmpl")

		if !isTemplate {
			// Não é template, passa adiante
			result = append(result, f)
			continue
		}

		// Renderizar template
		sourcePath := filepath.Join(f.SourceBase, f.RelPath)
		content, err := os.ReadFile(sourcePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read template %s: %w", sourcePath, err)
		}

		rendered, err := p.engine.Render(ctx, string(content), p.data)
		if err != nil {
			if errors.Is(err, template.ErrFileSkipped) {
				// Template chamou {{ skip }}, não gera arquivo
				continue
			}
			return nil, fmt.Errorf("failed to render template %s: %w", sourcePath, err)
		}

		// Calcular RelPath sem .tmpl
		relPath := f.RelPath
		if parts[len(parts)-1] == "tmpl" {
			// file.tmpl → file
			relPath = strings.TrimSuffix(relPath, ".tmpl")
		} else {
			// file.tmpl.ext → file.ext
			relPath = strings.TrimSuffix(relPath, ".tmpl"+filepath.Ext(relPath)) + filepath.Ext(relPath)
		}

		// Verificar se é multi-file
		multiFiles, isMulti := template.ParseMultiFile(rendered)

		if isMulti {
			// UM template → N files
			baseRelDir := filepath.Dir(relPath)
			baseName := strings.TrimSuffix(filepath.Base(relPath), filepath.Ext(relPath))

			// _index.tmpl → arquivos vão direto no dir
			// regular.tmpl → cria subdir
			if baseName != "_index" {
				baseRelDir = filepath.Join(baseRelDir, baseName)
			}

			for _, mf := range multiFiles {
				mfRelPath := filepath.Join(baseRelDir, mf.Name)

				// Escrever em temp
				tempFilePath := filepath.Join(p.tempDir, mfRelPath)
				if err := os.MkdirAll(filepath.Dir(tempFilePath), 0755); err != nil {
					return nil, err
				}
				if err := os.WriteFile(tempFilePath, []byte(mf.Content), mf.Mode); err != nil {
					return nil, err
				}

				result = append(result, File{
					SourceName: f.SourceName,
					RelPath:    mfRelPath,
					SourceBase: p.tempDir,
					TargetBase: f.TargetBase,
					Content:    []byte(mf.Content),
					Mode:       mf.Mode,
					Type:       TypeMultiFile,
					Priority:   f.Priority,
				})
			}
		} else {
			// UM template → 1 file
			// Escrever renderizado em temp
			tempFilePath := filepath.Join(p.tempDir, relPath)
			if err := os.MkdirAll(filepath.Dir(tempFilePath), 0755); err != nil {
				return nil, err
			}
			if err := os.WriteFile(tempFilePath, rendered, 0644); err != nil {
				return nil, err
			}

			result = append(result, File{
				SourceName: f.SourceName,
				RelPath:    relPath,
				SourceBase: p.tempDir,
				TargetBase: f.TargetBase,
				Content:    rendered,
				Mode:       0644,
				Type:       TypeTemplate,
				Priority:   f.Priority,
			})
		}
	}

	return result, nil
}
