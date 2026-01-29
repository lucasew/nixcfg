package config

import (
	"encoding/json"
	"fmt"
	"strings"
	"workspaced/pkg/common"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "get <key>",
			Short: "Get a configuration value (outputs JSON)",
			Long: `Get a configuration value using dot notation.

Examples:
  workspaced dispatch config get workspaces.www
  workspaced dispatch config get wallpaper.dir
  workspaced dispatch config get wallpaper

Outputs the value as JSON for easy parsing.`,
			Args: cobra.ExactArgs(1),
			RunE: func(c *cobra.Command, args []string) error {
				key := args[0]
				cfg, err := common.LoadConfig()
				if err != nil {
					return fmt.Errorf("failed to load config: %w", err)
				}

				result, err := getConfigValue(cfg, key)
				if err != nil {
					return err
				}

				// Encode to JSON
				jsonBytes, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to encode JSON: %w", err)
				}

				fmt.Println(string(jsonBytes))
				return nil
			},
		})
	})
}

func getConfigValue(cfg *common.GlobalConfig, key string) (any, error) {
	if key == "" {
		return cfg, nil
	}

	parts := strings.Split(key, ".")

	switch parts[0] {
	case "workspaces":
		if len(parts) == 2 {
			if val, ok := cfg.Workspaces[parts[1]]; ok {
				return val, nil
			}
			return nil, fmt.Errorf("workspace key not found: %s", parts[1])
		} else if len(parts) == 1 {
			return cfg.Workspaces, nil
		}

	case "wallpaper":
		if len(parts) == 2 {
			switch parts[1] {
			case "dir":
				return cfg.Wallpaper.Dir, nil
			case "default":
				return cfg.Wallpaper.Default, nil
			default:
				return nil, fmt.Errorf("unknown wallpaper field: %s", parts[1])
			}
		} else if len(parts) == 1 {
			return cfg.Wallpaper, nil
		}

	case "screenshot":
		if len(parts) == 2 {
			if parts[1] == "dir" {
				return cfg.Screenshot.Dir, nil
			}
			return nil, fmt.Errorf("unknown screenshot field: %s", parts[1])
		} else if len(parts) == 1 {
			return cfg.Screenshot, nil
		}

	case "hosts":
		if len(parts) == 2 {
			if val, ok := cfg.Hosts[parts[1]]; ok {
				return val, nil
			}
			return nil, fmt.Errorf("host key not found: %s", parts[1])
		} else if len(parts) == 1 {
			return cfg.Hosts, nil
		}

	case "backup":
		if len(parts) == 2 {
			switch parts[1] {
			case "rsyncnet_user":
				return cfg.Backup.RsyncnetUser, nil
			case "remote_path":
				return cfg.Backup.RemotePath, nil
			default:
				return nil, fmt.Errorf("unknown backup field: %s", parts[1])
			}
		} else if len(parts) == 1 {
			return cfg.Backup, nil
		}

	case "quicksync":
		if len(parts) == 2 {
			switch parts[1] {
			case "repo_dir":
				return cfg.QuickSync.RepoDir, nil
			case "remote_path":
				return cfg.QuickSync.RemotePath, nil
			default:
				return nil, fmt.Errorf("unknown quicksync field: %s", parts[1])
			}
		} else if len(parts) == 1 {
			return cfg.QuickSync, nil
		}

	case "browser":
		if len(parts) == 2 {
			switch parts[1] {
			case "default":
				return cfg.Browser.Default, nil
			case "webapp":
				return cfg.Browser.Engine, nil
			default:
				return nil, fmt.Errorf("unknown browser field: %s", parts[1])
			}
		} else if len(parts) == 1 {
			return cfg.Browser, nil
		}

	case "webapp":
		if len(parts) >= 2 {
			if wa, ok := cfg.Webapps[parts[1]]; ok {
				if len(parts) == 2 {
					return wa, nil
				}
				switch parts[2] {
				case "url":
					return wa.URL, nil
				case "profile":
					return wa.Profile, nil
				case "desktop_name":
					return wa.DesktopName, nil
				default:
					return nil, fmt.Errorf("unknown webapp field: %s", parts[2])
				}
			}
			return nil, fmt.Errorf("webapp not found: %s", parts[1])
		} else if len(parts) == 1 {
			return cfg.Webapps, nil
		}

	default:
		return nil, fmt.Errorf("unknown config section: %s", parts[0])
	}

	return nil, fmt.Errorf("invalid key format: %s", key)
}
