package wm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"sync"
	"workspaced/pkg/driver"
	"workspaced/pkg/driver/media"
	"workspaced/pkg/logging"
)

var wmMu sync.Mutex

// switchToWorkspace é a implementação interna sem lock para evitar deadlock
func switchToWorkspace(ctx context.Context, ws string, move bool) error {
	d, err := driver.Get[Driver](ctx)
	if err != nil {
		return err
	}
	return d.SwitchToWorkspace(ctx, ws, move)
}

// SwitchToWorkspace switches to the specified workspace.
func SwitchToWorkspace(ctx context.Context, ws string, move bool) error {
	wmMu.Lock()
	defer wmMu.Unlock()
	return switchToWorkspace(ctx, ws, move)
}

// ToggleScratchpad toggles the visibility of the scratchpad container.
func ToggleScratchpad(ctx context.Context) error {
	wmMu.Lock()
	defer wmMu.Unlock()
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
	wmMu.Lock()
	defer wmMu.Unlock()
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

	return switchToWorkspace(ctx, nextWS, move)
}

// RotateWorkspaces rotates the visible workspaces across all connected outputs.
func RotateWorkspaces(ctx context.Context) error {
	wmMu.Lock()
	defer wmMu.Unlock()

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

	if len(screens) <= 1 {
		return nil // Nothing to rotate
	}

	oldScreens := make([]string, len(screens))
	copy(oldScreens, screens)

	// Rotate screens list (A, B) -> (B, A)
	last := screens[len(screens)-1]
	screens = append([]string{last}, screens[:len(screens)-1]...)

	for i, fromScreen := range oldScreens {
		toScreen := screens[i]
		ws := workspaceScreens[fromScreen]

		// Move o workspace para o novo monitor.
		if err := d.MoveWorkspaceToOutput(ctx, ws, toScreen); err != nil {
			logging.ReportError(ctx, err)
		}
	}

	// Restaura o foco original
	if focusedWorkspace != "" {
		if err := switchToWorkspace(ctx, focusedWorkspace, false); err != nil {
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
