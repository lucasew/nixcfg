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
	"workspaced/pkg/drivers/audio"
	"workspaced/pkg/drivers/brightness"
	"workspaced/pkg/drivers/media"
)

type Workspace struct {
	Name    string `json:"name"`
	Focused bool   `json:"focused"`
	Output  string `json:"output"`
}

type Output struct {
	Name             string `json:"name"`
	CurrentWorkspace string `json:"current_workspace"`
}

func SwitchToWorkspace(ctx context.Context, num int, move bool) error {
	rpc := common.GetRPC(ctx)
	if move {
		return common.RunCmd(ctx, rpc, "move", "container", "to", "workspace", "number", strconv.Itoa(num)).Run()
	}
	return common.RunCmd(ctx, rpc, "workspace", "number", strconv.Itoa(num)).Run()
}

func ToggleScratchpad(ctx context.Context) error {
	rpc := common.GetRPC(ctx)
	return common.RunCmd(ctx, rpc, "scratchpad", "show").Run()
}

func ToggleScratchpadWithInfo(ctx context.Context) error {
	if err := ToggleScratchpad(ctx); err != nil {
		return err
	}
	audio.ShowStatus(ctx)
	brightness.ShowStatus(ctx)
	media.ShowStatus(ctx)
	return nil
}

func NextWorkspace(ctx context.Context, move bool) error {
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		runtimeDir = filepath.Join(os.TempDir(), fmt.Sprintf("workspaced-%d", os.Getuid()))
	}
	workspacedDir := filepath.Join(runtimeDir, "workspaced")
	os.MkdirAll(workspacedDir, 0700)

	wsFile := filepath.Join(workspacedDir, "last_ws")
	lastWS := 10
	if data, err := os.ReadFile(wsFile); err == nil {
		if val, err := strconv.Atoi(strings.TrimSpace(string(data))); err == nil {
			lastWS = val
		}
	}

	nextWS := lastWS + 1
	os.WriteFile(wsFile, []byte(strconv.Itoa(nextWS)), 0600)

	return SwitchToWorkspace(ctx, nextWS, move)
}

func RotateWorkspaces(ctx context.Context) error {
	rpc := common.GetRPC(ctx)

	// Get Workspaces
	out, err := common.RunCmd(ctx, rpc, "-t", "get_workspaces").Output()
	if err != nil {
		return err
	}
	var workspaces []Workspace
	json.Unmarshal(out, &workspaces)

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
	json.Unmarshal(out, &outputs)

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

		common.RunCmd(ctx, rpc, "workspace", "number", ws).Run()
		time.Sleep(100 * time.Millisecond)
		common.RunCmd(ctx, rpc, "move", "workspace", "to", "output", toScreen).Run()
		time.Sleep(100 * time.Millisecond)
	}

	for _, ws := range workspaceScreens {
		common.RunCmd(ctx, rpc, "workspace", "number", ws).Run()
		time.Sleep(100 * time.Millisecond)
	}

	if focusedWorkspace != "" {
		common.RunCmd(ctx, rpc, "workspace", "number", focusedWorkspace).Run()
	}

	return nil
}
