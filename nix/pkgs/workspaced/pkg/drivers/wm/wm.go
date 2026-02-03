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

	"workspaced/pkg/common"
	"workspaced/pkg/drivers/media"
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
// It uses common.GetRPC to determine whether to use swaymsg or i3-msg.
// If move is true, it moves the current container to that workspace instead of switching focus.
func SwitchToWorkspace(ctx context.Context, num int, move bool) error {
	rpc := common.GetRPC(ctx)
	if move {
		return common.RunCmd(ctx, rpc, "move", "container", "to", "workspace", "number", strconv.Itoa(num)).Run()
	}
	return common.RunCmd(ctx, rpc, "workspace", "number", strconv.Itoa(num)).Run()
}

// ToggleScratchpad toggles the visibility of the scratchpad container.
func ToggleScratchpad(ctx context.Context) error {
	rpc := common.GetRPC(ctx)
	return common.RunCmd(ctx, rpc, "scratchpad", "show").Run()
}

// ToggleScratchpadWithInfo toggles the scratchpad and shows a media status notification.
// This is useful for visual feedback when toggling the scratchpad.
func ToggleScratchpadWithInfo(ctx context.Context) error {
	if err := ToggleScratchpad(ctx); err != nil {
		return err
	}
	_ = media.ShowStatus(ctx)
	return nil
}

// NextWorkspace switches to (or moves the container to) the next available workspace.
//
// Mechanism:
// It maintains a monotonic counter in `$XDG_RUNTIME_DIR/workspaced/last_ws`.
// This persistence ensures that even if the daemon restarts, it doesn't recycle
// workspace numbers that might still be populated.
//
// The counter starts at 10 to leave single-digit workspaces (1-9) free for
// static/manual assignment (e.g., '1:www', '2:code').
func NextWorkspace(ctx context.Context, move bool) error {
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		runtimeDir = filepath.Join(os.TempDir(), fmt.Sprintf("workspaced-%d", os.Getuid()))
	}
	workspacedDir := filepath.Join(runtimeDir, "workspaced")
	_ = os.MkdirAll(workspacedDir, 0700)

	wsFile := filepath.Join(workspacedDir, "last_ws")
	lastWS := 10
	if data, err := os.ReadFile(wsFile); err == nil {
		if val, err := strconv.Atoi(strings.TrimSpace(string(data))); err == nil {
			lastWS = val
		}
	}

	nextWS := lastWS + 1
	_ = os.WriteFile(wsFile, []byte(strconv.Itoa(nextWS)), 0600)

	return SwitchToWorkspace(ctx, nextWS, move)
}

// RotateWorkspaces rotates the visible workspaces across all connected outputs.
// It effectively shifts the workspace on output A to output B, B to C, etc.
//
// The logic implements a physical "carousel" of workspaces:
//  1. Snapshots the current state (which workspace is on which output).
//  2. Calculates the rotation (Screen[i] gets Workspace from Screen[i-1]).
//  3. Sequentially moves workspaces.
//
// Note: Artificial delays (`time.Sleep`) are injected between IPC commands.
// This is necessary because Sway/i3 IPC is asynchronous and rapid commands
// can race, causing workspaces to land on the wrong output or focus to be lost.
func RotateWorkspaces(ctx context.Context) error {
	rpc := common.GetRPC(ctx)

	// Get Workspaces
	out, err := common.RunCmd(ctx, rpc, "-t", "get_workspaces").Output()
	if err != nil {
		return err
	}
	var workspaces []Workspace
	_ = json.Unmarshal(out, &workspaces)

	var focusedWorkspace string
	for _, w := range workspaces {
		if w.Focused {
			focusedWorkspace = w.Name
			break
		}
	}

	// Get Outputs
	out, err = common.RunCmd(ctx, rpc, "-t", "get_outputs").Output()
	if err != nil {
		return err
	}
	var outputs []Output
	_ = json.Unmarshal(out, &outputs)

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

		_ = common.RunCmd(ctx, rpc, "workspace", "number", ws).Run()
		time.Sleep(100 * time.Millisecond)
		_ = common.RunCmd(ctx, rpc, "move", "workspace", "to", "output", toScreen).Run()
		time.Sleep(100 * time.Millisecond)
	}

	for _, ws := range workspaceScreens {
		_ = common.RunCmd(ctx, rpc, "workspace", "number", ws).Run()
		time.Sleep(100 * time.Millisecond)
	}

	if focusedWorkspace != "" {
		_ = common.RunCmd(ctx, rpc, "workspace", "number", focusedWorkspace).Run()
	}

	return nil
}
