package template

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ProcessDotD processa um diretório .d.tmpl (concatenação de arquivos)
func (e *Engine) ProcessDotD(ctx context.Context, dirPath string, data interface{}) ([]byte, error) {
	// Check if directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return nil, nil // Empty content if directory doesn't exist
	}

	// Read directory entries
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	// Sort entries alphabetically (ReadDir already returns sorted)
	var fileNames []string
	for _, entry := range entries {
		if !entry.IsDir() {
			fileNames = append(fileNames, entry.Name())
		}
	}

	var result bytes.Buffer
	for _, fileName := range fileNames {
		filePath := filepath.Join(dirPath, fileName)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", fileName, err)
		}

		// Check if file is a template (ends with .tmpl)
		if strings.HasSuffix(fileName, ".tmpl") {
			rendered, err := e.Render(ctx, string(content), data)
			if err != nil {
				return nil, fmt.Errorf("failed to render template %s: %w", fileName, err)
			}
			result.Write(rendered)
		} else {
			result.Write(content)
		}

		// Add newline separator between files
		result.WriteString("\n")
	}

	return result.Bytes(), nil
}
