package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"workspaced/pkg/common"

	"github.com/BurntSushi/toml"
)

type Config struct {
	*common.GlobalConfig
	raw map[string]interface{}
}

func Load() (*Config, error) {
	commonCfg, err := common.LoadConfig()
	if err != nil {
		return nil, err
	}

	raw := make(map[string]interface{})
	home, _ := os.UserHomeDir()
	dotfiles, _ := common.GetDotfilesRoot()

	// 1. Merge dotfiles settings
	if dotfiles != "" {
		dotfilesSettingsPath := filepath.Join(dotfiles, "settings.toml")
		var tmp map[string]interface{}
		if _, err := toml.DecodeFile(dotfilesSettingsPath, &tmp); err == nil {
			for k, v := range tmp {
				raw[k] = v
			}
		}
	}

	// 2. Merge user settings
	userSettingsPath := filepath.Join(home, "settings.toml")
	var tmp map[string]interface{}
	if _, err := toml.DecodeFile(userSettingsPath, &tmp); err == nil {
		for k, v := range tmp {
			raw[k] = v
		}
	}

	return &Config{
		GlobalConfig: commonCfg,
		raw:          raw,
	}, nil
}

// UnmarshalKey retrieves a value from the raw configuration map by dot-notation key
// and unmarshals it into the provided target structure.
//
// It navigates the nested map structure using the dot-separated key (e.g. "desktop.wallpaper").
//
// Implementation Detail:
// Since the configuration is stored as a generic `map[string]interface{}`, we use
// a round-trip JSON serialization (Marshal -> Unmarshal) to convert the unstructured
// map data into the strongly-typed target structure `val`. This avoids manual type assertions
// and leverages the standard library's robust mapping logic.
func (c *Config) UnmarshalKey(key string, val interface{}) error {
	parts := strings.Split(key, ".")
	var current interface{} = c.raw

	for _, part := range parts {
		if m, ok := current.(map[string]interface{}); ok {
			val, ok := m[part]
			if !ok {
				return fmt.Errorf("key not found: %s", key)
			}
			current = val
		} else {
			return fmt.Errorf("key not found: %s", key)
		}
	}

	if current == nil {
		return fmt.Errorf("key not found: %s", key)
	}

	// Use json as a trick to unmarshal into the target type
	data, err := json.Marshal(current)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, val)
}
