package packagelockfile

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/jomadu/ai-resource-manager/internal/v4/core"
)

type PackageLockfile struct {
	Version string `json:"version"`
	Packages map[string]*PackageLockInfo `json:"packages"`
}

type PackageLockInfo struct {
	Version  string `json:"version"`
	Checksum string `json:"checksum"`
}

type FileManager struct {
	lockPath string
}

func NewFileManager() *FileManager {
	return &FileManager{lockPath: "arm-lock.json"}
}

func NewFileManagerWithPath(lockPath string) *FileManager {
	return &FileManager{lockPath: lockPath}
}

func (f *FileManager) GetPackageLockInfo(ctx context.Context, registryName, packageName string) (*PackageLockInfo, error) {
	lockfile, err := f.readLockFile()
	if err != nil {
		return nil, err
	}

	key := core.PackageKey(registryName, packageName)
	lockInfo, exists := lockfile.Packages[key]
	if !exists {
		return nil, errors.New("package not found")
	}

	return lockInfo, nil
}

func (f *FileManager) GetPackageLockfile(ctx context.Context) (*PackageLockfile, error) {
	return f.readLockFile()
}

func (f *FileManager) UpsertPackageLockInfo(ctx context.Context, registryName, packageName string, lockInfo *PackageLockInfo) error {
	var lockfile *PackageLockfile
	var err error

	// Try to read existing lockfile
	lockfile, err = f.readLockFile()
	if err != nil {
		// If file doesn't exist, create new lockfile
		if os.IsNotExist(err) {
			lockfile = &PackageLockfile{
				Version:  "1.0.0",
				Packages: make(map[string]*PackageLockInfo),
			}
		} else {
			return err
		}
	}

	// Upsert the package (create or update)
	key := core.PackageKey(registryName, packageName)
	lockfile.Packages[key] = lockInfo

	// Write lockfile back to disk
	return f.writeLockFile(lockfile)
}

func (f *FileManager) UpdatePackageName(ctx context.Context, registryName, packageName, newPackageName string) error {
	lockfile, err := f.readLockFile()
	if err != nil {
		return err
	}

	oldKey := core.PackageKey(registryName, packageName)
	newKey := core.PackageKey(registryName, newPackageName)

	lockInfo, exists := lockfile.Packages[oldKey]
	if !exists {
		return errors.New("package not found")
	}

	// Check if new package name already exists
	if _, exists := lockfile.Packages[newKey]; exists {
		return errors.New("package already exists")
	}

	// Rename the package
	lockfile.Packages[newKey] = lockInfo
	delete(lockfile.Packages, oldKey)

	return f.writeLockFile(lockfile)
}

func (f *FileManager) RemovePackageLockInfo(ctx context.Context, registryName, packageName string) error {
	lockfile, err := f.readLockFile()
	if err != nil {
		return err
	}

	key := core.PackageKey(registryName, packageName)
	_, exists := lockfile.Packages[key]
	if !exists {
		return errors.New("package not found")
	}

	// Remove the package
	delete(lockfile.Packages, key)

	// Delete lockfile if empty
	if f.isLockfileEmpty(lockfile) {
		return f.deleteLockFile()
	}

	return f.writeLockFile(lockfile)
}

func (f *FileManager) UpdateRegistryName(ctx context.Context, oldName, newName string) error {
	lockfile, err := f.readLockFile()
	if err != nil {
		return err
	}

	oldPrefix := oldName + "/"
	newPrefix := newName + "/"
	found := false

	// Check if any packages exist for old registry
	for key := range lockfile.Packages {
		if strings.HasPrefix(key, oldPrefix) {
			found = true
			break
		}
	}

	if !found {
		return errors.New("registry not found")
	}

	// Check if new registry name already exists
	for key := range lockfile.Packages {
		if strings.HasPrefix(key, newPrefix) {
			return errors.New("registry already exists")
		}
	}

	// Rename all packages from old registry to new registry
	for key, lockInfo := range lockfile.Packages {
		if strings.HasPrefix(key, oldPrefix) {
			regName, pkgName := core.ParsePackageKey(key)
			if regName == oldName {
				newKey := core.PackageKey(newName, pkgName)
				lockfile.Packages[newKey] = lockInfo
				delete(lockfile.Packages, key)
			}
		}
	}

	return f.writeLockFile(lockfile)
}

// readLockFile reads and parses the lockfile from disk
func (f *FileManager) readLockFile() (*PackageLockfile, error) {
	data, err := os.ReadFile(f.lockPath)
	if err != nil {
		// Return error if file doesn't exist (no auto-create)
		return nil, err
	}

	var lockfile PackageLockfile
	err = json.Unmarshal(data, &lockfile)
	if err != nil {
		return nil, err
	}

	// Initialize nil maps
	if lockfile.Packages == nil {
		lockfile.Packages = make(map[string]*PackageLockInfo)
	}

	return &lockfile, nil
}

// writeLockFile writes the lockfile to disk
func (f *FileManager) writeLockFile(lockfile *PackageLockfile) error {
	data, err := json.MarshalIndent(lockfile, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(f.lockPath, data, 0644)
}

// isLockfileEmpty checks if the lockfile has no packages
func (f *FileManager) isLockfileEmpty(lockfile *PackageLockfile) bool {
	return len(lockfile.Packages) == 0
}

// deleteLockFile removes the lockfile from disk
func (f *FileManager) deleteLockFile() error {
	return os.Remove(f.lockPath)
}