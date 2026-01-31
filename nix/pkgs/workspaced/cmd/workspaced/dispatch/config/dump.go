package config

import (
	"fmt"
	"workspaced/pkg/common"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

func GetDumpCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "dump",
		Short: "Dump the full merged configuration as TOML",
		Long: `Dump the complete merged configuration from all sources:
- Hardcoded defaults
- $DOTFILES/settings.toml
- ~/settings.toml

Outputs the result as TOML format.`,
		RunE: func(c *cobra.Command, args []string) error {
			cfg, err := common.LoadConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Encode to TOML format
			encoder := toml.NewEncoder(c.OutOrStdout())
			if err := encoder.Encode(cfg); err != nil {
				return fmt.Errorf("failed to encode TOML: %w", err)
			}

			return nil
		},
	}
}
