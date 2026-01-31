package apply

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"workspaced/pkg/config"
)

type DconfProvider struct{}

func (p *DconfProvider) Name() string {
	return "dconf"
}

func (p *DconfProvider) GetDesiredState(ctx context.Context) ([]DesiredState, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// Read desktop.raw.dconf section
	rawDconf := make(map[string]map[string]interface{})
	if err := cfg.UnmarshalKey("desktop.raw.dconf", &rawDconf); err != nil {
		// Section doesn't exist, that's ok
		rawDconf = make(map[string]map[string]interface{})
	}

	// Apply desktop.dark_mode if set
	var darkMode bool
	if err := cfg.UnmarshalKey("desktop.dark_mode", &darkMode); err == nil {
		if rawDconf["org/gnome/desktop/interface"] == nil {
			rawDconf["org/gnome/desktop/interface"] = make(map[string]interface{})
		}
		if darkMode {
			rawDconf["org/gnome/desktop/interface"]["color-scheme"] = "prefer-dark"
		} else {
			rawDconf["org/gnome/desktop/interface"]["color-scheme"] = "prefer-light"
		}
	}

	if len(rawDconf) == 0 {
		return nil, nil
	}

	// Generate dconf ini format (sorted for consistency)
	var sb strings.Builder
	paths := make([]string, 0, len(rawDconf))
	for path := range rawDconf {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	for _, path := range paths {
		settings := rawDconf[path]
		sb.WriteString(fmt.Sprintf("[%s]\n", path))

		keys := make([]string, 0, len(settings))
		for key := range settings {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			value := settings[key]
			sb.WriteString(fmt.Sprintf("%s=%s\n", key, formatDconfValue(value)))
		}
		sb.WriteString("\n")
	}

	tmpDir := filepath.Join(home, ".config", "workspaced", "generated")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return nil, err
	}

	dconfContent := sb.String()
	dconfFile := filepath.Join(tmpDir, "dconf.ini")
	if err := os.WriteFile(dconfFile, []byte(dconfContent), 0644); err != nil {
		return nil, err
	}

	// Apply dconf settings via dconf load
	if err := applyDconf(dconfFile); err != nil {
		fmt.Fprintf(os.Stderr, "warning: failed to apply dconf: %v\n", err)
	}

	// Use content hash as marker to track changes
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(dconfContent)))
	markerFile := filepath.Join(tmpDir, "dconf.marker")
	markerContent := filepath.Join(tmpDir, fmt.Sprintf("dconf-%s.hash", hash))

	// Create the hash file so the engine can find it
	if err := os.WriteFile(markerContent, []byte(hash), 0644); err != nil {
		return nil, err
	}

	return []DesiredState{
		{
			Target: markerFile,
			Source: markerContent,
			Mode:   0644,
		},
	}, nil
}

func applyDconf(iniFile string) error {
	cmd := exec.Command("dconf", "load", "/")
	file, err := os.Open(iniFile)
	if err != nil {
		return err
	}
	defer file.Close()

	cmd.Stdin = file
	return cmd.Run()
}

func formatDconfValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		return fmt.Sprintf("'%s'", val)
	case bool:
		return fmt.Sprintf("%t", val)
	case int, int64, float64:
		return fmt.Sprintf("%v", val)
	case []interface{}:
		// Handle arrays
		parts := make([]string, len(val))
		for i, item := range val {
			parts[i] = formatDconfValue(item)
		}
		return fmt.Sprintf("[%s]", strings.Join(parts, ", "))
	default:
		return fmt.Sprintf("%v", val)
	}
}
