package screenshot

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"time"
	"workspaced/pkg/exec"
	"workspaced/pkg/host"
	"workspaced/pkg/drivers/notification"
	"workspaced/pkg/config"
)

func Capture(ctx context.Context, area bool) (string, error) {
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

	rpc := host.GetRPC(ctx)
	if rpc == "swaymsg" {
		return captureWayland(ctx, path, area)
	}
	return captureX11(ctx, path, area)
}

func captureWayland(ctx context.Context, path string, area bool) (string, error) {
	if !exec.IsBinaryAvailable(ctx, "grim") {
		notifyMissing(ctx, "grim")
		return "", fmt.Errorf("grim not found")
	}

	args := []string{}
	if area {
		if !exec.IsBinaryAvailable(ctx, "slurp") {
			notifyMissing(ctx, "slurp")
			return "", fmt.Errorf("slurp not found")
		}
		out, err := exec.RunCmd(ctx, "slurp").Output()
		if err != nil {
			return "", err
		}
		args = append(args, "-g", string(out))
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

func captureX11(ctx context.Context, path string, area bool) (string, error) {
	if !exec.IsBinaryAvailable(ctx, "maim") {
		notifyMissing(ctx, "maim")
		return "", fmt.Errorf("maim not found")
	}

	args := []string{}
	if area {
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
