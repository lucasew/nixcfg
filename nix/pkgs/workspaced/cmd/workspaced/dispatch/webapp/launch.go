package webapp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
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
				cfg, err := config.Load()
				if err != nil {
					return err
				}

				var url string
				var wa WebappConfig
				var found bool

				if len(args) == 0 {
					// Use zenity to ask for URL
					out, err := exec.RunCmd(c.Context(), "zenity", "--entry", "--text=Link to be opened").Output()
					if err != nil {
						return nil // User cancelled
					}
					url = string(out)
				} else {
					target := args[0]
					var modCfg struct {
						Apps map[string]WebappConfig `json:"apps"`
					}
					if err := cfg.Module("webapp", &modCfg); err != nil {
						return fmt.Errorf("webapp module error: %w", err)
					}

					if app, ok := modCfg.Apps[target]; ok {
						wa = app
						url = app.URL
						found = true
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

func launchWebapp(ctx context.Context, url string, wa WebappConfig, engine string, isConfigured bool) error {
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
