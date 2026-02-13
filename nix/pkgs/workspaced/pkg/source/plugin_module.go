package source

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"workspaced/pkg/config"
	"workspaced/pkg/logging"
)

type ModuleScannerPlugin struct {
	baseDir  string
	cfg      *config.GlobalConfig
	priority int
}

func NewModuleScannerPlugin(baseDir string, cfg *config.GlobalConfig, priority int) *ModuleScannerPlugin {
	return &ModuleScannerPlugin{
		baseDir:  baseDir,
		cfg:      cfg,
		priority: priority,
	}
}

func (p *ModuleScannerPlugin) Name() string {
	return "module-scanner"
}

var presetBases = map[string]string{
	"home": "~",
	"etc":  "/etc",
	"usr":  "/usr",
	"root": "/",
	"var":  "/var",
	"bin":  "/usr/local/bin",
}

func (p *ModuleScannerPlugin) Process(ctx context.Context, files []File) ([]File, error) {
	logger := logging.GetLogger(ctx)

	entries, err := os.ReadDir(p.baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return files, nil
		}
		return nil, err
	}

	discovered := []File{}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		modName := entry.Name()
		modPath := filepath.Join(p.baseDir, modName)

		// Check if module is enabled
		modCfgRaw, ok := p.cfg.Modules[modName]
		if !ok {
			logger.Debug("module ignored (not in config)", "module", modName)
			continue
		}

		modCfg, ok := modCfgRaw.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("invalid config for module %q: expected map", modName)
		}

		enabled, _ := modCfg["enable"].(bool)
		if !enabled {
			logger.Debug("module disabled", "module", modName)
			continue
		}

		logger.Info("loading module", "module", modName)

		// Scan presets
		presets, err := os.ReadDir(modPath)
		if err != nil {
			return nil, err
		}

		for _, preset := range presets {
			if !preset.IsDir() {
				// Strict structure: no files in module root
				return nil, fmt.Errorf("strict structure violation: file %q found in module %q root (expected preset directory)", preset.Name(), modName)
			}

			presetName := preset.Name()
			targetBase, ok := presetBases[presetName]
			if !ok {
				return nil, fmt.Errorf("unknown preset %q in module %q", presetName, modName)
			}

			// Expand targetBase
			if targetBase == "~" {
				home, _ := os.UserHomeDir()
				targetBase = home
			}

			presetPath := filepath.Join(modPath, presetName)
			err := filepath.Walk(presetPath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					return nil
				}

				rel, _ := filepath.Rel(presetPath, path)
				discovered = append(discovered, File{
					SourceName: modName,
					RelPath:    rel,
					SourceBase: presetPath,
					TargetBase: targetBase,
					Type:       TypeStatic,
					Priority:   p.priority,
				})
				return nil
			})
			if err != nil {
				return nil, err
			}
		}
	}

	return append(files, discovered...), nil
}
