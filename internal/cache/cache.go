package cache

import (
	"context"
	"errors"

	"github.com/jomadu/ai-rules-manager/internal/arm"
)

// Cache provides local storage for registry data and ruleset files.
type Cache interface {
	ListVersions(ctx context.Context, registryKey, rulesetKey string) ([]string, error)
	Get(ctx context.Context, registryKey, rulesetKey, version string) ([]arm.File, error)
	Set(ctx context.Context, registryKey, rulesetKey, version string, files []arm.File) error
	InvalidateRegistry(ctx context.Context, registryKey string) error
	InvalidateRuleset(ctx context.Context, registryKey, rulesetKey string) error
	InvalidateVersion(ctx context.Context, registryKey, rulesetKey, version string) error
}

// FileCache implements filesystem-based caching.
type FileCache struct{}

// NewFileCache creates a new file-based cache.
func NewFileCache() *FileCache {
	return &FileCache{}
}

func (f *FileCache) ListVersions(ctx context.Context, registryKey, rulesetKey string) ([]string, error) {
	return nil, errors.New("not implemented")
}

func (f *FileCache) Get(ctx context.Context, registryKey, rulesetKey, version string) ([]arm.File, error) {
	return nil, errors.New("not implemented")
}

func (f *FileCache) Set(ctx context.Context, registryKey, rulesetKey, version string, files []arm.File) error {
	return errors.New("not implemented")
}

func (f *FileCache) InvalidateRegistry(ctx context.Context, registryKey string) error {
	return errors.New("not implemented")
}

func (f *FileCache) InvalidateRuleset(ctx context.Context, registryKey, rulesetKey string) error {
	return errors.New("not implemented")
}

func (f *FileCache) InvalidateVersion(ctx context.Context, registryKey, rulesetKey, version string) error {
	return errors.New("not implemented")
}
