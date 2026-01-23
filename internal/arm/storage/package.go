package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jomadu/ai-resource-manager/internal/arm/core"
)

// PackageCache handles package storage within a registry with per-package locking
type PackageCache struct {
	packagesDir string
}

// NewPackageCache creates package storage for given packages directory
func NewPackageCache(packagesDir string) *PackageCache {
	return &PackageCache{packagesDir: packagesDir}
}

// getPackageLock returns cross-process lock for specific package directory
func (p *PackageCache) getPackageLock(packageKey interface{}) (*FileLock, error) {
	hashedKey, err := GenerateKey(packageKey)
	if err != nil {
		return nil, err
	}
	packageDir := filepath.Join(p.packagesDir, hashedKey)
	return NewFileLock(packageDir), nil
}

// Package operations
func (p *PackageCache) SetPackageVersion(ctx context.Context, packageKey interface{}, version *core.Version, files []*core.File) error {
	lock, err := p.getPackageLock(packageKey)
	if err != nil {
		return err
	}
	if err := lock.Lock(ctx); err != nil {
		return err
	}
	defer func() { _ = lock.Unlock() }()

	hashedKey, err := GenerateKey(packageKey)
	if err != nil {
		return err
	}

	packageDir := filepath.Join(p.packagesDir, hashedKey)

	versionDir := filepath.Join(packageDir, fmt.Sprintf("v%d.%d.%d", version.Major, version.Minor, version.Patch))
	filesDir := filepath.Join(versionDir, "files")

	// Create directories
	if err := os.MkdirAll(filesDir, 0o755); err != nil {
		return err
	}

	// Store files
	for _, file := range files {
		filePath := filepath.Join(filesDir, file.Path)
		if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(filePath, file.Content, 0o644); err != nil {
			return err
		}
	}

	now := time.Now()

	// Update package metadata
	var packageMetadata interface{}
	if keyMap, ok := packageKey.(map[string]interface{}); ok {
		packageMetadata = keyMap
	} else {
		packageMetadata = packageKey
	}
	packageMetadataBytes, _ := json.Marshal(packageMetadata)
	_ = os.WriteFile(filepath.Join(packageDir, "metadata.json"), packageMetadataBytes, 0o644)

	// Update version metadata
	versionMetadata := struct {
		Version    core.Version `json:"version"`
		CreatedAt  time.Time    `json:"createdAt"`
		UpdatedAt  time.Time    `json:"updatedAt"`
		AccessedAt time.Time    `json:"accessedAt"`
	}{*version, now, now, now}
	versionMetadataBytes, _ := json.Marshal(versionMetadata)
	_ = os.WriteFile(filepath.Join(versionDir, "metadata.json"), versionMetadataBytes, 0o644)

	return nil
}

func (p *PackageCache) GetPackageVersion(ctx context.Context, packageKey interface{}, version *core.Version) ([]*core.File, error) {
	lock, err := p.getPackageLock(packageKey)
	if err != nil {
		return nil, err
	}
	if err := lock.Lock(ctx); err != nil {
		return nil, err
	}
	defer func() { _ = lock.Unlock() }()

	hashedKey, err := GenerateKey(packageKey)
	if err != nil {
		return nil, err
	}

	packageDir := filepath.Join(p.packagesDir, hashedKey)

	versionDir := filepath.Join(packageDir, fmt.Sprintf("v%d.%d.%d", version.Major, version.Minor, version.Patch))
	filesDir := filepath.Join(versionDir, "files")

	if _, err := os.Stat(filesDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("package version not found")
	}

	// Update access time
	versionMetadataPath := filepath.Join(versionDir, "metadata.json")
	if data, err := os.ReadFile(versionMetadataPath); err == nil {
		var metadata struct {
			Version    core.Version `json:"version"`
			CreatedAt  time.Time    `json:"createdAt"`
			UpdatedAt  time.Time    `json:"updatedAt"`
			AccessedAt time.Time    `json:"accessedAt"`
		}
		_ = json.Unmarshal(data, &metadata)
		metadata.AccessedAt = time.Now()
		updatedData, _ := json.Marshal(metadata)
		_ = os.WriteFile(versionMetadataPath, updatedData, 0o644)
	}

	// Read files
	var files []*core.File
	_ = filepath.Walk(filesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		relPath, _ := filepath.Rel(filesDir, path)
		content, _ := os.ReadFile(path)
		files = append(files, &core.File{
			Path:    relPath,
			Content: content,
			Size:    info.Size(),
		})
		return nil
	})

	return files, nil
}

func (p *PackageCache) ListPackageVersions(ctx context.Context, packageKey interface{}) ([]core.Version, error) {
	hashedKey, err := GenerateKey(packageKey)
	if err != nil {
		return nil, err
	}

	packageDir := filepath.Join(p.packagesDir, hashedKey)
	entries, err := os.ReadDir(packageDir)
	if err != nil {
		return nil, err
	}

	var versions []core.Version
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "v") {
			version, err := core.NewVersion(entry.Name())
			if err == nil {
				versions = append(versions, version)
			}
		}
	}

	return versions, nil
}

func (p *PackageCache) ListPackages(ctx context.Context) ([]interface{}, error) {
	entries, err := os.ReadDir(p.packagesDir)
	if err != nil {
		return nil, err
	}

	var packages []interface{}
	for _, entry := range entries {
		if entry.IsDir() {
			metadataPath := filepath.Join(p.packagesDir, entry.Name(), "metadata.json")
			if data, err := os.ReadFile(metadataPath); err == nil {
				var metadata interface{}
				if json.Unmarshal(data, &metadata) == nil {
					packages = append(packages, metadata)
				}
			}
		}
	}

	return packages, nil
}

// Cleanup operations
func (p *PackageCache) RemovePackageVersion(ctx context.Context, packageKey interface{}, version *core.Version) error {
	lock, err := p.getPackageLock(packageKey)
	if err != nil {
		return err
	}
	if err := lock.Lock(ctx); err != nil {
		return err
	}
	defer func() { _ = lock.Unlock() }()

	hashedKey, err := GenerateKey(packageKey)
	if err != nil {
		return err
	}

	packageDir := filepath.Join(p.packagesDir, hashedKey)

	versionDir := filepath.Join(packageDir, fmt.Sprintf("v%d.%d.%d", version.Major, version.Minor, version.Patch))
	return os.RemoveAll(versionDir)
}

func (p *PackageCache) RemovePackage(ctx context.Context, packageKey interface{}) error {
	lock, err := p.getPackageLock(packageKey)
	if err != nil {
		return err
	}
	if err := lock.Lock(ctx); err != nil {
		return err
	}
	defer func() { _ = lock.Unlock() }()

	hashedKey, err := GenerateKey(packageKey)
	if err != nil {
		return err
	}

	packageDir := filepath.Join(p.packagesDir, hashedKey)

	return os.RemoveAll(packageDir)
}

func (p *PackageCache) Remove(ctx context.Context) error {
	return os.RemoveAll(p.packagesDir)
}

func (p *PackageCache) RemoveOldVersions(ctx context.Context, maxAge time.Duration) error {
	entries, err := os.ReadDir(p.packagesDir)
	if err != nil {
		return err
	}

	cutoff := time.Now().Add(-maxAge)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		packageDir := filepath.Join(p.packagesDir, entry.Name())
		versionEntries, _ := os.ReadDir(packageDir)
		hasVersions := false
		for _, versionEntry := range versionEntries {
			if versionEntry.IsDir() && strings.HasPrefix(versionEntry.Name(), "v") {
				metadataPath := filepath.Join(packageDir, versionEntry.Name(), "metadata.json")
				if data, err := os.ReadFile(metadataPath); err == nil {
					var metadata struct {
						UpdatedAt time.Time `json:"updatedAt"`
					}
					if json.Unmarshal(data, &metadata) == nil && metadata.UpdatedAt.Before(cutoff) {
						_ = os.RemoveAll(filepath.Join(packageDir, versionEntry.Name()))
					} else {
						hasVersions = true
					}
				} else {
					hasVersions = true
				}
			}
		}
		// Remove package directory if no versions remain
		if !hasVersions {
			_ = os.RemoveAll(packageDir)
		}
	}
	return nil
}

func (p *PackageCache) RemoveUnusedVersions(ctx context.Context, maxTimeSinceLastAccess time.Duration) error {
	entries, err := os.ReadDir(p.packagesDir)
	if err != nil {
		return err
	}

	cutoff := time.Now().Add(-maxTimeSinceLastAccess)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		packageDir := filepath.Join(p.packagesDir, entry.Name())
		versionEntries, _ := os.ReadDir(packageDir)
		hasVersions := false
		for _, versionEntry := range versionEntries {
			if versionEntry.IsDir() && strings.HasPrefix(versionEntry.Name(), "v") {
				metadataPath := filepath.Join(packageDir, versionEntry.Name(), "metadata.json")
				if data, err := os.ReadFile(metadataPath); err == nil {
					var metadata struct {
						AccessedAt time.Time `json:"accessedAt"`
					}
					if json.Unmarshal(data, &metadata) == nil && metadata.AccessedAt.Before(cutoff) {
						_ = os.RemoveAll(filepath.Join(packageDir, versionEntry.Name()))
					} else {
						hasVersions = true
					}
				} else {
					hasVersions = true
				}
			}
		}
		// Remove package directory if no versions remain
		if !hasVersions {
			_ = os.RemoveAll(packageDir)
		}
	}
	return nil
}
