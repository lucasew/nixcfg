package webapp

import (
	"fmt"
	"sort"

	"workspaced/pkg/config"
	"workspaced/pkg/text"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "list",
			Short: "List configured webapps",
			RunE: func(c *cobra.Command, args []string) error {
				cfg, err := config.Load()
				if err != nil {
					return err
				}

				var modCfg struct {
					Apps map[string]config.WebappConfig `json:"apps"`
				}
				if err := cfg.Module("webapp", &modCfg); err != nil {
					return fmt.Errorf("webapp module error: %w", err)
				}

				if len(modCfg.Apps) == 0 {
					fmt.Println("No webapps configured.")
					return nil
				}

				fmt.Printf("%-20s %-30s %s\n", "NAME", "DISPLAY NAME", "URL")
				fmt.Println(text.ToTitleCase(fmt.Sprintf("%-20s %-30s %s", "----", "------------", "---")))

				// Sort names for stable output
				var names []string
				for name := range modCfg.Apps {
					names = append(names, name)
				}
				sort.Strings(names)

				for _, name := range names {
					app := modCfg.Apps[name]
					displayName := app.DesktopName
					if displayName == "" {
						displayName = text.ToTitleCase(name)
					}
					fmt.Printf("%-20s %-30s %s\n", name, displayName, app.URL)
				}
				return nil
			},
		})
	})
}
