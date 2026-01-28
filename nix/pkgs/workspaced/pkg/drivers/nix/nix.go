package nix

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"workspaced/pkg/common"
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

	return common.RunCmd(ctx, "nix-copy-closure", args...).Run()
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
