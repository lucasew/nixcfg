package launch

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"workspaced/pkg/config"
	"workspaced/pkg/driver"
	"workspaced/pkg/driver/dialog"
	"workspaced/pkg/env"
	execdriver "workspaced/pkg/driver/exec"
	"workspaced/pkg/executil"

	"github.com/spf13/cobra"
)

type WebappConfig struct {
	URL         string   `json:"url"`
	Profile     string   `json:"profile"`
	DesktopName string   `json:"desktop_name"`
	Icon        string   `json:"icon"`
	ExtraFlags  []string `json:"extra_flags"`
}

func newWebappCommand() *cobra.Command {
	var urlFlag string
	var profileFlag string
	var extraFlags []string

	cmd := &cobra.Command{
		Use:   "webapp",
		Short: "Launch a webapp",
		RunE: func(c *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			if urlFlag == "" {
				d, err := driver.Get[dialog.Prompter](c.Context())
				if err != nil {
					return fmt.Errorf("--url is required and no prompter available: %w", err)
				}
				urlFlag, err = d.Prompt(c.Context(), "Enter URL")
				if err != nil {
					return err
				}
			}

			if urlFlag == "" {
				return fmt.Errorf("--url is required")
			}

			wa := WebappConfig{
				URL:        urlFlag,
				Profile:    profileFlag,
				ExtraFlags: extraFlags,
			}

			return launchWebapp(c.Context(), urlFlag, wa, cfg.Browser.Engine)
		},
	}

	cmd.Flags().StringVarP(&urlFlag, "url", "u", "", "URL to open")
	cmd.Flags().StringVarP(&profileFlag, "profile", "p", "", "Browser profile name")
	cmd.Flags().StringSliceVarP(&extraFlags, "extra-flag", "e", nil, "Extra browser flags")

	return cmd
}

func launchWebapp(ctx context.Context, url string, wa WebappConfig, engine string) error {
	normalizedURL := env.NormalizeURL(url)
	args := []string{"--app=" + normalizedURL}

	if wa.Profile != "" {
		home, _ := os.UserHomeDir()
		profileDir := filepath.Join(home, ".config/workspaced/webapp/profiles", wa.Profile)
		args = append(args, "--user-data-dir="+profileDir)
	}

	if os.Getenv("WAYLAND_DISPLAY") != "" {
		args = append(args, "--enable-features=UseOzonePlatform", "--ozone-platform=wayland")
	}

	args = append(args, wa.ExtraFlags...)

	cmd := execdriver.MustRun(ctx, engine, args...)
	executil.InheritContextWriters(ctx, cmd)
	return cmd.Run()
}
