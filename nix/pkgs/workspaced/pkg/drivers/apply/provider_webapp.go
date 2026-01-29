package apply

import (
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"workspaced/pkg/common"
)

type WebappProvider struct{}

func (p *WebappProvider) Name() string {
	return "webapp"
}

func (p *WebappProvider) GetDesiredState(ctx context.Context) ([]DesiredState, error) {
	cfg, err := common.LoadConfig()
	if err != nil {
		return nil, err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	baseDir := filepath.Join(home, ".config/workspaced/webapp")
	shortcutsDir := filepath.Join(baseDir, "shortcuts")
	profilesDir := filepath.Join(baseDir, "profiles")

	for _, dir := range []string{shortcutsDir, profilesDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}

	desired := []DesiredState{}
	logger := common.GetLogger(ctx)

	for name, wa := range cfg.Webapps {
		normalizedURL := common.NormalizeURL(wa.URL)

		// 1. Manage Icon
		iconPath := filepath.Join(shortcutsDir, name+".png")
		if _, err := os.Stat(iconPath); os.IsNotExist(err) {
			logger.Info("downloading favicon", "webapp", name, "url", normalizedURL)
			if err := downloadAndEncodeFavicon(ctx, normalizedURL, iconPath); err != nil {
				logger.Error("failed to download favicon", "webapp", name, "error", err)
				// Continue without icon as it's optional
			}
		}

		// 2. Generate .desktop
		desktopName := wa.DesktopName
		if desktopName == "" {
			desktopName = common.ToTitleCase(name)
		}

		desktopContent := generateWebappDesktopFile(wa, cfg.Browser.Engine, name, desktopName, iconPath)
		desktopPath := filepath.Join(shortcutsDir, name+".desktop")
		if err := os.WriteFile(desktopPath, []byte(desktopContent), 0644); err != nil {
			return nil, err
		}

		// 3. Add to desired state
		target := filepath.Join(home, ".local/share/applications", "workspaced-webapp-"+name+".desktop")
		desired = append(desired, DesiredState{
			Target: target,
			Source: desktopPath,
		})
	}

	return desired, nil
}

func downloadAndEncodeFavicon(ctx context.Context, url, targetPath string) error {
	domain := url
	if strings.HasPrefix(domain, "https://") {
		domain = domain[8:]
	} else if strings.HasPrefix(domain, "http://") {
		domain = domain[7:]
	}
	parts := strings.Split(domain, "/")
	domain = parts[0]

	faviconURL := fmt.Sprintf("https://www.google.com/s2/favicons?sz=128&domain=%s", domain)

	req, err := http.NewRequestWithContext(ctx, "GET", faviconURL, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return err
	}

	out, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer out.Close()

	return png.Encode(out, img)
}

func generateWebappDesktopFile(wa common.WebappConfig, engine, id, name, iconPath string) string {
	url := common.NormalizeURL(wa.URL)
	args := []string{"--app=" + url}

	if wa.Profile != "" {
		home, _ := os.UserHomeDir()
		profileDir := filepath.Join(home, ".config/workspaced/webapp/profiles", wa.Profile)
		args = append(args, "--user-data-dir="+profileDir)
	}

	if os.Getenv("WAYLAND_DISPLAY") != "" {
		args = append(args, "--enable-features=UseOzonePlatform", "--ozone-platform=wayland")
	}

	args = append(args, wa.ExtraFlags...)

	return fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=%s
Icon=%s
Exec=%s %s
Categories=Network;WebBrowser;
StartupWMClass=%s
`, name, iconPath, engine, strings.Join(args, " "), id)
}
