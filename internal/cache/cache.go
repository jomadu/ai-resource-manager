package cache

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

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
type FileCache struct {
	cacheDir string
}

// NewFileCache creates a new file-based cache.
func NewFileCache() *FileCache {
	homeDir, _ := os.UserHomeDir()
	cacheDir := filepath.Join(homeDir, ".arm", "cache")
	return &FileCache{cacheDir: cacheDir}
}

// NewFileCacheWithDir creates a new file-based cache with custom directory.
func NewFileCacheWithDir(cacheDir string) *FileCache {
	return &FileCache{cacheDir: cacheDir}
}

func (f *FileCache) ListVersions(ctx context.Context, registryKey, rulesetKey string) ([]string, error) {
	rulesetDir := filepath.Join(f.cacheDir, "registries", registryKey, "rulesets", rulesetKey)
	entries, err := os.ReadDir(rulesetDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var versions []string
	for _, entry := range entries {
		if entry.IsDir() {
			versions = append(versions, entry.Name())
		}
	}
	sort.Strings(versions)
	return versions, nil
}

func (f *FileCache) Get(ctx context.Context, registryKey, rulesetKey, version string) ([]arm.File, error) {
	versionDir := filepath.Join(f.cacheDir, "registries", registryKey, "rulesets", rulesetKey, version)
	if _, err := os.Stat(versionDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("version %s not found in cache", version)
	}

	var files []arm.File
	err := filepath.WalkDir(versionDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(versionDir, path)
		if err != nil {
			return err
		}

		stat, err := os.Stat(path)
		if err != nil {
			return err
		}

		files = append(files, arm.File{
			Path:    filepath.ToSlash(relPath),
			Content: content,
			Size:    stat.Size(),
		})
		return nil
	})

	return files, err
}

func (f *FileCache) Set(ctx context.Context, registryKey, rulesetKey, version string, files []arm.File) error {
	versionDir := filepath.Join(f.cacheDir, "registries", registryKey, "rulesets", rulesetKey, version)
	if err := os.MkdirAll(versionDir, 0o755); err != nil {
		return err
	}

	for _, file := range files {
		filePath := filepath.Join(versionDir, filepath.FromSlash(file.Path))
		fileDir := filepath.Dir(filePath)
		if err := os.MkdirAll(fileDir, 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(filePath, file.Content, 0o644); err != nil {
			return err
		}
	}

	return nil
}

func (f *FileCache) InvalidateRegistry(ctx context.Context, registryKey string) error {
	registryDir := filepath.Join(f.cacheDir, "registries", registryKey)
	return os.RemoveAll(registryDir)
}

func (f *FileCache) InvalidateRuleset(ctx context.Context, registryKey, rulesetKey string) error {
	rulesetDir := filepath.Join(f.cacheDir, "registries", registryKey, "rulesets", rulesetKey)
	return os.RemoveAll(rulesetDir)
}

func (f *FileCache) InvalidateVersion(ctx context.Context, registryKey, rulesetKey, version string) error {
	versionDir := filepath.Join(f.cacheDir, "registries", registryKey, "rulesets", rulesetKey, version)
	return os.RemoveAll(versionDir)
}
