package host

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"workspaced/pkg/types"
)

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
