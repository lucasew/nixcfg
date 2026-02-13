package terminal

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"workspaced/pkg/driver"
	"workspaced/pkg/driver/dialog"

	"github.com/ktr0731/go-fuzzyfinder"
)

func init() {
	driver.Register[dialog.Chooser](&ChooserProvider{})
	driver.Register[dialog.Prompter](&PrompterProvider{})
	driver.Register[dialog.Confirmer](&ConfirmerProvider{})
}

type baseProvider struct{}

func (p *baseProvider) ID() string         { return "terminal" }
func (p *baseProvider) DefaultWeight() int { return 0 }

func (p *baseProvider) CheckCompatibility(ctx context.Context) error {
	// Sempre compatível, mas com peso 0 para ser fallback
	return nil
}

type ChooserProvider struct{ baseProvider }

func (p *ChooserProvider) Name() string                                    { return "Terminal (Fuzzy)" }
func (p *ChooserProvider) New(ctx context.Context) (dialog.Chooser, error) { return &Driver{}, nil }

type PrompterProvider struct{ baseProvider }

func (p *PrompterProvider) Name() string                                     { return "Terminal (Stdin)" }
func (p *PrompterProvider) New(ctx context.Context) (dialog.Prompter, error) { return &Driver{}, nil }

type ConfirmerProvider struct{ baseProvider }

func (p *ConfirmerProvider) Name() string                                      { return "Terminal (y/n)" }
func (p *ConfirmerProvider) New(ctx context.Context) (dialog.Confirmer, error) { return &Driver{}, nil }

// Driver implements Chooser, Prompter and Confirmer
type Driver struct{}

func (d *Driver) Choose(ctx context.Context, opts dialog.ChooseOptions) (*dialog.Item, error) {
	idx, err := fuzzyfinder.Find(
		opts.Items,
		func(i int) string {
			return opts.Items[i].Label
		},
	)
	if err != nil {
		if err == fuzzyfinder.ErrAbort {
			return nil, nil
		}
		return nil, err
	}
	return &opts.Items[idx], nil
}

func (d *Driver) Prompt(ctx context.Context, prompt string) (string, error) {
	fmt.Printf("%s: ", prompt)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text()), nil
	}
	return "", scanner.Err()
}

func (d *Driver) Confirm(ctx context.Context, message string) (bool, error) {
	fmt.Printf("%s [y/N]: ", message)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		text := strings.ToLower(strings.TrimSpace(scanner.Text()))
		return text == "y" || text == "yes", nil
	}
	return false, scanner.Err()
}

// Legacy compatibility for Driver interface
func (d *Driver) RunApp(ctx context.Context) error {
	return fmt.Errorf("RunApp not implemented for Terminal")
}
func (d *Driver) SwitchWindow(ctx context.Context) error {
	return fmt.Errorf("SwitchWindow not implemented for Terminal")
}
