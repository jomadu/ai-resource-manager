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
					Packages: map[string]map[string]*PackageLockInfo{
						"registry1": {
							"package1": {
								Version:  "1.2.3",
								Checksum: "sha256:abc123",
							},
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
			name: "error - registry not found",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version: "1.0.0",
					Packages: map[string]map[string]*PackageLockInfo{
						"registry1": {
							"package1": {
								Version:  "1.2.3",
								Checksum: "sha256:abc123",
							},
						},
					},
				}
				return createTestLockfile(t, lockfile)
			},
			registryName: "nonexistent-registry",
			packageName:  "package1",
			want:         nil,
			wantErr:      true,
			errContains:  "registry not found",
		},
		{
			name: "error - package not found",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version: "1.0.0",
					Packages: map[string]map[string]*PackageLockInfo{
						"registry1": {
							"package1": {
								Version:  "1.2.3",
								Checksum: "sha256:abc123",
							},
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

func TestFileManager_GetPackageLockfile(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		setupFile   func(t *testing.T) string
		want        *PackageLockfile
		wantErr     bool
		errContains string
	}{
		{
			name: "success - file exists",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version: "1.0.0",
					Packages: map[string]map[string]*PackageLockInfo{
						"registry1": {
							"package1": {
								Version:  "1.2.3",
								Checksum: "sha256:abc123",
							},
						},
					},
				}
				return createTestLockfile(t, lockfile)
			},
			want: &PackageLockfile{
				Version: "1.0.0",
				Packages: map[string]map[string]*PackageLockInfo{
					"registry1": {
						"package1": {
							Version:  "1.2.3",
							Checksum: "sha256:abc123",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error - file does not exist",
			setupFile: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "nonexistent.json")
			},
			want:        nil,
			wantErr:     true,
			errContains: "no such file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lockPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(lockPath)

			got, err := fm.GetPackageLockfile(ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetPackageLockfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetPackageLockfile() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("GetPackageLockfile() error = %v, should contain %v", err, tt.errContains)
				}
				return
			}

			if got == nil {
				t.Errorf("GetPackageLockfile() = nil, want %v", tt.want)
				return
			}

			if got.Version != tt.want.Version {
				t.Errorf("GetPackageLockfile() Version = %v, want %v", got.Version, tt.want.Version)
			}

			if len(got.Packages) != len(tt.want.Packages) {
				t.Errorf("GetPackageLockfile() Packages length = %v, want %v", len(got.Packages), len(tt.want.Packages))
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
		wantFile      *PackageLockfile
		wantErr       bool
		errContains   string
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
			wantFile: &PackageLockfile{
				Version: "1.0.0",
				Packages: map[string]map[string]*PackageLockInfo{
					"registry1": {
						"package1": {
							Version:  "1.2.3",
							Checksum: "sha256:abc123",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "success - add package to existing registry",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version: "1.0.0",
					Packages: map[string]map[string]*PackageLockInfo{
						"registry1": {
							"package1": {
								Version:  "1.2.3",
								Checksum: "sha256:abc123",
							},
						},
					},
				}
				return createTestLockfile(t, lockfile)
			},
			registryName: "registry1",
			packageName:  "package2",
			lockInfo: &PackageLockInfo{
				Version:  "2.0.0",
				Checksum: "sha256:def456",
			},
			wantFile: &PackageLockfile{
				Version: "1.0.0",
				Packages: map[string]map[string]*PackageLockInfo{
					"registry1": {
						"package1": {
							Version:  "1.2.3",
							Checksum: "sha256:abc123",
						},
						"package2": {
							Version:  "2.0.0",
							Checksum: "sha256:def456",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "success - update existing package",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version: "1.0.0",
					Packages: map[string]map[string]*PackageLockInfo{
						"registry1": {
							"package1": {
								Version:  "1.2.3",
								Checksum: "sha256:abc123",
							},
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
			wantFile: &PackageLockfile{
				Version: "1.0.0",
				Packages: map[string]map[string]*PackageLockInfo{
					"registry1": {
						"package1": {
							Version:  "1.3.0",
							Checksum: "sha256:xyz789",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "success - create new registry",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version: "1.0.0",
					Packages: map[string]map[string]*PackageLockInfo{
						"registry1": {
							"package1": {
								Version:  "1.2.3",
								Checksum: "sha256:abc123",
							},
						},
					},
				}
				return createTestLockfile(t, lockfile)
			},
			registryName: "registry2",
			packageName:  "package1",
			lockInfo: &PackageLockInfo{
				Version:  "2.0.0",
				Checksum: "sha256:def456",
			},
			wantFile: &PackageLockfile{
				Version: "1.0.0",
				Packages: map[string]map[string]*PackageLockInfo{
					"registry1": {
						"package1": {
							Version:  "1.2.3",
							Checksum: "sha256:abc123",
						},
					},
					"registry2": {
						"package1": {
							Version:  "2.0.0",
							Checksum: "sha256:def456",
						},
					},
				},
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

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpsertPackageLockInfo() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("UpsertPackageLockInfo() error = %v, should contain %v", err, tt.errContains)
				}
				return
			}

			got, err := fm.GetPackageLockfile(ctx)
			if err != nil {
				t.Fatalf("GetPackageLockfile() error = %v", err)
			}

			if got.Version != tt.wantFile.Version {
				t.Errorf("UpsertPackageLockInfo() Version = %v, want %v", got.Version, tt.wantFile.Version)
			}

			if len(got.Packages) != len(tt.wantFile.Packages) {
				t.Errorf("UpsertPackageLockInfo() Packages length = %v, want %v", len(got.Packages), len(tt.wantFile.Packages))
			}

			registry, exists := got.Packages[tt.registryName]
			if !exists {
				t.Errorf("UpsertPackageLockInfo() registry %v not found", tt.registryName)
				return
			}

			pkg, exists := registry[tt.packageName]
			if !exists {
				t.Errorf("UpsertPackageLockInfo() package %v not found in registry %v", tt.packageName, tt.registryName)
				return
			}

			if pkg.Version != tt.lockInfo.Version {
				t.Errorf("UpsertPackageLockInfo() package Version = %v, want %v", pkg.Version, tt.lockInfo.Version)
			}
			if pkg.Checksum != tt.lockInfo.Checksum {
				t.Errorf("UpsertPackageLockInfo() package Checksum = %v, want %v", pkg.Checksum, tt.lockInfo.Checksum)
			}
		})
	}
}

func TestFileManager_UpdatePackageName(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		setupFile     func(t *testing.T) string
		registryName  string
		packageName   string
		newPackageName string
		wantFile      *PackageLockfile
		wantErr       bool
		errContains   string
	}{
		{
			name: "success - rename package",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version: "1.0.0",
					Packages: map[string]map[string]*PackageLockInfo{
						"registry1": {
							"package1": {
								Version:  "1.2.3",
								Checksum: "sha256:abc123",
							},
						},
					},
				}
				return createTestLockfile(t, lockfile)
			},
			registryName:  "registry1",
			packageName:   "package1",
			newPackageName: "package1-renamed",
			wantFile: &PackageLockfile{
				Version: "1.0.0",
				Packages: map[string]map[string]*PackageLockInfo{
					"registry1": {
						"package1-renamed": {
							Version:  "1.2.3",
							Checksum: "sha256:abc123",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error - file does not exist",
			setupFile: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "nonexistent.json")
			},
			registryName:  "registry1",
			packageName:   "package1",
			newPackageName: "package2",
			wantFile:      nil,
			wantErr:       true,
			errContains:   "no such file",
		},
		{
			name: "error - registry not found",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version: "1.0.0",
					Packages: map[string]map[string]*PackageLockInfo{
						"registry1": {
							"package1": {
								Version:  "1.2.3",
								Checksum: "sha256:abc123",
							},
						},
					},
				}
				return createTestLockfile(t, lockfile)
			},
			registryName:  "nonexistent-registry",
			packageName:   "package1",
			newPackageName: "package2",
			wantFile:      nil,
			wantErr:       true,
			errContains:   "registry not found",
		},
		{
			name: "error - package not found",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version: "1.0.0",
					Packages: map[string]map[string]*PackageLockInfo{
						"registry1": {
							"package1": {
								Version:  "1.2.3",
								Checksum: "sha256:abc123",
							},
						},
					},
				}
				return createTestLockfile(t, lockfile)
			},
			registryName:  "registry1",
			packageName:   "nonexistent-package",
			newPackageName: "package2",
			wantFile:      nil,
			wantErr:       true,
			errContains:   "package not found",
		},
		{
			name: "error - new package name already exists",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version: "1.0.0",
					Packages: map[string]map[string]*PackageLockInfo{
						"registry1": {
							"package1": {
								Version:  "1.2.3",
								Checksum: "sha256:abc123",
							},
							"package2": {
								Version:  "2.0.0",
								Checksum: "sha256:def456",
							},
						},
					},
				}
				return createTestLockfile(t, lockfile)
			},
			registryName:  "registry1",
			packageName:   "package1",
			newPackageName: "package2",
			wantFile:      nil,
			wantErr:       true,
			errContains:   "package already exists",
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
				return
			}

			got, err := fm.GetPackageLockfile(ctx)
			if err != nil {
				t.Fatalf("GetPackageLockfile() error = %v", err)
			}

			registry, exists := got.Packages[tt.registryName]
			if !exists {
				t.Errorf("UpdatePackageName() registry %v not found", tt.registryName)
				return
			}

			if _, exists := registry[tt.packageName]; exists {
				t.Errorf("UpdatePackageName() old package name %v still exists", tt.packageName)
			}

			newPkg, exists := registry[tt.newPackageName]
			if !exists {
				t.Errorf("UpdatePackageName() new package name %v not found", tt.newPackageName)
				return
			}

			if newPkg.Version != tt.wantFile.Packages[tt.registryName][tt.newPackageName].Version {
				t.Errorf("UpdatePackageName() package Version = %v, want %v", newPkg.Version, tt.wantFile.Packages[tt.registryName][tt.newPackageName].Version)
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
		wantFile     *PackageLockfile
		wantFileExists bool
		wantErr      bool
		errContains  string
	}{
		{
			name: "success - remove package",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version: "1.0.0",
					Packages: map[string]map[string]*PackageLockInfo{
						"registry1": {
							"package1": {
								Version:  "1.2.3",
								Checksum: "sha256:abc123",
							},
							"package2": {
								Version:  "2.0.0",
								Checksum: "sha256:def456",
							},
						},
					},
				}
				return createTestLockfile(t, lockfile)
			},
			registryName: "registry1",
			packageName:  "package1",
			wantFile: &PackageLockfile{
				Version: "1.0.0",
				Packages: map[string]map[string]*PackageLockInfo{
					"registry1": {
						"package2": {
							Version:  "2.0.0",
							Checksum: "sha256:def456",
						},
					},
				},
			},
			wantFileExists: true,
			wantErr:        false,
		},
		{
			name: "success - remove last package in registry, cleanup empty registry",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version: "1.0.0",
					Packages: map[string]map[string]*PackageLockInfo{
						"registry1": {
							"package1": {
								Version:  "1.2.3",
								Checksum: "sha256:abc123",
							},
						},
					},
				}
				return createTestLockfile(t, lockfile)
			},
			registryName: "registry1",
			packageName:  "package1",
			wantFile:     nil,
			wantFileExists: false,
			wantErr:        false,
		},
		{
			name: "success - remove last package, delete lockfile",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version: "1.0.0",
					Packages: map[string]map[string]*PackageLockInfo{
						"registry1": {
							"package1": {
								Version:  "1.2.3",
								Checksum: "sha256:abc123",
							},
						},
					},
				}
				return createTestLockfile(t, lockfile)
			},
			registryName: "registry1",
			packageName:  "package1",
			wantFile:     nil,
			wantFileExists: false,
			wantErr:        false,
		},
		{
			name: "error - file does not exist",
			setupFile: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "nonexistent.json")
			},
			registryName: "registry1",
			packageName:  "package1",
			wantFile:     nil,
			wantFileExists: false,
			wantErr:        true,
			errContains:   "no such file",
		},
		{
			name: "error - registry not found",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version: "1.0.0",
					Packages: map[string]map[string]*PackageLockInfo{
						"registry1": {
							"package1": {
								Version:  "1.2.3",
								Checksum: "sha256:abc123",
							},
						},
					},
				}
				return createTestLockfile(t, lockfile)
			},
			registryName: "nonexistent-registry",
			packageName:  "package1",
			wantFile:     nil,
			wantFileExists: true,
			wantErr:        true,
			errContains:   "registry not found",
		},
		{
			name: "error - package not found",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version: "1.0.0",
					Packages: map[string]map[string]*PackageLockInfo{
						"registry1": {
							"package1": {
								Version:  "1.2.3",
								Checksum: "sha256:abc123",
							},
						},
					},
				}
				return createTestLockfile(t, lockfile)
			},
			registryName: "registry1",
			packageName:  "nonexistent-package",
			wantFile:     nil,
			wantFileExists: true,
			wantErr:        true,
			errContains:   "package not found",
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
				return
			}

			fileExists := fileExists(lockPath)
			if fileExists != tt.wantFileExists {
				t.Errorf("RemovePackageLockInfo() file exists = %v, want %v", fileExists, tt.wantFileExists)
			}

			if !tt.wantFileExists {
				return
			}

			got, err := fm.GetPackageLockfile(ctx)
			if err != nil {
				t.Fatalf("GetPackageLockfile() error = %v", err)
			}

			if len(got.Packages) != len(tt.wantFile.Packages) {
				t.Errorf("RemovePackageLockInfo() Packages length = %v, want %v", len(got.Packages), len(tt.wantFile.Packages))
			}

			if _, exists := got.Packages[tt.registryName]; exists && tt.wantFile != nil {
				if _, exists := got.Packages[tt.registryName][tt.packageName]; exists {
					t.Errorf("RemovePackageLockInfo() package %v still exists", tt.packageName)
				}
			}
		})
	}
}

func TestFileManager_UpdateRegistryName(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		setupFile     func(t *testing.T) string
		oldName       string
		newName       string
		wantFile      *PackageLockfile
		wantErr       bool
		errContains   string
	}{
		{
			name: "success - rename registry",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version: "1.0.0",
					Packages: map[string]map[string]*PackageLockInfo{
						"registry1": {
							"package1": {
								Version:  "1.2.3",
								Checksum: "sha256:abc123",
							},
							"package2": {
								Version:  "2.0.0",
								Checksum: "sha256:def456",
							},
						},
					},
				}
				return createTestLockfile(t, lockfile)
			},
			oldName: "registry1",
			newName: "registry1-renamed",
			wantFile: &PackageLockfile{
				Version: "1.0.0",
				Packages: map[string]map[string]*PackageLockInfo{
					"registry1-renamed": {
						"package1": {
							Version:  "1.2.3",
							Checksum: "sha256:abc123",
						},
						"package2": {
							Version:  "2.0.0",
							Checksum: "sha256:def456",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error - file does not exist",
			setupFile: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "nonexistent.json")
			},
			oldName:     "registry1",
			newName:     "registry2",
			wantFile:    nil,
			wantErr:     true,
			errContains: "no such file",
		},
		{
			name: "error - old registry not found",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version: "1.0.0",
					Packages: map[string]map[string]*PackageLockInfo{
						"registry1": {
							"package1": {
								Version:  "1.2.3",
								Checksum: "sha256:abc123",
							},
						},
					},
				}
				return createTestLockfile(t, lockfile)
			},
			oldName:     "nonexistent-registry",
			newName:     "registry2",
			wantFile:    nil,
			wantErr:     true,
			errContains: "registry not found",
		},
		{
			name: "error - new registry name already exists",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version: "1.0.0",
					Packages: map[string]map[string]*PackageLockInfo{
						"registry1": {
							"package1": {
								Version:  "1.2.3",
								Checksum: "sha256:abc123",
							},
						},
						"registry2": {
							"package1": {
								Version:  "2.0.0",
								Checksum: "sha256:def456",
							},
						},
					},
				}
				return createTestLockfile(t, lockfile)
			},
			oldName:     "registry1",
			newName:     "registry2",
			wantFile:    nil,
			wantErr:     true,
			errContains: "registry already exists",
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
				return
			}

			got, err := fm.GetPackageLockfile(ctx)
			if err != nil {
				t.Fatalf("GetPackageLockfile() error = %v", err)
			}

			if _, exists := got.Packages[tt.oldName]; exists {
				t.Errorf("UpdateRegistryName() old registry name %v still exists", tt.oldName)
			}

			newRegistry, exists := got.Packages[tt.newName]
			if !exists {
				t.Errorf("UpdateRegistryName() new registry name %v not found", tt.newName)
				return
			}

			if len(newRegistry) != len(tt.wantFile.Packages[tt.newName]) {
				t.Errorf("UpdateRegistryName() packages count = %v, want %v", len(newRegistry), len(tt.wantFile.Packages[tt.newName]))
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

