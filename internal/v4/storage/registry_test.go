package storage

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestNewRegistry(t *testing.T) {
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
			registry, err := NewRegistry(tt.registryKey)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewRegistry() expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("NewRegistry() unexpected error: %v", err)
				return
			}
			
			if registry == nil {
				t.Errorf("NewRegistry() returned nil registry")
				return
			}
			
			// Check registry directory exists
			registryDir := registry.GetRegistryDir()
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

func TestNewRegistryWithPath(t *testing.T) {
	tempDir := t.TempDir()
	registryKey := map[string]interface{}{"url": "https://github.com/user/repo", "type": "git"}
	
	registry, err := NewRegistryWithPath(tempDir, registryKey)
	if err != nil {
		t.Fatalf("NewRegistryWithPath() unexpected error: %v", err)
	}
	
	if registry == nil {
		t.Fatal("NewRegistryWithPath() returned nil registry")
	}
	
	// Check registry directory is under tempDir
	registryDir := registry.GetRegistryDir()
	if !filepath.HasPrefix(registryDir, tempDir) {
		t.Errorf("Registry directory not under tempDir: %s", registryDir)
	}
	
	// Check directory exists
	if _, err := os.Stat(registryDir); os.IsNotExist(err) {
		t.Errorf("Registry directory not created: %s", registryDir)
	}
}

func TestRegistry_GetPaths(t *testing.T) {
	tempDir := t.TempDir()
	registryKey := map[string]interface{}{"url": "https://github.com/user/repo", "type": "git"}
	
	registry, err := NewRegistryWithPath(tempDir, registryKey)
	if err != nil {
		t.Fatalf("NewRegistryWithPath() unexpected error: %v", err)
	}
	
	registryDir := registry.GetRegistryDir()
	repoDir := registry.GetRepoDir()
	packagesDir := registry.GetPackagesDir()
	
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

func TestRegistry_ConsistentKeys(t *testing.T) {
	tempDir := t.TempDir()
	registryKey := map[string]interface{}{"url": "https://github.com/user/repo", "type": "git"}
	
	// Create registry multiple times with same key
	registry1, err1 := NewRegistryWithPath(tempDir, registryKey)
	registry2, err2 := NewRegistryWithPath(tempDir, registryKey)
	
	if err1 != nil || err2 != nil {
		t.Fatalf("NewRegistryWithPath() unexpected errors: %v, %v", err1, err2)
	}
	
	// Should get same registry directory
	dir1 := registry1.GetRegistryDir()
	dir2 := registry2.GetRegistryDir()
	
	if dir1 != dir2 {
		t.Errorf("Same registry key produced different directories: %s vs %s", dir1, dir2)
	}
}

func TestRegistry_DifferentKeys(t *testing.T) {
	tempDir := t.TempDir()
	registryKey1 := map[string]interface{}{"url": "https://github.com/user/repo1", "type": "git"}
	registryKey2 := map[string]interface{}{"url": "https://github.com/user/repo2", "type": "git"}
	
	registry1, err1 := NewRegistryWithPath(tempDir, registryKey1)
	registry2, err2 := NewRegistryWithPath(tempDir, registryKey2)
	
	if err1 != nil || err2 != nil {
		t.Fatalf("NewRegistryWithPath() unexpected errors: %v, %v", err1, err2)
	}
	
	// Should get different registry directories
	dir1 := registry1.GetRegistryDir()
	dir2 := registry2.GetRegistryDir()
	
	if dir1 == dir2 {
		t.Errorf("Different registry keys produced same directory: %s", dir1)
	}
}