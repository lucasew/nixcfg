package template

import (
	"fmt"
	"os"
	"path/filepath"
	"workspaced/pkg/config"
	"workspaced/pkg/deployer"
	"workspaced/pkg/source"
	"workspaced/pkg/template"

	"github.com/spf13/cobra"
)

func getMaterializeCommand() *cobra.Command {
	var configPaths []string
	var sourcePaths []string
	var targetDir string

	cmd := &cobra.Command{
		Use:   "materialize",
		Short: "Materialize templates into a directory (low-level)",
		RunE: func(c *cobra.Command, args []string) error {
			ctx := c.Context()

			if targetDir == "" {
				return fmt.Errorf("--target is required")
			}

			// 1. Merge configs strictly
			cfg, err := config.LoadFiles(configPaths)
			if err != nil {
				return err
			}

			// 2. Setup pipeline
			engine := template.NewEngine(ctx)
			tempDir, err := os.MkdirTemp("", "workspaced-materialize-*")
			if err != nil {
				return err
			}
			defer os.RemoveAll(tempDir)

			pipeline := source.NewPipeline()

			// Add scanners for each source
			for i, srcPath := range sourcePaths {
				absSrc, err := filepath.Abs(srcPath)
				if err != nil {
					return err
				}
				scanner, err := source.NewScannerPlugin(source.ScannerConfig{
					Name:       fmt.Sprintf("source-%d", i),
					BaseDir:    absSrc,
					TargetBase: targetDir,
					Priority:   100,
				})
				if err != nil {
					return err
				}
				pipeline.AddPlugin(scanner)
			}

			// Add processors
			templatePlugin, _ := source.NewTemplateExpanderPlugin(engine, cfg, filepath.Join(tempDir, "templates"))
			pipeline.AddPlugin(templatePlugin)

			dotdPlugin, _ := source.NewDotDProcessorPlugin(engine, cfg, filepath.Join(tempDir, "dotd"))
			pipeline.AddPlugin(dotdPlugin)

			// Add strict conflict resolver
			pipeline.AddPlugin(source.NewStrictConflictResolverPlugin())

			// 3. Process
			files, err := pipeline.Run(ctx, nil)
			if err != nil {
				return err
			}

			// 4. Write to target
			executor := deployer.NewExecutor()
			actions := []deployer.Action{}
			for _, f := range files {
				actions = append(actions, deployer.Action{
					Type:   deployer.ActionCreate,
					Target: filepath.Join(f.TargetBase, f.RelPath),
					Desired: deployer.DesiredState{
						Source: filepath.Join(f.SourceBase, f.RelPath),
						Mode:   f.Mode,
					},
				})
			}

			// Execute without state tracking
			return executor.Execute(ctx, actions, &deployer.State{Files: make(map[string]deployer.ManagedInfo)})
		},
	}

	cmd.Flags().StringSliceVarP(&configPaths, "config", "c", nil, "Configuration file(s) to merge")
	cmd.Flags().StringSliceVarP(&sourcePaths, "source", "s", nil, "Source directory(ies) to scan")
	cmd.Flags().StringVarP(&targetDir, "target", "t", "", "Target directory to materialize files into")

	return cmd
}
