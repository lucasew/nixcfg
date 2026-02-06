package webapp

import (
	"fmt"

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
				cfg, err := config.LoadConfig()
				if err != nil {
					return err
				}

				if len(cfg.Webapps) == 0 {
					fmt.Println("No webapps configured.")
					return nil
				}

				fmt.Printf("%-20s %-30s %s\n", "NAME", "DISPLAY NAME", "URL")
				fmt.Println(text.ToTitleCase(fmt.Sprintf("%-20s %-30s %s", "----", "------------", "---")))

				for name, wa := range cfg.Webapps {
					displayName := wa.DesktopName
					if displayName == "" {
						displayName = text.ToTitleCase(name)
					}
					fmt.Printf("%-20s %-30s %s\n", name, displayName, wa.URL)
				}
				return nil
			},
		})
	})
}
