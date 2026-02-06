package sync

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	pkgExec "workspaced/pkg/exec"
	"workspaced/pkg/host"

	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Pull dotfiles changes and apply them",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			root, err := host.GetDotfilesRoot()
			if err != nil {
				return fmt.Errorf("failed to get dotfiles root: %w", err)
			}

			// 1. Git pull
			fmt.Println("==> Pulling dotfiles changes...")
			pullCmd := pkgExec.RunCmd(ctx, "git", "-C", root, "pull")
			pullCmd.Stdout = os.Stdout
			pullCmd.Stderr = os.Stderr
			if err := pullCmd.Run(); err != nil {
				return fmt.Errorf("git pull failed: %w", err)
			}

			// 2. Rebuild via shim and apply (WORKSPACED_REFRESH triggers rebuild, shim execs the new binary)
			fmt.Println("==> Rebuilding and applying...")
			bashPath, err := pkgExec.Which(ctx, "bash")
			if err != nil {
				return fmt.Errorf("bash not found: %w", err)
			}
			shimPath := filepath.Join(root, "bin/shim/workspaced")
			applyCmd := exec.CommandContext(ctx, bashPath, shimPath, "dispatch", "apply")
			applyCmd.Env = append(os.Environ(), "WORKSPACED_REFRESH=1")
			applyCmd.Stdout = os.Stdout
			applyCmd.Stderr = os.Stderr
			if err := applyCmd.Run(); err != nil {
				return fmt.Errorf("rebuild/apply failed: %w", err)
			}

			return nil
		},
	}
	return cmd
}
