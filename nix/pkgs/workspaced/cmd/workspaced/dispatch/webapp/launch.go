package webapp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"workspaced/pkg/env"
	"workspaced/pkg/exec"

	"workspaced/pkg/config"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "launch [name|url]",
			Short: "Launch a webapp",
			RunE: func(c *cobra.Command, args []string) error {
				cfg, err := config.LoadConfig()
				if err != nil {
					return err
				}

				var url string
				var wa config.WebappConfig
				var found bool

				if len(args) == 0 {
					// Use zenity to ask for URL
					out, err := exec.RunCmd(c.Context(), "zenity", "--entry", "--text=Link to be opened").Output()
					if err != nil {
						return nil // User cancelled
					}
					url = strings.TrimSpace(string(out))
				} else {
					target := args[0]
					if w, ok := cfg.Webapps[target]; ok {
						wa = w
						found = true
						url = wa.URL
					} else if modWebapp, ok := cfg.Modules["webapp"].(map[string]any); ok {
						if apps, ok := modWebapp["apps"].(map[string]any); ok {
							if appRaw, ok := apps[target].(map[string]any); ok {
								// Extract fields from map
								if u, ok := appRaw["url"].(string); ok {
									url = u
									wa.URL = u
								}
								if p, ok := appRaw["profile"].(string); ok {
									wa.Profile = p
								}
								if d, ok := appRaw["desktop_name"].(string); ok {
									wa.DesktopName = d
								}
								if i, ok := appRaw["icon"].(string); ok {
									wa.Icon = i
								}
								if f, ok := appRaw["extra_flags"].([]string); ok {
									wa.ExtraFlags = f
								}
								found = true
							}
						}
					}

					if !found {
						url = target
					}
				}

				if url == "" {
					return fmt.Errorf("no URL specified")
				}

				return launchWebapp(c.Context(), url, wa, cfg.Browser.Engine, found)
			},
		})
	})
}

func launchWebapp(ctx context.Context, url string, wa config.WebappConfig, engine string, isConfigured bool) error {
	normalizedURL := env.NormalizeURL(url)
	args := []string{"--app=" + normalizedURL}

	if isConfigured && wa.Profile != "" {
		home, _ := os.UserHomeDir()
		profileDir := filepath.Join(home, ".config/workspaced/webapp/profiles", wa.Profile)
		args = append(args, "--user-data-dir="+profileDir)
	}

	if os.Getenv("WAYLAND_DISPLAY") != "" {
		args = append(args, "--enable-features=UseOzonePlatform", "--ozone-platform=wayland")
	}

	if isConfigured {
		args = append(args, wa.ExtraFlags...)
	}

	cmd := exec.RunCmd(ctx, engine, args...)
	exec.InheritContextWriters(ctx, cmd)
	return cmd.Run()
}
