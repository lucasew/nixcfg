package screenshot

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"workspaced/pkg/common"
)

func Capture(ctx context.Context, area bool) (string, error) {
	config, err := common.LoadConfig()
	if err != nil {
		return "", err
	}

	dir := config.Screenshot.Dir
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create screenshot dir: %w", err)
	}

	timestamp := time.Now().Format("2006-01-27_15-04-05")
	filename := fmt.Sprintf("Screenshot_%s.png", timestamp)
	path := filepath.Join(dir, filename)

	rpc := common.GetRPC(ctx)
	if rpc == "swaymsg" {
		return captureWayland(ctx, path, area)
	}
	return captureX11(ctx, path, area)
}

func captureWayland(ctx context.Context, path string, area bool) (string, error) {
	if _, err := exec.LookPath("grim"); err != nil {
		notifyMissing(ctx, "grim")
		return "", fmt.Errorf("grim not found")
	}

	args := []string{}
	if area {
		if _, err := exec.LookPath("slurp"); err != nil {
			notifyMissing(ctx, "slurp")
			return "", fmt.Errorf("slurp not found")
		}
		out, err := common.RunCmd(ctx, "slurp").Output()
		if err != nil {
			return "", err
		}
		args = append(args, "-g", string(out))
	}
	args = append(args, path)

	if err := common.RunCmd(ctx, "grim", args...).Run(); err != nil {
		return "", err
	}

	// Copy to clipboard
	if _, err := exec.LookPath("wl-copy"); err == nil {
		exec.CommandContext(ctx, "sh", "-c", fmt.Sprintf("wl-copy < %s", path)).Run()
	}

	notifySaved(ctx, path)
	return path, nil
}

func captureX11(ctx context.Context, path string, area bool) (string, error) {
	if _, err := exec.LookPath("maim"); err != nil {
		notifyMissing(ctx, "maim")
		return "", fmt.Errorf("maim not found")
	}

	args := []string{}
	if area {
		args = append(args, "-s")
	}
	args = append(args, path)

	if err := common.RunCmd(ctx, "maim", args...).Run(); err != nil {
		return "", err
	}

	// Copy to clipboard
	if _, err := exec.LookPath("xclip"); err == nil {
		exec.CommandContext(ctx, "sh", "-c", fmt.Sprintf("xclip -selection clipboard -t image/png < %s", path)).Run()
	}

	notifySaved(ctx, path)
	return path, nil
}

func notifyMissing(ctx context.Context, tool string) {
	common.RunCmd(ctx, "notify-send", "-u", "critical", "Screenshot Error", fmt.Sprintf("Missing tool: %s", tool)).Run()
}

func notifySaved(ctx context.Context, path string) {
	common.RunCmd(ctx, "notify-send", "Screenshot Saved", path, "-i", "camera-photo").Run()
}
