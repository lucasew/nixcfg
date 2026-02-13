package source

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"workspaced/pkg/template"
)

// TemplateExpanderPlugin renderiza templates e expande multi-file
type TemplateExpanderPlugin struct {
	engine *template.Engine
	data   interface{}
}

// NewTemplateExpanderPlugin cria plugin de expansão de templates
func NewTemplateExpanderPlugin(engine *template.Engine, data interface{}) *TemplateExpanderPlugin {
	return &TemplateExpanderPlugin{
		engine: engine,
		data:   data,
	}
}

func (p *TemplateExpanderPlugin) Name() string {
	return "template-expander"
}

func (p *TemplateExpanderPlugin) Process(ctx context.Context, files []File) ([]File, error) {
	result := []File{}

	for _, f := range files {
		// Detectar se é template
		filename := filepath.Base(f.RelPath())
		parts := strings.Split(filename, ".")
		isTemplate := (len(parts) >= 2 && parts[len(parts)-1] == "tmpl") ||
			(len(parts) >= 3 && parts[len(parts)-2] == "tmpl")

		if !isTemplate {
			// Não é template, passa adiante
			result = append(result, f)
			continue
		}

		// Calcular RelPath sem .tmpl
		relPath := f.RelPath()
		if parts[len(parts)-1] == "tmpl" {
			// file.tmpl → file
			relPath = strings.TrimSuffix(relPath, ".tmpl")
		} else {
			// file.tmpl.ext → file.ext
			relPath = strings.TrimSuffix(relPath, ".tmpl"+filepath.Ext(relPath)) + filepath.Ext(relPath)
		}

		// Eagerly render to check if it's multi-file
		reader, err := f.Reader()
		if err != nil {
			return nil, fmt.Errorf("failed to read template source %s: %w", f.SourceInfo(), err)
		}
		srcContent, err := io.ReadAll(reader)
		reader.Close()
		if err != nil {
			return nil, err
		}

		rendered, err := p.engine.Render(ctx, string(srcContent), p.data)
		if err != nil {
			if errors.Is(err, template.ErrFileSkipped) {
				continue
			}
			return nil, fmt.Errorf("failed to render template %s: %w", f.SourceInfo(), err)
		}

		// Verificar se é multi-file
		multiFiles, isMulti := template.ParseMultiFile(rendered)

		if isMulti {
			// UM template → N files (EAGER)
			baseRelDir := filepath.Dir(relPath)
			baseName := strings.TrimSuffix(filepath.Base(relPath), filepath.Ext(relPath))

			if baseName != "_index" {
				baseRelDir = filepath.Join(baseRelDir, baseName)
			}

			for _, mf := range multiFiles {
				mfRelPath := filepath.Join(baseRelDir, mf.Name)
				result = append(result, &BufferFile{
					BasicFile: BasicFile{
						RelPathStr:    mfRelPath,
						TargetBaseDir: f.TargetBase(),
						FileMode:      mf.Mode,
						Info:          fmt.Sprintf("%s (multi:%s)", f.SourceInfo(), mf.Name),
						FileType:      TypeMultiFile,
					},
					Content: []byte(mf.Content),
				})
			}
		} else {
			// UM template → 1 file (LAZY)
			result = append(result, &TemplateFile{
				BasicFile: BasicFile{
					RelPathStr:    relPath,
					TargetBaseDir: f.TargetBase(),
					FileMode:      f.Mode(), // Usually templates produce non-exec files but we can keep source mode
					Info:          f.SourceInfo(),
					FileType:      TypeTemplate,
				},
				SourceFile: f,
				Engine:     p.engine,
				Data:       p.data,
				Context:    ctx,
			})
		}
	}

	return result, nil
}
