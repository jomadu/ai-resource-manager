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

// FileRegistryPackageCache implements filesystem-based package caching.
type FileRegistryPackageCache struct {
	registryKeyObj interface{}
	registryDir    string
}

// NewRegistryPackageCache creates a new registry-scoped package cache.
func NewRegistryPackageCache(registryKeyObj interface{}) (*FileRegistryPackageCache, error) {
	registryKey, err := GenerateKey(registryKeyObj)
	if err != nil {
		return nil, err
	}

	registriesDir := GetRegistriesDir()
	registryDir := filepath.Join(registriesDir, registryKey)

	return &FileRegistryPackageCache{
		registryKeyObj: registryKeyObj,
		registryDir:    registryDir,
	}, nil
}

func (f *FileRegistryPackageCache) ListVersions(ctx context.Context, keyObj interface{}) ([]string, error) {
	packageKey, err := GenerateKey(keyObj)
	if err != nil {
		return nil, err
	}

	packageDir := filepath.Join(f.registryDir, "packages", packageKey)
	entries, err := os.ReadDir(packageDir)
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

func (f *FileRegistryPackageCache) GetPackageVersion(ctx context.Context, keyObj interface{}, version string) ([]types.File, error) {
	registryKey, _ := GenerateKey(f.registryKeyObj)
	var files []types.File
	var walkErr error

	err := WithRegistryLock(registryKey, func() error {
		packageKey, err := GenerateKey(keyObj)
		if err != nil {
			return err
		}

		versionDir := filepath.Join(f.registryDir, "packages", packageKey, version)
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

func (f *FileRegistryPackageCache) SetPackageVersion(ctx context.Context, keyObj interface{}, version string, files []types.File) error {
	registryKey, _ := GenerateKey(f.registryKeyObj)

	return WithRegistryLock(registryKey, func() error {
		packageKey, err := GenerateKey(keyObj)
		if err != nil {
			return err
		}

		versionDir := filepath.Join(f.registryDir, "packages", packageKey, version)
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

// RegistryIndex tracks metadata for cached registry data.
type RegistryIndex struct {
	RegistryMetadata interface{}                  `json:"registryMetadata"`
	CreatedOn        time.Time                    `json:"createdOn"`
	LastUpdatedOn    time.Time                    `json:"lastUpdatedOn"`
	LastAccessedOn   time.Time                    `json:"lastAccessedOn"`
	Packages         map[string]PackageIndexEntry `json:"packages"`
}

// PackageIndexEntry tracks metadata for a cached package.
type PackageIndexEntry struct {
	PackageMetadata interface{}                  `json:"packageMetadata"`
	CreatedOn       time.Time                    `json:"createdOn"`
	LastUpdatedOn   time.Time                    `json:"lastUpdatedOn"`
	LastAccessedOn  time.Time                    `json:"lastAccessedOn"`
	Versions        map[string]VersionIndexEntry `json:"versions"`
}

// VersionIndexEntry tracks metadata for a cached version.
type VersionIndexEntry struct {
	CreatedOn      time.Time `json:"createdOn"`
	LastUpdatedOn  time.Time `json:"lastUpdatedOn"`
	LastAccessedOn time.Time `json:"lastAccessedOn"`
}

// loadIndex loads the registry index from disk, creating a new one if it doesn't exist.
func (f *FileRegistryPackageCache) loadIndex() (*RegistryIndex, error) {
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
			Packages:         make(map[string]PackageIndexEntry),
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
func (f *FileRegistryPackageCache) saveIndex(index *RegistryIndex) error {
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

// updateIndexOnAccess updates the index when a package version is accessed.
func (f *FileRegistryPackageCache) updateIndexOnAccess(keyObj interface{}, version string) error {
	packageKey, err := GenerateKey(keyObj)
	if err != nil {
		return err
	}

	index, err := f.loadIndex()
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	index.LastAccessedOn = now

	// Update existing package entry (should exist since we're accessing cached data)
	packageEntry, exists := index.Packages[packageKey]
	if exists {
		packageEntry.LastAccessedOn = now

		// Update existing version entry
		versionEntry, versionExists := packageEntry.Versions[version]
		if versionExists {
			versionEntry.LastAccessedOn = now
			packageEntry.Versions[version] = versionEntry
		}

		index.Packages[packageKey] = packageEntry
	}

	return f.saveIndex(index)
}

// updateIndexOnSet updates the index when a package version is cached.
func (f *FileRegistryPackageCache) updateIndexOnSet(keyObj interface{}, version string) error {
	packageKey, err := GenerateKey(keyObj)
	if err != nil {
		return err
	}

	index, err := f.loadIndex()
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	index.LastUpdatedOn = now

	// Update or create package entry
	packageEntry, exists := index.Packages[packageKey]
	if !exists {
		packageEntry = PackageIndexEntry{
			PackageMetadata: keyObj,
			CreatedOn:       now,
			LastUpdatedOn:   now,
			LastAccessedOn:  now,
			Versions:        make(map[string]VersionIndexEntry),
		}
	} else {
		packageEntry.LastUpdatedOn = now
	}

	// Update or create version entry
	versionEntry := VersionIndexEntry{
		CreatedOn:      now,
		LastUpdatedOn:  now,
		LastAccessedOn: now,
	}

	packageEntry.Versions[version] = versionEntry
	index.Packages[packageKey] = packageEntry

	return f.saveIndex(index)
}

// Cleanup removes cached versions that haven't been accessed within maxAge.
func (f *FileRegistryPackageCache) Cleanup(maxAge time.Duration) error {
	index, err := f.loadIndex()
	if err != nil {
		return err
	}

	cutoff := time.Now().UTC().Add(-maxAge)
	modified := false

	for packageKey, packageEntry := range index.Packages {
		for version, versionEntry := range packageEntry.Versions {
			if versionEntry.LastAccessedOn.Before(cutoff) {
				// Remove version directory
				versionDir := filepath.Join(f.registryDir, "packages", packageKey, version)
				if err := os.RemoveAll(versionDir); err != nil {
					return err
				}
				// Remove from index
				delete(packageEntry.Versions, version)
				modified = true
			}
		}
		// Update index entry
		index.Packages[packageKey] = packageEntry
	}

	if modified {
		return f.saveIndex(index)
	}
	return nil
}
