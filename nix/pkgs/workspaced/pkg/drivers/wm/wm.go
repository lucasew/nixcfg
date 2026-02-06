package wm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"workspaced/pkg/drivers/media"
	"workspaced/pkg/drivers/wm/api"
	"workspaced/pkg/drivers/wm/hyprland"
	"workspaced/pkg/drivers/wm/i3ipc"
	"workspaced/pkg/exec"
	"workspaced/pkg/logging"
)

// Re-export types for backward compatibility within the wm package if needed,
// but external packages should ideally use wm.Rect etc.
// We alias them here to keep the wm API surface the same.
type Rect = api.Rect
type Workspace = api.Workspace
type Output = api.Output
type Node = api.Node
type Driver = api.Driver

// GetDriver returns the appropriate WM driver for the current environment.
func GetDriver(ctx context.Context) (api.Driver, error) {
	rpc := exec.GetRPC(ctx)
	switch rpc {
	case "hyprctl":
		return &hyprland.Driver{}, nil
	case "swaymsg":
		return &i3ipc.Driver{Binary: "swaymsg"}, nil
	case "i3-msg":
		return &i3ipc.Driver{Binary: "i3-msg"}, nil
	}
	return nil, fmt.Errorf("no suitable WM driver found for RPC: %s", rpc)
}

// SwitchToWorkspace switches to the specified workspace number.
func SwitchToWorkspace(ctx context.Context, num int, move bool) error {
	d, err := GetDriver(ctx)
	if err != nil {
		return err
	}
	return d.SwitchToWorkspace(ctx, num, move)
}

// ToggleScratchpad toggles the visibility of the scratchpad container.
func ToggleScratchpad(ctx context.Context) error {
	d, err := GetDriver(ctx)
	if err != nil {
		return err
	}
	return d.ToggleScratchpad(ctx)
}

// ToggleScratchpadWithInfo toggles the scratchpad and shows a media status notification.
func ToggleScratchpadWithInfo(ctx context.Context) error {
	if err := ToggleScratchpad(ctx); err != nil {
		return err
	}
	if err := media.ShowStatus(ctx); err != nil {
		logging.ReportError(ctx, err)
	}
	return nil
}

// NextWorkspace switches to (or moves the container to) the next available workspace.
func NextWorkspace(ctx context.Context, move bool) error {
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		runtimeDir = filepath.Join(os.TempDir(), fmt.Sprintf("workspaced-%d", os.Getuid()))
	}
	workspacedDir := filepath.Join(runtimeDir, "workspaced")
	if err := os.MkdirAll(workspacedDir, 0700); err != nil {
		logging.ReportError(ctx, err)
	}

	wsFile := filepath.Join(workspacedDir, "last_ws")
	lastWS := 10
	if data, err := os.ReadFile(wsFile); err == nil {
		if val, err := strconv.Atoi(strings.TrimSpace(string(data))); err == nil {
			lastWS = val
		}
	}

	nextWS := lastWS + 1
	if err := os.WriteFile(wsFile, []byte(strconv.Itoa(nextWS)), 0600); err != nil {
		logging.ReportError(ctx, err)
	}

	return SwitchToWorkspace(ctx, nextWS, move)
}

// RotateWorkspaces rotates the visible workspaces across all connected outputs.
func RotateWorkspaces(ctx context.Context) error {
	d, err := GetDriver(ctx)
	if err != nil {
		return err
	}

	workspaces, err := d.GetWorkspaces(ctx)
	if err != nil {
		return err
	}

	var focusedWorkspace string
	for _, w := range workspaces {
		if w.Focused {
			focusedWorkspace = w.Name
			break
		}
	}

	outputs, err := d.GetOutputs(ctx)
	if err != nil {
		return err
	}

	var screens []string
	workspaceScreens := make(map[string]string)

	for _, o := range outputs {
		if o.CurrentWorkspace != "" {
			screens = append(screens, o.Name)
			workspaceScreens[o.Name] = o.CurrentWorkspace
		}
	}

	if len(screens) == 0 {
		return fmt.Errorf("no screens found")
	}

	oldScreens := make([]string, len(screens))
	copy(oldScreens, screens)

	// Rotate screens
	last := screens[len(screens)-1]
	screens = append([]string{last}, screens[:len(screens)-1]...)

	for i, fromScreen := range oldScreens {
		toScreen := screens[i]
		ws := workspaceScreens[fromScreen]

		if err := d.SwitchToWorkspace(ctx, parseWS(ws), false); err != nil {
			logging.ReportError(ctx, err)
		}
		time.Sleep(100 * time.Millisecond)

		rpc := exec.GetRPC(ctx)
		if rpc == "hyprctl" {
			if err := exec.RunCmd(ctx, rpc, "dispatch", "moveworkspacetomonitor", ws, toScreen).Run(); err != nil {
				logging.ReportError(ctx, err)
			}
		} else {
			if err := exec.RunCmd(ctx, rpc, "move", "workspace", "to", "output", toScreen).Run(); err != nil {
				logging.ReportError(ctx, err)
			}
		}
		time.Sleep(100 * time.Millisecond)
	}

	for _, ws := range workspaceScreens {
		if err := d.SwitchToWorkspace(ctx, parseWS(ws), false); err != nil {
			logging.ReportError(ctx, err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	if focusedWorkspace != "" {
		if err := d.SwitchToWorkspace(ctx, parseWS(focusedWorkspace), false); err != nil {
			logging.ReportError(ctx, err)
		}
	}

	return nil
}

func parseWS(ws string) int {
	val, _ := strconv.Atoi(ws)
	return val
}

// GetFocusedOutput returns the name and geometry of the currently focused output.
func GetFocusedOutput(ctx context.Context) (string, *api.Rect, error) {
	d, err := GetDriver(ctx)
	if err != nil {
		return "", nil, err
	}
	return d.GetFocusedOutput(ctx)
}

// GetFocusedWindowRect returns the geometry of the currently focused window.
func GetFocusedWindowRect(ctx context.Context) (*api.Rect, error) {
	d, err := GetDriver(ctx)
	if err != nil {
		return nil, err
	}
	return d.GetFocusedWindowRect(ctx)
}
