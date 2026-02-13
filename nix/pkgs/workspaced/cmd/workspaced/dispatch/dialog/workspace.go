package dialog

import (
	"sort"
	"strconv"
	"workspaced/pkg/driver"
	"workspaced/pkg/driver/dialog"
	"workspaced/pkg/driver/wm"

	"workspaced/pkg/config"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		cmd := &cobra.Command{
			Use:   "workspace",
			Short: "Workspace switcher",
			RunE: func(c *cobra.Command, args []string) error {
				move, _ := c.Flags().GetBool("move")
				cfg, err := config.LoadConfig()
				if err != nil {
					return err
				}

				var items []dialog.Item
				var keys []string
				for k := range cfg.Workspaces {
					keys = append(keys, k)
				}
				sort.Strings(keys)

				for _, k := range keys {
					items = append(items, dialog.Item{
						Label: k,
						Value: strconv.Itoa(cfg.Workspaces[k]),
					})
				}

				d, err := driver.Get[dialog.Driver](c.Context())
				if err != nil {
					return err
				}

				selected, err := d.Choose(c.Context(), dialog.Options{
					Prompt: "Workspace",
					Items:  items,
				})
				if err != nil {
					return err
				}

				if selected == nil {
					return nil
				}

				return wm.SwitchToWorkspace(c.Context(), selected.Value, move)
			},
		}
		cmd.Flags().Bool("move", false, "Move container to workspace")
		parent.AddCommand(cmd)
	})
}
