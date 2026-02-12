package source

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ScannerPlugin descobre arquivos em um diretório
type ScannerPlugin struct {
	name       string
	baseDir    string
	targetBase string
	priority   int
}

// ScannerConfig configura um ScannerPlugin
type ScannerConfig struct {
	Name       string // Nome identificador
	BaseDir    string // Diretório fonte (onde estão os arquivos)
	TargetBase string // Base path para targets (opcional, default: $HOME)
	Priority   int    // Priority em conflitos
}

// NewScannerPlugin cria um plugin scanner
func NewScannerPlugin(cfg ScannerConfig) (*ScannerPlugin, error) {
	if cfg.Name == "" {
		return nil, fmt.Errorf("scanner name is required")
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

	return &ScannerPlugin{
		name:       cfg.Name,
		baseDir:    baseDir,
		targetBase: targetBase,
		priority:   cfg.Priority,
	}, nil
}

func (p *ScannerPlugin) Name() string {
	return fmt.Sprintf("scanner:%s", p.name)
}

func (p *ScannerPlugin) Process(ctx context.Context, files []File) ([]File, error) {
	discovered := []File{}

	err := filepath.Walk(p.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip diretórios por enquanto (serão processados depois)
		if info.IsDir() {
			return nil
		}

		// Calcular caminho relativo
		rel, err := filepath.Rel(p.baseDir, path)
		if err != nil {
			return err
		}

		discovered = append(discovered, File{
			SourceName: p.name,
			RelPath:    rel,
			SourceBase: p.baseDir,
			TargetBase: p.targetBase,
			Type:       TypeStatic, // Será determinado por outros plugins
			Mode:       0,          // 0 = symlink por padrão
			Priority:   p.priority,
		})

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan directory: %w", err)
	}

	// Append aos arquivos existentes
	return append(files, discovered...), nil
}
