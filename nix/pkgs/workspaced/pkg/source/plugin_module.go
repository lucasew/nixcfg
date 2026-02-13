package source

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"workspaced/pkg/config"
	"workspaced/pkg/logging"

	"github.com/xeipuuv/gojsonschema"
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

func (p *ModuleScannerPlugin) validateConfig(ctx context.Context, modName string, modPath string, modCfg map[string]any) error {
	logger := logging.GetLogger(ctx)
	schemaPath := filepath.Join(modPath, "schema.json")
	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		return nil // No schema, skip validation
	}

	logger.Debug("validating module config", "module", modName, "schema", schemaPath)

	absSchemaPath, err := filepath.Abs(schemaPath)
	if err != nil {
		return err
	}

	// Wrapper schema that adds 'enable' property and requires it
	wrapperSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"enable": map[string]any{"type": "boolean"},
		},
		"required": []string{"enable"},
		"allOf": []map[string]any{
			{"$ref": "file://" + absSchemaPath},
		},
	}

	schemaLoader := gojsonschema.NewGoLoader(wrapperSchema)
	documentLoader := gojsonschema.NewGoLoader(modCfg)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("failed to validate module %q: %w", modName, err)
	}

	if !result.Valid() {
		var errs string
		for _, desc := range result.Errors() {
			errs += fmt.Sprintf("- %s\n", desc)
		}
		return fmt.Errorf("config validation failed for module %q:\n%s", modName, errs)
	}

	logger.Debug("module config valid", "module", modName)
	return nil
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

		// Check if module configuration exists
		modCfgRaw, ok := p.cfg.Modules[modName]
		if !ok {
			modCfgRaw = make(map[string]any)
		}

		modCfg, ok := modCfgRaw.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("invalid config for module %q: expected map, got %T", modName, modCfgRaw)
		}

		// Validate config if schema exists
		if err := p.validateConfig(ctx, modName, modPath, modCfg); err != nil {
			return nil, err
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
				name := preset.Name()
				if name == "schema.json" || name == "module.toml" || name == "defaults.toml" {
					continue
				}
				return nil, fmt.Errorf("strict structure violation: file %q found in module %q root (expected preset directory or module meta files)", name, modName)
			}

			presetName := preset.Name()
			targetBase, ok := presetBases[presetName]
			if !ok {
				return nil, fmt.Errorf("unknown preset %q in module %q", presetName, modName)
			}

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
				discovered = append(discovered, &StaticFile{
					BasicFile: BasicFile{
						RelPathStr:    rel,
						TargetBaseDir: targetBase,
						FileMode:      info.Mode(),
						Info:          fmt.Sprintf("module:%s (%s/%s)", modName, presetName, rel),
						FileType:      TypeStatic,
					},
					AbsPath: path,
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
