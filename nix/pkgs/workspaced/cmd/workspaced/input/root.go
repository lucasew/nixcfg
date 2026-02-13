package input

import (
	"fmt"
	"workspaced/pkg/driver/dialog"
	"workspaced/pkg/registry"

	"github.com/spf13/cobra"
)

var Registry registry.CommandRegistry

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "input",
		Short: "Interactive user input commands",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "text [prompt]",
		Short: "Ask for a text string",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			prompt := "Input"
			if len(args) > 0 {
				prompt = args[0]
			}
			res, err := dialog.Prompt(c.Context(), prompt)
			if err != nil {
				return err
			}
			fmt.Println(res)
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "confirm [message]",
		Short: "Ask for a yes/no confirmation",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			msg := "Confirm?"
			if len(args) > 0 {
				msg = args[0]
			}
			ok, err := dialog.Confirm(c.Context(), msg)
			if err != nil {
				return err
			}
			if !ok {
				return fmt.Errorf("cancelled")
			}
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "choose [prompt] [options...]",
		Short: "Select an item from a list",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			prompt := args[0]
			var items []dialog.Item
			for _, arg := range args[1:] {
				items = append(items, dialog.Item{Label: arg, Value: arg})
			}
			res, err := dialog.Choose(c.Context(), dialog.ChooseOptions{
				Prompt: prompt,
				Items:  items,
			})
			if err != nil {
				return err
			}
			if res != nil {
				fmt.Println(res.Value)
			}
			return nil
		},
	})

	return Registry.GetCommand(cmd)
}
