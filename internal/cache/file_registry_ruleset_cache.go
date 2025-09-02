package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

// FileRegistryRulesetCache implements filesystem-based ruleset caching.
type FileRegistryRulesetCache struct {
	registryKeyObj interface{}
	registryDir    string
}

// NewRegistryRulesetCache creates a new registry-scoped ruleset cache.
func NewRegistryRulesetCache(registryKeyObj interface{}) (*FileRegistryRulesetCache, error) {
	registryKey, err := GenerateKey(registryKeyObj)
	if err != nil {
		return nil, err
	}

	cacheDir := GetCacheDir()
	registryDir := filepath.Join(cacheDir, registryKey)

	return &FileRegistryRulesetCache{
		registryKeyObj: registryKeyObj,
		registryDir:    registryDir,
	}, nil
}

func (f *FileRegistryRulesetCache) ListVersions(ctx context.Context, keyObj interface{}) ([]string, error) {
	rulesetKey, err := GenerateKey(keyObj)
	if err != nil {
		return nil, err
	}

	rulesetDir := filepath.Join(f.registryDir, "rulesets", rulesetKey)
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

func (f *FileRegistryRulesetCache) GetRulesetVersion(ctx context.Context, keyObj interface{}, version string) ([]types.File, error) {
	registryKey, _ := GenerateKey(f.registryKeyObj)
	var files []types.File
	var walkErr error

	err := WithRegistryLock(registryKey, func() error {
		rulesetKey, err := GenerateKey(keyObj)
		if err != nil {
			return err
		}

		versionDir := filepath.Join(f.registryDir, "rulesets", rulesetKey, version)
		if _, statErr := os.Stat(versionDir); os.IsNotExist(statErr) {
			return fmt.Errorf("version %s not found in cache", version)
		}

		walkErr = filepath.WalkDir(versionDir, func(path string, d fs.DirEntry, err error) error {
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

		// Update index on access
		if walkErr == nil {
			_ = f.updateIndexOnAccess(keyObj, version)
		}

		return walkErr
	})

	if err != nil {
		return nil, err
	}
	return files, walkErr
}

func (f *FileRegistryRulesetCache) SetRulesetVersion(ctx context.Context, keyObj interface{}, version string, files []types.File) error {
	registryKey, _ := GenerateKey(f.registryKeyObj)

	return WithRegistryLock(registryKey, func() error {
		rulesetKey, err := GenerateKey(keyObj)
		if err != nil {
			return err
		}

		versionDir := filepath.Join(f.registryDir, "rulesets", rulesetKey, version)
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

		// Update index on set
		_ = f.updateIndexOnSet(keyObj, version)

		return nil
	})
}

func (f *FileRegistryRulesetCache) InvalidateRuleset(ctx context.Context, rulesetKey string) error {
	rulesetDir := filepath.Join(f.registryDir, "rulesets", rulesetKey)
	return os.RemoveAll(rulesetDir)
}

func (f *FileRegistryRulesetCache) InvalidateVersion(ctx context.Context, rulesetKey, version string) error {
	versionDir := filepath.Join(f.registryDir, "rulesets", rulesetKey, version)
	return os.RemoveAll(versionDir)
}

// RegistryIndex tracks metadata for cached registry data.
type RegistryIndex struct {
	RegistryMetadata interface{}                  `json:"registry_metadata"`
	CreatedOn        time.Time                    `json:"created_on"`
	LastUpdatedOn    time.Time                    `json:"last_updated_on"`
	LastAccessedOn   time.Time                    `json:"last_accessed_on"`
	Rulesets         map[string]RulesetIndexEntry `json:"rulesets"`
}

// RulesetIndexEntry tracks metadata for a cached ruleset.
type RulesetIndexEntry struct {
	RulesetMetadata interface{}                  `json:"ruleset_metadata"`
	CreatedOn       time.Time                    `json:"created_on"`
	LastUpdatedOn   time.Time                    `json:"last_updated_on"`
	LastAccessedOn  time.Time                    `json:"last_accessed_on"`
	Versions        map[string]VersionIndexEntry `json:"versions"`
}

// VersionIndexEntry tracks metadata for a cached version.
type VersionIndexEntry struct {
	CreatedOn      time.Time `json:"created_on"`
	LastUpdatedOn  time.Time `json:"last_updated_on"`
	LastAccessedOn time.Time `json:"last_accessed_on"`
}

// loadIndex loads the registry index from disk, creating a new one if it doesn't exist.
func (f *FileRegistryRulesetCache) loadIndex() (*RegistryIndex, error) {
	indexPath := filepath.Join(f.registryDir, "index.json")

	data, err := os.ReadFile(indexPath)
	if os.IsNotExist(err) {
		// Create new index
		now := time.Now().UTC()
		return &RegistryIndex{
			RegistryMetadata: f.registryKeyObj,
			CreatedOn:        now,
			LastUpdatedOn:    now,
			LastAccessedOn:   now,
			Rulesets:         make(map[string]RulesetIndexEntry),
		}, nil
	}
	if err != nil {
		return nil, err
	}

	var index RegistryIndex
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, err
	}

	return &index, nil
}

// saveIndex saves the registry index to disk.
func (f *FileRegistryRulesetCache) saveIndex(index *RegistryIndex) error {
	if err := os.MkdirAll(f.registryDir, 0o755); err != nil {
		return err
	}

	indexPath := filepath.Join(f.registryDir, "index.json")
	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(indexPath, data, 0o644)
}

// updateIndexOnAccess updates the index when a ruleset version is accessed.
func (f *FileRegistryRulesetCache) updateIndexOnAccess(keyObj interface{}, version string) error {
	rulesetKey, err := GenerateKey(keyObj)
	if err != nil {
		return err
	}

	index, err := f.loadIndex()
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	index.LastAccessedOn = now

	// Update existing ruleset entry (should exist since we're accessing cached data)
	rulesetEntry, exists := index.Rulesets[rulesetKey]
	if exists {
		rulesetEntry.LastAccessedOn = now

		// Update existing version entry
		versionEntry, versionExists := rulesetEntry.Versions[version]
		if versionExists {
			versionEntry.LastAccessedOn = now
			rulesetEntry.Versions[version] = versionEntry
		}

		index.Rulesets[rulesetKey] = rulesetEntry
	}

	return f.saveIndex(index)
}

// updateIndexOnSet updates the index when a ruleset version is cached.
func (f *FileRegistryRulesetCache) updateIndexOnSet(keyObj interface{}, version string) error {
	rulesetKey, err := GenerateKey(keyObj)
	if err != nil {
		return err
	}

	index, err := f.loadIndex()
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	index.LastUpdatedOn = now

	// Update or create ruleset entry
	rulesetEntry, exists := index.Rulesets[rulesetKey]
	if !exists {
		rulesetEntry = RulesetIndexEntry{
			RulesetMetadata: keyObj,
			CreatedOn:       now,
			LastUpdatedOn:   now,
			LastAccessedOn:  now,
			Versions:        make(map[string]VersionIndexEntry),
		}
	} else {
		rulesetEntry.LastUpdatedOn = now
	}

	// Update or create version entry
	versionEntry := VersionIndexEntry{
		CreatedOn:      now,
		LastUpdatedOn:  now,
		LastAccessedOn: now,
	}

	rulesetEntry.Versions[version] = versionEntry
	index.Rulesets[rulesetKey] = rulesetEntry

	return f.saveIndex(index)
}

// CleanupOldVersions removes cached versions that haven't been accessed within maxAge.
func (f *FileRegistryRulesetCache) CleanupOldVersions(ctx context.Context, maxAge time.Duration) error {
	index, err := f.loadIndex()
	if err != nil {
		return err
	}

	cutoff := time.Now().UTC().Add(-maxAge)
	modified := false

	for rulesetKey, rulesetEntry := range index.Rulesets {
		for version, versionEntry := range rulesetEntry.Versions {
			if versionEntry.LastAccessedOn.Before(cutoff) {
				// Remove version directory
				versionDir := filepath.Join(f.registryDir, "rulesets", rulesetKey, version)
				if err := os.RemoveAll(versionDir); err != nil {
					return err
				}
				// Remove from index
				delete(rulesetEntry.Versions, version)
				modified = true
			}
		}
		// Update index entry
		index.Rulesets[rulesetKey] = rulesetEntry
	}

	if modified {
		return f.saveIndex(index)
	}
	return nil
}
