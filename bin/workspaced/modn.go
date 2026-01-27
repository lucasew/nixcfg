package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

var modnCmd = &cobra.Command{
	Use:   "modn",
	Short: "Rotate workspaces across outputs",
	Run: func(cmd *cobra.Command, args []string) {
		runOrRoute("modn", args, RunModn)
	},
}

type Workspace struct {
	Name    string `json:"name"`
	Focused bool   `json:"focused"`
	Output  string `json:"output"`
}

type Output struct {
	Name             string `json:"name"`
	CurrentWorkspace string `json:"current_workspace"`
}

func getRPC() string {
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		return "swaymsg"
	}
	return "i3-msg"
}

func RunModn() (string, error) {
	rpc := getRPC()

	// Get Workspaces
	out, err := exec.Command(rpc, "-t", "get_workspaces").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get workspaces: %w", err)
	}
	var workspaces []Workspace
	if err := json.Unmarshal(out, &workspaces); err != nil {
		return "", fmt.Errorf("failed to parse workspaces: %w", err)
	}

	var focusedWorkspace string
	for _, w := range workspaces {
		if w.Focused {
			focusedWorkspace = w.Name
			break
		}
	}

	// Get Outputs
	out, err = exec.Command(rpc, "-t", "get_outputs").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get outputs: %w", err)
	}
	var outputs []Output
	if err := json.Unmarshal(out, &outputs); err != nil {
		return "", fmt.Errorf("failed to parse outputs: %w", err)
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
		return "No screens found", nil
	}

	oldScreens := make([]string, len(screens))
	copy(oldScreens, screens)

	// Rotate screens: insert last at beginning (0)
	last := screens[len(screens)-1]
	screens = append([]string{last}, screens[:len(screens)-1]...)

	// Perform moves
	for i, fromScreen := range oldScreens {
		toScreen := screens[i]
		ws := workspaceScreens[fromScreen]

		// i3/sway logic: focus workspace, then move it to output
		exec.Command(rpc, "workspace", "number", ws).Run()
		time.Sleep(100 * time.Millisecond)
		exec.Command(rpc, "move", "workspace", "to", "output", toScreen).Run()
		time.Sleep(100 * time.Millisecond)
	}

	// Refocus workspaces to clean up
	for _, ws := range workspaceScreens {
		exec.Command(rpc, "workspace", "number", ws).Run()
		time.Sleep(100 * time.Millisecond)
	}

	if focusedWorkspace != "" {
		exec.Command(rpc, "workspace", "number", focusedWorkspace).Run()
	}

	return "Rotated workspaces", nil
}
