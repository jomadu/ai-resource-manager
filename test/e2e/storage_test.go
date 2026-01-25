package e2e

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jomadu/ai-resource-manager/test/e2e/helpers"
)

func TestPackageCachedAfterFirstInstall(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)
	storageDir := filepath.Join(os.Getenv("HOME"), ".arm", "storage")

	// Setup: Create test Git repository
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.Commit("Initial commit")
	repo.Tag("v1.0.0")

	// Setup: Add registry and sink
	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
	arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")

	// Install package
	arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "cursor-rules")

	// Verify storage directory exists
	helpers.AssertDirExists(t, storageDir)

	// Verify registries directory exists
	registriesDir := filepath.Join(storageDir, "registries")
	helpers.AssertDirExists(t, registriesDir)

	// Find the registry directory (it's a hash)
	entries, err := os.ReadDir(registriesDir)
	if err != nil {
		t.Fatalf("failed to read registries directory: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected at least one registry in storage")
	}

	// Verify registry metadata.json exists
	registryDir := filepath.Join(registriesDir, entries[0].Name())
	registryMetadata := filepath.Join(registryDir, "metadata.json")
	helpers.AssertFileExists(t, registryMetadata)

	// Verify packages directory exists
	packagesDir := filepath.Join(registryDir, "packages")
	helpers.AssertDirExists(t, packagesDir)

	// Verify at least one package cached
	packageEntries, err := os.ReadDir(packagesDir)
	if err != nil {
		t.Fatalf("failed to read packages directory: %v", err)
	}
	if len(packageEntries) == 0 {
		t.Fatal("expected at least one package in cache")
	}

	// Verify package metadata.json exists
	packageDir := filepath.Join(packagesDir, packageEntries[0].Name())
	packageMetadata := filepath.Join(packageDir, "metadata.json")
	helpers.AssertFileExists(t, packageMetadata)

	// Verify version directory exists
	versionEntries, err := os.ReadDir(packageDir)
	if err != nil {
		t.Fatalf("failed to read package directory: %v", err)
	}
	
	var versionDir string
	for _, entry := range versionEntries {
		if entry.IsDir() {
			versionDir = filepath.Join(packageDir, entry.Name())
			break
		}
	}
	if versionDir == "" {
		t.Fatal("expected version directory in package cache")
	}

	// Verify version metadata.json exists
	versionMetadata := filepath.Join(versionDir, "metadata.json")
	helpers.AssertFileExists(t, versionMetadata)

	// Verify files directory exists
	filesDir := filepath.Join(versionDir, "files")
	helpers.AssertDirExists(t, filesDir)

	// Verify files exist in cache
	fileCount := helpers.CountFilesRecursive(t, filesDir)
	if fileCount == 0 {
		t.Error("expected cached files in storage")
	}
}

func TestCacheReusedOnSecondInstall(t *testing.T) {
	workDir1 := t.TempDir()
	workDir2 := t.TempDir()
	arm1 := helpers.NewARMRunner(t, workDir1)
	arm2 := helpers.NewARMRunner(t, workDir2)

	// Setup: Create test Git repository
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.Commit("Initial commit")
	repo.Tag("v1.0.0")

	repoURL := "file://" + repoDir

	// First project: Install package
	arm1.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
	arm1.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")
	arm1.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "cursor-rules")

	// Get storage directory
	storageDir := filepath.Join(os.Getenv("HOME"), ".arm", "storage")
	registriesDir := filepath.Join(storageDir, "registries")

	// Find registry directory
	entries, err := os.ReadDir(registriesDir)
	if err != nil {
		t.Fatalf("failed to read registries directory: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected registry in storage")
	}
	registryDir := filepath.Join(registriesDir, entries[0].Name())

	// Get modification time of cached files
	packagesDir := filepath.Join(registryDir, "packages")
	packageEntries, err := os.ReadDir(packagesDir)
	if err != nil {
		t.Fatalf("failed to read packages directory: %v", err)
	}
	if len(packageEntries) == 0 {
		t.Fatal("expected package in cache")
	}
	packageDir := filepath.Join(packagesDir, packageEntries[0].Name())

	// Get version directory
	versionEntries, err := os.ReadDir(packageDir)
	if err != nil {
		t.Fatalf("failed to read package directory: %v", err)
	}
	var versionDir string
	for _, entry := range versionEntries {
		if entry.IsDir() {
			versionDir = filepath.Join(packageDir, entry.Name())
			break
		}
	}
	if versionDir == "" {
		t.Fatal("expected version directory")
	}

	filesDir := filepath.Join(versionDir, "files")
	fileInfo1, err := os.Stat(filesDir)
	if err != nil {
		t.Fatalf("failed to stat files directory: %v", err)
	}
	modTime1 := fileInfo1.ModTime()

	// Wait a bit to ensure different timestamps if files are recreated
	time.Sleep(100 * time.Millisecond)

	// Second project: Install same package
	arm2.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
	arm2.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")
	arm2.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "cursor-rules")

	// Verify cache was reused (modification time unchanged)
	fileInfo2, err := os.Stat(filesDir)
	if err != nil {
		t.Fatalf("failed to stat files directory after second install: %v", err)
	}
	modTime2 := fileInfo2.ModTime()

	if !modTime1.Equal(modTime2) {
		t.Errorf("cache was not reused: mod time changed from %v to %v", modTime1, modTime2)
	}

	// Verify both projects have compiled files
	sinkDir1 := filepath.Join(workDir1, ".cursor/rules")
	sinkDir2 := filepath.Join(workDir2, ".cursor/rules")
	
	if helpers.CountFilesRecursive(t, sinkDir1) == 0 {
		t.Error("first project should have compiled files")
	}
	if helpers.CountFilesRecursive(t, sinkDir2) == 0 {
		t.Error("second project should have compiled files")
	}
}

func TestCacheKeyGenerationWithPatterns(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)

	// Setup: Create test Git repository with multiple files
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	repo.WriteFile("security.yml", helpers.SecurityRuleset)
	repo.WriteFile("clean-code.yml", helpers.MinimalRuleset)
	repo.Commit("Initial commit")
	repo.Tag("v1.0.0")

	// Use a unique registry name to avoid cache pollution from other tests
	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry-patterns")
	arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")

	// Install with include pattern
	arm.MustRun("install", "ruleset", "--include", "security.yml", "test-registry-patterns/test-ruleset@1.0.0", "cursor-rules")

	// Get storage directory
	storageDir := filepath.Join(os.Getenv("HOME"), ".arm", "storage")
	registriesDir := filepath.Join(storageDir, "registries")

	entries, err := os.ReadDir(registriesDir)
	if err != nil {
		t.Fatalf("failed to read registries directory: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected registry in storage")
	}

	// Find the registry for our test (by checking metadata)
	var testRegistryDir string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		registryDir := filepath.Join(registriesDir, entry.Name())
		metadataPath := filepath.Join(registryDir, "metadata.json")
		
		data, err := os.ReadFile(metadataPath)
		if err != nil {
			continue
		}
		
		var regMeta map[string]interface{}
		if err := json.Unmarshal(data, &regMeta); err != nil {
			continue
		}
		
		// Check if this is our test registry
		if url, ok := regMeta["url"].(string); ok && url == repoURL {
			testRegistryDir = registryDir
			break
		}
	}
	
	if testRegistryDir == "" {
		t.Fatal("could not find test registry in storage")
	}

	// Verify package metadata includes pattern information
	packagesDir := filepath.Join(testRegistryDir, "packages")
	packageEntries, err := os.ReadDir(packagesDir)
	if err != nil {
		t.Fatalf("failed to read packages directory: %v", err)
	}
	if len(packageEntries) == 0 {
		t.Fatal("expected package in cache")
	}

	// Find the package with include pattern
	foundPackageWithPattern := false
	for _, pkgEntry := range packageEntries {
		packageDir := filepath.Join(packagesDir, pkgEntry.Name())
		packageMetadata := filepath.Join(packageDir, "metadata.json")

		data, err := os.ReadFile(packageMetadata)
		if err != nil {
			continue
		}

		var metadata map[string]interface{}
		if err := json.Unmarshal(data, &metadata); err != nil {
			continue
		}

		// Check if this package has the include pattern we specified
		includeField, ok := metadata["include"]
		if !ok || includeField == nil {
			continue
		}
		
		includeArray, ok := includeField.([]interface{})
		if !ok || len(includeArray) == 0 {
			continue
		}
		
		// Check if it contains "security.yml"
		for _, pattern := range includeArray {
			if patternStr, ok := pattern.(string); ok && patternStr == "security.yml" {
				foundPackageWithPattern = true
				break
			}
		}
		
		if foundPackageWithPattern {
			break
		}
	}

	if !foundPackageWithPattern {
		t.Error("expected to find package with include pattern 'security.yml'")
	}
}

func TestCleanCacheWithDefaultAge(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)

	// Setup: Create test Git repository
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.Commit("Initial commit")
	repo.Tag("v1.0.0")

	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
	arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")
	arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "cursor-rules")

	// Get storage directory
	storageDir := filepath.Join(os.Getenv("HOME"), ".arm", "storage")
	registriesDir := filepath.Join(storageDir, "registries")

	// Verify cache exists
	entries, err := os.ReadDir(registriesDir)
	if err != nil {
		t.Fatalf("failed to read registries directory: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected registry in storage before clean")
	}

	// Clean cache with default age (7 days) - should not remove recent cache
	arm.MustRun("clean", "cache")

	// Verify cache still exists (it's recent)
	entries, err = os.ReadDir(registriesDir)
	if err != nil {
		t.Fatalf("failed to read registries directory after clean: %v", err)
	}
	if len(entries) == 0 {
		t.Error("cache should not be removed with default age (recent install)")
	}
}

func TestCleanCacheWithMaxAge(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)

	// Setup: Create test Git repository
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.Commit("Initial commit")
	repo.Tag("v1.0.0")

	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
	arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")
	arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "cursor-rules")

	// Get storage directory
	storageDir := filepath.Join(os.Getenv("HOME"), ".arm", "storage")
	registriesDir := filepath.Join(storageDir, "registries")

	// Verify cache exists
	entries, err := os.ReadDir(registriesDir)
	if err != nil {
		t.Fatalf("failed to read registries directory: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected registry in storage before clean")
	}

	// Clean cache with max-age of 0 seconds - should remove all cache
	arm.MustRun("clean", "cache", "--max-age", "0s")

	// Verify version directories were removed (registry directory may still exist)
	entries, err = os.ReadDir(registriesDir)
	if err != nil {
		t.Fatalf("failed to read registries directory after clean: %v", err)
	}
	
	// Check if any version directories still exist
	hasVersions := false
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		registryDir := filepath.Join(registriesDir, entry.Name())
		packagesDir := filepath.Join(registryDir, "packages")
		if _, err := os.Stat(packagesDir); os.IsNotExist(err) {
			continue
		}
		
		packageEntries, err := os.ReadDir(packagesDir)
		if err != nil {
			continue
		}
		
		for _, pkgEntry := range packageEntries {
			if !pkgEntry.IsDir() {
				continue
			}
			packageDir := filepath.Join(packagesDir, pkgEntry.Name())
			versionEntries, err := os.ReadDir(packageDir)
			if err != nil {
				continue
			}
			
			for _, verEntry := range versionEntries {
				if verEntry.IsDir() && strings.HasPrefix(verEntry.Name(), "v") {
					hasVersions = true
					break
				}
			}
			if hasVersions {
				break
			}
		}
		if hasVersions {
			break
		}
	}
	
	if hasVersions {
		t.Error("version directories should be removed with max-age 0s")
	}
}

func TestCleanCacheWithNuke(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)

	// Setup: Create test Git repository
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.Commit("Initial commit")
	repo.Tag("v1.0.0")

	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
	arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")
	arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "cursor-rules")

	// Get storage directory
	storageDir := filepath.Join(os.Getenv("HOME"), ".arm", "storage")

	// Verify storage exists
	helpers.AssertDirExists(t, storageDir)

	// Clean cache with --nuke - should remove entire storage directory
	arm.MustRun("clean", "cache", "--nuke")

	// Verify storage directory was removed
	if _, err := os.Stat(storageDir); !os.IsNotExist(err) {
		t.Error("storage directory should be removed with --nuke")
	}
}

func TestCacheStructure(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)

	// Setup: Create test Git repository
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.Commit("Initial commit")
	repo.Tag("v1.0.0")

	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
	arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")
	arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "cursor-rules")

	// Get storage directory
	storageDir := filepath.Join(os.Getenv("HOME"), ".arm", "storage")
	registriesDir := filepath.Join(storageDir, "registries")

	// Find registry directory
	entries, err := os.ReadDir(registriesDir)
	if err != nil {
		t.Fatalf("failed to read registries directory: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected registry in storage")
	}
	registryDir := filepath.Join(registriesDir, entries[0].Name())

	// Verify registry metadata.json
	registryMetadata := filepath.Join(registryDir, "metadata.json")
	helpers.AssertFileExists(t, registryMetadata)

	data, err := os.ReadFile(registryMetadata)
	if err != nil {
		t.Fatalf("failed to read registry metadata: %v", err)
	}

	var regMeta map[string]interface{}
	if err := json.Unmarshal(data, &regMeta); err != nil {
		t.Fatalf("failed to parse registry metadata: %v", err)
	}

	// Verify required fields
	if _, ok := regMeta["url"]; !ok {
		t.Error("registry metadata missing url field")
	}
	if _, ok := regMeta["type"]; !ok {
		t.Error("registry metadata missing type field")
	}

	// Verify packages directory structure
	packagesDir := filepath.Join(registryDir, "packages")
	helpers.AssertDirExists(t, packagesDir)

	packageEntries, err := os.ReadDir(packagesDir)
	if err != nil {
		t.Fatalf("failed to read packages directory: %v", err)
	}
	if len(packageEntries) == 0 {
		t.Fatal("expected package in cache")
	}

	packageDir := filepath.Join(packagesDir, packageEntries[0].Name())

	// Verify package metadata.json
	packageMetadata := filepath.Join(packageDir, "metadata.json")
	helpers.AssertFileExists(t, packageMetadata)

	// Verify version directory
	versionEntries, err := os.ReadDir(packageDir)
	if err != nil {
		t.Fatalf("failed to read package directory: %v", err)
	}
	
	var versionDir string
	for _, entry := range versionEntries {
		if entry.IsDir() {
			versionDir = filepath.Join(packageDir, entry.Name())
			break
		}
	}
	if versionDir == "" {
		t.Fatal("expected version directory")
	}

	// Verify version metadata.json
	versionMetadata := filepath.Join(versionDir, "metadata.json")
	helpers.AssertFileExists(t, versionMetadata)

	data, err = os.ReadFile(versionMetadata)
	if err != nil {
		t.Fatalf("failed to read version metadata: %v", err)
	}

	var verMeta map[string]interface{}
	if err := json.Unmarshal(data, &verMeta); err != nil {
		t.Fatalf("failed to parse version metadata: %v", err)
	}

	// Verify timestamp fields exist
	if _, ok := verMeta["accessedAt"]; !ok {
		t.Error("version metadata missing accessedAt field")
	}
	if _, ok := verMeta["createdAt"]; !ok {
		t.Error("version metadata missing createdAt field")
	}
	if _, ok := verMeta["updatedAt"]; !ok {
		t.Error("version metadata missing updatedAt field")
	}

	// Verify files directory
	filesDir := filepath.Join(versionDir, "files")
	helpers.AssertDirExists(t, filesDir)

	// Verify files exist
	fileCount := helpers.CountFilesRecursive(t, filesDir)
	if fileCount == 0 {
		t.Error("expected files in cache")
	}
}
