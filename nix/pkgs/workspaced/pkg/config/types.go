package config

import (
	"maps"
)

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

	return result
}
