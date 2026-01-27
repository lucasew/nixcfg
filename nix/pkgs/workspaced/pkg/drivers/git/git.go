package git

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"workspaced/pkg/common"
	"workspaced/pkg/drivers/notification"
)

func QuickSync(ctx context.Context) error {
	config, err := common.LoadConfig()
	if err != nil {
		return err
	}

	logger := common.GetLogger(ctx)
	repoDir := config.QuickSync.RepoDir
	entries, err := os.ReadDir(repoDir)
	if err != nil {
		return fmt.Errorf("failed to read repo dir %s: %w", repoDir, err)
	}

	for _, entry := range entries {
		// Check for context cancellation
		if err := ctx.Err(); err != nil {
			return err
		}

		if entry.IsDir() {
			repoPath := filepath.Join(repoDir, entry.Name())
			if _, err := os.Stat(filepath.Join(repoPath, ".git")); err == nil {
				logger.Info("syncing repository", "repo", entry.Name())
				if err := SyncRepo(ctx, repoPath); err != nil {
					logger.Error("failed to sync repo", "repo", entry.Name(), "error", err)
					n := &notification.Notification{
						Title:   "Sincronização Falhou",
						Message: fmt.Sprintf("Conflito ou erro em %s. Intervenção manual necessária.", entry.Name()),
						Urgency: "critical",
						Icon:    "dialog-warning",
					}
					n.Notify(ctx)
				}
			}
		}
	}

	return nil
}

func SyncRepo(ctx context.Context, path string) error {
	hostname, _ := os.Hostname()
	logger := common.GetLogger(ctx)

	// git add -A
	logger.Info("git add", "path", path)
	if err := common.RunCmd(ctx, "git", "-C", path, "add", "-A").Run(); err != nil {
		return fmt.Errorf("git add failed: %w", err)
	}

	// git commit -sm "backup checkpoint <host>"
	// Check if there are changes to commit
	if err := common.RunCmd(ctx, "git", "-C", path, "diff-index", "HEAD", "--exit-code").Run(); err != nil {
		commitMsg := fmt.Sprintf("backup checkpoint %s", hostname)
		logger.Info("git commit", "path", path, "msg", commitMsg)
		if err := common.RunCmd(ctx, "git", "-C", path, "commit", "-sm", commitMsg).Run(); err != nil {
			return fmt.Errorf("git commit failed: %w", err)
		}
	}

	// git pull --rebase
	logger.Info("git pull --rebase", "path", path)
	if err := common.RunCmd(ctx, "git", "-C", path, "pull", "--rebase").Run(); err != nil {
		common.RunCmd(ctx, "git", "-C", path, "rebase", "--abort").Run()
		return fmt.Errorf("git pull rebase failed (conflict?): %w", err)
	}

	// git push
	logger.Info("git push", "path", path)
	if err := common.RunCmd(ctx, "git", "-C", path, "push").Run(); err != nil {
		return fmt.Errorf("git push failed: %w", err)
	}

	return nil
}
