package packagelockfile

import (
	"context"
	"encoding/json"
	"errors"
	"os"
)

type PackageLockfile struct {
	Version string `json:"version"`
	Packages map[string]map[string]*PackageLockInfo `json:"packages"`
}

type PackageLockInfo struct {
	Version  string `json:"version"`
	Checksum string `json:"checksum"`
}

type Manager interface {
	GetPackageLockInfo(ctx context.Context, registryName, packageName string) (*PackageLockInfo, error)
	GetPackageLockfile(ctx context.Context) (*PackageLockfile, error)
	UpsertPackageLockInfo(ctx context.Context, registryName, packageName string, lockInfo *PackageLockInfo) error
	UpdatePackageName(ctx context.Context, registryName, packageName, newPackageName string) error
	RemovePackageLockInfo(ctx context.Context, registryName, packageName string) error
	UpdateRegistryName(ctx context.Context, oldName, newName string) error
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

	registry, exists := lockfile.Packages[registryName]
	if !exists {
		return nil, errors.New("registry not found")
	}

	lockInfo, exists := registry[packageName]
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
				Packages: make(map[string]map[string]*PackageLockInfo),
			}
		} else {
			return err
		}
	}

	// Create registry if it doesn't exist
	if lockfile.Packages[registryName] == nil {
		lockfile.Packages[registryName] = make(map[string]*PackageLockInfo)
	}

	// Upsert the package (create or update)
	lockfile.Packages[registryName][packageName] = lockInfo

	// Write lockfile back to disk
	return f.writeLockFile(lockfile)
}

func (f *FileManager) UpdatePackageName(ctx context.Context, registryName, packageName, newPackageName string) error {
	lockfile, err := f.readLockFile()
	if err != nil {
		return err
	}

	registry, exists := lockfile.Packages[registryName]
	if !exists {
		return errors.New("registry not found")
	}

	lockInfo, exists := registry[packageName]
	if !exists {
		return errors.New("package not found")
	}

	// Check if new package name already exists
	if _, exists := registry[newPackageName]; exists {
		return errors.New("package already exists")
	}

	// Rename the package
	registry[newPackageName] = lockInfo
	delete(registry, packageName)

	return f.writeLockFile(lockfile)
}

func (f *FileManager) RemovePackageLockInfo(ctx context.Context, registryName, packageName string) error {
	lockfile, err := f.readLockFile()
	if err != nil {
		return err
	}

	registry, exists := lockfile.Packages[registryName]
	if !exists {
		return errors.New("registry not found")
	}

	_, exists = registry[packageName]
	if !exists {
		return errors.New("package not found")
	}

	// Remove the package
	delete(registry, packageName)

	// Remove empty registry
	if len(registry) == 0 {
		delete(lockfile.Packages, registryName)
	}

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

	oldRegistry, exists := lockfile.Packages[oldName]
	if !exists {
		return errors.New("registry not found")
	}

	// Check if new registry name already exists
	if _, exists := lockfile.Packages[newName]; exists {
		return errors.New("registry already exists")
	}

	// Move all packages from old registry to new registry
	lockfile.Packages[newName] = oldRegistry
	delete(lockfile.Packages, oldName)

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
		lockfile.Packages = make(map[string]map[string]*PackageLockInfo)
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