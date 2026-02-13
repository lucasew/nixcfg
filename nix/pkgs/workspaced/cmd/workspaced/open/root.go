package open

import (
	"fmt"
	"workspaced/pkg/config"
	"workspaced/pkg/driver/opener"
	"workspaced/pkg/driver/terminal"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open [target]",
		Short: "Open a file, URL or webapp",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) == 0 {
				return c.Help()
			}
			return opener.Open(c.Context(), args[0])
		},
	}

	var urlFlag string
	var profileFlag string
	var extraFlags []string

	webappCmd := &cobra.Command{
		Use:   "webapp [name]",
		Short: "Launch a configured webapp",
		RunE: func(c *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			var wa opener.WebappConfig

			if len(args) > 0 {
				name := args[0]
				var modCfg struct {
					Apps map[string]opener.WebappConfig `json:"apps"`
				}
				if err := cfg.Module("webapp", &modCfg); err == nil {
					if app, ok := modCfg.Apps[name]; ok {
						wa = app
					}
				}
			}

			// Override with flags
			if urlFlag != "" {
				wa.URL = urlFlag
			}
			if profileFlag != "" {
				wa.Profile = profileFlag
			}
			if len(extraFlags) > 0 {
				wa.ExtraFlags = append(wa.ExtraFlags, extraFlags...)
			}

			if wa.URL == "" {
				return fmt.Errorf("URL is required (use name or --url)")
			}

			return opener.OpenWebapp(c.Context(), wa)
		},
	}

	webappCmd.Flags().StringVarP(&urlFlag, "url", "u", "", "URL to open")
	webappCmd.Flags().StringVarP(&profileFlag, "profile", "p", "", "Browser profile name")
	webappCmd.Flags().StringSliceVarP(&extraFlags, "extra-flag", "e", nil, "Extra browser flags")

	cmd.AddCommand(webappCmd)

	cmd.AddCommand(&cobra.Command{
		Use:   "terminal [cmd...]",
		Short: "Launch the preferred terminal",
		RunE: func(c *cobra.Command, args []string) error {
			opts := terminal.Options{
				Title: "Terminal",
			}
			if len(args) > 0 {
				opts.Command = args[0]
				opts.Args = args[1:]
			}
			return terminal.Open(c.Context(), opts)
		},
	})

	return cmd
}
