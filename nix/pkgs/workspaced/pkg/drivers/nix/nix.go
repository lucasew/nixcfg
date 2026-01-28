package nix

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"workspaced/pkg/common"
	"workspaced/pkg/types"
)

type Direction int

const (
	To Direction = iota
	From
)

func ResolveFlakePath(ctx context.Context, repo string) (string, error) {
	if repo == "" || repo == "." || repo == "," {
		root, err := common.GetDotfilesRoot()
		if err != nil {
			return "", err
		}
		repo = root
	}

	out, err := common.RunCmd(ctx, "nix", "flake", "metadata", repo, "--json").Output()
	if err != nil {
		return "", fmt.Errorf("failed to resolve flake metadata for %s: %w", repo, err)
	}

	var meta struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(out, &meta); err != nil {
		return "", fmt.Errorf("failed to parse flake metadata: %w", err)
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

	target := fmt.Sprintf("%s#%s", flake, hostname)
	args := []string{action, "--flake", target}

	cmdName := "nixos-rebuild"
	var cmd *exec.Cmd
	if os.Getuid() != 0 {
		cmd = common.RunCmd(ctx, "sudo", append([]string{cmdName}, args...)...)
	} else {
		cmd = common.RunCmd(ctx, cmdName, args...)
	}
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
