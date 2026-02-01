package common

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
)

type CASWriter struct {
	tempFile *os.File
	writer   io.Writer
	hasher   io.Writer
	dir      string
}

func NewCASWriter() (*CASWriter, error) {
	dataDir, err := GetUserDataDir()
	if err != nil {
		return nil, err
	}
	genDir := filepath.Join(dataDir, "generated")
	if err := os.MkdirAll(genDir, 0755); err != nil {
		return nil, err
	}

	tempFile, err := os.CreateTemp(genDir, ".tmp-*")
	if err != nil {
		return nil, err
	}

	hasher := sha256.New()
	writer := io.MultiWriter(tempFile, hasher)

	return &CASWriter{
		tempFile: tempFile,
		writer:   writer,
		hasher:   hasher,
		dir:      genDir,
	}, nil
}

func (c *CASWriter) Write(p []byte) (n int, err error) {
	return c.writer.Write(p)
}

func (c *CASWriter) Seal() (string, error) {
	c.tempFile.Close()
	hash := c.hasher.(interface{ Sum(b []byte) []byte }).Sum(nil)
	hashStr := hex.EncodeToString(hash)
	finalPath := filepath.Join(c.dir, hashStr)

	if _, err := os.Stat(finalPath); os.IsNotExist(err) {
		if err := os.Rename(c.tempFile.Name(), finalPath); err != nil {
			return "", err
		}
	} else {
		os.Remove(c.tempFile.Name())
	}

	return finalPath, nil
}
