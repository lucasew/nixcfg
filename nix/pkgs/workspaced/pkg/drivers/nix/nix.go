package nix

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"workspaced/pkg/common"
	"workspaced/pkg/drivers/notification"
	"workspaced/pkg/drivers/sudo"
	"workspaced/pkg/types"
)

var buildCache sync.Map // key: sourcePath#attribute, value: resultPath

type Direction int

const (
	To Direction = iota
	From
)

func parseFlakeRef(ref string) (repo string, item string) {
	parts := strings.SplitN(ref, "#", 2)
	repo = parts[0]
	if len(parts) > 1 {
		item = parts[1]
	}
	return
}

func ResolveFlakePath(ctx context.Context, repo string) (string, error) {
	if repo == "" || repo == "." || repo == "," {
		root, err := common.GetDotfilesRoot()
		if err != nil {
			return "", err
		}
		repo = root
	}

	// Use nix flake archive to ensure the source is in the Nix store and get its path
	out, err := common.RunCmd(ctx, "nix", "flake", "archive", repo, "--json").Output()
	if err != nil {
		return "", fmt.Errorf("failed to archive flake %s to store: %w", repo, err)
	}

	var meta struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(out, &meta); err != nil {
		return "", fmt.Errorf("failed to parse flake archive output: %w", err)
	}

	return meta.Path, nil
}

func CopyClosure(ctx context.Context, target string, path string, direction Direction) error {
	args := []string{}
	if direction == To {
		args = append(args, "-s", "--to", target, path)
	} else {
		args = append(args, "--from", target, path)
	}

	cmd := common.RunCmd(ctx, "nix-copy-closure", args...)
	common.InheritContextWriters(ctx, cmd)
	return cmd.Run()
}

func GetRemoteCacheDir(ctx context.Context, target string) (string, error) {
	// Sentinel: Use XDG_RUNTIME_DIR (wiped on reboot) or fallback to user cache
	script := `echo "${XDG_RUNTIME_DIR:-${XDG_CACHE_HOME:-$HOME/.cache}}/rbuild-outputs"`
	out, err := common.RunCmd(ctx, "ssh", target, script).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func RemoteBuild(ctx context.Context, ref string, target string, copyBack bool) (string, error) {
	logger := common.GetLogger(ctx)

	if target == "" {
		target = os.Getenv("NIX_RBUILD_TARGET")
		if target == "" {
			target = "whiterun"
		}
	}

	n := &notification.Notification{
		Title: "Nix Remote Build",
		Icon:  "nix-snowflake",
	}

	updateProgress := func(msg string, prog float64) {
		n.Message = msg
		n.Progress = prog
		_ = n.Notify(ctx)
		logger.Info(msg, "progress", prog)
	}

	// 1. Resolve source
	updateProgress("Resolvendo metadados do flake...", 0.1)
	repo, item := parseFlakeRef(ref)

	sourcePath, err := ResolveFlakePath(ctx, repo)
	if err != nil {
		return "", err
	}

	// 2. Sync source to target
	updateProgress(fmt.Sprintf("Sincronizando fontes para %s...", target), 0.3)
	if err := CopyClosure(ctx, target, sourcePath, To); err != nil {
		return "", fmt.Errorf("failed to copy source to %s: %w", target, err)
	}

	// 3. Remote build
	updateProgress("Compilando no servidor remoto...", 0.6)
	remoteCache, err := GetRemoteCacheDir(ctx, target)
	if err != nil {
		return "", fmt.Errorf("failed to get remote cache dir: %w", err)
	}

	buildID := make([]byte, 8)
	_, _ = rand.Read(buildID)
	uuid := fmt.Sprintf("%x", buildID)
	outLink := fmt.Sprintf("%s/%s", remoteCache, uuid)

	buildCmd := "nix build"
	// if useNom {
	// 	buildCmd = "nom build"
	// }

	safeRef := fmt.Sprintf("%s#%s", sourcePath, item)
	remoteArgs := []string{
		target, "-t",
		"mkdir", "-p", remoteCache, "&&",
		buildCmd, fmt.Sprintf("%q", safeRef), "--out-link", outLink, "--show-trace",
	}

	cmdBuild := common.RunCmd(ctx, "ssh", remoteArgs...)
	common.InheritContextWriters(ctx, cmdBuild)
	if err := cmdBuild.Run(); err != nil {
		return "", fmt.Errorf("remote build failed: %w", err)
	}

	// Get result path
	out, err := common.RunCmd(ctx, "ssh", target, "realpath", outLink).Output()
	if err != nil {
		return "", fmt.Errorf("failed to resolve result path: %w", err)
	}
	resultPath := strings.TrimSpace(string(out))

	// 4. Copy back
	if copyBack {
		updateProgress("Sincronizando resultado de volta...", 0.9)
		if err := CopyClosure(ctx, target, resultPath, From); err != nil {
			return "", fmt.Errorf("failed to copy result from %s: %w", target, err)
		}
	}

	updateProgress("Build concluído com sucesso.", 1.0)
	return resultPath, nil
}

func Build(ctx context.Context, ref string, useCache bool) (string, error) {
	logger := common.GetLogger(ctx)

	repo, item := parseFlakeRef(ref)

	// Resolve the source path to a store path
	sourcePath, err := ResolveFlakePath(ctx, repo)
	if err != nil {
		return "", err
	}

	cacheKey := fmt.Sprintf("%s#%s", sourcePath, item)
	if useCache {
		if val, ok := buildCache.Load(cacheKey); ok {
			resultPath := val.(string)
			if _, err := os.Stat(resultPath); err == nil {
				logger.Debug("build cache hit", "ref", ref, "path", resultPath)
				return resultPath, nil
			}
			buildCache.Delete(cacheKey)
		}
	}

	logger.Info("performing nix build", "ref", ref)
	// We use the store path of the source to ensure deterministic build and avoid re-evaluation if not needed
	cmd := common.RunCmd(ctx, "nix", "build", fmt.Sprintf("%s#%s", sourcePath, item), "--no-link", "--print-out-paths")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("nix build failed: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	resultPath := lines[0]
	// If multiple paths, try to find the one with bin/
	for _, line := range lines {
		if info, err := os.Stat(filepath.Join(line, "bin")); err == nil && info.IsDir() {
			resultPath = line
			break
		}
	}

	if useCache {
		buildCache.Store(cacheKey, resultPath)
	}

	return resultPath, nil
}

func Rebuild(ctx context.Context, action string, flake string) error {
	hostname := common.GetHostname()
	if flake == "" || flake == "." || flake == "," {
		root, err := common.GetDotfilesRoot()
		if err != nil {
			return err
		}
		flake = root
	}

	if common.IsInStore() {
		flake = "github:lucasew/nixcfg"
	}

	// Check if we are on a known node
	supportedNodes := []string{"riverwood", "whiterun", "ravenrock", "atomicpi", "recovery"}
	isSupported := false
	for _, node := range supportedNodes {
		if hostname == node {
			isSupported = true
			break
		}
	}

	if !isSupported {
		return fmt.Errorf("hostname %s is not a supported NixOS node for rebuild", hostname)
	}

	var target string
	if strings.HasPrefix(flake, "/nix/store/") {
		target = flake
	} else {
		target = fmt.Sprintf("%s#%s", flake, hostname)
	}
	args := []string{action, "--flake", target}

	cmdName := "nixos-rebuild"

	if os.Getuid() != 0 {
		isDaemon := false
		if val := ctx.Value(types.DaemonModeKey); val != nil {
			isDaemon = val.(bool)
		}
		if isDaemon {
			return sudo.Enqueue(ctx, &types.SudoCommand{
				Slug:    "rebuild",
				Command: cmdName,
				Args:    args,
			})
		}
		cmd := common.RunCmd(ctx, "sudo", append([]string{cmdName}, args...)...)
		common.InheritContextWriters(ctx, cmd)
		return cmd.Run()
	} else {
		cmd := common.RunCmd(ctx, cmdName, args...)
		common.InheritContextWriters(ctx, cmd)
		return cmd.Run()
	}
}

func HomeManagerSwitch(ctx context.Context, action string, flake string) error {
	if flake == "" || flake == "." || flake == "," {
		root, err := common.GetDotfilesRoot()
		if err != nil {
			return err
		}
		flake = root
	}

	if common.IsInStore() {
		flake = "github:lucasew/nixcfg"
	}

	if strings.HasPrefix(flake, "/nix/store/") {
		activatePath := filepath.Join(flake, "activate")
		if _, err := os.Stat(activatePath); err == nil {
			cmd := common.RunCmd(ctx, activatePath)
			common.InheritContextWriters(ctx, cmd)
			return cmd.Run()
		}
	}

	target := fmt.Sprintf("%s#main", flake)
	args := []string{action, "--flake", target}

	cmd := common.RunCmd(ctx, "home-manager", args...)
	common.InheritContextWriters(ctx, cmd)
	return cmd.Run()
}

func GetFlakeOutput(ctx context.Context, flake, output string) (string, error) {
	cmd := common.RunCmd(ctx, "nix", "build", fmt.Sprintf("%s#%s", flake, output), "--no-link", "--print-out-paths")
	if stderr, ok := ctx.Value(types.StderrKey).(io.Writer); ok {
		cmd.Stderr = stderr
	}
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
