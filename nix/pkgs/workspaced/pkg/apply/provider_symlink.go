package apply

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"workspaced/pkg/config"
	"workspaced/pkg/env"
	"workspaced/pkg/icons"
	"workspaced/pkg/logging"
	"workspaced/pkg/text"
)

var ErrFileSkipped = errors.New("file skipped")

const (
	markerFileStart = "<<<WORKSPACED_FILE:"
	markerFileEnd   = "<<<WORKSPACED_ENDFILE>>>"
)

func makeFuncMap(ctx context.Context) template.FuncMap {
	return template.FuncMap{
	"skip": func() (string, error) {
		return "", ErrFileSkipped
	},
	"dotfiles": func() (string, error) {
		return env.GetDotfilesRoot()
	},
	"userDataDir": func() (string, error) {
		return env.GetUserDataDir()
	},
	"file": func(name string, mode ...string) string {
		perm := "0644"
		if len(mode) > 0 {
			perm = mode[0]
		}
		return fmt.Sprintf("%s%s:%s>>>\n", markerFileStart, name, perm)
	},
	"endfile": func() string {
		return fmt.Sprintf("\n%s\n", markerFileEnd)
	},
	// String functions
	"split": func(s, sep string) []string {
		return strings.Split(s, sep)
	},
	"join": func(arr []string, sep string) string {
		return strings.Join(arr, sep)
	},
	"trimSpace": strings.TrimSpace,
	"replace": func(s, old, new string) string {
		return strings.ReplaceAll(s, old, new)
	},
	// Array/slice functions
	"list": func(items ...interface{}) []interface{} {
		return items
	},
	"last": func(arr []string) string {
		if len(arr) == 0 {
			return ""
		}
		return arr[len(arr)-1]
	},
	// Logic helpers
	"default": func(def interface{}, val interface{}) interface{} {
		if val == nil || val == "" {
			return def
		}
		return val
	},
	"ternary": func(condition bool, trueVal, falseVal interface{}) interface{} {
		if condition {
			return trueVal
		}
		return falseVal
	},
	// Webapp helpers
	"favicon": func(url string) (string, error) {
		return getFavicon(ctx, url)
	},
	"isWayland": func() bool {
		return os.Getenv("WAYLAND_DISPLAY") != ""
	},
	"titleCase": func(s string) string {
		return text.ToTitleCase(s)
	},
	"normalizeURL": func(url string) string {
		return env.NormalizeURL(url)
	},
	}
}

func getFavicon(ctx context.Context, url string) (string, error) {
	iconPath, err := icons.GetIconPath(ctx, url)
	if err != nil {
		logger := logging.GetLogger(ctx)
		logger.Error("failed to get favicon", "url", url, "error", err)
		// Return fallback icon
		return "applications-internet", nil
	}
	return iconPath, nil
}

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

		// Check if file is a template
		// Supports: colors.tmpl.toml → colors.toml  OR  colors.tmpl → colors
		filename := filepath.Base(rel)
		parts := strings.Split(filename, ".")
		isTemplate := (len(parts) >= 2 && parts[len(parts)-1] == "tmpl") ||
		              (len(parts) >= 3 && parts[len(parts)-2] == "tmpl")

		if isTemplate {
			// Read and render template
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			rendered, err := renderTemplate(ctx, string(content), cfg)
			if err != nil {
				if errors.Is(err, ErrFileSkipped) {
					// Skip this file - don't add to desired state
					return nil
				}
				return err
			}

			// Remove .tmpl from filename
			var newFilename string
			if parts[len(parts)-1] == "tmpl" {
				// file.tmpl → file
				newFilename = strings.Join(parts[:len(parts)-1], ".")
			} else {
				// file.tmpl.ext → file.ext
				newParts := append(parts[:len(parts)-2], parts[len(parts)-1])
				newFilename = strings.Join(newParts, ".")
			}

			// Check if this is a multi-file template
			if multiFiles, isMulti := parseMultiFile(rendered); isMulti {
				// Multi-file: template name becomes directory, each file inside
				dir := filepath.Dir(rel)
				baseDir := filepath.Join(dir, newFilename)

				for _, mf := range multiFiles {
					// Write each file to temp
					tempFilePath := filepath.Join(renderedDir, baseDir, mf.name)
					if err := os.MkdirAll(filepath.Dir(tempFilePath), 0755); err != nil {
						return err
					}
					if err := os.WriteFile(tempFilePath, []byte(mf.content), mf.mode); err != nil {
						return err
					}

					// Add to desired state
					targetPath := filepath.Join(home, baseDir, mf.name)
					desired = append(desired, DesiredState{
						Target: targetPath,
						Source: tempFilePath,
						Mode:   mf.mode,
					})
				}
			} else {
				// Single-file: works as before
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
			}
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

func renderTemplate(ctx context.Context, content string, cfg *config.Config) ([]byte, error) {
	tmpl, err := template.New("config").Funcs(makeFuncMap(ctx)).Parse(content)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, cfg); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type multiFile struct {
	name    string
	mode    os.FileMode
	content string
}

func parseMultiFile(rendered []byte) ([]multiFile, bool) {
	content := string(rendered)
	if !strings.Contains(content, markerFileStart) {
		return nil, false
	}

	var files []multiFile
	parts := strings.Split(content, markerFileStart)

	for i, part := range parts {
		if i == 0 {
			// Skip content before first marker
			continue
		}

		// Parse header: filename:mode>>>
		headerEnd := strings.Index(part, ">>>")
		if headerEnd == -1 {
			continue
		}
		header := part[:headerEnd]
		rest := part[headerEnd+3:]

		// Split header
		headerParts := strings.SplitN(header, ":", 2)
		if len(headerParts) != 2 {
			continue
		}
		filename := headerParts[0]
		modeStr := headerParts[1]

		// Find end marker
		endIdx := strings.Index(rest, markerFileEnd)
		if endIdx == -1 {
			// No end marker, take rest of content
			endIdx = len(rest)
		}

		fileContent := strings.TrimSpace(rest[:endIdx])

		// Parse mode
		var mode os.FileMode = 0644
		if modeStr != "" {
			if parsed, err := strconv.ParseUint(modeStr, 8, 32); err == nil {
				mode = os.FileMode(parsed)
			}
		}

		files = append(files, multiFile{
			name:    filename,
			mode:    mode,
			content: fileContent,
		})
	}

	return files, len(files) > 0
}
