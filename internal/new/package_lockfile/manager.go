package packagelockfile

import (
	"context"
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
	SetPackageLockInfo(ctx context.Context, registryName, packageName string, lockInfo *PackageLockInfo) error
	UpdatePackageName(ctx context.Context, registryName, packageName, newPackageName string) error
	RemovePackageLockInfo(ctx context.Context, registryName, packageName string) error
	UpdateRegistryName(ctx context.Context, oldName, newName string) error
	UpdatePackageLockfileVersion(ctx context.Context, version string) error
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
	// TODO: Implement
	return nil, nil
}

func (f *FileManager) GetPackagesLockInfo(ctx context.Context) (map[string]map[string]*PackageLockInfo, error) {
	// TODO: Implement
	return nil, nil
}

func (f *FileManager) AddPackageLockInfo(ctx context.Context, registryName, packageName string, lockInfo *PackageLockInfo) error {
	// TODO: Implement
	return nil
}

func (f *FileManager) UpdatePackageLockInfoName(ctx context.Context, registryName, packageName, newPackageName string) error {
	// TODO: Implement
	return nil
}

func (f *FileManager) UpdatePackageLockInfo(ctx context.Context, registryName, packageName string, lockInfo *PackageLockInfo) error {
	// TODO: Implement
	return nil
}

func (f *FileManager) RemovePackageLockInfo(ctx context.Context, registryName, packageName string) error {
	// TODO: Implement
	return nil
}