package wm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"workspaced/pkg/driver"
	"workspaced/pkg/driver/media"
	"workspaced/pkg/logging"
)

// SwitchToWorkspace switches to the specified workspace.
func SwitchToWorkspace(ctx context.Context, ws string, move bool) error {
	d, err := driver.Get[Driver](ctx)
	if err != nil {
		return err
	}
	return d.SwitchToWorkspace(ctx, ws, move)
}

// ToggleScratchpad toggles the visibility of the scratchpad container.
func ToggleScratchpad(ctx context.Context) error {
	d, err := driver.Get[Driver](ctx)
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

	nextWS := strconv.Itoa(lastWS + 1)
	if err := os.WriteFile(wsFile, []byte(nextWS), 0600); err != nil {
		logging.ReportError(ctx, err)
	}

	return SwitchToWorkspace(ctx, nextWS, move)
}

// RotateWorkspaces rotates the visible workspaces across all connected outputs.
func RotateWorkspaces(ctx context.Context) error {
	d, err := driver.Get[Driver](ctx)
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

		if err := d.SwitchToWorkspace(ctx, ws, false); err != nil {
			logging.ReportError(ctx, err)
		}
		time.Sleep(100 * time.Millisecond)

		if err := d.MoveWorkspaceToOutput(ctx, ws, toScreen); err != nil {
			logging.ReportError(ctx, err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	for _, ws := range workspaceScreens {
		if err := d.SwitchToWorkspace(ctx, ws, false); err != nil {
			logging.ReportError(ctx, err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	if focusedWorkspace != "" {
		if err := d.SwitchToWorkspace(ctx, focusedWorkspace, false); err != nil {
			logging.ReportError(ctx, err)
		}
	}

	return nil
}

// GetFocusedOutput returns the name and geometry of the currently focused output.
func GetFocusedOutput(ctx context.Context) (string, *Rect, error) {
	d, err := driver.Get[Driver](ctx)
	if err != nil {
		return "", nil, err
	}
	return d.GetFocusedOutput(ctx)
}

// GetFocusedWindowRect returns the geometry of the currently focused window.
func GetFocusedWindowRect(ctx context.Context) (*Rect, error) {
	d, err := driver.Get[Driver](ctx)
	if err != nil {
		return nil, err
	}
	return d.GetFocusedWindowRect(ctx)
}
