package config

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"workspaced/pkg/env"

	"github.com/BurntSushi/toml"
)

type Config struct {
	*GlobalConfig
	raw map[string]interface{}
}

// GlobalConfig represents the schema of the settings.toml file.
type GlobalConfig struct {
	Workspaces map[string]int            `toml:"workspaces"`
	Desktop    DesktopConfig             `toml:"desktop"`
	Screenshot ScreenshotConfig          `toml:"screenshot"`
	Hosts      map[string]HostConfig     `toml:"hosts"`
	Backup     BackupConfig              `toml:"backup"`
	QuickSync  QuickSyncConfig           `toml:"quicksync"`
	Browser    BrowserConfig             `toml:"browser"`
	Webapps    map[string]WebappConfig   `toml:"webapp"`
	LazyTools  map[string]LazyToolConfig `toml:"lazy_tools"`
	Palette    PaletteConfig             `toml:"palette"`
}

type DesktopConfig struct {
	DarkMode  bool            `toml:"dark_mode"`
	Wallpaper WallpaperConfig `toml:"wallpaper"`
}

type LazyToolConfig struct {
	Version string   `toml:"version"`
	Pkg     string   `toml:"pkg"`    // Optional: mise ref (e.g. github:owner/repo)
	Global  bool     `toml:"global"` // Whether to put in global PATH
	Alias   string   `toml:"alias"`  // Binary name if Bins is empty
	Bins    []string `toml:"bins"`   // List of binaries
}

type HostConfig struct {
	MAC         string `toml:"mac"`
	TailscaleIP string `toml:"tailscale_ip"`
	ZerotierIP  string `toml:"zerotier_ip"`
	Port        int    `toml:"port"`
	User        string `toml:"user"`
}

type BrowserConfig struct {
	Default string `toml:"default"`
	Engine  string `toml:"webapp"`
}

type WebappConfig struct {
	URL         string   `toml:"url"`
	Profile     string   `toml:"profile"`
	DesktopName string   `toml:"desktop_name"`
	Icon        string   `toml:"icon"`
	ExtraFlags  []string `toml:"extra_flags"`
}

type WallpaperConfig struct {
	Dir     string `toml:"dir"`
	Default string `toml:"default"`
}

type ScreenshotConfig struct {
	Dir string `toml:"dir"`
}

type BackupConfig struct {
	RsyncnetUser string `toml:"rsyncnet_user"`
	RemotePath   string `toml:"remote_path"`
}

type QuickSyncConfig struct {
	RepoDir    string `toml:"repo_dir"`
	RemotePath string `toml:"remote_path"`
}

type PaletteConfig struct {
	Base00 string `toml:"base00" json:"base00"`
	Base01 string `toml:"base01" json:"base01"`
	Base02 string `toml:"base02" json:"base02"`
	Base03 string `toml:"base03" json:"base03"`
	Base04 string `toml:"base04" json:"base04"`
	Base05 string `toml:"base05" json:"base05"`
	Base06 string `toml:"base06" json:"base06"`
	Base07 string `toml:"base07" json:"base07"`
	Base08 string `toml:"base08" json:"base08"`
	Base09 string `toml:"base09" json:"base09"`
	Base0A string `toml:"base0A" json:"base0A"`
	Base0B string `toml:"base0B" json:"base0B"`
	Base0C string `toml:"base0C" json:"base0C"`
	Base0D string `toml:"base0D" json:"base0D"`
	Base0E string `toml:"base0E" json:"base0E"`
	Base0F string `toml:"base0F" json:"base0F"`
}

func (p PaletteConfig) Get(key string) string {
	v := reflect.ValueOf(p)
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if tag := field.Tag.Get("toml"); tag == key {
			return v.Field(i).String()
		}
	}
	return ""
}

// Merge returns a new PaletteConfig with values from other overriding non-empty values
func (p PaletteConfig) Merge(other PaletteConfig) PaletteConfig {
	result := p
	if other.Base00 != "" {
		result.Base00 = other.Base00
	}
	if other.Base01 != "" {
		result.Base01 = other.Base01
	}
	if other.Base02 != "" {
		result.Base02 = other.Base02
	}
	if other.Base03 != "" {
		result.Base03 = other.Base03
	}
	if other.Base04 != "" {
		result.Base04 = other.Base04
	}
	if other.Base05 != "" {
		result.Base05 = other.Base05
	}
	if other.Base06 != "" {
		result.Base06 = other.Base06
	}
	if other.Base07 != "" {
		result.Base07 = other.Base07
	}
	if other.Base08 != "" {
		result.Base08 = other.Base08
	}
	if other.Base09 != "" {
		result.Base09 = other.Base09
	}
	if other.Base0A != "" {
		result.Base0A = other.Base0A
	}
	if other.Base0B != "" {
		result.Base0B = other.Base0B
	}
	if other.Base0C != "" {
		result.Base0C = other.Base0C
	}
	if other.Base0D != "" {
		result.Base0D = other.Base0D
	}
	if other.Base0E != "" {
		result.Base0E = other.Base0E
	}
	if other.Base0F != "" {
		result.Base0F = other.Base0F
	}
	return result
}

// Merge returns a new WallpaperConfig with values from other overriding non-empty values
func (w WallpaperConfig) Merge(other WallpaperConfig) WallpaperConfig {
	result := w
	if other.Dir != "" {
		result.Dir = other.Dir
	}
	if other.Default != "" {
		result.Default = other.Default
	}
	return result
}

// Merge returns a new ScreenshotConfig with values from other overriding non-empty values
func (s ScreenshotConfig) Merge(other ScreenshotConfig) ScreenshotConfig {
	result := s
	if other.Dir != "" {
		result.Dir = other.Dir
	}
	return result
}

// Merge returns a new BackupConfig with values from other overriding non-empty values
func (b BackupConfig) Merge(other BackupConfig) BackupConfig {
	result := b
	if other.RsyncnetUser != "" {
		result.RsyncnetUser = other.RsyncnetUser
	}
	if other.RemotePath != "" {
		result.RemotePath = other.RemotePath
	}
	return result
}

// Merge returns a new QuickSyncConfig with values from other overriding non-empty values
func (q QuickSyncConfig) Merge(other QuickSyncConfig) QuickSyncConfig {
	result := q
	if other.RepoDir != "" {
		result.RepoDir = other.RepoDir
	}
	if other.RemotePath != "" {
		result.RemotePath = other.RemotePath
	}
	return result
}

// Merge returns a new BrowserConfig with values from other overriding non-empty values
func (b BrowserConfig) Merge(other BrowserConfig) BrowserConfig {
	result := b
	if other.Default != "" {
		result.Default = other.Default
	}
	if other.Engine != "" {
		result.Engine = other.Engine
	}
	return result
}

// Merge returns a new GlobalConfig with values from other overriding non-empty values.
// For maps (Workspaces, Hosts, Webapps), keys are merged additively.
func (g GlobalConfig) Merge(other GlobalConfig) GlobalConfig {
	result := g

	// Deep copy maps to avoid aliasing
	if result.Workspaces == nil {
		result.Workspaces = make(map[string]int)
	} else {
		// Copy the map
		newWorkspaces := make(map[string]int, len(result.Workspaces))
		maps.Copy(newWorkspaces, result.Workspaces)
		result.Workspaces = newWorkspaces
	}

	if result.Hosts == nil {
		result.Hosts = make(map[string]HostConfig)
	} else {
		// Copy the map
		newHosts := make(map[string]HostConfig, len(result.Hosts))
		maps.Copy(newHosts, result.Hosts)
		result.Hosts = newHosts
	}

	if result.Webapps == nil {
		result.Webapps = make(map[string]WebappConfig)
	} else {
		// Copy the map
		newWebapps := make(map[string]WebappConfig, len(result.Webapps))
		maps.Copy(newWebapps, result.Webapps)
		result.Webapps = newWebapps
	}

	if result.LazyTools == nil {
		result.LazyTools = make(map[string]LazyToolConfig)
	} else {
		newLazyTools := make(map[string]LazyToolConfig, len(result.LazyTools))
		maps.Copy(newLazyTools, result.LazyTools)
		result.LazyTools = newLazyTools
	}

	// Merge Workspaces map (additive, override on conflict)
	maps.Copy(result.Workspaces, other.Workspaces)

	// Merge Hosts map (additive, override on conflict)
	maps.Copy(result.Hosts, other.Hosts)

	// Merge Webapps map (additive, override on conflict)
	maps.Copy(result.Webapps, other.Webapps)

	// Merge LazyTools map
	maps.Copy(result.LazyTools, other.LazyTools)

	// Merge nested configs using their Merge methods
	result.Desktop.Wallpaper = result.Desktop.Wallpaper.Merge(other.Desktop.Wallpaper)
	result.Screenshot = result.Screenshot.Merge(other.Screenshot)
	result.Backup = result.Backup.Merge(other.Backup)
	result.QuickSync = result.QuickSync.Merge(other.QuickSync)
	result.Browser = result.Browser.Merge(other.Browser)
	result.Palette = result.Palette.Merge(other.Palette)

	return result
}

func Load() (*Config, error) {
	commonCfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	raw := make(map[string]interface{})
	home, _ := os.UserHomeDir()
	dotfiles, _ := env.GetDotfilesRoot()

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

// LoadConfig reads the global configuration using a layered approach:
// 1. Start with hardcoded defaults
// 2. Merge with $DOTFILES/settings.toml (if exists)
// 3. Merge with ~/settings.toml (if exists)
// It also expands environment variables and tilde (~) in paths.
func LoadConfig() (*GlobalConfig, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dotfiles, _ := env.GetDotfilesRoot()

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

	// 2. Load and merge base config from $DOTFILES/settings.toml
	if dotfiles != "" {
		dotfilesSettingsPath := filepath.Join(dotfiles, "settings.toml")
		if _, err := os.Stat(dotfilesSettingsPath); err == nil {
			var dotfilesConfig GlobalConfig
			if _, err := toml.DecodeFile(dotfilesSettingsPath, &dotfilesConfig); err == nil {
				config = config.Merge(dotfilesConfig)
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
	}

	// 4. Expand paths
	config.Desktop.Wallpaper.Dir = env.ExpandPath(config.Desktop.Wallpaper.Dir)
	config.Desktop.Wallpaper.Default = env.ExpandPath(config.Desktop.Wallpaper.Default)
	config.Screenshot.Dir = env.ExpandPath(config.Screenshot.Dir)
	config.QuickSync.RepoDir = env.ExpandPath(config.QuickSync.RepoDir)

	return &config, nil
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

	// Use json as a trick to unmarshal into the target type
	data, err := json.Marshal(current)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, val)
}
