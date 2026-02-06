package wm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"workspaced/pkg/drivers/media"
	"workspaced/pkg/exec"
	"workspaced/pkg/logging"
)

// Workspace represents a workspace as returned by Sway/i3 IPC.
type Workspace struct {
	Name    string `json:"name"`
	Focused bool   `json:"focused"`
	Output  string `json:"output"`
}

// Output represents a display output as returned by Sway/i3 IPC.
type Output struct {
	Name             string `json:"name"`
	CurrentWorkspace string `json:"current_workspace"`
}

// SwitchToWorkspace switches to the specified workspace number.
// It uses exec.GetRPC to determine whether to use swaymsg or i3-msg.
// If move is true, it moves the current container to that workspace instead of switching focus.
func SwitchToWorkspace(ctx context.Context, num int, move bool) error {
	rpc := exec.GetRPC(ctx)
	if move {
		return exec.RunCmd(ctx, rpc, "move", "container", "to", "workspace", "number", strconv.Itoa(num)).Run()
	}
	return exec.RunCmd(ctx, rpc, "workspace", "number", strconv.Itoa(num)).Run()
}

// ToggleScratchpad toggles the visibility of the scratchpad container.
func ToggleScratchpad(ctx context.Context) error {
	rpc := exec.GetRPC(ctx)
	return exec.RunCmd(ctx, rpc, "scratchpad", "show").Run()
}

// ToggleScratchpadWithInfo toggles the scratchpad and shows a media status notification.
// This is useful for visual feedback when toggling the scratchpad.
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
// It maintains a counter in XDG_RUNTIME_DIR/workspaced/last_ws to ensure unique,
// monotonically increasing workspace numbers until the system (or runtime dir) resets.
//
// The counter starts at 10 and increments with each call.
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
// It effectively shifts the workspace on output A to output B, B to C, etc.
// This creates a "carousel" effect for workspaces.
//
// The function:
// 1. Fetches current workspaces and outputs via IPC.
// 2. Maps outputs to their currently visible workspace.
// 3. Rotates the list of screens.
// 4. Moves workspaces to their new target screens.
// 5. Restores focus to the originally focused workspace.
func RotateWorkspaces(ctx context.Context) error {
	rpc := exec.GetRPC(ctx)

	// Get Workspaces
	out, err := exec.RunCmd(ctx, rpc, "-t", "get_workspaces").Output()
	if err != nil {
		return err
	}
	var workspaces []Workspace
	if err := json.Unmarshal(out, &workspaces); err != nil {
		logging.ReportError(ctx, err)
		return err
	}

	var focusedWorkspace string
	for _, w := range workspaces {
		if w.Focused {
			focusedWorkspace = w.Name
			break
		}
	}

	// Get Outputs
	out, err = exec.RunCmd(ctx, rpc, "-t", "get_outputs").Output()
	if err != nil {
		return err
	}
	var outputs []Output
	if err := json.Unmarshal(out, &outputs); err != nil {
		logging.ReportError(ctx, err)
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

		if err := exec.RunCmd(ctx, rpc, "workspace", "number", ws).Run(); err != nil {
			logging.ReportError(ctx, err)
		}
		time.Sleep(100 * time.Millisecond)
		if err := exec.RunCmd(ctx, rpc, "move", "workspace", "to", "output", toScreen).Run(); err != nil {
			logging.ReportError(ctx, err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	for _, ws := range workspaceScreens {
		if err := exec.RunCmd(ctx, rpc, "workspace", "number", ws).Run(); err != nil {
			logging.ReportError(ctx, err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	if focusedWorkspace != "" {
		if err := exec.RunCmd(ctx, rpc, "workspace", "number", focusedWorkspace).Run(); err != nil {
			logging.ReportError(ctx, err)
		}
	}

	return nil
}
