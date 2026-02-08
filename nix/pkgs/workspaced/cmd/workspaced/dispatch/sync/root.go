package sync

import (
	"fmt"
	"os"
	"path/filepath"
	"workspaced/pkg/env"
	"workspaced/pkg/exec"

	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	var rebuildOnly bool

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Pull dotfiles changes and apply them",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			root, err := env.GetDotfilesRoot()
			if err != nil {
				return fmt.Errorf("failed to get dotfiles root: %w", err)
			}

			// 1. Git pull
			fmt.Println("==> Pulling dotfiles changes...")
			pullCmd := exec.RunCmd(ctx, "git", "-C", root, "pull")
			pullCmd.Stdout = os.Stdout
			pullCmd.Stderr = os.Stderr
			if err := pullCmd.Run(); err != nil {
				return fmt.Errorf("git pull failed: %w", err)
			}

			// 2. Rebuild (and optionally apply)
			if rebuildOnly {
				fmt.Println("==> Rebuilding only...")
				bashPath, err := exec.Which(ctx, "bash")
				if err != nil {
					return fmt.Errorf("bash not found: %w", err)
				}
				shimPath := filepath.Join(root, "bin/shim/workspaced")
				rebuildCmd := exec.RunCmd(ctx, bashPath, shimPath, "--version")
				rebuildCmd.Env = append(os.Environ(), "WORKSPACED_REFRESH=1")
				rebuildCmd.Stdout = os.Stdout
				rebuildCmd.Stderr = os.Stderr
				if err := rebuildCmd.Run(); err != nil {
					return fmt.Errorf("rebuild failed: %w", err)
				}
			} else {
				fmt.Println("==> Rebuilding and applying...")
				bashPath, err := exec.Which(ctx, "bash")
				if err != nil {
					return fmt.Errorf("bash not found: %w", err)
				}
				shimPath := filepath.Join(root, "bin/shim/workspaced")
				applyCmd := exec.RunCmd(ctx, bashPath, shimPath, "dispatch", "apply")
				applyCmd.Env = append(os.Environ(), "WORKSPACED_REFRESH=1")
				applyCmd.Stdout = os.Stdout
				applyCmd.Stderr = os.Stderr
				if err := applyCmd.Run(); err != nil {
					return fmt.Errorf("rebuild/apply failed: %w", err)
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&rebuildOnly, "rebuild-only", false, "Only rebuild, skip apply")
	return cmd
}
