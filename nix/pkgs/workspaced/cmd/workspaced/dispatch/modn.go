package dispatch

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
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

func getRPC(env []string) string {
	for _, e := range env {
		if strings.HasPrefix(e, "WAYLAND_DISPLAY=") {
			return "swaymsg"
		}
	}
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		return "swaymsg"
	}
	return "i3-msg"
}

func init() {
	Command.AddCommand(modnCmd)
}

var modnCmd = &cobra.Command{
	Use:   "modn",
	Short: "Rotate workspaces across outputs",
	RunE: func(c *cobra.Command, args []string) error {
		ctx := c.Context()
		var env []string
		if ctx != nil {
			env, _ = ctx.Value("env").([]string)
		}

		rpc := getRPC(env)

		// Get Workspaces
		cmd := runCmd(c, rpc, "-t", "get_workspaces")
		out, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("failed to get workspaces: %w", err)
		}
		var workspaces []Workspace
		if err := json.Unmarshal(out, &workspaces); err != nil {
			return fmt.Errorf("failed to parse workspaces: %w", err)
		}

		var focusedWorkspace string
		for _, w := range workspaces {
			if w.Focused {
				focusedWorkspace = w.Name
				break
			}
		}

		// Get Outputs
		cmd = runCmd(c, rpc, "-t", "get_outputs")
		out, err = cmd.Output()
		if err != nil {
			return fmt.Errorf("failed to get outputs: %w", err)
		}
		var outputs []Output
		if err := json.Unmarshal(out, &outputs); err != nil {
			return fmt.Errorf("failed to parse outputs: %w", err)
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

		// Rotate screens: insert last at beginning (0)
		last := screens[len(screens)-1]
		screens = append([]string{last}, screens[:len(screens)-1]...)

		// Perform moves
		for i, fromScreen := range oldScreens {
			toScreen := screens[i]
			ws := workspaceScreens[fromScreen]

			// i3/sway logic: focus workspace, then move it to output
			runCmd(c, rpc, "workspace", "number", ws).Run()
			time.Sleep(100 * time.Millisecond)
			runCmd(c, rpc, "move", "workspace", "to", "output", toScreen).Run()
			time.Sleep(100 * time.Millisecond)
		}

		// Refocus workspaces to clean up
		for _, ws := range workspaceScreens {
			runCmd(c, rpc, "workspace", "number", ws).Run()
			time.Sleep(100 * time.Millisecond)
		}

		if focusedWorkspace != "" {
			runCmd(c, rpc, "workspace", "number", focusedWorkspace).Run()
		}

		fmt.Println("Rotated workspaces")
		return nil
	},
}
