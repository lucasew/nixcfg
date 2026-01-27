package common

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"workspaced/pkg/types"

	"github.com/BurntSushi/toml"
)

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

type GlobalConfig struct {
	Workspaces map[string]int   `toml:"workspaces"`
	Wallpaper  WallpaperConfig  `toml:"wallpaper"`
	Screenshot ScreenshotConfig `toml:"screenshot"`
}

type WallpaperConfig struct {
	Dir     string `toml:"dir"`
	Default string `toml:"default"`
}

type ScreenshotConfig struct {
	Dir string `toml:"dir"`
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
	}

	if _, err := os.Stat(settingsPath); err == nil {
		if _, err := toml.DecodeFile(settingsPath, config); err != nil {
			return config, err
		}
	}

	// Expand paths
	config.Wallpaper.Dir = ExpandPath(config.Wallpaper.Dir)
	config.Screenshot.Dir = ExpandPath(config.Screenshot.Dir)

	return config, nil
}

func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return os.ExpandEnv(path)
}
