package apply

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"workspaced/pkg/common"
	"workspaced/pkg/config"
)

type LazyShimProvider struct{}

func (p *LazyShimProvider) Name() string {
	return "lazyshim"
}

func (p *LazyShimProvider) GetDesiredState(ctx context.Context) ([]DesiredState, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	dataDir, err := common.GetUserDataDir()
	if err != nil {
		return nil, err
	}

	shimDir := filepath.Join(dataDir, "shim")
	globalDir := filepath.Join(shimDir, "global")
	lazyDir := filepath.Join(shimDir, "lazy")

	for _, dir := range []string{globalDir, lazyDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}

	desired := []DesiredState{}

	// 1. Generate 'x' shim
	xContent := fmt.Sprintf(`#!/usr/bin/env bash
export PATH="%s:$PATH"
exec "$@"
`, lazyDir)
	xSource, err := p.materialize(xContent)
	if err != nil {
		return nil, err
	}
	desired = append(desired, DesiredState{
		Target: filepath.Join(globalDir, "x"),
		Source: xSource,
		Mode:   0755,
	})

	// 2. Generate tool shims
	for name, tool := range cfg.LazyTools {
		pkg := tool.Pkg
		if pkg == "" {
			pkg = name
		}

		bins := tool.Bins
		if len(bins) == 0 {
			binName := tool.Alias
			if binName == "" {
				// Clean up name (e.g. github:owner/repo -> repo)
				parts := strings.Split(name, ":")
				binName = parts[len(parts)-1]
				parts = strings.Split(binName, "/")
				binName = parts[len(parts)-1]
			}
			bins = []string{binName}
		}

		for _, bin := range bins {
			content := fmt.Sprintf(`#!/usr/bin/env bash
exec mise exec %s@%s -- %s "$@"
`, pkg, tool.Version, bin)
			source, err := p.materialize(content)
			if err != nil {
				return nil, err
			}

			// Always in lazy
			desired = append(desired, DesiredState{
				Target: filepath.Join(lazyDir, bin),
				Source: source,
				Mode:   0755,
			})

			// If global, also in global
			if tool.Global {
				// Global shim calls 'x' to ensure lazy shims are available if the tool needs them
				globalContent := fmt.Sprintf(`#!/usr/bin/env bash
exec x %s "$@"
`, bin)
				globalSource, err := p.materialize(globalContent)
				if err != nil {
					return nil, err
				}
				desired = append(desired, DesiredState{
					Target: filepath.Join(globalDir, bin),
					Source: globalSource,
					Mode:   0755,
				})
			}
		}
	}

	return desired, nil
}

func (p *LazyShimProvider) materialize(content string) (string, error) {
	cw, err := common.NewCASWriter()
	if err != nil {
		return "", err
	}
	if _, err := fmt.Fprint(cw, content); err != nil {
		return "", err
	}
	return cw.Seal()
}
