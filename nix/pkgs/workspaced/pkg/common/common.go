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
	"slices"
	"strings"
	"workspaced/pkg/types"

	"github.com/BurntSushi/toml"
)

// EssentialPaths defines the list of directories that must be present in the PATH
// for the application to function correctly on NixOS.
// These typically include locations for wrapped binaries and current system software.
var EssentialPaths = []string{"/run/wrappers/bin", "/run/current-system/sw/bin"}

func init() {
	if home, err := os.UserHomeDir(); err == nil {
		EssentialPaths = append(EssentialPaths, filepath.Join(home, ".nix-profile/bin"))
	}
	newPath := strings.Split(os.Getenv("PATH"), ":")

	for _, path := range EssentialPaths {
		if !slices.Contains(newPath, path) {
			newPath = append([]string{path}, newPath...)
		}
	}
	os.Setenv("PATH", strings.Join(newPath, ":"))
}

// GetLogger retrieves the logger instance from the context.
// It returns the default slog logger if no logger is found in the context.
func GetLogger(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(types.LoggerKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

// ChannelLogHandler is a custom slog.Handler that broadcasts log records to a channel.
// This is used to stream server-side logs to the client via the daemon connection.
type ChannelLogHandler struct {
	Out    chan<- types.StreamPacket
	Parent slog.Handler
	Ctx    context.Context
}

// Enabled reports whether the handler handles records at the given level.
func (h *ChannelLogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

// Handle processes a log record, marshals it to JSON, and sends it as a StreamPacket.
// It also delegates to the parent handler if one is configured.
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

// WithAttrs returns a new ChannelLogHandler with the given attributes added.
func (h *ChannelLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ChannelLogHandler{Out: h.Out, Parent: h.Parent.WithAttrs(attrs), Ctx: h.Ctx}
}

// WithGroup returns a new ChannelLogHandler with the given group name.
func (h *ChannelLogHandler) WithGroup(name string) slog.Handler {
	return &ChannelLogHandler{Out: h.Out, Parent: h.Parent.WithGroup(name), Ctx: h.Ctx}
}

// RunCmd creates an exec.Cmd with environment variables injected from the context.
// It ensures that the PATH includes EssentialPaths.
func RunCmd(ctx context.Context, name string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, name, args...)
	return cmd
}


// InheritContextWriters configures the command's Stdout and Stderr to write to the writers
// stored in the context, allowing output capture or redirection.
func InheritContextWriters(ctx context.Context, cmd *exec.Cmd) {
	if stdout, ok := ctx.Value(types.StdoutKey).(io.Writer); ok {
		cmd.Stdout = stdout
	}
	if stderr, ok := ctx.Value(types.StderrKey).(io.Writer); ok {
		cmd.Stderr = stderr
	}
}

// GetRPC determines the appropriate IPC command for the current window manager.
// It checks for WAYLAND_DISPLAY to decide between "swaymsg" (Wayland) and "i3-msg" (X11).
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

// GetDotfilesRoot locates the root directory of the dotfiles repository.
// It checks standard locations (~/.dotfiles, /etc/.dotfiles).
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

// GetHostname returns the current system hostname.
func GetHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}

// IsRiverwood checks if the current host is "riverwood".
func IsRiverwood() bool {
	return GetHostname() == "riverwood"
}

// IsWhiterun checks if the current host is "whiterun".
func IsWhiterun() bool {
	return GetHostname() == "whiterun"
}

// IsPhone checks if the environment suggests we are running on a phone (Termux).
func IsPhone() bool {
	return os.Getenv("TERMUX_VERSION") != ""
}

// IsBinaryAvailable checks if a command exists in the PATH.
func IsBinaryAvailable(ctx context.Context, name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// IsInStore checks if the dotfiles root is located inside the Nix store.
func IsInStore() bool {
	root, err := GetDotfilesRoot()
	if err != nil {
		return false
	}
	return strings.HasPrefix(root, "/nix/store")
}

// IsNixOS checks if the system is NixOS by verifying the existence of /etc/NIXOS.
func IsNixOS() bool {
	_, err := os.Stat("/etc/NIXOS")
	return err == nil
}

// GlobalConfig represents the schema of the settings.toml file.
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

// LoadConfig reads the global configuration from ~/settings.toml.
// If the file is missing or invalid, it falls back to hardcoded defaults.
// It also expands environment variables and tilde (~) in paths.
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

// ExpandPath expands the tilde (~) to the user's home directory
// and expands environment variables (e.g. $HOME, ${VAR}) in the path.
func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return os.ExpandEnv(path)
}
