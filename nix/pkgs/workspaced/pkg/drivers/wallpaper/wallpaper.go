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
	"workspaced/pkg/common"
	"workspaced/pkg/config"
)

func SetStatic(ctx context.Context, path string) error {
	logger := common.GetLogger(ctx)
	if path == "" {
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}
		wallpaperDir := cfg.Desktop.Wallpaper.Dir
		files, _ := filepath.Glob(filepath.Join(wallpaperDir, "*"))
		if len(files) == 0 {
			return fmt.Errorf("no wallpapers found in %s", wallpaperDir)
		}
		path = files[rand.Intn(len(files))]
	}

	logger.Info("setting wallpaper", "path", path)

	rpc := common.GetRPC(ctx)
	if rpc == "swaymsg" {
		return common.RunCmd(ctx, "systemd-run", "--user", "-u", "wallpaper-change", "--collect", "swaybg", "-i", path).Run()
	}
	return common.RunCmd(ctx, "systemd-run", "--user", "-u", "wallpaper-change", "--collect", "feh", "--bg-fill", path).Run()
}

func SetAnimated(ctx context.Context, path string) error {
	// Simple wrapper for the video logic
	// Since it uses xrandr and xwinwrap, it's mostly X11
	// We'll just run it via common.RunCmd
	return common.RunCmd(ctx, "sh", "-c", fmt.Sprintf("sd wall video %s", path)).Run()
}

type APODResponse struct {
	HDURL string `json:"hdurl"`
	URL   string `json:"url"`
}

func SetAPOD(ctx context.Context) error {
	logger := common.GetLogger(ctx)
	apiKey := os.Getenv("NASA_API_KEY")
	if apiKey == "" {
		apiKey = "DEMO_KEY"
	}

	logger.Info("fetching NASA Astronomy Picture of the Day")
	resp, err := http.Get(fmt.Sprintf("https://api.nasa.gov/planetary/apod?api_key=%s", apiKey))
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	var apod APODResponse
	if err := json.NewDecoder(resp.Body).Decode(&apod); err != nil {
		return err
	}

	url := apod.HDURL
	if url == "" {
		url = apod.URL
	}

	home, _ := os.UserHomeDir()
	cacheDir := filepath.Join(home, ".cache/workspaced")
	_ = os.MkdirAll(cacheDir, 0755)
	outPath := filepath.Join(cacheDir, "apod.jpg")

	out, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	imgResp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() { _ = imgResp.Body.Close() }()

	if _, err := io.Copy(out, imgResp.Body); err != nil {
		return err
	}

	return SetStatic(ctx, outPath)
}
