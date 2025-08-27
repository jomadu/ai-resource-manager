package cache

import (
	"os"
	"path/filepath"

	"github.com/jomadu/ai-rules-manager/pkg/registry"
)

// RulesetCache provides registry-agnostic ruleset content caching
type RulesetCache interface {
	// CacheContent stores extracted content by commit hash
	CacheContent(rulesetKey, commit string, files []registry.File) error
	// GetCachedContent retrieves previously cached content
	GetCachedContent(rulesetKey, commit string) ([]registry.File, error)
	// HasCachedContent checks if content exists in cache
	HasCachedContent(rulesetKey, commit string) bool
}

// FileRulesetCache implements RulesetCache for file-based storage
type FileRulesetCache struct {
	basePath string
}

func NewFileRulesetCache(basePath string) *FileRulesetCache {
	return &FileRulesetCache{basePath: basePath}
}

func (f *FileRulesetCache) CacheContent(rulesetKey, commit string, files []registry.File) error {
	cachePath := filepath.Join(f.basePath, "registries", rulesetKey, commit)
	if err := os.MkdirAll(cachePath, 0755); err != nil {
		return err
	}
	for _, file := range files {
		filePath := filepath.Join(cachePath, file.Path)
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		if err := os.WriteFile(filePath, file.Content, 0644); err != nil {
			return err
		}
	}
	return nil
}

func (f *FileRulesetCache) GetCachedContent(rulesetKey, commit string) ([]registry.File, error) {
	cachePath := filepath.Join(f.basePath, "registries", rulesetKey, commit)
	var files []registry.File
	err := filepath.Walk(cachePath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		relPath, _ := filepath.Rel(cachePath, path)
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		files = append(files, registry.File{
			Path:    relPath,
			Content: content,
			Size:    info.Size(),
		})
		return nil
	})
	if os.IsNotExist(err) {
		return nil, nil
	}
	return files, err
}

func (f *FileRulesetCache) HasCachedContent(rulesetKey, commit string) bool {
	cachePath := filepath.Join(f.basePath, "registries", rulesetKey, commit)
	_, err := os.Stat(cachePath)
	return err == nil
}