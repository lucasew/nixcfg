package apply

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"workspaced/pkg/config"
	"workspaced/pkg/env"
)

type SymlinkProvider struct{}

func (p *SymlinkProvider) Name() string {
	return "symlink"
}

func (p *SymlinkProvider) GetDesiredState(ctx context.Context) ([]DesiredState, error) {
	dotfiles, err := env.GetDotfilesRoot()
	if err != nil {
		return nil, err
	}
	configDir := filepath.Join(dotfiles, "config")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return nil, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// Create directory for rendered templates
	renderedDir := filepath.Join(home, ".config/workspaced/generated/templates")
	if err := os.MkdirAll(renderedDir, 0755); err != nil {
		return nil, err
	}

	desired := []DesiredState{}
	err = filepath.Walk(configDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(configDir, path)
		if err != nil {
			return err
		}

		// Check if file is a template (.tmpl as second-to-last component)
		// Example: colors.tmpl.toml → colors.toml
		filename := filepath.Base(rel)
		parts := strings.Split(filename, ".")
		isTemplate := len(parts) >= 3 && parts[len(parts)-2] == "tmpl"

		if isTemplate {
			// Read and render template
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			rendered, err := renderTemplate(string(content), cfg)
			if err != nil {
				return err
			}

			// Remove .tmpl from filename
			newParts := append(parts[:len(parts)-2], parts[len(parts)-1])
			newFilename := strings.Join(newParts, ".")

			// Write rendered content to temp file
			renderedPath := filepath.Join(renderedDir, rel)
			renderedPath = strings.ReplaceAll(renderedPath, ".tmpl", "")
			if err := os.MkdirAll(filepath.Dir(renderedPath), 0755); err != nil {
				return err
			}
			if err := os.WriteFile(renderedPath, rendered, 0644); err != nil {
				return err
			}

			dir := filepath.Dir(rel)
			targetRel := filepath.Join(dir, newFilename)

			desired = append(desired, DesiredState{
				Target: filepath.Join(home, targetRel),
				Source: renderedPath,
				Mode:   0644,
			})
		} else {
			// Regular symlink
			desired = append(desired, DesiredState{
				Target: filepath.Join(home, rel),
				Source: path,
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return desired, nil
}

func renderTemplate(content string, cfg *config.Config) ([]byte, error) {
	tmpl, err := template.New("config").Parse(content)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, cfg); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
