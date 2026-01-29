package webapp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"workspaced/pkg/common"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "launch [name|url]",
			Short: "Launch a webapp",
			RunE: func(c *cobra.Command, args []string) error {
				cfg, err := common.LoadConfig()
				if err != nil {
					return err
				}

				var url string
				var wa common.WebappConfig
				var found bool

				if len(args) == 0 {
					// Use zenity to ask for URL
					out, err := common.RunCmd(c.Context(), "zenity", "--entry", "--text=Link to be opened").Output()
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
					} else {
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

func launchWebapp(ctx context.Context, url string, wa common.WebappConfig, engine string, isConfigured bool) error {
	normalizedURL := common.NormalizeURL(url)
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

	cmd := common.RunCmd(ctx, engine, args...)
	common.InheritContextWriters(ctx, cmd)
	return cmd.Run()
}
