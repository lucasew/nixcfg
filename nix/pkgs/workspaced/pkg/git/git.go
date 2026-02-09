package git

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"workspaced/pkg/config"
	"workspaced/pkg/exec"
	"workspaced/pkg/logging"
	"workspaced/pkg/notification"
)

func QuickSync(ctx context.Context) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	logger := logging.GetLogger(ctx)
	repoDir := cfg.QuickSync.RepoDir
	entries, err := os.ReadDir(repoDir)
	if err != nil {
		return fmt.Errorf("failed to read repo dir %s: %w", repoDir, err)
	}

	var repos []string
	for _, entry := range entries {
		if entry.IsDir() {
			repoPath := filepath.Join(repoDir, entry.Name())
			if _, err := os.Stat(filepath.Join(repoPath, ".git")); err == nil {
				repos = append(repos, entry.Name())
			}
		}
	}

	total := len(repos)
	if total == 0 {
		return nil
	}

	n := &notification.Notification{
		Title: "Sincronização Git",
		Icon:  "git",
	}

	for i, repoName := range repos {
		// Check for context cancellation
		if err := ctx.Err(); err != nil {
			return err
		}

		repoPath := filepath.Join(repoDir, repoName)
		n.Message = fmt.Sprintf("Sincronizando %s...", repoName)
		n.Progress = float64(i) / float64(total)
		_ = n.Notify(ctx)

		logger.Info("syncing repository", "repo", repoName)
		if err := SyncRepo(ctx, repoPath); err != nil {
			logger.Error("failed to sync repo", "repo", repoName, "error", err)
			errN := &notification.Notification{
				Title:   "Sincronização Falhou",
				Message: fmt.Sprintf("Conflito ou erro em %s. Intervenção manual necessária.", repoName),
				Urgency: "critical",
				Icon:    "dialog-warning",
			}
			_ = errN.Notify(ctx)
		}
	}

	n.Message = "Sincronização concluída."
	n.Progress = 1.0
	_ = n.Notify(ctx)

	return nil
}

func SyncRepo(ctx context.Context, path string) error {
	hostname, _ := os.Hostname()
	logger := logging.GetLogger(ctx)

	// git add -A
	logger.Info("git add", "path", path)
	if err := exec.RunCmd(ctx, "git", "-C", path, "add", "-A").Run(); err != nil {
		return fmt.Errorf("git add failed: %w", err)
	}

	// git commit -sm "backup checkpoint <host>"
	// Check if there are changes to commit
	if err := exec.RunCmd(ctx, "git", "-C", path, "diff-index", "HEAD", "--exit-code").Run(); err != nil {
		commitMsg := fmt.Sprintf("backup checkpoint %s", hostname)
		logger.Info("git commit", "path", path, "msg", commitMsg)
		if err := exec.RunCmd(ctx, "git", "-C", path, "commit", "-sm", commitMsg).Run(); err != nil {
			return fmt.Errorf("git commit failed: %w", err)
		}
	}

	// git pull --rebase
	logger.Info("git pull --rebase", "path", path)
	if err := exec.RunCmd(ctx, "git", "-C", path, "pull", "--rebase").Run(); err != nil {
		_ = exec.RunCmd(ctx, "git", "-C", path, "rebase", "--abort").Run()
		return fmt.Errorf("git pull rebase failed (conflict?): %w", err)
	}

	// git push
	logger.Info("git push", "path", path)
	if err := exec.RunCmd(ctx, "git", "-C", path, "push").Run(); err != nil {
		return fmt.Errorf("git push failed: %w", err)
	}

	return nil
}
