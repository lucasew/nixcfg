package i3ipc

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	dapi "workspaced/pkg/api"
	"workspaced/pkg/driver"
	"workspaced/pkg/exec"
	"workspaced/pkg/wm/api"
)

func init() {
	driver.Register[api.Driver](&SwayProvider{})
	driver.Register[api.Driver](&I3Provider{})
}

type SwayProvider struct{}

func (p *SwayProvider) Name() string { return "Sway" }
func (p *SwayProvider) CheckCompatibility(ctx context.Context) error {
	if exec.GetEnv(ctx, "WAYLAND_DISPLAY") == "" {
		return fmt.Errorf("%w: WAYLAND_DISPLAY not set", driver.ErrIncompatible)
	}
	if !exec.IsBinaryAvailable(ctx, "swaymsg") {
		return fmt.Errorf("%w: swaymsg not found", driver.ErrIncompatible)
	}
	return nil
}

func (p *SwayProvider) New(ctx context.Context) (api.Driver, error) {
	return &Driver{Binary: "swaymsg"}, nil
}

type I3Provider struct{}

func (p *I3Provider) Name() string { return "i3" }
func (p *I3Provider) CheckCompatibility(ctx context.Context) error {
	if exec.GetEnv(ctx, "DISPLAY") == "" {
		return fmt.Errorf("%w: DISPLAY not set", driver.ErrIncompatible)
	}
	if !exec.IsBinaryAvailable(ctx, "i3-msg") {
		return fmt.Errorf("%w: i3-msg not found", driver.ErrIncompatible)
	}
	return nil
}

func (p *I3Provider) New(ctx context.Context) (api.Driver, error) {
	return &Driver{Binary: "i3-msg"}, nil
}

type Driver struct {
	Binary string
}

func (d *Driver) MoveWorkspaceToOutput(ctx context.Context, workspace string, output string) error {
	return exec.RunCmd(ctx, d.Binary, "workspace", workspace, "move", "workspace", "to", "output", output).Run()
}

func (d *Driver) SwitchToWorkspace(ctx context.Context, num int, move bool) error {
	if move {
		return exec.RunCmd(ctx, d.Binary, "move", "container", "to", "workspace", "number", strconv.Itoa(num)).Run()
	}
	return exec.RunCmd(ctx, d.Binary, "workspace", "number", strconv.Itoa(num)).Run()
}

func (d *Driver) ToggleScratchpad(ctx context.Context) error {
	return exec.RunCmd(ctx, d.Binary, "scratchpad", "show").Run()
}

func (d *Driver) GetOutputs(ctx context.Context) ([]api.Output, error) {
	out, err := exec.RunCmd(ctx, d.Binary, "-t", "get_outputs").Output()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", dapi.ErrIPC, err)
	}
	var outputs []api.Output
	if err := json.Unmarshal(out, &outputs); err != nil {
		return nil, fmt.Errorf("%w: %w", dapi.ErrIPC, err)
	}
	return outputs, nil
}

func (d *Driver) GetWorkspaces(ctx context.Context) ([]api.Workspace, error) {
	out, err := exec.RunCmd(ctx, d.Binary, "-t", "get_workspaces").Output()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", dapi.ErrIPC, err)
	}
	var workspaces []api.Workspace
	if err := json.Unmarshal(out, &workspaces); err != nil {
		return nil, fmt.Errorf("%w: %w", dapi.ErrIPC, err)
	}
	return workspaces, nil
}

func (d *Driver) GetFocusedOutput(ctx context.Context) (string, *api.Rect, error) {
	outputs, err := d.GetOutputs(ctx)
	if err != nil {
		return "", nil, err
	}

	for _, o := range outputs {
		if o.Focused {
			return o.Name, &o.Rect, nil
		}
	}

	workspaces, err := d.GetWorkspaces(ctx)
	if err != nil {
		return "", nil, err
	}

	var focusedOutputName string
	for _, w := range workspaces {
		if w.Focused {
			focusedOutputName = w.Output
			break
		}
	}

	if focusedOutputName != "" {
		for _, o := range outputs {
			if o.Name == focusedOutputName {
				return o.Name, &o.Rect, nil
			}
		}
	}

	return "", nil, dapi.ErrNoFocusedOutput
}

func (d *Driver) GetFocusedWindowRect(ctx context.Context) (*api.Rect, error) {
	out, err := exec.RunCmd(ctx, d.Binary, "-t", "get_tree").Output()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", dapi.ErrIPC, err)
	}

	var root api.Node
	if err := json.Unmarshal(out, &root); err != nil {
		return nil, fmt.Errorf("%w: %w", dapi.ErrIPC, err)
	}

	found := findFocusedNode(&root)
	if found != nil {
		return &found.Rect, nil
	}

	return nil, dapi.ErrNoFocusedWindow
}

func findFocusedNode(node *api.Node) *api.Node {
	if node.Focused {
		return node
	}
	for _, n := range node.Nodes {
		if found := findFocusedNode(n); found != nil {
			return found
		}
	}
	for _, n := range node.FloatingNodes {
		if found := findFocusedNode(n); found != nil {
			return found
		}
	}
	return nil
}
