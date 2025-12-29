package packagelockfile

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileManager_GetPackageLockInfo(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		setupFile     func(t *testing.T) string
		registryName  string
		packageName   string
		want          *PackageLockInfo
		wantErr       bool
		errContains   string
	}{
		{
			name: "success - package exists",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version: "1.0.0",
					Packages: map[string]*PackageLockInfo{
						"registry1/package1": {
							Version:  "1.2.3",
							Checksum: "sha256:abc123",
						},
					},
				}
				return createTestLockfile(t, lockfile)
			},
			registryName: "registry1",
			packageName:  "package1",
			want: &PackageLockInfo{
				Version:  "1.2.3",
				Checksum: "sha256:abc123",
			},
			wantErr: false,
		},
		{
			name: "error - file does not exist",
			setupFile: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "nonexistent.json")
			},
			registryName: "registry1",
			packageName:  "package1",
			want:         nil,
			wantErr:      true,
			errContains:  "no such file",
		},
		{
			name: "error - package not found",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version: "1.0.0",
					Packages: map[string]*PackageLockInfo{
						"registry1/package1": {
							Version:  "1.2.3",
							Checksum: "sha256:abc123",
						},
					},
				}
				return createTestLockfile(t, lockfile)
			},
			registryName: "registry1",
			packageName:  "nonexistent-package",
			want:         nil,
			wantErr:      true,
			errContains:  "package not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lockPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(lockPath)

			got, err := fm.GetPackageLockInfo(ctx, tt.registryName, tt.packageName)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetPackageLockInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetPackageLockInfo() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("GetPackageLockInfo() error = %v, should contain %v", err, tt.errContains)
				}
				return
			}

			if got == nil {
				t.Errorf("GetPackageLockInfo() = nil, want %v", tt.want)
				return
			}

			if got.Version != tt.want.Version {
				t.Errorf("GetPackageLockInfo() Version = %v, want %v", got.Version, tt.want.Version)
			}
			if got.Checksum != tt.want.Checksum {
				t.Errorf("GetPackageLockInfo() Checksum = %v, want %v", got.Checksum, tt.want.Checksum)
			}
		})
	}
}

func TestFileManager_UpsertPackageLockInfo(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		setupFile     func(t *testing.T) string
		registryName  string
		packageName   string
		lockInfo      *PackageLockInfo
		wantErr       bool
	}{
		{
			name: "success - create new lockfile",
			setupFile: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "new-lock.json")
			},
			registryName: "registry1",
			packageName:  "package1",
			lockInfo: &PackageLockInfo{
				Version:  "1.2.3",
				Checksum: "sha256:abc123",
			},
			wantErr: false,
		},
		{
			name: "success - update existing package",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version: "1.0.0",
					Packages: map[string]*PackageLockInfo{
						"registry1/package1": {
							Version:  "1.2.3",
							Checksum: "sha256:abc123",
						},
					},
				}
				return createTestLockfile(t, lockfile)
			},
			registryName: "registry1",
			packageName:  "package1",
			lockInfo: &PackageLockInfo{
				Version:  "1.3.0",
				Checksum: "sha256:xyz789",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lockPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(lockPath)

			err := fm.UpsertPackageLockInfo(ctx, tt.registryName, tt.packageName, tt.lockInfo)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpsertPackageLockInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				got, err := fm.GetPackageLockInfo(ctx, tt.registryName, tt.packageName)
				if err != nil {
					t.Fatalf("GetPackageLockInfo() error = %v", err)
				}

				if got.Version != tt.lockInfo.Version {
					t.Errorf("UpsertPackageLockInfo() Version = %v, want %v", got.Version, tt.lockInfo.Version)
				}
				if got.Checksum != tt.lockInfo.Checksum {
					t.Errorf("UpsertPackageLockInfo() Checksum = %v, want %v", got.Checksum, tt.lockInfo.Checksum)
				}
			}
		})
	}
}

func TestFileManager_UpdatePackageName(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		setupFile      func(t *testing.T) string
		registryName   string
		packageName    string
		newPackageName string
		wantErr        bool
		errContains    string
	}{
		{
			name: "success - rename package",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version: "1.0.0",
					Packages: map[string]*PackageLockInfo{
						"registry1/package1": {
							Version:  "1.2.3",
							Checksum: "sha256:abc123",
						},
					},
				}
				return createTestLockfile(t, lockfile)
			},
			registryName:   "registry1",
			packageName:    "package1",
			newPackageName: "package1-renamed",
			wantErr:        false,
		},
		{
			name: "error - package not found",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version: "1.0.0",
					Packages: map[string]*PackageLockInfo{
						"registry1/package1": {
							Version:  "1.2.3",
							Checksum: "sha256:abc123",
						},
					},
				}
				return createTestLockfile(t, lockfile)
			},
			registryName:   "registry1",
			packageName:    "nonexistent-package",
			newPackageName: "package2",
			wantErr:        true,
			errContains:    "package not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lockPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(lockPath)

			err := fm.UpdatePackageName(ctx, tt.registryName, tt.packageName, tt.newPackageName)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdatePackageName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdatePackageName() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("UpdatePackageName() error = %v, should contain %v", err, tt.errContains)
				}
			}
		})
	}
}

func TestFileManager_RemovePackageLockInfo(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		setupFile    func(t *testing.T) string
		registryName string
		packageName  string
		wantErr      bool
		errContains  string
	}{
		{
			name: "success - remove package",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version: "1.0.0",
					Packages: map[string]*PackageLockInfo{
						"registry1/package1": {
							Version:  "1.2.3",
							Checksum: "sha256:abc123",
						},
						"registry1/package2": {
							Version:  "2.0.0",
							Checksum: "sha256:def456",
						},
					},
				}
				return createTestLockfile(t, lockfile)
			},
			registryName: "registry1",
			packageName:  "package1",
			wantErr:      false,
		},
		{
			name: "error - package not found",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version: "1.0.0",
					Packages: map[string]*PackageLockInfo{
						"registry1/package1": {
							Version:  "1.2.3",
							Checksum: "sha256:abc123",
						},
					},
				}
				return createTestLockfile(t, lockfile)
			},
			registryName: "registry1",
			packageName:  "nonexistent-package",
			wantErr:      true,
			errContains:  "package not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lockPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(lockPath)

			err := fm.RemovePackageLockInfo(ctx, tt.registryName, tt.packageName)

			if (err != nil) != tt.wantErr {
				t.Errorf("RemovePackageLockInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("RemovePackageLockInfo() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("RemovePackageLockInfo() error = %v, should contain %v", err, tt.errContains)
				}
			}
		})
	}
}

func TestFileManager_UpdateRegistryName(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		setupFile   func(t *testing.T) string
		oldName     string
		newName     string
		wantErr     bool
		errContains string
	}{
		{
			name: "success - rename registry",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version: "1.0.0",
					Packages: map[string]*PackageLockInfo{
						"registry1/package1": {
							Version:  "1.2.3",
							Checksum: "sha256:abc123",
						},
						"registry1/package2": {
							Version:  "2.0.0",
							Checksum: "sha256:def456",
						},
					},
				}
				return createTestLockfile(t, lockfile)
			},
			oldName: "registry1",
			newName: "registry1-renamed",
			wantErr: false,
		},
		{
			name: "error - registry not found",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version: "1.0.0",
					Packages: map[string]*PackageLockInfo{
						"registry1/package1": {
							Version:  "1.2.3",
							Checksum: "sha256:abc123",
						},
					},
				}
				return createTestLockfile(t, lockfile)
			},
			oldName:     "nonexistent-registry",
			newName:     "registry2",
			wantErr:     true,
			errContains: "registry not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lockPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(lockPath)

			err := fm.UpdateRegistryName(ctx, tt.oldName, tt.newName)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateRegistryName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateRegistryName() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("UpdateRegistryName() error = %v, should contain %v", err, tt.errContains)
				}
			}
		})
	}
}

// Helper functions

func createTestLockfile(t *testing.T, lockfile *PackageLockfile) string {
	tmpDir := t.TempDir()
	lockPath := filepath.Join(tmpDir, "test-lock.json")

	data, err := json.MarshalIndent(lockfile, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal lockfile: %v", err)
	}

	err = os.WriteFile(lockPath, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write lockfile: %v", err)
	}

	return lockPath
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}