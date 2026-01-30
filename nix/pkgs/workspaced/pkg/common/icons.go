package common

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"image/draw"
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
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return "", err
	}

	img = makeBackgroundTransparent(img)
	img = cropToContentSquare(img)

	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = out.Close() }()

	if err := png.Encode(out, img); err != nil {
		return "", err
	}

	return path, nil
}

func makeBackgroundTransparent(img image.Image) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Detect target color at (0,0)
	targetColor := img.At(bounds.Min.X, bounds.Min.Y)
	tr, tg, tb, ta := targetColor.RGBA()

	// If background is already transparent, do nothing
	if ta == 0 {
		return img
	}

	// Create mutable image
	dst := image.NewNRGBA(bounds)
	draw.Draw(dst, bounds, img, bounds.Min, draw.Src)

	type point struct{ x, y int }
	queue := make([]point, 0)
	visited := make([]bool, width*height)

	isTarget := func(x, y int) bool {
		r, g, b, a := dst.At(x, y).RGBA()
		return r == tr && g == tg && b == tb && a == ta
	}

	enqueue := func(x, y int) {
		idx := (y-bounds.Min.Y)*width + (x - bounds.Min.X)
		if !visited[idx] && isTarget(x, y) {
			visited[idx] = true
			queue = append(queue, point{x, y})
		}
	}

	// Seed from all edges
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		enqueue(x, bounds.Min.Y)
		enqueue(x, bounds.Max.Y-1)
	}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		enqueue(bounds.Min.X, y)
		enqueue(bounds.Max.X-1, y)
	}

	// Seed from center if it matches target color
	centerX, centerY := bounds.Min.X+width/2, bounds.Min.Y+height/2
	if isTarget(centerX, centerY) {
		enqueue(centerX, centerY)
	}

	// BFS Flood Fill
	head := 0
	for head < len(queue) {
		p := queue[head]
		head++

		// Set transparent
		dst.SetNRGBA(p.x, p.y, color.NRGBA{0, 0, 0, 0})

		// Check neighbors
		if p.x > bounds.Min.X {
			enqueue(p.x-1, p.y)
		}
		if p.x < bounds.Max.X-1 {
			enqueue(p.x+1, p.y)
		}
		if p.y > bounds.Min.Y {
			enqueue(p.x, p.y-1)
		}
		if p.y < bounds.Max.Y-1 {
			enqueue(p.x, p.y+1)
		}
	}

	return dst
}

func cropToContentSquare(img image.Image) image.Image {
	bounds := img.Bounds()
	minX, minY := bounds.Max.X, bounds.Max.Y
	maxX, maxY := bounds.Min.X, bounds.Min.Y
	found := false

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a > 0 {
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
				found = true
			}
		}
	}

	if !found {
		return img
	}

	w := maxX - minX + 1
	h := maxY - minY + 1
	size := w
	if h > w {
		size = h
	}

	dst := image.NewNRGBA(image.Rect(0, 0, size, size))
	// No need to fill with transparency as NewNRGBA is initialized with zeros (transparent)

	offsetX := (size - w) / 2
	offsetY := (size - h) / 2

	draw.Draw(dst, image.Rect(offsetX, offsetY, offsetX+w, offsetY+h), img, image.Point{minX, minY}, draw.Src)
	return dst
}
