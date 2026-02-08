package nix

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"workspaced/pkg/nix"
	"workspaced/pkg/exec"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		cmd := &cobra.Command{
			Use:                "run <ref> [args...]",
			Short:              "Builds a package with RAM caching and runs it locally",
			Args:               cobra.MinimumNArgs(1),
			DisableFlagParsing: true,
			RunE: func(cmd *cobra.Command, args []string) error {
				if len(args) == 0 {
					return fmt.Errorf("no flake reference provided")
				}
				ctx := cmd.Context()
				ref := args[0]
				runArgs := args[1:]

				if len(runArgs) > 0 && runArgs[0] == "--" {
					runArgs = runArgs[1:]
				}

				// Parse ref: repo#item/binary
				parts := strings.Split(ref, "#")
				repo := parts[0]
				item := ""
				if len(parts) > 1 {
					item = parts[1]
				}

				binary := ""
				if strings.Contains(item, "/") {
					itemParts := strings.Split(item, "/")
					item = itemParts[0]
					binary = itemParts[1]
				}

				// Build with cache
				resultPath, err := nix.Build(ctx, repo+"#"+item, true)
				if err != nil {
					return err
				}

				// Find binary
				binDir := filepath.Join(resultPath, "bin")
				if binary == "" {
					entries, err := os.ReadDir(binDir)
					if err != nil || len(entries) == 0 {
						return fmt.Errorf("no binary found in %s", binDir)
					}
					binary = entries[0].Name()
				}

				binPath := filepath.Join(binDir, binary)
				if _, err := os.Stat(binPath); err != nil {
					entries, _ := os.ReadDir(binDir)
					for _, entry := range entries {
						if strings.Contains(entry.Name(), binary) {
							binPath = filepath.Join(binDir, entry.Name())
							break
						}
					}
				}

				// Run
				ec := exec.RunCmd(ctx, binPath, runArgs...)
				exec.InheritContextWriters(ctx, ec)
				ec.Stdin = os.Stdin
				return ec.Run()
			},
		}
		parent.AddCommand(cmd)
	})
}
