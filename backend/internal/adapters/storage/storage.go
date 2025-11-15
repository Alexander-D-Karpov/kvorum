package storage

import (
	"context"
	"io"
	"time"
)

type Storage interface {
	Upload(ctx context.Context, key string, r io.Reader, contentType string) error
	Download(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	GetPresignedURL(ctx context.Context, key string, ttl time.Duration) (string, error)
}

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{basePath: basePath}
}

func (s *LocalStorage) Upload(ctx context.Context, key string, r io.Reader, contentType string) error {
	return nil
}

func (s *LocalStorage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	return nil, nil
}

func (s *LocalStorage) Delete(ctx context.Context, key string) error {
	return nil
}

func (s *LocalStorage) GetPresignedURL(ctx context.Context, key string, ttl time.Duration) (string, error) {
	return "", nil
}
