package packagelockfile

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileManager_GetRulesetLockInfo(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		setupFile   func(t *testing.T) string
		packageKey  string
		want        *PackageLockInfo
		wantErr     bool
		errContains string
	}{
		{
			name: "success - ruleset exists",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version: 1,
					Rulesets: map[string]*PackageLockInfo{
						"registry1/package1@1.2.3": {
							Integrity: "sha256-abc123",
						},
					},
				}
				return createTestLockfile(t, lockfile)
			},
			packageKey: "registry1/package1@1.2.3",
			want: &PackageLockInfo{
				Integrity: "sha256-abc123",
			},
			wantErr: false,
		},
		{
			name: "error - file does not exist",
			setupFile: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "nonexistent.json")
			},
			packageKey:  "registry1/package1@1.2.3",
			want:        nil,
			wantErr:     true,
			errContains: "no such file",
		},
		{
			name: "error - ruleset not found",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version:  1,
					Rulesets: map[string]*PackageLockInfo{},
				}
				return createTestLockfile(t, lockfile)
			},
			packageKey:  "registry1/nonexistent@1.0.0",
			want:        nil,
			wantErr:     true,
			errContains: "ruleset not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lockPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(lockPath)

			got, err := fm.GetRulesetLockInfo(ctx, tt.packageKey)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetRulesetLockInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetRulesetLockInfo() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("GetRulesetLockInfo() error = %v, should contain %v", err, tt.errContains)
				}
				return
			}

			if got == nil {
				t.Errorf("GetRulesetLockInfo() = nil, want %v", tt.want)
				return
			}

			if got.Integrity != tt.want.Integrity {
				t.Errorf("GetRulesetLockInfo() Integrity = %v, want %v", got.Integrity, tt.want.Integrity)
			}
		})
	}
}

func TestFileManager_UpsertRulesetLockInfo(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		setupFile  func(t *testing.T) string
		packageKey string
		lockInfo   *PackageLockInfo
		wantErr    bool
	}{
		{
			name: "success - create new lockfile",
			setupFile: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "new-lock.json")
			},
			packageKey: "registry1/package1@1.2.3",
			lockInfo: &PackageLockInfo{
				Integrity: "sha256-abc123",
			},
			wantErr: false,
		},
		{
			name: "success - update existing ruleset",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version: 1,
					Rulesets: map[string]*PackageLockInfo{
						"registry1/package1@1.2.3": {
							Integrity: "sha256-abc123",
						},
					},
				}
				return createTestLockfile(t, lockfile)
			},
			packageKey: "registry1/package1@1.3.0",
			lockInfo: &PackageLockInfo{
				Integrity: "sha256-xyz789",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lockPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(lockPath)

			err := fm.UpsertRulesetLockInfo(ctx, tt.packageKey, tt.lockInfo)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpsertRulesetLockInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				got, err := fm.GetRulesetLockInfo(ctx, tt.packageKey)
				if err != nil {
					t.Fatalf("GetRulesetLockInfo() error = %v", err)
				}

				if got.Integrity != tt.lockInfo.Integrity {
					t.Errorf("UpsertRulesetLockInfo() Integrity = %v, want %v", got.Integrity, tt.lockInfo.Integrity)
				}
			}
		})
	}
}

func TestFileManager_RemoveRulesetLockInfo(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		setupFile   func(t *testing.T) string
		packageKey  string
		wantErr     bool
		errContains string
	}{
		{
			name: "success - remove ruleset",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version: 1,
					Rulesets: map[string]*PackageLockInfo{
						"registry1/package1@1.2.3": {
							Integrity: "sha256-abc123",
						},
						"registry1/package2@2.0.0": {
							Integrity: "sha256-def456",
						},
					},
				}
				return createTestLockfile(t, lockfile)
			},
			packageKey: "registry1/package1@1.2.3",
			wantErr:    false,
		},
		{
			name: "error - ruleset not found",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version:  1,
					Rulesets: map[string]*PackageLockInfo{},
				}
				return createTestLockfile(t, lockfile)
			},
			packageKey:  "registry1/nonexistent@1.0.0",
			wantErr:     true,
			errContains: "ruleset not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lockPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(lockPath)

			err := fm.RemoveRulesetLockInfo(ctx, tt.packageKey)

			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveRulesetLockInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("RemoveRulesetLockInfo() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("RemoveRulesetLockInfo() error = %v, should contain %v", err, tt.errContains)
				}
			}
		})
	}
}

func TestFileManager_GetPromptsetLockInfo(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		setupFile   func(t *testing.T) string
		packageKey  string
		want        *PackageLockInfo
		wantErr     bool
		errContains string
	}{
		{
			name: "success - promptset exists",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version: 1,
					Promptsets: map[string]*PackageLockInfo{
						"registry1/prompts@2.1.0": {
							Integrity: "sha256-xyz789",
						},
					},
				}
				return createTestLockfile(t, lockfile)
			},
			packageKey: "registry1/prompts@2.1.0",
			want: &PackageLockInfo{
				Integrity: "sha256-xyz789",
			},
			wantErr: false,
		},
		{
			name: "error - promptset not found",
			setupFile: func(t *testing.T) string {
				lockfile := &PackageLockfile{
					Version:    1,
					Promptsets: map[string]*PackageLockInfo{},
				}
				return createTestLockfile(t, lockfile)
			},
			packageKey:  "registry1/nonexistent@1.0.0",
			want:        nil,
			wantErr:     true,
			errContains: "promptset not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lockPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(lockPath)

			got, err := fm.GetPromptsetLockInfo(ctx, tt.packageKey)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetPromptsetLockInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetPromptsetLockInfo() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("GetPromptsetLockInfo() error = %v, should contain %v", err, tt.errContains)
				}
				return
			}

			if got.Integrity != tt.want.Integrity {
				t.Errorf("GetPromptsetLockInfo() Integrity = %v, want %v", got.Integrity, tt.want.Integrity)
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

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}