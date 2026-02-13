package source

import (
	"context"
	"path/filepath"
	"sort"
	"strings"
	"workspaced/pkg/template"
)

// DotDProcessorPlugin processa diretórios .d.tmpl (concatenação)
type DotDProcessorPlugin struct {
	engine *template.Engine
	data   interface{}
}

// NewDotDProcessorPlugin cria plugin de processamento .d.tmpl
func NewDotDProcessorPlugin(engine *template.Engine, data interface{}) *DotDProcessorPlugin {
	return &DotDProcessorPlugin{
		engine: engine,
		data:   data,
	}
}

func (p *DotDProcessorPlugin) Name() string {
	return "dotd-processor"
}

func (p *DotDProcessorPlugin) Process(ctx context.Context, files []File) ([]File, error) {
	result := []File{}
	dotdGroups := make(map[string][]File) // Key: Final Target Path, Value: component files

	// 1. Group files by their intended concatenated target
	for _, f := range files {
		// A file is part of a DotD if its RelPath contains ".d.tmpl/"
		rel := f.RelPath()
		if strings.Contains(rel, ".d.tmpl/") || strings.Contains(rel, ".d.tmpl"+string(filepath.Separator)) {
			// Extract the target path by removing the ".d.tmpl" part and the filename
			parts := strings.Split(rel, ".d.tmpl")
			targetRelPath := parts[0] // e.g. ".bashrc" if RelPath was ".bashrc.d.tmpl/00-base.sh"
			finalTargetPath := filepath.Join(f.TargetBase(), targetRelPath)
			dotdGroups[finalTargetPath] = append(dotdGroups[finalTargetPath], f)
		} else {
			result = append(result, f)
		}
	}

	// 2. Process each group
	// To keep it deterministic, sort targetPaths
	targets := make([]string, 0, len(dotdGroups))
	for t := range dotdGroups {
		targets = append(targets, t)
	}
	sort.Strings(targets)

	for _, targetPath := range targets {
		groupFiles := dotdGroups[targetPath]
		// Sort components by RelPath to ensure deterministic concatenation
		sort.Slice(groupFiles, func(i, j int) bool {
			return groupFiles[i].RelPath() < groupFiles[j].RelPath()
		})

		// Use info from the first file for bases
		first := groupFiles[0]
		relPath, _ := filepath.Rel(first.TargetBase(), targetPath)

		result = append(result, &ConcatenatedFile{
			BasicFile: BasicFile{
				RelPathStr:    relPath,
				TargetBaseDir: first.TargetBase(),
				FileMode:      0644,
				Info:          "concatenated:" + targetPath,
				FileType:      TypeDotD,
			},
			Components: groupFiles,
		})
	}

	return result, nil
}
