package common

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// GetBinaryHash returns the SHA256 hash of the current executable
func GetBinaryHash() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	file, err := os.Open(exePath)
	if err != nil {
		return "", fmt.Errorf("failed to open executable: %w", err)
	}
	defer func() { _ = file.Close() }()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to hash executable: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
