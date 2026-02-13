package dialog

import (
	"context"
	"workspaced/pkg/driver"
)

// Choose allows selecting an item from a list. It tries graphical choosers first.
func Choose(ctx context.Context, opts ChooseOptions) (*Item, error) {
	d, err := driver.Get[Chooser](ctx)
	if err != nil {
		return nil, err
	}
	return d.Choose(ctx, opts)
}

// Prompt asks for a simple text input.
func Prompt(ctx context.Context, prompt string) (string, error) {
	d, err := driver.Get[Prompter](ctx)
	if err != nil {
		return "", err
	}
	return d.Prompt(ctx, prompt)
}

// Confirm asks a yes/no question.
func Confirm(ctx context.Context, message string) (bool, error) {
	d, err := driver.Get[Confirmer](ctx)
	if err != nil {
		return false, err
	}
	return d.Confirm(ctx, message)
}
