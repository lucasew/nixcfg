package config

import (
	"os"
	"path/filepath"
	"workspaced/pkg/common"

	"github.com/BurntSushi/toml"
)

// loadDefault reads the global configuration using a layered approach:
// 1. Start with hardcoded defaults
// 2. Merge with $DOTFILES/settings.toml (if exists)
// 3. Merge with ~/settings.toml (if exists)
// It also expands environment variables and tilde (~) in paths.
func loadDefault() (*GlobalConfig, error) {
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
	config.Desktop.Wallpaper.Dir = common.ExpandPath(config.Desktop.Wallpaper.Dir)
	config.Desktop.Wallpaper.Default = common.ExpandPath(config.Desktop.Wallpaper.Default)
	config.Screenshot.Dir = common.ExpandPath(config.Screenshot.Dir)
	config.QuickSync.RepoDir = common.ExpandPath(config.QuickSync.RepoDir)

	return &config, nil
}
