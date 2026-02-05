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
	*GlobalConfig
	raw map[string]interface{}
}

func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dotfiles, _ := common.GetDotfilesRoot()

	// 1. Start with hardcoded defaults
	config := GlobalConfig{
		Workspaces: map[string]int{
			"www":  1,
			"meet": 2,
		},
		Desktop: DesktopConfig{
			Wallpaper: WallpaperConfig{
				Dir: filepath.Join(dotfiles, "assets/wallpapers"),
			},
		},
		Screenshot: ScreenshotConfig{
			Dir: filepath.Join(home, "Pictures/Screenshots"),
		},
		Backup: BackupConfig{
			RsyncnetUser: "de3163@de3163.rsync.net",
			RemotePath:   "backup/lucasew",
		},
		QuickSync: QuickSyncConfig{
			RepoDir:    filepath.Join(home, ".personal"),
			RemotePath: "/data2/home/de3163/git-personal",
		},
		Hosts: make(map[string]HostConfig),
		Browser: BrowserConfig{
			Default: "zen",
			Engine:  "brave",
		},
		Webapps:   make(map[string]WebappConfig),
		LazyTools: make(map[string]LazyToolConfig),
	}

	// Raw map for UnmarshalKey
	raw := make(map[string]interface{})

	// 2. Load and merge base config from $DOTFILES/settings.toml
	if dotfiles != "" {
		dotfilesSettingsPath := filepath.Join(dotfiles, "settings.toml")
		if _, err := os.Stat(dotfilesSettingsPath); err == nil {
			var dotfilesConfig GlobalConfig
			if _, err := toml.DecodeFile(dotfilesSettingsPath, &dotfilesConfig); err == nil {
				config = config.Merge(dotfilesConfig)
			}

			// Merge into raw
			var tmp map[string]interface{}
			if _, err := toml.DecodeFile(dotfilesSettingsPath, &tmp); err == nil {
				for k, v := range tmp {
					raw[k] = v
				}
			}
		}
	}

	// 3. Load and merge user config from ~/settings.toml
	userSettingsPath := filepath.Join(home, "settings.toml")
	if _, err := os.Stat(userSettingsPath); err == nil {
		var userConfig GlobalConfig
		if _, err := toml.DecodeFile(userSettingsPath, &userConfig); err == nil {
			config = config.Merge(userConfig)
		}

		// Merge into raw
		var tmp map[string]interface{}
		if _, err := toml.DecodeFile(userSettingsPath, &tmp); err == nil {
			for k, v := range tmp {
				raw[k] = v
			}
		}
	}

	// 4. Expand paths
	config.Desktop.Wallpaper.Dir = common.ExpandPath(config.Desktop.Wallpaper.Dir)
	config.Desktop.Wallpaper.Default = common.ExpandPath(config.Desktop.Wallpaper.Default)
	config.Screenshot.Dir = common.ExpandPath(config.Screenshot.Dir)
	config.QuickSync.RepoDir = common.ExpandPath(config.QuickSync.RepoDir)

	return &Config{
		GlobalConfig: &config,
		raw:          raw,
	}, nil
}

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

	// Use json as a hack to unmarshal into the target type from map[string]interface{}
	// This avoids manual type assertion hell but incurs a performance cost.
	data, err := json.Marshal(current)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, val)
}
