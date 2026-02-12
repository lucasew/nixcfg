package hyprland

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
	driver.Register[api.Driver](&Provider{})
}

type Provider struct{}

func (p *Provider) Name() string { return "Hyprland" }

func (p *Provider) CheckCompatibility(ctx context.Context) error {
	if exec.GetEnv(ctx, "HYPRLAND_INSTANCE_SIGNATURE") == "" {
		return fmt.Errorf("%w: HYPRLAND_INSTANCE_SIGNATURE not set", driver.ErrIncompatible)
	}
	if !exec.IsBinaryAvailable(ctx, "hyprctl") {
		return fmt.Errorf("%w: hyprctl not found", driver.ErrIncompatible)
	}
	return nil
}

func (p *Provider) New(ctx context.Context) (api.Driver, error) {
	return &Driver{}, nil
}

type Driver struct{}

func (d *Driver) MoveWorkspaceToOutput(ctx context.Context, workspace string, output string) error {
	return exec.RunCmd(ctx, "hyprctl", "dispatch", "moveworkspacetomonitor", workspace, output).Run()
}

func (d *Driver) SwitchToWorkspace(ctx context.Context, num int, move bool) error {
	cmd := "workspace"
	if move {
		cmd = "movetoworkspace"
	}
	return exec.RunCmd(ctx, "hyprctl", "dispatch", cmd, strconv.Itoa(num)).Run()
}

func (d *Driver) ToggleScratchpad(ctx context.Context) error {
	return exec.RunCmd(ctx, "hyprctl", "dispatch", "togglespecialworkspace").Run()
}

func (d *Driver) GetOutputs(ctx context.Context) ([]api.Output, error) {
	out, err := exec.RunCmd(ctx, "hyprctl", "monitors", "-j").Output()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", dapi.ErrIPC, err)
	}
	var monitors []struct {
		Name            string `json:"name"`
		Focused         bool   `json:"focused"`
		X               int    `json:"x"`
		Y               int    `json:"y"`
		Width           int    `json:"width"`
		Height          int    `json:"height"`
		ActiveWorkspace struct {
			Name string `json:"name"`
		} `json:"activeWorkspace"`
	}
	if err := json.Unmarshal(out, &monitors); err != nil {
		return nil, fmt.Errorf("%w: %w", dapi.ErrIPC, err)
	}
	var outputs []api.Output
	for _, m := range monitors {
		outputs = append(outputs, api.Output{
			Name:             m.Name,
			Focused:          m.Focused,
			CurrentWorkspace: m.ActiveWorkspace.Name,
			Rect:             api.Rect{X: m.X, Y: m.Y, Width: m.Width, Height: m.Height},
		})
	}
	return outputs, nil
}

func (d *Driver) GetWorkspaces(ctx context.Context) ([]api.Workspace, error) {
	out, err := exec.RunCmd(ctx, "hyprctl", "workspaces", "-j").Output()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", dapi.ErrIPC, err)
	}
	var workspaces []struct {
		Name    string `json:"name"`
		Monitor string `json:"monitor"`
	}
	if err := json.Unmarshal(out, &workspaces); err != nil {
		return nil, fmt.Errorf("%w: %w", dapi.ErrIPC, err)
	}

	activeWSOut, _ := exec.RunCmd(ctx, "hyprctl", "activeworkspace", "-j").Output()
	var activeWS struct {
		Name string `json:"name"`
	}
	_ = json.Unmarshal(activeWSOut, &activeWS)

	var result []api.Workspace
	for _, w := range workspaces {
		result = append(result, api.Workspace{
			Name:    w.Name,
			Output:  w.Monitor,
			Focused: w.Name == activeWS.Name,
		})
	}
	return result, nil
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
	return "", nil, dapi.ErrNoFocusedOutput
}

func (d *Driver) GetFocusedWindowRect(ctx context.Context) (*api.Rect, error) {
	out, err := exec.RunCmd(ctx, "hyprctl", "activewindow", "-j").Output()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", dapi.ErrIPC, err)
	}
	var win struct {
		At   []int `json:"at"`
		Size []int `json:"size"`
	}
	if err := json.Unmarshal(out, &win); err != nil {
		return nil, fmt.Errorf("%w: %w", dapi.ErrIPC, err)
	}
	if len(win.At) != 2 || len(win.Size) != 2 {
		return nil, fmt.Errorf("%w: invalid hyprland active window geometry", dapi.ErrIPC)
	}
	return &api.Rect{X: win.At[0], Y: win.At[1], Width: win.Size[0], Height: win.Size[1]}, nil
}
