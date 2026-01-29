package common

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

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

func GetIconPath(ctx context.Context, url string) (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	iconsDir := filepath.Join(configDir, "webapp/icons")
	if err := os.MkdirAll(iconsDir, 0755); err != nil {
		return "", err
	}

	normalized := NormalizeURL(url)
	hash := sha256.Sum256([]byte(normalized))
	hashStr := hex.EncodeToString(hash[:])
	path := filepath.Join(iconsDir, hashStr+".png")

	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	logger := GetLogger(ctx)
	logger.Info("downloading favicon", "url", normalized, "target", path)

	domain := normalized
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
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return "", err
	}

	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if err := png.Encode(out, img); err != nil {
		return "", err
	}

	return path, nil
}
