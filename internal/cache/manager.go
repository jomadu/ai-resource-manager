package cache

import (
	"context"
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
		if entry.IsDir() {
			registryKey := entry.Name()
			cache := &FileRegistryRulesetCache{
				registryKey: registryKey,
				baseDir:     filepath.Join(m.cacheDir, registryKey),
			}
			if err := cache.CleanupOldVersions(ctx, maxAge); err != nil {
				return err
			}
		}
	}

	return nil
}
