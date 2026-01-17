package packagelockfile

import "context"

// Manager manages package lock file operations
type Manager interface {
	GetDependencyLock(ctx context.Context, registry, packageName, version string) (*DependencyLockConfig, error)
	GetLockFile(ctx context.Context) (*LockFile, error)
	UpsertDependencyLock(ctx context.Context, registry, packageName, version string, config *DependencyLockConfig) error
	RemoveDependencyLock(ctx context.Context, registry, packageName, version string) error
	UpdateRegistryName(ctx context.Context, oldName, newName string) error
}
