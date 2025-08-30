package cache

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

// FileRegistryRulesetCache implements filesystem-based ruleset caching.
type FileRegistryRulesetCache struct {
	registryKey string
	baseDir     string
}

// NewRegistryRulesetCache creates a new registry-scoped ruleset cache.
func NewRegistryRulesetCache(registryKey string) *FileRegistryRulesetCache {
	homeDir, _ := os.UserHomeDir()
	baseDir := filepath.Join(homeDir, ".arm", "cache", "registries", registryKey)

	return &FileRegistryRulesetCache{
		registryKey: registryKey,
		baseDir:     baseDir,
	}
}

func (f *FileRegistryRulesetCache) ListVersions(ctx context.Context, rulesetKey string) ([]string, error) {
	rulesetDir := filepath.Join(f.baseDir, "rulesets", rulesetKey)
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

func (f *FileRegistryRulesetCache) GetRulesetVersion(ctx context.Context, rulesetKey, version string) ([]types.File, error) {
	versionDir := filepath.Join(f.baseDir, "rulesets", rulesetKey, version)
	if _, err := os.Stat(versionDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("version %s not found in cache", version)
	}

	var files []types.File
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

		files = append(files, types.File{
			Path:    filepath.ToSlash(relPath),
			Content: content,
			Size:    stat.Size(),
		})
		return nil
	})

	return files, err
}

func (f *FileRegistryRulesetCache) SetRulesetVersion(ctx context.Context, rulesetKey, version string, files []types.File) error {
	versionDir := filepath.Join(f.baseDir, "rulesets", rulesetKey, version)
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

func (f *FileRegistryRulesetCache) InvalidateRuleset(ctx context.Context, rulesetKey string) error {
	rulesetDir := filepath.Join(f.baseDir, "rulesets", rulesetKey)
	return os.RemoveAll(rulesetDir)
}

func (f *FileRegistryRulesetCache) InvalidateVersion(ctx context.Context, rulesetKey, version string) error {
	versionDir := filepath.Join(f.baseDir, "rulesets", rulesetKey, version)
	return os.RemoveAll(versionDir)
}
