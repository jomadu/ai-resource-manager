package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Manager provides cache management operations across all registries.
type Manager struct{}

// NewManager creates a new cache manager.
func NewManager() *Manager {
	return &Manager{}
}

// CleanupOldVersions removes old cached versions across all registries.
func (m *Manager) CleanupOldVersions(ctx context.Context, maxAge time.Duration) error {
	cacheDir := GetRegistriesDir()
	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No cache directory exists
		}
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		baseDir := filepath.Join(cacheDir, entry.Name())

		// Try to read registry metadata from existing index
		registryKeyObj, err := m.readRegistryMetadata(baseDir)
		if err != nil {
			// Remove corrupted registry and its lock
			registryKey := entry.Name()
			_ = os.RemoveAll(baseDir)
			_ = os.Remove(filepath.Join(cacheDir, ".locks", registryKey+".lock"))
			continue
		}

		cache, err := NewRegistryRulesetCache(registryKeyObj)
		if err != nil {
			// Remove registry with invalid metadata and its lock
			registryKey := entry.Name()
			_ = os.RemoveAll(baseDir)
			_ = os.Remove(filepath.Join(cacheDir, ".locks", registryKey+".lock"))
			continue
		}

		if err := cache.CleanupOldVersions(ctx, maxAge); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) readRegistryMetadata(baseDir string) (interface{}, error) {
	indexPath := filepath.Join(baseDir, "index.json")

	data, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, err
	}

	var index RegistryIndex
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, err
	}

	if index.RegistryMetadata == nil {
		return nil, fmt.Errorf("missing registry metadata")
	}

	return index.RegistryMetadata, nil
}
