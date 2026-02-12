package template

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/template"
	"workspaced/pkg/env"
	"workspaced/pkg/icons"
	"workspaced/pkg/logging"
	"workspaced/pkg/text"
)

// ErrFileSkipped é retornado quando um template chama {{ skip }}
var ErrFileSkipped = fmt.Errorf("file skipped")

// makeFuncMap cria o FuncMap padrão do workspaced
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
		// Filesystem helpers
		"readDir": func(path string) ([]string, error) {
			entries, err := os.ReadDir(path)
			if err != nil {
				return nil, err
			}
			var names []string
			for _, entry := range entries {
				if !entry.IsDir() {
					names = append(names, entry.Name())
				}
			}
			return names, nil
		},
		"isPhone": func() bool {
			return env.IsPhone()
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
