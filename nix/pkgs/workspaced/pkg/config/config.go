package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"workspaced/pkg/driver"
	"workspaced/pkg/env"

	"github.com/BurntSushi/toml"
)

type Config struct {
	*GlobalConfig
	raw map[string]interface{}
}

type GlobalConfig struct {
	Workspaces map[string]int            `toml:"workspaces" json:"workspaces"`
	Desktop    DesktopConfig             `toml:"desktop" json:"desktop"`
	Screenshot ScreenshotConfig          `toml:"screenshot" json:"screenshot"`
	Hosts      map[string]HostConfig     `toml:"hosts" json:"hosts"`
	Backup     BackupConfig              `toml:"backup" json:"backup"`
	QuickSync  QuickSyncConfig           `toml:"quicksync" json:"quicksync"`
	Browser    BrowserConfig             `toml:"browser" json:"browser"`
	LazyTools  map[string]LazyToolConfig `toml:"lazy_tools" json:"lazy_tools"`
	Palette    PaletteConfig             `toml:"palette" json:"palette"`
	Fonts      FontsConfig               `toml:"fonts" json:"fonts"`
	Modules    map[string]any            `toml:"modules" json:"modules"`
	Drivers    map[string]map[string]int `toml:"drivers" json:"drivers"`
}

type DesktopConfig struct {
	DarkMode  bool            `toml:"dark_mode" json:"dark_mode"`
	Wallpaper WallpaperConfig `toml:"wallpaper" json:"wallpaper"`
}

type LazyToolConfig struct {
	Version string   `toml:"version" json:"version"`
	Pkg     string   `toml:"pkg" json:"pkg"`
	Global  bool     `toml:"global" json:"global"`
	Alias   string   `toml:"alias" json:"alias"`
	Bins    []string `toml:"bins" json:"bins"`
}

type HostConfig struct {
	MAC         string `toml:"mac" json:"mac"`
	TailscaleIP string `toml:"tailscale_ip" json:"tailscale_ip"`
	ZerotierIP  string `toml:"zerotier_ip" json:"zerotier_ip"`
	Port        int    `toml:"port" json:"port"`
	User        string `toml:"user" json:"user"`
}

type BrowserConfig struct {
	Default string `toml:"default" json:"default"`
	Engine  string `toml:"webapp" json:"webapp"`
}

type WallpaperConfig struct {
	Dir     string `toml:"dir" json:"dir"`
	Default string `toml:"default" json:"default"`
}

type ScreenshotConfig struct {
	Dir string `toml:"dir" json:"dir"`
}

type BackupConfig struct {
	RsyncnetUser string `toml:"rsyncnet_user" json:"rsyncnet_user"`
	RemotePath   string `toml:"remote_path" json:"remote_path"`
}

type QuickSyncConfig struct {
	RepoDir    string `toml:"repo_dir" json:"repo_dir"`
	RemotePath string `toml:"remote_path" json:"remote_path"`
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
	Serif     string `toml:"serif" json:"serif"`
	SansSerif string `toml:"sans_serif" json:"sans_serif"`
	Monospace string `toml:"monospace" json:"monospace"`
	Emoji     string `toml:"emoji" json:"emoji"`
}

func (c *Config) Module(name string, target interface{}) error {
	return c.UnmarshalKey("modules."+name, target)
}

func (c *Config) UnmarshalKey(key string, val interface{}) error {
	parts := strings.Split(key, ".")
	var current interface{} = c.raw
	for _, part := range parts {
		if mRaw, ok := current.(map[string]any); ok {
			v, ok := mRaw[part]
			if !ok {
				return fmt.Errorf("key %q not found in config", key)
			}
			current = v
		} else if mRaw, ok := current.(map[string]interface{}); ok {
			v, ok := mRaw[part]
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
	data, err := json.Marshal(current)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, val)
}

func Load() (*Config, error) {
	home, _ := os.UserHomeDir()
	dotfiles, _ := env.GetDotfilesRoot()
	gCfg, err := LoadConfigBase()
	if err != nil {
		return nil, err
	}
	structToMap := func(s interface{}) map[string]any {
		data, _ := json.Marshal(s)
		var res map[string]any
		json.Unmarshal(data, &res)
		return res
	}
	rawMerged := structToMap(gCfg)
	userConfigs := []string{}
	if dotfiles != "" {
		userConfigs = append(userConfigs, filepath.Join(dotfiles, "settings.toml"))
	}
	userConfigs = append(userConfigs, filepath.Join(home, "settings.toml"))
	enabledModules := make(map[string]bool)
	for _, path := range userConfigs {
		if _, err := os.Stat(path); err == nil {
			var currentRaw map[string]any
			if _, err := toml.DecodeFile(path, &currentRaw); err == nil {
				if modsRaw, ok := currentRaw["modules"]; ok {
					modsVal := reflect.ValueOf(modsRaw)
					if modsVal.Kind() == reflect.Map {
						for _, modKey := range modsVal.MapKeys() {
							mVal := modsVal.MapIndex(modKey)
							if mVal.Kind() == reflect.Interface {
								mVal = mVal.Elem()
							}
							if mVal.Kind() == reflect.Map {
								eVal := mVal.MapIndex(reflect.ValueOf("enable"))
								if eVal.IsValid() {
									if eVal.Kind() == reflect.Interface {
										eVal = eVal.Elem()
									}
									if eVal.Kind() == reflect.Bool && eVal.Bool() {
										enabledModules[modKey.String()] = true
									}
								}
							}
						}
					}
				}
			}
		}
	}
	modulesDir := filepath.Join(dotfiles, "modules")
	moduleMeta := make(map[string]ModuleMetadata)
	if entries, err := os.ReadDir(modulesDir); err == nil {
		defaultsRaw := make(map[string]any)
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			name := entry.Name()
			if !enabledModules[name] {
				continue
			}
			modPath := filepath.Join(modulesDir, name)
			metaPath := filepath.Join(modPath, "module.toml")
			if _, err := os.Stat(metaPath); err == nil {
				var meta struct {
					Module ModuleMetadata `toml:"module"`
				}
				if _, err := toml.DecodeFile(metaPath, &meta); err != nil {
					return nil, err
				}
				moduleMeta[name] = meta.Module
			}
			defaultsPath := filepath.Join(modPath, "defaults.toml")
			if _, err := os.Stat(defaultsPath); err == nil {
				var currentDefaults map[string]any
				if _, err := toml.DecodeFile(defaultsPath, &currentDefaults); err == nil {
					var driversDefaults map[string]any
					if d, ok := currentDefaults["drivers"]; ok {
						if dMap, ok := d.(map[string]any); ok {
							driversDefaults = dMap
							delete(currentDefaults, "drivers")
						}
					}
					wrapped := map[string]any{"modules": map[string]any{name: currentDefaults}}
					MergeStrict(defaultsRaw, wrapped, false)
					if driversDefaults != nil {
						wrappedDrivers := map[string]any{"drivers": driversDefaults}
						MergeStrict(defaultsRaw, wrappedDrivers, false)
					}
				}
			}
		}
		if err := validateDependencies(enabledModules, moduleMeta); err != nil {
			return nil, err
		}
		MergeStrict(rawMerged, defaultsRaw, true)
	}
	for _, path := range userConfigs {
		if _, err := os.Stat(path); err == nil {
			var currentRaw map[string]any
			if _, err := toml.DecodeFile(path, &currentRaw); err != nil {
				if err := MergeStrict(rawMerged, currentRaw, true); err != nil {
					return nil, err
				}
			}
		}
	}
	finalGCfg := &GlobalConfig{}
	data, _ := json.Marshal(rawMerged)
	json.Unmarshal(data, finalGCfg)
	if modBase16, ok := finalGCfg.Modules["base16"].(map[string]any); ok {
		data, _ := json.Marshal(modBase16)
		json.Unmarshal(data, &finalGCfg.Palette)
	}
	if modFonts, ok := finalGCfg.Modules["fontconfig"].(map[string]any); ok {
		data, _ := json.Marshal(modFonts)
		json.Unmarshal(data, &finalGCfg.Fonts)
	}
	finalGCfg.Desktop.Wallpaper.Dir = env.ExpandPath(finalGCfg.Desktop.Wallpaper.Dir)
	finalGCfg.Desktop.Wallpaper.Default = env.ExpandPath(finalGCfg.Desktop.Wallpaper.Default)
	finalGCfg.Screenshot.Dir = env.ExpandPath(finalGCfg.Screenshot.Dir)
	finalGCfg.QuickSync.RepoDir = env.ExpandPath(finalGCfg.QuickSync.RepoDir)
	if err := driver.SetWeights(finalGCfg.Drivers); err != nil {
		return nil, err
	}
	return &Config{GlobalConfig: finalGCfg, raw: rawMerged}, nil
}

func LoadConfig() (*GlobalConfig, error) {
	cfg, err := Load()
	if err != nil {
		return nil, err
	}
	return cfg.GlobalConfig, nil
}

func LoadConfigBase() (*GlobalConfig, error) {
	home, _ := os.UserHomeDir()
	dotfiles, _ := env.GetDotfilesRoot()
	return &GlobalConfig{
		Workspaces: map[string]int{"www": 1, "meet": 2},
		Desktop:    DesktopConfig{Wallpaper: WallpaperConfig{Dir: filepath.Join(dotfiles, "assets/wallpapers")}},
		Screenshot: ScreenshotConfig{Dir: filepath.Join(home, "Pictures/Screenshots")},
		Backup:     BackupConfig{RsyncnetUser: "de3163@de3163.rsync.net", RemotePath: "backup/lucasew"},
		QuickSync:  QuickSyncConfig{RepoDir: filepath.Join(home, ".personal"), RemotePath: "/data2/home/de3163/git-personal"},
		Hosts:      make(map[string]HostConfig),
		Browser:    BrowserConfig{Default: "zen", Engine: "brave"},
		LazyTools:  make(map[string]LazyToolConfig),
		Modules:    make(map[string]any),
		Drivers:    make(map[string]map[string]int),
	}, nil
}

type ModuleMetadata struct {
	Requires   []string `toml:"requires"`
	Recommends []string `toml:"recommends"`
}

func validateDependencies(enabled map[string]bool, meta map[string]ModuleMetadata) error {
	for name := range enabled {
		m, ok := meta[name]
		if !ok {
			continue
		}
		for _, req := range m.Requires {
			if !enabled[req] {
				return fmt.Errorf("module %q requires %q, but it is not enabled", name, req)
			}
		}
		for _, rec := range m.Recommends {
			if !enabled[rec] {
				fmt.Printf("Warning: module %q recommends %q, but it is not enabled\n", name, rec)
			}
		}
	}
	deps := make(map[string][]string)
	for name, m := range meta {
		deps[name] = m.Requires
	}
	return detectCycles(deps)
}

func detectCycles(deps map[string][]string) error {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	var check func(node string) error
	check = func(node string) error {
		visited[node] = true
		recStack[node] = true
		for _, neighbor := range deps[node] {
			if !visited[neighbor] {
				if err := check(neighbor); err != nil {
					return err
				}
			} else if recStack[neighbor] {
				return fmt.Errorf("circular dependency detected involving module %q", neighbor)
			}
		}
		recStack[node] = false
		return nil
	}
	for node := range deps {
		if !visited[node] {
			if err := check(node); err != nil {
				return err
			}
		}
	}
	return nil
}

func MergeStrict(dst, src map[string]any, allowSubstitution bool) error {
	for k, v := range src {
		if v == nil {
			continue
		}
		if existing, ok := dst[k]; ok && existing != nil {
			if reflect.TypeOf(v).Kind() == reflect.Slice || reflect.TypeOf(v).Kind() == reflect.Array {
				return fmt.Errorf("lists are forbidden in strict config (key: %s)", k)
			}
			vVal := reflect.ValueOf(v)
			eVal := reflect.ValueOf(existing)
			if vVal.Kind() == reflect.Map && eVal.Kind() == reflect.Map {
				vMap := make(map[string]any)
				for _, key := range vVal.MapKeys() {
					vMap[key.String()] = vVal.MapIndex(key).Interface()
				}
				eMap := make(map[string]any)
				for _, key := range eVal.MapKeys() {
					eMap[key.String()] = eVal.MapIndex(key).Interface()
				}
				if err := MergeStrict(eMap, vMap, allowSubstitution); err != nil {
					return err
				}
				dst[k] = eMap
				continue
			}
			if !reflect.DeepEqual(existing, v) {
				if allowSubstitution {
					dst[k] = v
				} else {
					return fmt.Errorf("substitution forbidden: key %q already has value %v, cannot overwrite with %v", k, existing, v)
				}
			}
		} else {
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
		Drivers:    make(map[string]map[string]int),
	}
	mergedRaw := make(map[string]any)
	for _, path := range paths {
		var currentRaw map[string]any
		if _, err := toml.DecodeFile(path, &currentRaw); err != nil {
			return nil, fmt.Errorf("failed to decode %s: %w", path, err)
		}
		if err := MergeStrict(mergedRaw, currentRaw, true); err != nil {
			return nil, fmt.Errorf("strict merge failed for %s: %w", path, err)
		}
	}
	data, err := json.Marshal(mergedRaw)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return config, nil
}
