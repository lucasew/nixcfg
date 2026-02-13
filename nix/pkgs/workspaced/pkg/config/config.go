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
	LazyTools  map[string]LazyToolConfig `toml:"lazy_tools"`
	Palette    PaletteConfig             `toml:"palette"`
	Fonts      FontsConfig               `toml:"fonts"`
	Modules    map[string]any            `toml:"modules"`
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
	// Base24 extended colors
	Base10 string `toml:"base10,omitempty" json:"base10,omitempty"`
	Base11 string `toml:"base11,omitempty" json:"base11,omitempty"`
	Base12 string `toml:"base12,omitempty" json:"base12,omitempty"`
	Base13 string `toml:"base13,omitempty" json:"base13,omitempty"`
	Base14 string `toml:"base14,omitempty" json:"base14,omitempty"`
	Base15 string `toml:"base15,omitempty" json:"base15,omitempty"`
	Base16 string `toml:"base16,omitempty" json:"base16,omitempty"`
	Base17 string `toml:"base17,omitempty" json:"base17,omitempty"`
}

type FontsConfig struct {
	Serif     string `toml:"serif"`
	SansSerif string `toml:"sans_serif"`
	Monospace string `toml:"monospace"`
	Emoji     string `toml:"emoji"`
}

func (f FontsConfig) Merge(other FontsConfig) FontsConfig {
	result := f
	if other.Serif != "" {
		result.Serif = other.Serif
	}
	if other.SansSerif != "" {
		result.SansSerif = other.SansSerif
	}
	if other.Monospace != "" {
		result.Monospace = other.Monospace
	}
	if other.Emoji != "" {
		result.Emoji = other.Emoji
	}
	return result
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
	if other.Base10 != "" {
		result.Base10 = other.Base10
	}
	if other.Base11 != "" {
		result.Base11 = other.Base11
	}
	if other.Base12 != "" {
		result.Base12 = other.Base12
	}
	if other.Base13 != "" {
		result.Base13 = other.Base13
	}
	if other.Base14 != "" {
		result.Base14 = other.Base14
	}
	if other.Base15 != "" {
		result.Base15 = other.Base15
	}
	if other.Base16 != "" {
		result.Base16 = other.Base16
	}
	if other.Base17 != "" {
		result.Base17 = other.Base17
	}

	// Auto-fill base24 extended colors by repeating accent colors if empty
	if result.Base10 == "" && result.Base08 != "" {
		result.Base10 = result.Base08
	}
	if result.Base11 == "" && result.Base09 != "" {
		result.Base11 = result.Base09
	}
	if result.Base12 == "" && result.Base0A != "" {
		result.Base12 = result.Base0A
	}
	if result.Base13 == "" && result.Base0B != "" {
		result.Base13 = result.Base0B
	}
	if result.Base14 == "" && result.Base0C != "" {
		result.Base14 = result.Base0C
	}
	if result.Base15 == "" && result.Base0D != "" {
		result.Base15 = result.Base0D
	}
	if result.Base16 == "" && result.Base0E != "" {
		result.Base16 = result.Base0E
	}
	if result.Base17 == "" && result.Base0F != "" {
		result.Base17 = result.Base0F
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
// For maps (Workspaces, Hosts), keys are merged additively.
func (g GlobalConfig) Merge(other GlobalConfig) (GlobalConfig, error) {
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

	if result.LazyTools == nil {
		result.LazyTools = make(map[string]LazyToolConfig)
	} else {
		newLazyTools := make(map[string]LazyToolConfig, len(result.LazyTools))
		maps.Copy(newLazyTools, result.LazyTools)
		result.LazyTools = newLazyTools
	}

	if result.Modules == nil {
		result.Modules = make(map[string]any)
	} else {
		newModules := make(map[string]any, len(result.Modules))
		maps.Copy(newModules, result.Modules)
		result.Modules = newModules
	}

	// Merge Workspaces map (additive, override on conflict)
	maps.Copy(result.Workspaces, other.Workspaces)

	// Merge Hosts map (additive, override on conflict)
	maps.Copy(result.Hosts, other.Hosts)

	// Merge LazyTools map
	maps.Copy(result.LazyTools, other.LazyTools)

	// Merge Modules map strictly
	if err := MergeStrict(result.Modules, other.Modules); err != nil {
		return result, err
	}

	// Merge nested configs using their Merge methods
	result.Desktop.Wallpaper = result.Desktop.Wallpaper.Merge(other.Desktop.Wallpaper)
	result.Screenshot = result.Screenshot.Merge(other.Screenshot)
	result.Backup = result.Backup.Merge(other.Backup)
	result.QuickSync = result.QuickSync.Merge(other.QuickSync)
	result.Browser = result.Browser.Merge(other.Browser)
	result.Palette = result.Palette.Merge(other.Palette)
	result.Fonts = result.Fonts.Merge(other.Fonts)

	return result, nil
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
		LazyTools: make(map[string]LazyToolConfig),
		Modules:   make(map[string]any),
	}

	// 2. Load and merge base config from $DOTFILES/settings.toml
	if dotfiles != "" {
		dotfilesSettingsPath := filepath.Join(dotfiles, "settings.toml")
		if _, err := os.Stat(dotfilesSettingsPath); err == nil {
			var dotfilesConfig GlobalConfig
			if _, err := toml.DecodeFile(dotfilesSettingsPath, &dotfilesConfig); err == nil {
				var err error
				config, err = config.Merge(dotfilesConfig)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	// 3. Load and merge user config from ~/settings.toml
	userSettingsPath := filepath.Join(home, "settings.toml")
	if _, err := os.Stat(userSettingsPath); err == nil {
		var userConfig GlobalConfig
		if _, err := toml.DecodeFile(userSettingsPath, &userConfig); err == nil {
			var err error
			config, err = config.Merge(userConfig)
			if err != nil {
				return nil, err
			}
		}
	}

	// 4. Expand paths
	config.Desktop.Wallpaper.Dir = env.ExpandPath(config.Desktop.Wallpaper.Dir)
	config.Desktop.Wallpaper.Default = env.ExpandPath(config.Desktop.Wallpaper.Default)
	config.Screenshot.Dir = env.ExpandPath(config.Screenshot.Dir)
	config.QuickSync.RepoDir = env.ExpandPath(config.QuickSync.RepoDir)

	return &config, nil
}

// UnmarshalKey extracts a sub-config by key path (e.g. "modules.webapp") into a target struct.
func (c *Config) UnmarshalKey(key string, val interface{}) error {
	parts := strings.Split(key, ".")
	var current interface{} = c.raw

	for _, part := range parts {
		if m, ok := current.(map[string]interface{}); ok {
			v, ok := m[part]
			if !ok {
				return fmt.Errorf("key %q not found in config", key)
			}
			current = v
		} else if m, ok := current.(map[string]any); ok { // Handle both map types
			v, ok := m[part]
			if !ok {
				return fmt.Errorf("key %q not found in config", key)
			}
			current = v
		} else {
			return fmt.Errorf("key %q not found or not a map", key)
		}
	}

	if current == nil {
		return fmt.Errorf("value for key %q is nil", key)
	}

	// Use json as a trick to unmarshal into the target type
	data, err := json.Marshal(current)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, val)
}

// Module extracts configuration for a specific module into a target struct.
func (c *Config) Module(name string, target interface{}) error {
	return c.UnmarshalKey("modules."+name, target)
}

func MergeStrict(dst, src map[string]any) error {
	for k, v := range src {
		if v == nil {
			continue
		}
		if existing, ok := dst[k]; ok && existing != nil {
			// Check for lists
			if reflect.TypeOf(v).Kind() == reflect.Slice || reflect.TypeOf(v).Kind() == reflect.Array {
				return fmt.Errorf("lists are forbidden in strict config (key: %s)", k)
			}

			// If both are maps, recurse
			if vMap, ok := v.(map[string]any); ok {
				if existingMap, ok := existing.(map[string]any); ok {
					if err := MergeStrict(existingMap, vMap); err != nil {
						return err
					}
					continue
				}
			}

			// If atomic and different, error (no substitution)
			if !reflect.DeepEqual(existing, v) {
				return fmt.Errorf("substitution forbidden: key %q already has value %v, cannot overwrite with %v", k, existing, v)
			}
		} else {
			// New key
			if reflect.TypeOf(v).Kind() == reflect.Slice || reflect.TypeOf(v).Kind() == reflect.Array {
				return fmt.Errorf("lists are forbidden in strict config (key: %s)", k)
			}
			dst[k] = v
		}
	}
	return nil
}

func LoadFiles(paths []string) (*GlobalConfig, error) {
	config := &GlobalConfig{
		Workspaces: make(map[string]int),
		Hosts:      make(map[string]HostConfig),
		LazyTools:  make(map[string]LazyToolConfig),
		Modules:    make(map[string]any),
	}

	mergedRaw := make(map[string]any)

	for _, path := range paths {
		var currentRaw map[string]any
		if _, err := toml.DecodeFile(path, &currentRaw); err != nil {
			return nil, fmt.Errorf("failed to decode %s: %w", path, err)
		}
		if err := MergeStrict(mergedRaw, currentRaw); err != nil {
			return nil, fmt.Errorf("strict merge failed for %s: %w", path, err)
		}
	}

	// Use JSON trick to unmarshal the merged map into the struct
	data, err := json.Marshal(mergedRaw)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}
