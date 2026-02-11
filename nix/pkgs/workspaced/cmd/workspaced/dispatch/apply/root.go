package apply

import (
	"fmt"
	"os"
	"path/filepath"
	"workspaced/pkg/apply"
	"workspaced/pkg/env"
	"workspaced/pkg/exec"
	"workspaced/pkg/logging"
	"workspaced/pkg/nix"

	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply [action]",
		Short: "Declaratively apply system and user configurations",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := logging.GetLogger(ctx)

			action := "switch"
			if len(args) > 0 {
				action = args[0]
			}

			dryRun, _ := cmd.Flags().GetBool("dry-run")

			// 1. Load state
			state, err := apply.LoadState()
			if err != nil {
				return err
			}

			// 2. Collect desired state
			providers := []apply.Provider{
				&apply.DconfProvider{},
				&apply.SymlinkProvider{},
			}

			desired := []apply.DesiredState{}
			for _, p := range providers {
				d, err := p.GetDesiredState(ctx)
				if err != nil {
					return fmt.Errorf("provider %s failed: %w", p.Name(), err)
				}
				desired = append(desired, d...)
			}

			// 3. Plan
			actions, err := apply.Plan(ctx, desired, state)
			if err != nil {
				return err
			}

			// 4. Show and execute
			hasChanges := false
			for _, a := range actions {
				if a.Type != apply.ActionNoop {
					hasChanges = true
					cmd.Printf("[%s] %s\n", a.Type, a.Target)
					if a.Type == apply.ActionUpdate || a.Type == apply.ActionCreate {
						cmd.Printf("      -> %s\n", a.Desired.Source)
					}
				}
			}

			if !hasChanges {
				logger.Info("no file changes needed")
			} else if dryRun {
				logger.Info("dry-run: skipping file execution")
			} else {
				if err := apply.Execute(ctx, actions, state); err != nil {
					return err
				}
				if err := apply.SaveState(state); err != nil {
					return err
				}

				// Reload GTK theme if not on Termux
				if !env.IsPhone() {
					home, _ := os.UserHomeDir()
					dummyTheme := home + "/.local/share/themes/dummy"
					if _, err := os.Stat(dummyTheme); err == nil {
						// Switch to dummy and back to force GTK reload
						if err := exec.RunCmd(ctx, "dconf", "write", "/org/gnome/desktop/interface/gtk-theme", "'dummy'").Run(); err != nil {
							logger.Warn("failed to switch to dummy theme", "error", err)
						}
						if err := exec.RunCmd(ctx, "dconf", "write", "/org/gnome/desktop/interface/gtk-theme", "'base16'").Run(); err != nil {
							logger.Warn("failed to switch to base16 theme", "error", err)
						}
					}
				}
			}

			// 5. System specific hooks
			if env.IsNixOS() {
				logger.Info("running NixOS rebuild", "action", action)
				if dryRun {
					logger.Info("dry-run: skipping nixos-rebuild")
				} else {
					flake := ""
					if env.IsRiverwood() {
						logger.Info("performing remote build for riverwood")
						hostname := env.GetHostname()
						ref := fmt.Sprintf(".#nixosConfigurations.%s.config.system.build.toplevel", hostname)
						result, err := nix.RemoteBuild(ctx, ref, "whiterun", true)
						if err != nil {
							return fmt.Errorf("remote build failed: %w", err)
						}
						flake = result
					}
					if err := nix.Rebuild(ctx, action, flake); err != nil {
						return err
					}
				}
			}

			return nil
		},
	}
	cmd.Flags().BoolP("dry-run", "d", false, "Only show what would be done")
	return cmd
}
