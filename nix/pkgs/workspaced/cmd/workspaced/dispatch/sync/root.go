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

			// 2. Determine command to run
			var shimArgs []string
			var actionMsg string
			if rebuildOnly {
				shimArgs = []string{"--version"}
				actionMsg = "==> Rebuilding only..."
			} else {
				shimArgs = []string{"dispatch", "apply"}
				actionMsg = "==> Rebuilding and applying..."
			}

			// 3. Execute rebuild (and optionally apply)
			fmt.Println(actionMsg)
			bashPath, err := exec.Which(ctx, "bash")
			if err != nil {
				return fmt.Errorf("bash not found: %w", err)
			}
			shimPath := filepath.Join(root, "bin/shim/workspaced")
			shimCmd := exec.RunCmd(ctx, bashPath, append([]string{shimPath}, shimArgs...)...)
			shimCmd.Env = append(os.Environ(), "WORKSPACED_REFRESH=1")
			shimCmd.Stdout = os.Stdout
			shimCmd.Stderr = os.Stderr
			if err := shimCmd.Run(); err != nil {
				return fmt.Errorf("command failed: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&rebuildOnly, "rebuild-only", false, "Only rebuild, skip apply")
	return cmd
}
