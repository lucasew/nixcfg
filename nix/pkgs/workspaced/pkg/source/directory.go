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

// DirectorySource escaneia um diretório de templates
type DirectorySource struct {
	name       string
	baseDir    string
	targetBase string // Base path para targets (default: $HOME)
	priority   int
	engine     *template.Engine
	data       interface{}
	tempDir    string // Diretório para arquivos renderizados
}

// DirectorySourceConfig configura um DirectorySource
type DirectorySourceConfig struct {
	Name       string          // Nome identificador da source
	BaseDir    string          // Diretório fonte (onde estão os templates)
	TargetBase string          // Base path para targets (opcional, default: $HOME)
	Priority   int             // Priority em conflitos
	Engine     *template.Engine // Template engine (opcional, cria padrão se nil)
	Data       interface{}     // Dados para renderizar templates
	TempDir    string          // Dir para temporários (opcional, usa padrão)
}

// NewDirectorySource cria uma nova source de diretório
func NewDirectorySource(ctx context.Context, cfg DirectorySourceConfig) (*DirectorySource, error) {
	// Validações
	if cfg.Name == "" {
		return nil, fmt.Errorf("source name is required")
	}
	if cfg.BaseDir == "" {
		return nil, fmt.Errorf("base directory is required")
	}

	// Expand paths
	baseDir := os.ExpandEnv(cfg.BaseDir)
	if strings.HasPrefix(baseDir, "~/") {
		home, _ := os.UserHomeDir()
		baseDir = filepath.Join(home, baseDir[2:])
	}

	// Verificar se diretório existe
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("base directory does not exist: %s", baseDir)
	}

	// Target base (default: $HOME)
	targetBase := cfg.TargetBase
	if targetBase == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		targetBase = home
	}

	// Template engine (default: novo com padrões workspaced)
	engine := cfg.Engine
	if engine == nil {
		engine = template.NewEngine(ctx)
	}

	// Temp dir (default: $HOME/.config/workspaced/generated/<source-name>)
	tempDir := cfg.TempDir
	if tempDir == "" {
		home, _ := os.UserHomeDir()
		tempDir = filepath.Join(home, ".config/workspaced/generated", cfg.Name)
	}

	// Criar temp dir
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	return &DirectorySource{
		name:       cfg.Name,
		baseDir:    baseDir,
		targetBase: targetBase,
		priority:   cfg.Priority,
		engine:     engine,
		data:       cfg.Data,
		tempDir:    tempDir,
	}, nil
}

func (s *DirectorySource) Name() string {
	return s.name
}

func (s *DirectorySource) Priority() int {
	return s.priority
}

func (s *DirectorySource) Scan(ctx context.Context) ([]File, error) {
	var files []File

	err := filepath.Walk(s.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calcular caminho relativo
		rel, err := filepath.Rel(s.baseDir, path)
		if err != nil {
			return err
		}

		// Processar .d.tmpl directories (concatenação)
		if info.IsDir() && strings.HasSuffix(info.Name(), ".d.tmpl") {
			dotdFiles, err := s.processDotD(ctx, path, rel)
			if err != nil {
				return fmt.Errorf("failed to process .d.tmpl directory %s: %w", path, err)
			}
			files = append(files, dotdFiles...)
			return filepath.SkipDir
		}

		// Skip outros diretórios
		if info.IsDir() {
			return nil
		}

		// Processar arquivo
		fileList, err := s.processFile(ctx, path, rel)
		if err != nil {
			if errors.Is(err, template.ErrFileSkipped) {
				// Template chamou {{ skip }}, ignorar silenciosamente
				return nil
			}
			return err
		}

		files = append(files, fileList...)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// processDotD processa diretório .d.tmpl (concatenação)
func (s *DirectorySource) processDotD(ctx context.Context, dirPath, relPath string) ([]File, error) {
	// Renderizar concatenação
	concatenated, err := s.engine.ProcessDotD(ctx, dirPath, s.data)
	if err != nil {
		return nil, err
	}

	// Nome do arquivo de destino (remove .d.tmpl)
	targetName := strings.TrimSuffix(filepath.Base(relPath), ".d.tmpl")
	dir := filepath.Dir(relPath)
	targetRel := filepath.Join(dir, targetName)

	// Escrever arquivo renderizado em temp
	renderedPath := filepath.Join(s.tempDir, targetRel)
	if err := os.MkdirAll(filepath.Dir(renderedPath), 0755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(renderedPath, concatenated, 0644); err != nil {
		return nil, err
	}

	return []File{
		{
			SourceName: s.name,
			SourcePath: renderedPath,
			TargetPath: filepath.Join(s.targetBase, targetRel),
			Type:       TypeDotD,
			Content:    concatenated,
			Mode:       0644,
			Priority:   s.priority,
		},
	}, nil
}

// processFile processa um arquivo individual
func (s *DirectorySource) processFile(ctx context.Context, path, relPath string) ([]File, error) {
	// Detectar se é template
	// Suporta: file.tmpl → file  OU  file.tmpl.ext → file.ext
	filename := filepath.Base(relPath)
	parts := strings.Split(filename, ".")
	isTemplate := (len(parts) >= 2 && parts[len(parts)-1] == "tmpl") ||
		(len(parts) >= 3 && parts[len(parts)-2] == "tmpl")

	if !isTemplate {
		// Arquivo regular - criar symlink
		return []File{
			{
				SourceName: s.name,
				SourcePath: path,
				TargetPath: filepath.Join(s.targetBase, relPath),
				Type:       TypeSymlink,
				Mode:       0, // 0 = symlink
				Priority:   s.priority,
			},
		}, nil
	}

	// É template - renderizar
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	rendered, err := s.engine.Render(ctx, string(content), s.data)
	if err != nil {
		return nil, err
	}

	// Calcular nome de destino (remove .tmpl)
	var targetFilename string
	if parts[len(parts)-1] == "tmpl" {
		// file.tmpl → file
		targetFilename = strings.Join(parts[:len(parts)-1], ".")
	} else {
		// file.tmpl.ext → file.ext
		newParts := append(parts[:len(parts)-2], parts[len(parts)-1])
		targetFilename = strings.Join(newParts, ".")
	}

	dir := filepath.Dir(relPath)

	// Verificar se é multi-file template
	if multiFiles, isMulti := template.ParseMultiFile(rendered); isMulti {
		return s.processMultiFile(multiFiles, dir, targetFilename)
	}

	// Template simples
	targetRel := filepath.Join(dir, targetFilename)
	renderedPath := filepath.Join(s.tempDir, targetRel)
	if err := os.MkdirAll(filepath.Dir(renderedPath), 0755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(renderedPath, rendered, 0644); err != nil {
		return nil, err
	}

	return []File{
		{
			SourceName: s.name,
			SourcePath: renderedPath,
			TargetPath: filepath.Join(s.targetBase, targetRel),
			Type:       TypeTemplate,
			Content:    rendered,
			Mode:       0644,
			Priority:   s.priority,
		},
	}, nil
}

// processMultiFile processa template multi-file
func (s *DirectorySource) processMultiFile(multiFiles []template.MultiFile, dir, baseName string) ([]File, error) {
	var files []File

	// Determinar diretório base
	// _index.tmpl → arquivos vão direto no dir
	// regular.tmpl → cria subdiretório com o nome
	var baseDir string
	if baseName == "_index" {
		baseDir = dir
	} else {
		baseDir = filepath.Join(dir, baseName)
	}

	for _, mf := range multiFiles {
		targetRel := filepath.Join(baseDir, mf.Name)

		// Escrever em temp
		tempFilePath := filepath.Join(s.tempDir, targetRel)
		if err := os.MkdirAll(filepath.Dir(tempFilePath), 0755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(tempFilePath, []byte(mf.Content), mf.Mode); err != nil {
			return nil, err
		}

		files = append(files, File{
			SourceName: s.name,
			SourcePath: tempFilePath,
			TargetPath: filepath.Join(s.targetBase, targetRel),
			Type:       TypeMultiFile,
			Content:    []byte(mf.Content),
			Mode:       mf.Mode,
			Priority:   s.priority,
		})
	}

	return files, nil
}
