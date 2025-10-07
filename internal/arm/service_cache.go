package arm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jomadu/ai-rules-manager/internal/cache"
)

func (a *ArmService) CleanCacheWithAge(ctx context.Context, maxAge time.Duration) error {
	cacheDir := cache.GetRegistriesDir()
	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		if os.IsNotExist(err) {
			a.ui.Success("No cache directory found - nothing to clean")
			return nil
		}
		return fmt.Errorf("failed to read cache directory: %w", err)
	}

	cleanedCount := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		baseDir := filepath.Join(cacheDir, entry.Name())

		// Try to read registry metadata from existing index
		registryKeyObj, err := a.readRegistryMetadata(baseDir)
		if err != nil {
			// Remove corrupted registry and its lock
			registryKey := entry.Name()
			_ = os.RemoveAll(baseDir)
			_ = os.Remove(filepath.Join(cacheDir, ".locks", registryKey+".lock"))
			cleanedCount++
			continue
		}

		packageCache, err := cache.NewRegistryPackageCache(registryKeyObj)
		if err != nil {
			// Remove registry with invalid metadata and its lock
			registryKey := entry.Name()
			_ = os.RemoveAll(baseDir)
			_ = os.Remove(filepath.Join(cacheDir, ".locks", registryKey+".lock"))
			cleanedCount++
			continue
		}

		if err := packageCache.Cleanup(maxAge); err != nil {
			return fmt.Errorf("failed to clean cache for registry %s: %w", entry.Name(), err)
		}
		cleanedCount++
	}

	a.ui.Success(fmt.Sprintf("Cache cleaned: processed %d registries, removed versions older than %v", cleanedCount, maxAge))
	return nil
}

func (a *ArmService) NukeCache(ctx context.Context) error {
	cacheDir := cache.GetCacheDir()
	err := os.RemoveAll(cacheDir)
	if err != nil {
		return fmt.Errorf("failed to remove cache directory: %w", err)
	}
	a.ui.Success("Cache directory removed successfully")
	return nil
}

// readRegistryMetadata reads registry metadata from the index file
func (a *ArmService) readRegistryMetadata(baseDir string) (interface{}, error) {
	indexPath := filepath.Join(baseDir, "index.json")

	data, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, err
	}

	var index cache.RegistryIndex
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, err
	}

	if index.RegistryMetadata == nil {
		return nil, fmt.Errorf("missing registry metadata")
	}

	return index.RegistryMetadata, nil
}
