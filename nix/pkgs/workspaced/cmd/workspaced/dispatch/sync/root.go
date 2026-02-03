package sync

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
			pullCmd := common.RunCmd(ctx, "git", "-C", root, "pull")
			pullCmd.Stdout = os.Stdout
			pullCmd.Stderr = os.Stderr
			if err := pullCmd.Run(); err != nil {
				return fmt.Errorf("git pull failed: %w", err)
			}

			// 2. Rebuild workspaced unconditionally
			fmt.Println("==> Rebuilding workspaced...")
			if err := rebuildWorkspaced(ctx, root); err != nil {
				fmt.Printf("Warning: workspaced rebuild failed: %v\n", err)
				// Non-fatal, continue
			}

			// 3. Workspaced dispatch apply (using newly built binary)
			fmt.Println("==> Applying configuration...")
			newBinary := filepath.Join(os.Getenv("HOME"), ".local/share/workspaced/bin/workspaced")
			applyCmd := exec.CommandContext(ctx, newBinary, "dispatch", "apply")
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

// rebuildWorkspaced rebuilds workspaced from source unconditionally
func rebuildWorkspaced(ctx context.Context, dotfilesRoot string) error {
	sourceDir := filepath.Join(dotfilesRoot, "nix/pkgs/workspaced")

	// Build workspaced using mise exec, matching the shim logic
	buildCmd := exec.CommandContext(ctx, "bash", "-c", fmt.Sprintf(`
		SOURCE_DIR="%s"
		GO_VERSION="$(sed -n 's/^go = "\(.*\)"/\1/p' "$SOURCE_DIR/mise.toml")"

		# Activate mise so that mise exec works
		if [ -f "$HOME/.local/bin/mise" ]; then
			eval "$("$HOME"/.local/bin/mise activate bash)"
		fi

		mkdir -p ~/.local/share/workspaced/bin
		cd "$SOURCE_DIR" || exit 1

		BUILD_ID="$(date +%%s)"
		mise exec "go@$GO_VERSION" -- env CGO_ENABLED=0 go build -v -ldflags "-X workspaced/pkg/common.BuildID=$BUILD_ID" -o ~/.local/share/workspaced/bin/workspaced ./cmd/workspaced
	`, sourceDir))

	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	return buildCmd.Run()
}
