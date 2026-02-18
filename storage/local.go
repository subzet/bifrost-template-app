package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	baseDir string
	baseURL string
}

func NewLocalStorage(baseDir, baseURL string) (*LocalStorage, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}
	return &LocalStorage{baseDir: baseDir, baseURL: baseURL}, nil
}

func (s *LocalStorage) Upload(_ context.Context, key string, r io.Reader, _ int64, _ string) (string, error) {
	path := filepath.Join(s.baseDir, filepath.FromSlash(key))
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", fmt.Errorf("mkdir: %w", err)
	}
	f, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer f.Close()
	if _, err := io.Copy(f, r); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}
	return s.baseURL + "/" + key, nil
}
