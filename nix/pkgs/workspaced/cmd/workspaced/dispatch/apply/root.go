package apply

import (
	"fmt"
	"workspaced/pkg/apply"
	"workspaced/pkg/drivers/nix"
	"workspaced/pkg/env"
	"workspaced/pkg/logging"

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

			engine := apply.NewEngine(nil)

			// 1. Load state
			state, err := engine.LoadState()
			if err != nil {
				return err
			}

			// 2. Collect desired state
			providers := []apply.Provider{
				&apply.ProfileProvider{},
				&apply.DconfProvider{},
				&apply.DesktopProvider{},
				&apply.SymlinkProvider{},
				&apply.TermuxProvider{},
				&apply.WebappProvider{},
				&apply.LazyShimProvider{},
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
			actions, err := engine.Plan(ctx, desired, state)
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
				if err := engine.Execute(ctx, actions, state); err != nil {
					return err
				}
				if err := engine.SaveState(state); err != nil {
					return err
				}
			}

			// 5. Home Manager hook
			if !env.IsPhone() {
				logger.Info("running home-manager rebuild")
				if dryRun {
					logger.Info("dry-run: skipping home-manager rebuild")
				} else {
					flake := ""
					if env.IsRiverwood() {
						logger.Info("performing remote build for home-manager on riverwood")
						ref := ".#homeConfigurations.main.activationPackage"
						result, err := nix.RemoteBuild(ctx, ref, "whiterun", true)
						if err != nil {
							return fmt.Errorf("remote build for home-manager failed: %w", err)
						}
						flake = result
					}
					if err := nix.HomeManagerSwitch(ctx, action, flake); err != nil {
						return err
					}
				}
			}

			// 6. System specific hooks
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
