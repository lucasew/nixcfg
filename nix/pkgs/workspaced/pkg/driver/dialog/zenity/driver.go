package zenity

import (
	"context"
	"fmt"
	"strings"
	"workspaced/pkg/driver"
	"workspaced/pkg/driver/dialog"
	execdriver "workspaced/pkg/driver/exec"
	"workspaced/pkg/executil"
)

func init() {
	driver.Register[dialog.Prompter](&PrompterProvider{})
	driver.Register[dialog.Confirmer](&ConfirmerProvider{})
}

type baseProvider struct{}

func (p *baseProvider) ID() string         { return "zenity" }
func (p *baseProvider) DefaultWeight() int { return driver.DefaultWeight }

func (p *baseProvider) CheckCompatibility(ctx context.Context) error {
	if executil.GetEnv(ctx, "DISPLAY") == "" && executil.GetEnv(ctx, "WAYLAND_DISPLAY") == "" {
		return fmt.Errorf("%w: neither DISPLAY nor WAYLAND_DISPLAY set", driver.ErrIncompatible)
	}
	if !execdriver.IsBinaryAvailable(ctx, "zenity") {
		return fmt.Errorf("%w: zenity not found", driver.ErrIncompatible)
	}
	return nil
}

type PrompterProvider struct{ baseProvider }

func (p *PrompterProvider) Name() string                                     { return "Zenity (Prompt)" }
func (p *PrompterProvider) New(ctx context.Context) (dialog.Prompter, error) { return &Driver{}, nil }

type ConfirmerProvider struct{ baseProvider }

func (p *ConfirmerProvider) Name() string                                      { return "Zenity (Confirm)" }
func (p *ConfirmerProvider) New(ctx context.Context) (dialog.Confirmer, error) { return &Driver{}, nil }

type Driver struct{}

func (d *Driver) Prompt(ctx context.Context, prompt string) (string, error) {
	out, err := execdriver.MustRun(ctx, "zenity", "--entry", "--text", prompt).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (d *Driver) Confirm(ctx context.Context, message string) (bool, error) {
	err := execdriver.MustRun(ctx, "zenity", "--question", "--text", message).Run()
	if err != nil {
		// Zenity returns non-zero if No is selected
		return false, nil
	}
	return true, nil
}
