package apply

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"workspaced/pkg/common"
	"workspaced/pkg/config"
)

type WebappProvider struct{}

func (p *WebappProvider) Name() string {
	return "webapp"
}

func (p *WebappProvider) GetDesiredState(ctx context.Context) ([]DesiredState, error) {
	cfg, err := config.Load()
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

	// 0. Add generic launcher
	borderlessPath := filepath.Join(shortcutsDir, "borderless-browser.desktop")
	borderlessContent := `[Desktop Entry]
Type=Application
Name=Borderless Browser
Icon=applications-internet
Exec=workspaced dispatch webapp launch
Categories=Network;WebBrowser;
`
	if err := os.WriteFile(borderlessPath, []byte(borderlessContent), 0644); err != nil {
		return nil, err
	}
	desired = append(desired, DesiredState{
		Target: filepath.Join(home, ".local/share/applications", "workspaced-webapp-borderless-browser.desktop"),
		Source: borderlessPath,
	})

	for name, wa := range cfg.Webapps {
		// 1. Manage Icon
		var iconToUse string
		if wa.Icon != "" {
			iconToUse = wa.Icon
		} else {
			var err error
			iconToUse, err = common.GetIconPath(ctx, wa.URL)
			if err != nil {
				logger.Error("failed to get icon path", "webapp", name, "url", wa.URL, "error", err)
				// Continue without icon as it's optional
			}
		}

		// 2. Generate .desktop
		desktopName := wa.DesktopName
		if desktopName == "" {
			desktopName = common.ToTitleCase(name)
		}

		desktopContent := generateWebappDesktopFile(wa, cfg.Browser.Engine, name, desktopName, iconToUse)
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

func generateWebappDesktopFile(wa config.WebappConfig, engine, id, name, iconPath string) string {
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
