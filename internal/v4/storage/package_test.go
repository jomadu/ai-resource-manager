package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jomadu/ai-resource-manager/internal/v4/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPackageCache(t *testing.T) {
	tempDir := t.TempDir()
	packagesDir := filepath.Join(tempDir, "packages")
	
	// Create packages directory
	err := os.MkdirAll(packagesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create packages directory: %v", err)
	}
	
	pkg := NewPackageCache(packagesDir)
	if pkg == nil {
		t.Fatal("NewPackageCache() returned nil")
	}
}

func TestPackageCache_SetPackageVersion(t *testing.T) {
	tempDir := t.TempDir()
	packagesDir := filepath.Join(tempDir, "packages")
	pkg := NewPackageCache(packagesDir)
	
	packageKey := map[string]interface{}{
		"name":     "clean-code",
		"includes": []string{"**/*.yml"},
		"excludes": []string{"**/test/**"},
	}
	version := core.Version{Major: 1, Minor: 0, Patch: 0}
	files := []*core.File{
		{Path: "rules.yml", Content: []byte("rule: value"), Size: 11},
		{Path: "nested/rule.yml", Content: []byte("nested: rule"), Size: 12},
	}
	
	ctx := context.Background()
	err := pkg.SetPackageVersion(ctx, packageKey, version, files)
	if err != nil {
		t.Errorf("SetPackageVersion() unexpected error: %v", err)
	}
	
	// Verify package directory was created
	packageHash, _ := GenerateKey(packageKey)
	packageDir := filepath.Join(packagesDir, packageHash)
	if _, err := os.Stat(packageDir); os.IsNotExist(err) {
		t.Errorf("Package directory not created: %s", packageDir)
	}
	
	// Verify version directory was created
	versionDir := filepath.Join(packageDir, fmt.Sprintf("v%d.%d.%d", version.Major, version.Minor, version.Patch))
	if _, err := os.Stat(versionDir); os.IsNotExist(err) {
		t.Errorf("Version directory not created: %s", versionDir)
	}
	
	// Verify files were stored
	filesDir := filepath.Join(versionDir, "files")
	for _, file := range files {
		filePath := filepath.Join(filesDir, file.Path)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("File not stored: %s", filePath)
		}
		
		// Verify file content
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Errorf("Failed to read file %s: %v", filePath, err)
		}
		if string(content) != string(file.Content) {
			t.Errorf("File content mismatch for %s: got %s, want %s", 
				filePath, string(content), string(file.Content))
		}
	}
	
	// Verify package metadata.json was created
	packageMetadataPath := filepath.Join(packageDir, "metadata.json")
	if _, err := os.Stat(packageMetadataPath); os.IsNotExist(err) {
		t.Errorf("Package metadata.json not created: %s", packageMetadataPath)
	}
	
	// Verify version metadata.json was created
	versionMetadataPath := filepath.Join(versionDir, "metadata.json")
	if _, err := os.Stat(versionMetadataPath); os.IsNotExist(err) {
		t.Errorf("Version metadata.json not created: %s", versionMetadataPath)
	}
}

func TestPackageCache_GetPackageVersion(t *testing.T) {
	tempDir := t.TempDir()
	packagesDir := filepath.Join(tempDir, "packages")
	pkg := NewPackageCache(packagesDir)
	
	packageKey := map[string]interface{}{
		"name":     "clean-code",
		"includes": []string{"**/*.yml"},
	}
	version := core.Version{Major: 1, Minor: 0, Patch: 0}
	originalFiles := []*core.File{
		{Path: "rules.yml", Content: []byte("rule: value"), Size: 11},
		{Path: "nested/rule.yml", Content: []byte("nested: rule"), Size: 12},
	}
	
	ctx := context.Background()
	
	// First set the package version
	err := pkg.SetPackageVersion(ctx, packageKey, version, originalFiles)
	if err != nil {
		t.Fatalf("SetPackageVersion() failed: %v", err)
	}
	
	// Then get it back
	retrievedFiles, err := pkg.GetPackageVersion(ctx, packageKey, version)
	if err != nil {
		t.Errorf("GetPackageVersion() unexpected error: %v", err)
	}
	
	if len(retrievedFiles) != len(originalFiles) {
		t.Errorf("GetPackageVersion() returned %d files, want %d", 
			len(retrievedFiles), len(originalFiles))
	}
	
	// Verify file contents match (order independent)
	fileMap := make(map[string]*core.File)
	for _, file := range retrievedFiles {
		fileMap[file.Path] = file
	}
	
	for _, originalFile := range originalFiles {
		retrievedFile, exists := fileMap[originalFile.Path]
		if !exists {
			t.Errorf("Missing file: %s", originalFile.Path)
			continue
		}
		
		if string(retrievedFile.Content) != string(originalFile.Content) {
			t.Errorf("File content mismatch for %s: got %s, want %s",
				retrievedFile.Path, string(retrievedFile.Content), string(originalFile.Content))
		}
	}
}

func TestPackageCache_GetPackageVersion_NotFound(t *testing.T) {
	tempDir := t.TempDir()
	packagesDir := filepath.Join(tempDir, "packages")
	pkg := NewPackageCache(packagesDir)
	
	packageKey := map[string]interface{}{"name": "nonexistent"}
	version := core.Version{Major: 1, Minor: 0, Patch: 0}
	
	ctx := context.Background()
	files, err := pkg.GetPackageVersion(ctx, packageKey, version)
	
	if err == nil {
		t.Errorf("GetPackageVersion() expected error for nonexistent package")
	}
	
	if files != nil {
		t.Errorf("GetPackageVersion() returned files for nonexistent package")
	}
}

func TestPackageCache_ListPackageVersions(t *testing.T) {
	tempDir := t.TempDir()
	packagesDir := filepath.Join(tempDir, "packages")
	pkg := NewPackageCache(packagesDir)
	
	packageKey := map[string]interface{}{"name": "clean-code"}
	versions := []core.Version{
		{Major: 1, Minor: 0, Patch: 0},
		{Major: 1, Minor: 1, Patch: 0},
		{Major: 2, Minor: 0, Patch: 0},
	}
	files := []*core.File{{Path: "test.yml", Content: []byte("test"), Size: 4}}
	
	ctx := context.Background()
	
	// Set multiple versions
	for _, version := range versions {
		err := pkg.SetPackageVersion(ctx, packageKey, version, files)
		if err != nil {
			t.Fatalf("SetPackageVersion() failed for %v: %v", version, err)
		}
	}
	
	// List versions
	listedVersions, err := pkg.ListPackageVersions(ctx, packageKey)
	if err != nil {
		t.Errorf("ListPackageVersions() unexpected error: %v", err)
	}
	
	if len(listedVersions) != len(versions) {
		t.Errorf("ListPackageVersions() returned %d versions, want %d", 
			len(listedVersions), len(versions))
	}
	
	// Verify all versions are present (order may vary)
	versionMap := make(map[core.Version]bool)
	for _, v := range listedVersions {
		versionMap[v] = true
	}
	
	for _, expectedVersion := range versions {
		if !versionMap[expectedVersion] {
			t.Errorf("ListPackageVersions() missing version: %v", expectedVersion)
		}
	}
}

func TestPackageCache_ListPackages(t *testing.T) {
	tempDir := t.TempDir()
	packagesDir := filepath.Join(tempDir, "packages")
	pkg := NewPackageCache(packagesDir)
	
	packageKeys := []interface{}{
		map[string]interface{}{"name": "clean-code"},
		map[string]interface{}{"name": "security"},
		map[string]interface{}{"name": "typescript"},
	}
	version := core.Version{Major: 1, Minor: 0, Patch: 0}
	files := []*core.File{{Path: "test.yml", Content: []byte("test"), Size: 4}}
	
	ctx := context.Background()
	
	// Set multiple packages
	for _, packageKey := range packageKeys {
		err := pkg.SetPackageVersion(ctx, packageKey, version, files)
		if err != nil {
			t.Fatalf("SetPackageVersion() failed: %v", err)
		}
	}
	
	// List packages
	listedPackages, err := pkg.ListPackages(ctx)
	if err != nil {
		t.Errorf("ListPackages() unexpected error: %v", err)
	}
	
	if len(listedPackages) != len(packageKeys) {
		t.Errorf("ListPackages() returned %d packages, want %d", 
			len(listedPackages), len(packageKeys))
	}
	
	// Verify packages contain expected data (basic check)
	for _, pkg := range listedPackages {
		if pkg == nil {
			t.Errorf("ListPackages() returned nil package")
		}
	}
}

func TestPackageCache_RemovePackageVersion(t *testing.T) {
	tempDir := t.TempDir()
	packagesDir := filepath.Join(tempDir, "packages")
	pkg := NewPackageCache(packagesDir)
	
	packageKey := map[string]interface{}{"name": "clean-code"}
	version := core.Version{Major: 1, Minor: 0, Patch: 0}
	files := []*core.File{{Path: "test.yml", Content: []byte("test"), Size: 4}}
	
	ctx := context.Background()
	
	// Set package version
	err := pkg.SetPackageVersion(ctx, packageKey, version, files)
	if err != nil {
		t.Fatalf("SetPackageVersion() failed: %v", err)
	}
	
	// Remove package version
	err = pkg.RemovePackageVersion(ctx, packageKey, version)
	if err != nil {
		t.Errorf("RemovePackageVersion() unexpected error: %v", err)
	}
	
	// Verify version is gone
	_, err = pkg.GetPackageVersion(ctx, packageKey, version)
	if err == nil {
		t.Errorf("GetPackageVersion() should fail after RemovePackageVersion()")
	}
}

func TestPackageCache_RemovePackage(t *testing.T) {
	tempDir := t.TempDir()
	packagesDir := filepath.Join(tempDir, "packages")
	pkg := NewPackageCache(packagesDir)
	
	packageKey := map[string]interface{}{"name": "clean-code"}
	versions := []core.Version{
		{Major: 1, Minor: 0, Patch: 0},
		{Major: 1, Minor: 1, Patch: 0},
	}
	files := []*core.File{{Path: "test.yml", Content: []byte("test"), Size: 4}}
	
	ctx := context.Background()
	
	// Set multiple versions
	for _, version := range versions {
		err := pkg.SetPackageVersion(ctx, packageKey, version, files)
		if err != nil {
			t.Fatalf("SetPackageVersion() failed: %v", err)
		}
	}
	
	// Remove entire package
	err := pkg.RemovePackage(ctx, packageKey)
	if err != nil {
		t.Errorf("RemovePackage() unexpected error: %v", err)
	}
	
	// Verify all versions are gone
	for _, version := range versions {
		_, err = pkg.GetPackageVersion(ctx, packageKey, version)
		if err == nil {
			t.Errorf("GetPackageVersion() should fail after RemovePackage() for version %v", version)
		}
	}
}

func TestPackageCache_RemoveOldVersionsByTimestamp(t *testing.T) {
	tempDir := t.TempDir()
	pkg := NewPackageCache(tempDir)
	ctx := context.Background()

	packageKey := "test-package"
	files := []*core.File{{Path: "test.txt", Content: []byte("content")}}

	// Store old version
	oldVersion := core.Version{Major: 1, Minor: 0, Patch: 0}
	err := pkg.SetPackageVersion(ctx, packageKey, oldVersion, files)
	require.NoError(t, err)

	// Wait to create age difference
	time.Sleep(50 * time.Millisecond)

	// Store new version
	newVersion := core.Version{Major: 1, Minor: 1, Patch: 0}
	err = pkg.SetPackageVersion(ctx, packageKey, newVersion, files)
	require.NoError(t, err)

	// Remove versions older than 25ms (should remove old version only)
	err = pkg.RemoveOldVersions(ctx, 25*time.Millisecond)
	require.NoError(t, err)

	// Check that only new version remains
	versions, err := pkg.ListPackageVersions(ctx, packageKey)
	require.NoError(t, err)
	assert.Len(t, versions, 1)
	assert.Equal(t, newVersion, versions[0])

	// Verify package still exists
	packages, err := pkg.ListPackages(ctx)
	require.NoError(t, err)
	assert.Len(t, packages, 1)
}

func TestPackageCache_RemoveUnusedVersionsByAccess(t *testing.T) {
	tempDir := t.TempDir()
	pkg := NewPackageCache(tempDir)
	ctx := context.Background()

	packageKey := "test-package"
	files := []*core.File{{Path: "test.txt", Content: []byte("content")}}

	// Store two versions
	version1 := core.Version{Major: 1, Minor: 0, Patch: 0}
	version2 := core.Version{Major: 1, Minor: 1, Patch: 0}
	
	err := pkg.SetPackageVersion(ctx, packageKey, version1, files)
	require.NoError(t, err)
	err = pkg.SetPackageVersion(ctx, packageKey, version2, files)
	require.NoError(t, err)

	// Wait to create access time difference
	time.Sleep(50 * time.Millisecond)

	// Access version2 only
	_, err = pkg.GetPackageVersion(ctx, packageKey, version2)
	require.NoError(t, err)

	// Remove versions not accessed in last 25ms (should remove version1 only)
	err = pkg.RemoveUnusedVersions(ctx, 25*time.Millisecond)
	require.NoError(t, err)

	// Check that only version2 remains
	versions, err := pkg.ListPackageVersions(ctx, packageKey)
	require.NoError(t, err)
	assert.Len(t, versions, 1)
	assert.Equal(t, version2, versions[0])

	// Verify package still exists
	packages, err := pkg.ListPackages(ctx)
	require.NoError(t, err)
	assert.Len(t, packages, 1)
}

func TestPackageCache_RemoveAllVersionsRemovesPackage(t *testing.T) {
	tempDir := t.TempDir()
	pkg := NewPackageCache(tempDir)
	ctx := context.Background()

	packageKey := "test-package"
	files := []*core.File{{Path: "test.txt", Content: []byte("content")}}

	// Store version
	version := core.Version{Major: 1, Minor: 0, Patch: 0}
	err := pkg.SetPackageVersion(ctx, packageKey, version, files)
	require.NoError(t, err)

	// Wait to create age
	time.Sleep(50 * time.Millisecond)

	// Remove all old versions (should remove package entirely)
	err = pkg.RemoveOldVersions(ctx, 25*time.Millisecond)
	require.NoError(t, err)

	// Verify package is completely gone
	packages, err := pkg.ListPackages(ctx)
	require.NoError(t, err)
	assert.Len(t, packages, 0)
}

// Helper function for reading version metadata
func readVersionMetadata(t *testing.T, baseDir string, packageKey interface{}, version core.Version) struct {
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	AccessedAt time.Time `json:"accessedAt"`
} {
	hashedKey, err := GenerateKey(packageKey)
	require.NoError(t, err)
	
	versionDir := fmt.Sprintf("v%d.%d.%d", version.Major, version.Minor, version.Patch)
	metadataPath := filepath.Join(baseDir, hashedKey, versionDir, "metadata.json")
	data, err := os.ReadFile(metadataPath)
	require.NoError(t, err)
	
	var metadata struct {
		CreatedAt  time.Time `json:"createdAt"`
		UpdatedAt  time.Time `json:"updatedAt"`
		AccessedAt time.Time `json:"accessedAt"`
	}
	err = json.Unmarshal(data, &metadata)
	require.NoError(t, err)
	
	return metadata
}