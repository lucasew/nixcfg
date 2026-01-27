package workspace

import (
	"encoding/json"
	"fmt"
	"time"
	"workspaced/cmd/workspaced/dispatch/common"
	"workspaced/cmd/workspaced/dispatch/types"

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

var modnCmd = &cobra.Command{
	Use:   "modn",
	Short: "Rotate workspaces across outputs",
	RunE: func(c *cobra.Command, args []string) error {
		ctx := c.Context()
		var env []string
		if ctx != nil {
			env, _ = ctx.Value(types.EnvKey).([]string)
		}

		rpc := common.GetRPC(env)

		// Get Workspaces
		cmd := common.RunCmd(c, rpc, "-t", "get_workspaces")
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
		cmd = common.RunCmd(c, rpc, "-t", "get_outputs")
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
			common.RunCmd(c, rpc, "workspace", "number", ws).Run()
			time.Sleep(100 * time.Millisecond)
			common.RunCmd(c, rpc, "move", "workspace", "to", "output", toScreen).Run()
			time.Sleep(100 * time.Millisecond)
		}

		// Refocus workspaces to clean up
		for _, ws := range workspaceScreens {
			common.RunCmd(c, rpc, "workspace", "number", ws).Run()
			time.Sleep(100 * time.Millisecond)
		}

		if focusedWorkspace != "" {
			common.RunCmd(c, rpc, "workspace", "number", focusedWorkspace).Run()
		}

		fmt.Println("Rotated workspaces")
		return nil
	},
}

func init() {
	Command.AddCommand(modnCmd)
}
