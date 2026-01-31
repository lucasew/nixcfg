package config

import (
	"fmt"
	"workspaced/pkg/config"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "colors",
			Short: "Output ANSI escape sequences for the current color palette",
			RunE: func(c *cobra.Command, args []string) error {
				cfg, err := config.Load()
				if err != nil {
					return err
				}

				var desktop map[string]interface{}
				if err := cfg.UnmarshalKey("desktop", &desktop); err != nil {
					return err
				}

				palette, ok := desktop["palette"].(map[string]interface{})
				if !ok {
					return fmt.Errorf("palette not found")
				}

				base16, ok := palette["base16"].(map[string]interface{})
				if !ok {
					return fmt.Errorf("base16 palette not found")
				}

				get := func(key string) string {
					val, _ := base16[key].(string)
					return val
				}

				colors := []string{
					get("base00"), get("base08"), get("base0B"), get("base0A"),
					get("base0D"), get("base0E"), get("base0C"), get("base05"),
					get("base03"), get("base08"), get("base0B"), get("base0A"),
					get("base0D"), get("base0E"), get("base0C"), get("base07"),
				}

				for i, color := range colors {
					if color != "" {
						fmt.Printf("printf \"\\033]4;%d;#%s\\033\\\\\"\n", i, color)
					}
				}

				fg := get("base05")
				bg := get("base00")
				if fg != "" {
					fmt.Printf("printf \"\\033]10;#%s\\033\\\\\"\n", fg)
					fmt.Printf("printf \"\\033]12;#%s\\033\\\\\"\n", fg)
				}
				if bg != "" {
					fmt.Printf("printf \"\\033]11;#%s\\033\\\\\"\n", bg)
				}

				return nil
			},
		})
	})
}
