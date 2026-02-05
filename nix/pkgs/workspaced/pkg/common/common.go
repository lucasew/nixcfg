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
)

var (
	ErrCommandNotFound = fmt.Errorf("command not found")
)

// EssentialPaths defines the list of directories that must be present in the PATH
// for the application to function correctly on NixOS.
// These typically include locations for wrapped binaries and current system software.
var EssentialPaths = []string{"/run/wrappers/bin", "/run/current-system/sw/bin"}

func init() {
	if home, err := os.UserHomeDir(); err == nil {
		EssentialPaths = append(EssentialPaths, filepath.Join(home, ".nix-profile/bin"))
	}
	if root, err := GetDotfilesRoot(); err == nil && root != "" {
		EssentialPaths = append(EssentialPaths, filepath.Join(root, "bin/shim"))
	}
	if dataDir, err := GetUserDataDir(); err == nil && dataDir != "" {
		EssentialPaths = append(EssentialPaths, filepath.Join(dataDir, "shim/global"))
	}
	newPath := strings.Split(os.Getenv("PATH"), ":")

	for _, path := range EssentialPaths {
		if !slices.Contains(newPath, path) {
			newPath = append([]string{path}, newPath...)
		}
	}
	if err := os.Setenv("PATH", strings.Join(newPath, ":")); err != nil {
		panic(err)
	}
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
// It uses the custom Which implementation to avoid SIGSYS errors on Android/Termux.
func RunCmd(ctx context.Context, name string, args ...string) *exec.Cmd {
	// Resolve the full path using our custom Which to avoid SIGSYS on Android
	fullPath, err := Which(ctx, name)
	if err != nil {
		// If Which fails, fall back to the original name
		// This allows exec.CommandContext to handle the error properly
		fullPath = name
	}
	cmd := exec.CommandContext(ctx, fullPath, args...)
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

// GetUserDataDir returns the path to the user data directory for workspaced (~/.local/share/workspaced)
func GetUserDataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(home, ".local/share/workspaced")
	if err := os.MkdirAll(path, 0755); err != nil {
		return "", err
	}
	return path, nil
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

// Which locates a command in the PATH without using os/exec.LookPath
// to avoid SIGSYS errors on Android/Go 1.24.
func Which(ctx context.Context, name string) (string, error) {
	if filepath.IsAbs(name) {
		if _, err := os.Stat(name); err == nil {
			return name, nil
		}
		return "", fmt.Errorf("file not found: %s", name)
	}

	path := os.Getenv("PATH")
	if env, ok := ctx.Value(types.EnvKey).([]string); ok {
		for _, e := range env {
			if strings.HasPrefix(e, "PATH=") {
				path = strings.TrimPrefix(e, "PATH=")
				break
			}
		}
	}

	for _, dir := range filepath.SplitList(path) {
		fullPath := filepath.Join(dir, name)
		if info, err := os.Stat(fullPath); err == nil && !info.IsDir() {
			return fullPath, nil
		}
	}
	return "", fmt.Errorf("%w: %s", ErrCommandNotFound, name)
}

// IsBinaryAvailable checks if a command exists in the PATH using the internal Which implementation.
func IsBinaryAvailable(ctx context.Context, name string) bool {
	_, err := Which(ctx, name)
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

// ExpandPath expands the tilde (~) to the user's home directory
// and expands environment variables (e.g. $HOME, ${VAR}) in the path.
func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return os.ExpandEnv(path)
}

func ToTitleCase(s string) string {
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")
	words := strings.Fields(s)
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
		}
	}
	return strings.Join(words, " ")
}

func NormalizeURL(url string) string {
	if strings.HasPrefix(url, "file://") || strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	}
	if strings.HasPrefix(url, "/") || strings.HasPrefix(url, "~/") {
		return "file://" + ExpandPath(url)
	}
	return "https://" + url
}
