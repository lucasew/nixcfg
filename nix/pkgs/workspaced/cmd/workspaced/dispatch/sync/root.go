package sync

import (
	"fmt"
	"os"
	"os/exec"
	"workspaced/pkg/common"

	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Pull dotfiles changes and apply them",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			root, err := common.GetDotfilesRoot()
			if err != nil {
				return fmt.Errorf("failed to get dotfiles root: %w", err)
			}

			// Find git binary without using os/exec.LookPath to avoid SIGSYS on Android/Go 1.24
			gitPath, err := common.Which(ctx, "git")
			if err != nil {
				return fmt.Errorf("git not found in PATH: %w", err)
			}

			// 1. Git pull
			fmt.Println("==> Pulling dotfiles changes...")
			pullCmd := exec.CommandContext(ctx, gitPath, "-C", root, "pull")
			pullCmd.Stdout = os.Stdout
			pullCmd.Stderr = os.Stderr
			if err := pullCmd.Run(); err != nil {
				return fmt.Errorf("git pull failed: %w", err)
			}

			// 2. Workspaced dispatch apply
			fmt.Println("==> Applying configuration...")
			self, err := os.Executable()
			if err != nil {
				self = "workspaced"
			}
			applyCmd := exec.CommandContext(ctx, self, "dispatch", "apply")
			applyCmd.Stdout = os.Stdout
			applyCmd.Stderr = os.Stderr
			if err := applyCmd.Run(); err != nil {
				return fmt.Errorf("apply failed: %w", err)
			}

			return nil
		},
	}
	return cmd
}
