package wallpaper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"workspaced/pkg/api"
	"workspaced/pkg/config"
	"workspaced/pkg/driver"
	"workspaced/pkg/exec"
	"workspaced/pkg/logging"
)

func SetStatic(ctx context.Context, path string) error {
	logger := logging.GetLogger(ctx)
	if path == "" {
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}
		wallpaperDir := cfg.Desktop.Wallpaper.Dir
		files, err := filepath.Glob(filepath.Join(wallpaperDir, "*"))
		if err != nil {
			return fmt.Errorf("error when listing wallpaper candidates: %w", err)
		}
		if len(files) == 0 {
			return fmt.Errorf("%w: wallpapers in %s", api.ErrNoTargetFound, wallpaperDir)
		}
		path = files[rand.Intn(len(files))]
	}

	logger.Info("setting wallpaper", "path", path)

	// Stop existing wallpaper-change service if it exists
	stopCmd := exec.RunCmd(ctx, "systemctl", "--user", "stop", "wallpaper-change.service")
	_ = stopCmd.Run() // Ignore errors if service doesn't exist

	d, err := driver.Get[Driver](ctx)
	if err != nil {
		return err
	}
	return d.SetStatic(ctx, path)
}

func SetAnimated(ctx context.Context, path string) error {
	// Simple wrapper for the video logic
	// Since it uses xrandr and xwinwrap, it's mostly X11
	// We'll just run it via exec.RunCmd
	return exec.RunCmd(ctx, "sh", "-c", fmt.Sprintf("sd wall video %s", path)).Run()
}

type APODResponse struct {
	HDURL string `json:"hdurl"`
	URL   string `json:"url"`
}

func SetAPOD(ctx context.Context) error {
	logger := logging.GetLogger(ctx)
	apiKey := os.Getenv("NASA_API_KEY")
	if apiKey == "" {
		apiKey = "DEMO_KEY"
	}

	logger.Info("fetching NASA Astronomy Picture of the Day")
	resp, err := http.Get(fmt.Sprintf("https://api.nasa.gov/planetary/apod?api_key=%s", apiKey))
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logging.ReportError(ctx, err)
		}
	}()

	var apod APODResponse
	if err := json.NewDecoder(resp.Body).Decode(&apod); err != nil {
		return err
	}

	url := apod.HDURL
	if url == "" {
		url = apod.URL
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	cacheDir := filepath.Join(home, ".cache/workspaced")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return err
	}
	outPath := filepath.Join(cacheDir, "apod.jpg")

	out, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := out.Close(); err != nil {
			logging.ReportError(ctx, err)
		}
	}()

	imgResp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		if err := imgResp.Body.Close(); err != nil {
			logging.ReportError(ctx, err)
		}
	}()

	if _, err := io.Copy(out, imgResp.Body); err != nil {
		return err
	}

	return SetStatic(ctx, outPath)
}
