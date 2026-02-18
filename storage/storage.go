package storage

import (
	"context"
	"io"
)

type Storage interface {
	Upload(ctx context.Context, key string, r io.Reader, size int64, contentType string) (string, error)
}

type noopStorage struct{}

func (n *noopStorage) Upload(_ context.Context, _ string, _ io.Reader, _ int64, _ string) (string, error) {
	return "", nil
}

func Noop() Storage {
	return &noopStorage{}
}
