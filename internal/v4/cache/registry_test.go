package cache

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestNewRegistryCache(t *testing.T) {
	tests := []struct {
		name        string
		registryKey interface{}
		wantErr     bool
	}{
		{
			name:        "simple registry key",
			registryKey: map[string]interface{}{"url": "https://github.com/user/repo", "type": "git"},
			wantErr:     false,
		},
		{
			name:        "complex registry key",
			registryKey: map[string]interface{}{"url": "https://gitlab.com", "type": "gitlab", "project_id": "123"},
			wantErr:     false,
		},
		{
			name:        "nil registry key",
			registryKey: nil,
			wantErr:     false,
		},
		{
			name:        "unmarshalable registry key",
			registryKey: make(chan int),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache, err := NewRegistryCache(tt.registryKey)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewRegistryCache() expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("NewRegistryCache() unexpected error: %v", err)
				return
			}
			
			if cache == nil {
				t.Errorf("NewRegistryCache() returned nil cache")
				return
			}
			
			// Check registry directory exists
			registryDir := cache.GetRegistryDir()
			if registryDir == "" {
				t.Errorf("GetRegistryDir() returned empty path")
			}
			
			if _, err := os.Stat(registryDir); os.IsNotExist(err) {
				t.Errorf("Registry directory not created: %s", registryDir)
			}
			
			// Check metadata.json exists
			metadataPath := filepath.Join(registryDir, "metadata.json")
			if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
				t.Errorf("Registry metadata.json not created: %s", metadataPath)
			}
			
			// Cleanup
			os.RemoveAll(registryDir)
		})
	}
}

func TestNewRegistryCacheWithPath(t *testing.T) {
	tempDir := t.TempDir()
	registryKey := map[string]interface{}{"url": "https://github.com/user/repo", "type": "git"}
	
	cache, err := NewRegistryCacheWithPath(tempDir, registryKey)
	if err != nil {
		t.Fatalf("NewRegistryCacheWithPath() unexpected error: %v", err)
	}
	
	if cache == nil {
		t.Fatal("NewRegistryCacheWithPath() returned nil cache")
	}
	
	// Check registry directory is under tempDir
	registryDir := cache.GetRegistryDir()
	if !filepath.HasPrefix(registryDir, tempDir) {
		t.Errorf("Registry directory not under tempDir: %s", registryDir)
	}
	
	// Check directory exists
	if _, err := os.Stat(registryDir); os.IsNotExist(err) {
		t.Errorf("Registry directory not created: %s", registryDir)
	}
}

func TestRegistryCache_GetPaths(t *testing.T) {
	tempDir := t.TempDir()
	registryKey := map[string]interface{}{"url": "https://github.com/user/repo", "type": "git"}
	
	cache, err := NewRegistryCacheWithPath(tempDir, registryKey)
	if err != nil {
		t.Fatalf("NewRegistryCacheWithPath() unexpected error: %v", err)
	}
	
	registryDir := cache.GetRegistryDir()
	repoDir := cache.GetRepoDir()
	packagesDir := cache.GetPackagesDir()
	
	// Check paths are correct
	expectedRepoDir := filepath.Join(registryDir, "repo")
	if repoDir != expectedRepoDir {
		t.Errorf("GetRepoDir() = %s, want %s", repoDir, expectedRepoDir)
	}
	
	expectedPackagesDir := filepath.Join(registryDir, "packages")
	if packagesDir != expectedPackagesDir {
		t.Errorf("GetPackagesDir() = %s, want %s", packagesDir, expectedPackagesDir)
	}
}

func TestRegistryCache_ConsistentKeys(t *testing.T) {
	tempDir := t.TempDir()
	registryKey := map[string]interface{}{"url": "https://github.com/user/repo", "type": "git"}
	
	// Create cache multiple times with same key
	cache1, err1 := NewRegistryCacheWithPath(tempDir, registryKey)
	cache2, err2 := NewRegistryCacheWithPath(tempDir, registryKey)
	
	if err1 != nil || err2 != nil {
		t.Fatalf("NewRegistryCacheWithPath() unexpected errors: %v, %v", err1, err2)
	}
	
	// Should get same registry directory
	dir1 := cache1.GetRegistryDir()
	dir2 := cache2.GetRegistryDir()
	
	if dir1 != dir2 {
		t.Errorf("Same registry key produced different directories: %s vs %s", dir1, dir2)
	}
}

func TestRegistryCache_DifferentKeys(t *testing.T) {
	tempDir := t.TempDir()
	registryKey1 := map[string]interface{}{"url": "https://github.com/user/repo1", "type": "git"}
	registryKey2 := map[string]interface{}{"url": "https://github.com/user/repo2", "type": "git"}
	
	cache1, err1 := NewRegistryCacheWithPath(tempDir, registryKey1)
	cache2, err2 := NewRegistryCacheWithPath(tempDir, registryKey2)
	
	if err1 != nil || err2 != nil {
		t.Fatalf("NewRegistryCacheWithPath() unexpected errors: %v, %v", err1, err2)
	}
	
	// Should get different registry directories
	dir1 := cache1.GetRegistryDir()
	dir2 := cache2.GetRegistryDir()
	
	if dir1 == dir2 {
		t.Errorf("Different registry keys produced same directory: %s", dir1)
	}
}

func TestRegistryCache_UpdateAccessTime(t *testing.T) {
	tempDir := t.TempDir()
	registryKey := map[string]interface{}{"url": "https://github.com/user/repo", "type": "git"}
	
	cache, err := NewRegistryCacheWithPath(tempDir, registryKey)
	if err != nil {
		t.Fatalf("NewRegistryCacheWithPath() unexpected error: %v", err)
	}
	
	ctx := context.Background()
	err = cache.UpdateAccessTime(ctx)
	if err != nil {
		t.Errorf("UpdateAccessTime() unexpected error: %v", err)
	}
	
	// TODO: verify metadata.json was updated with new access time
}

func TestRegistryCache_UpdateUpdatedTime(t *testing.T) {
	tempDir := t.TempDir()
	registryKey := map[string]interface{}{"url": "https://github.com/user/repo", "type": "git"}
	
	cache, err := NewRegistryCacheWithPath(tempDir, registryKey)
	if err != nil {
		t.Fatalf("NewRegistryCacheWithPath() unexpected error: %v", err)
	}
	
	ctx := context.Background()
	err = cache.UpdateUpdatedTime(ctx)
	if err != nil {
		t.Errorf("UpdateUpdatedTime() unexpected error: %v", err)
	}
	
	// TODO: verify metadata.json was updated with new updated time
}