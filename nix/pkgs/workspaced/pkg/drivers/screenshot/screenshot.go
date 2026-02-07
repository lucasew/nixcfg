package screenshot

import (
	"context"
	"fmt"
	"image"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
	"time"

	"workspaced/pkg/config"
	"workspaced/pkg/drivers/clipboard"
	"workspaced/pkg/drivers/notification"
	"workspaced/pkg/wm"
	"workspaced/pkg/exec"
	"workspaced/pkg/logging"
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
	if rpc == "swaymsg" || rpc == "hyprctl" {
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
			return "", fmt.Errorf("screenshot: slurp binary not found in PATH")
		}
		out, err := exec.RunCmd(ctx, "slurp").CombinedOutput()
		if err != nil {
			// Se o usuário cancelar a seleção, o slurp sai com status 1.
			// Não queremos tratar isso como um erro catastrófico.
			if len(out) == 0 {
				return "", nil
			}
			return "", fmt.Errorf("slurp failed: %w (output: %q)", err, string(out))
		}
		args = append(args, "-g", strings.TrimSpace(string(out)))
	}

	args = append(args, path)

	cmdGrim := exec.RunCmd(ctx, "grim", args...)
	if out, err := cmdGrim.CombinedOutput(); err != nil {
		return "", fmt.Errorf("grim failed to capture: %w (output: %q, command: grim %s)", err, string(out), strings.Join(args, " "))
	}

	copyFileToClipboard(ctx, path)

	notifySaved(ctx, path, target)
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
		out, err := exec.RunCmd(ctx, "maim", "-s").CombinedOutput()
		if err != nil {
			if len(out) == 0 {
				return "", nil // Provável cancelamento
			}
			return "", fmt.Errorf("maim selection failed: %w (output: %q)", err, string(out))
		}
		args = append(args, "-g", strings.TrimSpace(string(out)))
	}

	args = append(args, path)

	cmdMaim := exec.RunCmd(ctx, "maim", args...)
	if out, err := cmdMaim.CombinedOutput(); err != nil {
		return "", fmt.Errorf("maim failed to capture: %w (output: %q, command: maim %s)", err, string(out), strings.Join(args, " "))
	}

	copyFileToClipboard(ctx, path)

	notifySaved(ctx, path, target)
	return path, nil
}

func copyFileToClipboard(ctx context.Context, path string) {
	file, err := os.Open(path)
	if err != nil {
		logging.ReportError(ctx, fmt.Errorf("failed to open screenshot for clipboard: %w", err))
		return
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		logging.ReportError(ctx, fmt.Errorf("failed to decode screenshot for clipboard: %w", err))
		return
	}

	if err := clipboard.WriteImage(ctx, img); err != nil {
		logging.ReportError(ctx, err)
	}
}

func notifyMissing(ctx context.Context, tool string) {
	n := &notification.Notification{
		Title:   "Screenshot Error",
		Message: fmt.Sprintf("Missing tool: %s", tool),
		Urgency: "critical",
	}
	if err := n.Notify(ctx); err != nil {
		logging.ReportError(ctx, err)
	}
}

func notifySaved(ctx context.Context, path string, target Target) {
	strategy := "Unknown"
	switch target {
	case All:
		strategy = "All screens"
	case Output:
		strategy = "Current monitor"
	case Window:
		strategy = "Current window"
	case Selection:
		strategy = "Selected area"
	}

	n := &notification.Notification{
		Title:   fmt.Sprintf("Screenshot Saved (%s)", strategy),
		Message: path,
		Icon:    "camera-photo",
	}
	if err := n.Notify(ctx); err != nil {
		logging.ReportError(ctx, err)
	}
}
