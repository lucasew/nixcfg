package common

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"workspaced/pkg/types"

	"github.com/BurntSushi/toml"
)

var EssentialPaths = []string{"/run/wrappers/bin", "/run/current-system/sw/bin"}

func GetLogger(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(types.LoggerKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

type ChannelLogHandler struct {
	Out    chan<- types.StreamPacket
	Parent slog.Handler
	Ctx    context.Context
}

func (h *ChannelLogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (h *ChannelLogHandler) Handle(ctx context.Context, r slog.Record) error {
	entry := types.LogEntry{
		Level:   r.Level.String(),
		Message: r.Message,
		Attrs:   make(map[string]any),
	}
	r.Attrs(func(a slog.Attr) bool {
		entry.Attrs[a.Key] = a.Value.Any()
		return true
	})
	payload, _ := json.Marshal(entry)

	select {
	case h.Out <- types.StreamPacket{Type: "log", Payload: payload}:
	case <-h.Ctx.Done():
		return h.Ctx.Err()
	}

	if h.Parent != nil {
		return h.Parent.Handle(ctx, r)
	}
	return nil
}

func (h *ChannelLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ChannelLogHandler{Out: h.Out, Parent: h.Parent.WithAttrs(attrs), Ctx: h.Ctx}
}

func (h *ChannelLogHandler) WithGroup(name string) slog.Handler {
	return &ChannelLogHandler{Out: h.Out, Parent: h.Parent.WithGroup(name), Ctx: h.Ctx}
}

func RunCmd(ctx context.Context, name string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, name, args...)
	if env, ok := ctx.Value(types.EnvKey).([]string); ok {
		cmd.Env = EnsureEssentialPaths(env)
	}
	return cmd
}

func EnsureEssentialPaths(env []string) []string {
	var newEnv []string
	foundPath := false
	for _, e := range env {
		if strings.HasPrefix(e, "PATH=") {
			foundPath = true
			currentPath := strings.TrimPrefix(e, "PATH=")
			pathParts := strings.Split(currentPath, ":")

			// Check which essential paths are missing and prepend them
			for i := len(EssentialPaths) - 1; i >= 0; i-- {
				p := EssentialPaths[i]
				missing := true
				for _, cp := range pathParts {
					if cp == p {
						missing = false
						break
					}
				}
				if missing {
					pathParts = append([]string{p}, pathParts...)
				}
			}
			newEnv = append(newEnv, "PATH="+strings.Join(pathParts, ":"))
		} else {
			newEnv = append(newEnv, e)
		}
	}
	if !foundPath {
		newEnv = append(newEnv, "PATH="+strings.Join(EssentialPaths, ":")+":"+os.Getenv("PATH"))
	}
	return newEnv
}

func InheritContextWriters(ctx context.Context, cmd *exec.Cmd) {
	if stdout, ok := ctx.Value(types.StdoutKey).(io.Writer); ok {
		cmd.Stdout = stdout
	}
	if stderr, ok := ctx.Value(types.StderrKey).(io.Writer); ok {
		cmd.Stderr = stderr
	}
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

func GetHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}

func IsRiverwood() bool {
	return GetHostname() == "riverwood"
}

func IsWhiterun() bool {
	return GetHostname() == "whiterun"
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

func IsNixOS() bool {
	_, err := os.Stat("/etc/NIXOS")
	return err == nil
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
