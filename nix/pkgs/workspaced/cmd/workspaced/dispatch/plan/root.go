package plan

import (
	"fmt"
	"os"
	"path/filepath"
	"workspaced/pkg/env"
	"workspaced/pkg/exec"

	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plan",
		Short: "Rebuild and show what would be applied",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			root, err := env.GetDotfilesRoot()
			if err != nil {
				return fmt.Errorf("failed to get dotfiles root: %w", err)
			}

			// 1. Determine command to run
			shimArgs := []string{"dispatch", "apply", "--dry-run"}
			fmt.Println("==> Rebuilding and planning...")

			// 2. Execute rebuild and dry-run apply
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

	return cmd
}
