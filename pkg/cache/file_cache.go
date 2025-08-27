package cache

import (
	"os"
	"path/filepath"
	"time"
)

// Cache provides content-addressable storage with TTL
type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, content []byte, ttl time.Duration) error
	Delete(key string) error
	Clear() error
}

// FileCache implements Cache for file-based storage
type FileCache struct {
	basePath string
}

func NewFileCache(basePath string) *FileCache {
	return &FileCache{basePath: basePath}
}

func (f *FileCache) Get(key string) ([]byte, error) {
	filePath := filepath.Join(f.basePath, key)
	data, err := os.ReadFile(filePath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	return data, err
}

func (f *FileCache) Set(key string, content []byte, ttl time.Duration) error {
	filePath := filepath.Join(f.basePath, key)
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(filePath, content, 0644)
}

func (f *FileCache) Delete(key string) error {
	filePath := filepath.Join(f.basePath, key)
	return os.Remove(filePath)
}

func (f *FileCache) Clear() error {
	return os.RemoveAll(f.basePath)
}