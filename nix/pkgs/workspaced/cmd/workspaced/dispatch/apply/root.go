package apply

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"workspaced/pkg/apply"
	"workspaced/pkg/config"
	"workspaced/pkg/deployer"
	"workspaced/pkg/dotfiles"
	"workspaced/pkg/env"
	"workspaced/pkg/exec"
	"workspaced/pkg/logging"
	"workspaced/pkg/nix"
	"workspaced/pkg/source"
	"workspaced/pkg/template"

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

			// Carregar configuração
			cfg, err := config.LoadConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Obter dotfiles root
			dotfilesRoot, err := env.GetDotfilesRoot()
			if err != nil {
				return fmt.Errorf("failed to get dotfiles root: %w", err)
			}

			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to get home directory: %w", err)
			}

			// Criar template engine compartilhada
			engine := template.NewEngine(ctx)

			// Configurar pipeline de plugins
			configDir := filepath.Join(dotfilesRoot, "config")
			pipeline := source.NewPipeline()

			// 1. Provider dconf (legacy)
			pipeline.AddPlugin(source.NewProviderPlugin(&apply.DconfProvider{}, 50))

			// 2. Scanner - descobre arquivos em config/
			if _, err := os.Stat(configDir); err == nil {
				scanner, err := source.NewScannerPlugin(source.ScannerConfig{
					Name:       "legacy-config",
					BaseDir:    configDir,
					TargetBase: home,
					Priority:   50, // Legacy has lower priority than modules
				})
				if err != nil {
					return fmt.Errorf("failed to create scanner: %w", err)
				}
				pipeline.AddPlugin(scanner)
			}

			// 2.5 Modules Scanner
			modulesDir := filepath.Join(dotfilesRoot, "modules")
			if _, err := os.Stat(modulesDir); err == nil {
				pipeline.AddPlugin(source.NewModuleScannerPlugin(modulesDir, cfg, 100))
			}

			// 3. TemplateExpander - renderiza .tmpl (inclui multi-file)
			pipeline.AddPlugin(source.NewTemplateExpanderPlugin(engine, cfg))

			// 4. DotDProcessor - concatena .d.tmpl/
			pipeline.AddPlugin(source.NewDotDProcessorPlugin(engine, cfg))

			// 5. StrictConflictResolver - garante unicidade total
			pipeline.AddPlugin(source.NewStrictConflictResolverPlugin())

			// StateStore
			stateStore, err := deployer.NewFileStateStore("~/.config/workspaced/state.json")
			if err != nil {
				return fmt.Errorf("failed to create state store: %w", err)
			}

			// Hooks
			hooks := []dotfiles.Hook{
				// Hook para reload GTK theme
				&dotfiles.FuncHook{
					AfterFn: func(ctx context.Context, actions []deployer.Action, execErr error) error {
						if execErr != nil {
							return nil // Não executar se houve erro
						}
						if env.IsPhone() {
							return nil // Não executar em phone
						}

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
						return nil
					},
				},
			}

			// Criar manager com pipeline
			mgr, err := dotfiles.NewManager(dotfiles.Config{
				Pipeline:   pipeline,
				StateStore: stateStore,
				Hooks:      hooks,
			})
			if err != nil {
				return fmt.Errorf("failed to create manager: %w", err)
			}

			// Aplicar configurações
			result, err := mgr.Apply(ctx, dotfiles.ApplyOptions{
				DryRun: dryRun,
			})
			if err != nil {
				return err
			}

			// Mostrar resultado
			if result.FilesCreated > 0 || result.FilesUpdated > 0 || result.FilesDeleted > 0 {
				for _, a := range result.Actions {
					if a.Type != deployer.ActionNoop {
						cmd.Printf("[%s] %s\n", a.Type, a.Target)
						if a.Type == deployer.ActionUpdate || a.Type == deployer.ActionCreate {
							cmd.Printf("      -> %s\n", a.Desired.File.SourceInfo())
						}
					}
				}
				cmd.Printf("\nSummary: %d created, %d updated, %d deleted\n", result.FilesCreated, result.FilesUpdated, result.FilesDeleted)
			}

			// NixOS rebuild (se aplicável)
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
						nixResult, err := nix.RemoteBuild(ctx, ref, "whiterun", true)
						if err != nil {
							return fmt.Errorf("remote build failed: %w", err)
						}
						flake = nixResult
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
