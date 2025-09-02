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
type Manager struct {
	cacheDir string
}

// NewManager creates a new cache manager.
func NewManager() *Manager {
	homeDir, _ := os.UserHomeDir()
	cacheDir := filepath.Join(homeDir, ".arm", "cache", "registries")
	return &Manager{cacheDir: cacheDir}
}

// CleanupOldVersions removes old cached versions across all registries.
func (m *Manager) CleanupOldVersions(ctx context.Context, maxAge time.Duration) error {
	entries, err := os.ReadDir(m.cacheDir)
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

		baseDir := filepath.Join(m.cacheDir, entry.Name())

		// Try to read registry metadata from existing index
		registryKeyObj, err := m.readRegistryMetadata(baseDir)
		if err != nil {
			// Nuke corrupted registry - it will repopulate later
			if err := os.RemoveAll(baseDir); err != nil {
				return err
			}
			continue
		}

		cache, err := NewRegistryRulesetCache(registryKeyObj)
		if err != nil {
			// Nuke registry with invalid metadata
			if err := os.RemoveAll(baseDir); err != nil {
				return err
			}
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
