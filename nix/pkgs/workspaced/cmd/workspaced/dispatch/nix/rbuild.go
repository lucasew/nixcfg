package nix

import (
	"crypto/rand"
	"fmt"
	"os"
	"strings"
	"workspaced/pkg/common"
	"workspaced/pkg/drivers/nix"
	"workspaced/pkg/drivers/notification"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		var target string
		var copyBack bool
		var useNom bool

		cmd := &cobra.Command{
			Use:   "rbuild <ref>",
			Short: "Performs a remote build of a Nix flake reference",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				ctx := cmd.Context()
				logger := common.GetLogger(ctx)
				ref := args[0]

				if target == "" {
					target = os.Getenv("NIX_RBUILD_TARGET")
					if target == "" {
						target = "whiterun"
					}
				}

				n := &notification.Notification{
					Title: "Nix Remote Build",
					Icon:  "nix-snowflake",
				}

				updateProgress := func(msg string, prog float64) {
					n.Message = msg
					n.Progress = prog
					n.Notify(ctx)
					logger.Info(msg, "progress", prog)
				}

				// 1. Resolve source
				updateProgress("Resolvendo metadados do flake...", 0.1)
				parts := strings.Split(ref, "#")
				repo := parts[0]
				item := ""
				if len(parts) > 1 {
					item = parts[1]
				}

				sourcePath, err := nix.ResolveFlakePath(ctx, repo)
				if err != nil {
					return err
				}

				// 2. Sync source to target
				updateProgress(fmt.Sprintf("Sincronizando fontes para %s...", target), 0.3)
				if err := nix.CopyClosure(ctx, target, sourcePath, nix.To); err != nil {
					return fmt.Errorf("failed to copy source to %s: %w", target, err)
				}

				// 3. Remote build
				updateProgress("Compilando no servidor remoto...", 0.6)
				remoteCache, err := nix.GetRemoteCacheDir(ctx, target)
				if err != nil {
					return fmt.Errorf("failed to get remote cache dir: %w", err)
				}

				buildID := make([]byte, 8)
				rand.Read(buildID)
				uuid := fmt.Sprintf("%x", buildID)
				outLink := fmt.Sprintf("%s/%s", remoteCache, uuid)

				buildCmd := "nix build"
				if useNom {
					buildCmd = "nom build"
				}

				safeRef := fmt.Sprintf("%s#%s", sourcePath, item)
				remoteArgs := []string{
					target, "-t",
					"mkdir", "-p", remoteCache, "&&",
					buildCmd, fmt.Sprintf("%q", safeRef), "--out-link", outLink, "--show-trace",
				}

				cmdBuild := common.RunCmd(ctx, "ssh", remoteArgs...)
				common.InheritContextWriters(ctx, cmdBuild)
				if err := cmdBuild.Run(); err != nil {
					return fmt.Errorf("remote build failed: %w", err)
				}

				// Get result path
				out, err := common.RunCmd(ctx, "ssh", target, "realpath", outLink).Output()
				if err != nil {
					return fmt.Errorf("failed to resolve result path: %w", err)
				}
				resultPath := strings.TrimSpace(string(out))

				// 4. Copy back
				if copyBack {
					updateProgress("Sincronizando resultado de volta...", 0.9)
					if err := nix.CopyClosure(ctx, target, resultPath, nix.From); err != nil {
						return fmt.Errorf("failed to copy result from %s: %w", target, err)
					}
				}

				updateProgress("Build concluído com sucesso.", 1.0)
				cmd.Println(resultPath)

				return nil
			},
		}

		cmd.Flags().StringVarP(&target, "target", "t", "", "Remote host to build on (default: whiterun)")
		cmd.Flags().BoolVar(&copyBack, "copy-back", true, "Copy result back to local store")
		cmd.Flags().BoolVar(&useNom, "nom", true, "Use nix-output-monitor (nom)")

		parent.AddCommand(cmd)
	})
}
