package common

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"workspaced/pkg/types"

	"github.com/BurntSushi/toml"
)

func GetLogger(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(types.LoggerKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

func RunCmd(ctx context.Context, name string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, name, args...)
	if env, ok := ctx.Value(types.EnvKey).([]string); ok {
		cmd.Env = env
	}
	return cmd
}

func GetRPC(ctx context.Context) string {
	if env, ok := ctx.Value(types.EnvKey).([]string); ok {
		for _, e := range env {
			if strings.HasPrefix(e, "WAYLAND_DISPLAY=") {
				return "swaymsg"
			}
		}
	}
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		return "swaymsg"
	}
	return "i3-msg"
}

func GetDotfilesRoot() (string, error) {
	home, err := os.UserHomeDir()
	if err == nil {
		path := filepath.Join(home, ".dotfiles")
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			return path, nil
		}
	}
	// Fallback to /etc/.dotfiles
	path := "/etc/.dotfiles"
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return path, nil
	}
	return "", fmt.Errorf("could not find dotfiles root")
}

func IsRiverwood() bool {
	hostname, _ := os.Hostname()
	return hostname == "riverwood"
}

func IsWhiterun() bool {
	hostname, _ := os.Hostname()
	return hostname == "whiterun"
}

func IsPhone() bool {
	return os.Getenv("TERMUX_VERSION") != ""
}

func IsBinaryAvailable(ctx context.Context, name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func IsInStore() bool {
	root, err := GetDotfilesRoot()
	if err != nil {
		return false
	}
	return strings.HasPrefix(root, "/nix/store")
}

type GlobalConfig struct {
	Workspaces map[string]int    `toml:"workspaces"`
	Wallpaper  WallpaperConfig   `toml:"wallpaper"`
	Screenshot ScreenshotConfig  `toml:"screenshot"`
	Hosts      map[string]string `toml:"hosts"`
	Backup     BackupConfig      `toml:"backup"`
	QuickSync  QuickSyncConfig   `toml:"quicksync"`
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

func LoadConfig() (*GlobalConfig, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	settingsPath := filepath.Join(home, "settings.toml")

	dotfiles, _ := GetDotfilesRoot()

	config := &GlobalConfig{
		Workspaces: map[string]int{
			"www":  1,
			"meet": 2,
		},
		Wallpaper: WallpaperConfig{
			Dir: filepath.Join(dotfiles, "assets/wallpapers"),
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
		Hosts: map[string]string{
			"whiterun": "a8:a1:59:9c:ab:32",
		},
	}

	if _, err := os.Stat(settingsPath); err == nil {
		if _, err := toml.DecodeFile(settingsPath, config); err != nil {
			return config, err
		}
	}

	// Expand paths
	config.Wallpaper.Dir = ExpandPath(config.Wallpaper.Dir)
	config.Screenshot.Dir = ExpandPath(config.Screenshot.Dir)
	config.QuickSync.RepoDir = ExpandPath(config.QuickSync.RepoDir)

	return config, nil
}

func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return os.ExpandEnv(path)
}
