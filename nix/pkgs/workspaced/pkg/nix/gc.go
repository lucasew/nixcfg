package nix

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"workspaced/pkg/sudo"
	"workspaced/pkg/exec"
	"workspaced/pkg/logging"
	"workspaced/pkg/types"
)

func CleanupProfiles(ctx context.Context) error {
	logger := logging.GetLogger(ctx)
	logger.Info("Searching for old Nix profiles to cleanup...")

	baseDir := "/nix/var/nix/profiles"

	dirsToScan := []string{
		baseDir,
		filepath.Join(baseDir, "per-user/root"),
		filepath.Join(baseDir, "per-user", os.Getenv("USER")),
	}

	if user := os.Getenv("USER"); user == "" {
		// Fallback if USER env is not set
		if home, err := os.UserHomeDir(); err == nil {
			dirsToScan = append(dirsToScan, filepath.Join(baseDir, "per-user", filepath.Base(home)))
		}
	}

	var filesToRemove []string

	for _, dir := range dirsToScan {
		if _, err := os.Stat(dir); err != nil {
			continue
		}

		entries, err := os.ReadDir(dir)
		if err != nil {
			logger.Error("failed to read directory", "dir", dir, "error", err)
			continue
		}

		// Find "master" links (e.g., 'system', 'profile')
		masterLinks := make(map[string]string)
		for _, entry := range entries {
			if entry.Type()&os.ModeSymlink != 0 {
				name := entry.Name()
				// Profiles usually don't have hyphens in the base name,
				// or they are specifically 'system', 'default', 'profile', 'home-manager'
				if !strings.Contains(name, "-") {
					target, err := os.Readlink(filepath.Join(dir, name))
					if err == nil {
						masterLinks[name] = target
					}
				}
			}
		}

		for master, target := range masterLinks {
			prefix := master + "-"
			for _, entry := range entries {
				name := entry.Name()
				if strings.HasPrefix(name, prefix) && strings.HasSuffix(name, "-link") {
					if name == target {
						continue
					}
					filesToRemove = append(filesToRemove, filepath.Join(dir, name))
				}
			}
		}
	}

	if len(filesToRemove) == 0 {
		logger.Info("No old profiles found to cleanup.")
		return nil
	}

	sort.Strings(filesToRemove)

	logger.Info(fmt.Sprintf("Found %d old profile links to remove.", len(filesToRemove)))

	if os.Getuid() != 0 {
		return sudo.Enqueue(ctx, &types.SudoCommand{
			Slug:    "nix-gc-cleanup",
			Command: "rm",
			Args:    filesToRemove,
		})
	} else {
		// If running interactively, we can just run sudo directly
		cmd := exec.RunCmd(ctx, "rm", filesToRemove...)
		exec.InheritContextWriters(ctx, cmd)
		return cmd.Run()
	}
}
