package packagelockfile

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileManager_UpsertDependencyLock(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		key     string
		config  *DependencyLockConfig
		wantErr bool
	}{
		{
			name: "create new lockfile",
			key:  "myregistry/clean-code@1.2.0",
			config: &DependencyLockConfig{
				Integrity: "sha256-abc123",
			},
			wantErr: false,
		},
		{
			name: "update existing dependency",
			key:  "myregistry/clean-code@1.3.0",
			config: &DependencyLockConfig{
				Integrity: "sha256-xyz789",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lockPath := filepath.Join(t.TempDir(), "test-lock.json")
			fm := NewFileManagerWithPath(lockPath)

			err := fm.UpsertDependencyLock(ctx, tt.key, tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpsertDependencyLock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify the dependency was stored
				got, err := fm.GetDependencyLock(ctx, tt.key)
				if err != nil {
					t.Fatalf("GetDependencyLock() error = %v", err)
				}
				if got.Integrity != tt.config.Integrity {
					t.Errorf("GetDependencyLock() Integrity = %v, want %v", got.Integrity, tt.config.Integrity)
				}
			}
		})
	}
}

func TestFileManager_GetDependencyLock(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		setupFile   func(t *testing.T) string
		key         string
		want        *DependencyLockConfig
		wantErr     bool
		errContains string
	}{
		{
			name: "dependency exists",
			setupFile: func(t *testing.T) string {
				lockfile := &LockFile{
					Version: 1,
					Dependencies: map[string]DependencyLockConfig{
						"myregistry/clean-code@1.2.0": {
							Integrity: "sha256-abc123",
						},
					},
				}
				return createTestLockfile(t, lockfile)
			},
			key: "myregistry/clean-code@1.2.0",
			want: &DependencyLockConfig{
				Integrity: "sha256-abc123",
			},
			wantErr: false,
		},
		{
			name: "file does not exist",
			setupFile: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "nonexistent.json")
			},
			key:         "myregistry/clean-code@1.2.0",
			want:        nil,
			wantErr:     true,
			errContains: "no such file",
		},
		{
			name: "dependency not found",
			setupFile: func(t *testing.T) string {
				lockfile := &LockFile{
					Version:      1,
					Dependencies: map[string]DependencyLockConfig{},
				}
				return createTestLockfile(t, lockfile)
			},
			key:         "myregistry/nonexistent@1.0.0",
			want:        nil,
			wantErr:     true,
			errContains: "dependency not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lockPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(lockPath)

			got, err := fm.GetDependencyLock(ctx, tt.key)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetDependencyLock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetDependencyLock() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("GetDependencyLock() error = %v, should contain %v", err, tt.errContains)
				}
				return
			}

			if got == nil {
				t.Errorf("GetDependencyLock() = nil, want %v", tt.want)
				return
			}

			if got.Integrity != tt.want.Integrity {
				t.Errorf("GetDependencyLock() Integrity = %v, want %v", got.Integrity, tt.want.Integrity)
			}
		})
	}
}

func TestFileManager_RemoveDependencyLock(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		setupFile   func(t *testing.T) string
		key         string
		wantErr     bool
		errContains string
	}{
		{
			name: "remove dependency",
			setupFile: func(t *testing.T) string {
				lockfile := &LockFile{
					Version: 1,
					Dependencies: map[string]DependencyLockConfig{
						"myregistry/clean-code@1.2.0": {
							Integrity: "sha256-abc123",
						},
						"myregistry/other@2.0.0": {
							Integrity: "sha256-def456",
						},
					},
				}
				return createTestLockfile(t, lockfile)
			},
			key:     "myregistry/clean-code@1.2.0",
			wantErr: false,
		},
		{
			name: "dependency not found",
			setupFile: func(t *testing.T) string {
				lockfile := &LockFile{
					Version:      1,
					Dependencies: map[string]DependencyLockConfig{},
				}
				return createTestLockfile(t, lockfile)
			},
			key:         "myregistry/nonexistent@1.0.0",
			wantErr:     true,
			errContains: "dependency not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lockPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(lockPath)

			err := fm.RemoveDependencyLock(ctx, tt.key)

			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveDependencyLock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("RemoveDependencyLock() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("RemoveDependencyLock() error = %v, should contain %v", err, tt.errContains)
				}
			}
		})
	}
}

func TestFileManager_UpdateRegistryName(t *testing.T) {
	ctx := context.Background()

	lockfile := &LockFile{
		Version: 1,
		Dependencies: map[string]DependencyLockConfig{
			"oldregistry/clean-code@1.2.0": {
				Integrity: "sha256-abc123",
			},
			"otherregistry/other@2.0.0": {
				Integrity: "sha256-def456",
			},
		},
	}
	lockPath := createTestLockfile(t, lockfile)
	fm := NewFileManagerWithPath(lockPath)

	err := fm.UpdateRegistryName(ctx, "oldregistry", "newregistry")
	if err != nil {
		t.Fatalf("UpdateRegistryName() error = %v", err)
	}

	// Verify old key is gone and new key exists
	_, err = fm.GetDependencyLock(ctx, "oldregistry/clean-code@1.2.0")
	if err == nil {
		t.Error("Expected old key to be removed")
	}

	got, err := fm.GetDependencyLock(ctx, "newregistry/clean-code@1.2.0")
	if err != nil {
		t.Fatalf("GetDependencyLock() error = %v", err)
	}
	if got.Integrity != "sha256-abc123" {
		t.Errorf("GetDependencyLock() Integrity = %v, want %v", got.Integrity, "sha256-abc123")
	}

	// Verify other registry unchanged
	other, err := fm.GetDependencyLock(ctx, "otherregistry/other@2.0.0")
	if err != nil {
		t.Fatalf("GetDependencyLock() error = %v", err)
	}
	if other.Integrity != "sha256-def456" {
		t.Errorf("GetDependencyLock() Integrity = %v, want %v", other.Integrity, "sha256-def456")
	}
}

// Helper functions

func createTestLockfile(t *testing.T, lockfile *LockFile) string {
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

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}