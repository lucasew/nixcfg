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

			// 1. Git pull
			fmt.Println("==> Pulling dotfiles changes...")
			pullCmd := exec.CommandContext(ctx, "git", "-C", root, "pull", "--really")
			pullCmd.Stdout = os.Stdout
			pullCmd.Stderr = os.Stderr
			if err := pullCmd.Run(); err != nil {
				return fmt.Errorf("git pull failed: %w", err)
			}

			// 2. Workspaced dispatch apply
			fmt.Println("==> Applying configuration...")
			applyCmd := exec.CommandContext(ctx, "workspaced", "dispatch", "apply")
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
