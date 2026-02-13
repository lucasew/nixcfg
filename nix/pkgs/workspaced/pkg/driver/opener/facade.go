package opener

import (
	"context"
	"os"
	"path/filepath"
	"workspaced/pkg/config"
	"workspaced/pkg/driver"
	"workspaced/pkg/env"
	"workspaced/pkg/exec"
)

// WebappConfig is used for passing parameters to OpenWebapp
type WebappConfig struct {
	URL        string
	Profile    string
	ExtraFlags []string
}

// Open opens a generic target (file or URL) using the available opener driver.
func Open(ctx context.Context, target string) error {
	d, err := driver.Get[Driver](ctx)
	if err != nil {
		return err
	}
	return d.Open(ctx, target)
}

// OpenWebapp launches a URL as a webapp using the configured browser engine.
func OpenWebapp(ctx context.Context, wa WebappConfig) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	engine := cfg.Browser.Engine
	normalizedURL := env.NormalizeURL(wa.URL)
	args := []string{"--app=" + normalizedURL}

	if wa.Profile != "" {
		home, _ := os.UserHomeDir()
		profileDir := filepath.Join(home, ".config/workspaced/webapp/profiles", wa.Profile)
		args = append(args, "--user-data-dir="+profileDir)
	}

	if os.Getenv("WAYLAND_DISPLAY") != "" {
		args = append(args, "--enable-features=UseOzonePlatform", "--ozone-platform=wayland")
	}

	args = append(args, wa.ExtraFlags...)

	cmd := exec.RunCmd(ctx, engine, args...)
	exec.InheritContextWriters(ctx, cmd)
	return cmd.Run()
}
