package screenshot

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"workspaced/pkg/config"
	"workspaced/pkg/drivers/notification"
	"workspaced/pkg/drivers/wm"
	"workspaced/pkg/exec"
)

type Target int

const (
	All Target = iota
	Output
	Window
	Selection
)

func Capture(ctx context.Context, target Target) (string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}

	dir := cfg.Screenshot.Dir
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create screenshot dir: %w", err)
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("Screenshot_%s.png", timestamp)
	path := filepath.Join(dir, filename)

	rpc := exec.GetRPC(ctx)
	if rpc == "swaymsg" {
		return captureWayland(ctx, path, target)
	}
	return captureX11(ctx, path, target)
}

func captureWayland(ctx context.Context, path string, target Target) (string, error) {
	if !exec.IsBinaryAvailable(ctx, "grim") {
		notifyMissing(ctx, "grim")
		return "", fmt.Errorf("grim not found")
	}

	args := []string{}
	switch target {
	case All:
		// Default behavior of grim is to capture all outputs
	case Output:
		name, _, err := wm.GetFocusedOutput(ctx)
		if err != nil {
			return "", fmt.Errorf("failed to get focused output: %w", err)
		}
		args = append(args, "-o", name)
	case Window:
		rect, err := wm.GetFocusedWindowRect(ctx)
		if err != nil {
			return "", fmt.Errorf("failed to get focused window rect: %w", err)
		}
		args = append(args, "-g", fmt.Sprintf("%d,%d %dx%d", rect.X, rect.Y, rect.Width, rect.Height))
	case Selection:
		if !exec.IsBinaryAvailable(ctx, "slurp") {
			notifyMissing(ctx, "slurp")
			return "", fmt.Errorf("slurp not found")
		}
		out, err := exec.RunCmd(ctx, "slurp").Output()
		if err != nil {
			return "", err
		}
		args = append(args, "-g", strings.TrimSpace(string(out)))
	}

	args = append(args, path)

	if err := exec.RunCmd(ctx, "grim", args...).Run(); err != nil {
		return "", err
	}

	// Copy to clipboard
	if exec.IsBinaryAvailable(ctx, "wl-copy") {
		_ = exec.RunCmd(ctx, "sh", "-c", fmt.Sprintf("wl-copy < %s", path)).Run()
	}

	notifySaved(ctx, path)
	return path, nil
}

func captureX11(ctx context.Context, path string, target Target) (string, error) {
	if !exec.IsBinaryAvailable(ctx, "maim") {
		notifyMissing(ctx, "maim")
		return "", fmt.Errorf("maim not found")
	}

	args := []string{}
	switch target {
	case All:
		// Default behavior of maim is to capture the root window (all screens combined)
	case Output:
		_, rect, err := wm.GetFocusedOutput(ctx)
		if err != nil {
			return "", fmt.Errorf("failed to get focused output: %w", err)
		}
		args = append(args, "-g", fmt.Sprintf("%dx%d+%d+%d", rect.Width, rect.Height, rect.X, rect.Y))
	case Window:
		rect, err := wm.GetFocusedWindowRect(ctx)
		if err != nil {
			return "", fmt.Errorf("failed to get focused window rect: %w", err)
		}
		args = append(args, "-g", fmt.Sprintf("%dx%d+%d+%d", rect.Width, rect.Height, rect.X, rect.Y))
	case Selection:
		args = append(args, "-s")
	}

	args = append(args, path)

	if err := exec.RunCmd(ctx, "maim", args...).Run(); err != nil {
		return "", err
	}

	// Copy to clipboard
	if exec.IsBinaryAvailable(ctx, "xclip") {
		_ = exec.RunCmd(ctx, "sh", "-c", fmt.Sprintf("xclip -selection clipboard -t image/png < %s", path)).Run()
	}

	notifySaved(ctx, path)
	return path, nil
}

func notifyMissing(ctx context.Context, tool string) {
	n := &notification.Notification{
		Title:   "Screenshot Error",
		Message: fmt.Sprintf("Missing tool: %s", tool),
		Urgency: "critical",
	}
	_ = n.Notify(ctx)
}

func notifySaved(ctx context.Context, path string) {
	n := &notification.Notification{
		Title:   "Screenshot Saved",
		Message: path,
		Icon:    "camera-photo",
	}
	_ = n.Notify(ctx)
}
