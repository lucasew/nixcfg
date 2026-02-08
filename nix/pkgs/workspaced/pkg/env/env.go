package env

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"workspaced/pkg/api"
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
	return "", api.ErrDotfilesRootNotFound
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

// GetConfigDir returns the path to the user config directory for workspaced (~/.config/workspaced)
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".config/workspaced")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
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

func NormalizeURL(url string) string {
	if strings.HasPrefix(url, "file://") || strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	}
	if strings.HasPrefix(url, "/") || strings.HasPrefix(url, "~/") {
		return "file://" + ExpandPath(url)
	}
	return "https://" + url
}
