package apply

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"workspaced/pkg/config"
)

type DesktopProvider struct{}

func (p *DesktopProvider) Name() string {
	return "desktop"
}

type DesktopScript struct {
	Name        string `toml:"name"`
	Description string `toml:"description"`
	Icon        string `toml:"icon"`
	Exec        string `toml:"exec"`
	Terminal    bool   `toml:"terminal"`
	Categories  string `toml:"categories"`
}

func downloadIcon(ctx context.Context, iconURL string, iconsDir string) (string, error) {
	// Generate icon filename from URL hash
	hash := sha256.Sum256([]byte(iconURL))
	iconFilename := fmt.Sprintf("icon-%x.png", hash[:8])
	iconPath := filepath.Join(iconsDir, iconFilename)

	// Check if already downloaded
	if _, err := os.Stat(iconPath); err == nil {
		return iconPath, nil
	}

	// Download the icon
	req, err := http.NewRequestWithContext(ctx, "GET", iconURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download icon: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download icon: status %d", resp.StatusCode)
	}

	// Save icon
	out, err := os.Create(iconPath)
	if err != nil {
		return "", fmt.Errorf("failed to create icon file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		os.Remove(iconPath)
		return "", fmt.Errorf("failed to save icon: %w", err)
	}

	return iconPath, nil
}

func (p *DesktopProvider) GetDesiredState(ctx context.Context) ([]DesiredState, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// Get desktop_scripts section from settings.toml
	desktopScripts := make(map[string]DesktopScript)
	if err := cfg.UnmarshalKey("desktop_scripts", &desktopScripts); err != nil {
		// Section doesn't exist or is empty, that's ok
		return nil, nil
	}

	tmpDir := filepath.Join(home, ".config", "workspaced", "generated", "desktop")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return nil, err
	}

	iconsDir := filepath.Join(home, ".local", "share", "icons", "workspaced")
	if err := os.MkdirAll(iconsDir, 0755); err != nil {
		return nil, err
	}

	desired := []DesiredState{}

	for scriptID, script := range desktopScripts {
		iconPath := script.Icon

		// If icon is a URL, download it
		if strings.HasPrefix(script.Icon, "http://") || strings.HasPrefix(script.Icon, "https://") {
			downloaded, err := downloadIcon(ctx, script.Icon, iconsDir)
			if err != nil {
				// Log error but continue with original URL as fallback
				fmt.Fprintf(os.Stderr, "warning: failed to download icon for %s: %v\n", scriptID, err)
			} else {
				iconPath = downloaded
			}
		}

		desktopContent := fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=%s
Comment=%s
Icon=%s
Exec=%s
Terminal=%t
Categories=%s
`,
			script.Name,
			script.Description,
			iconPath,
			script.Exec,
			script.Terminal,
			script.Categories,
		)

		tmpFile := filepath.Join(tmpDir, scriptID+".desktop")
		if err := os.WriteFile(tmpFile, []byte(desktopContent), 0644); err != nil {
			return nil, err
		}

		targetPath := filepath.Join(home, ".local", "share", "applications", scriptID+".desktop")
		desired = append(desired, DesiredState{
			Target: targetPath,
			Source: tmpFile,
			Mode:   0644,
		})
	}

	return desired, nil
}
