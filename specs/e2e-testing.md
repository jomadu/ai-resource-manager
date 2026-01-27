# E2E Testing

## Job to be Done
Validate ARM functionality through comprehensive end-to-end tests covering all user workflows and edge cases.

## Activities
1. Test installation workflows (install, update, upgrade, uninstall)
2. Test registry types (Git, GitLab, Cloudsmith)
3. Test compilation for all tools (Cursor, Amazon Q, Copilot, Markdown)
4. Test version resolution (semver, branches, constraints)
5. Test cache management and cleanup
6. Test authentication and authorization
7. Test pattern filtering and archive extraction
8. Test manifest and lock file management
9. Test integrity verification
10. Test multi-sink installations

## Acceptance Criteria
- [x] 56 E2E tests across 14 test suites covering all workflows
- [x] Tests use isolated temporary directories
- [x] Tests don't pollute user's ARM directories
- [x] Tests create local Git repositories for reproducibility
- [x] Tests verify file system state (files, directories, JSON)
- [x] Tests verify manifest and lock file correctness
- [x] Tests verify compiled output for each tool
- [x] Tests verify cache structure and metadata
- [x] Tests verify cleanup behavior (empty directories, index files)
- [x] All tests pass reliably (100% pass rate)

## Test Suites

### install_test.go
- TestRulesetInstallation - Basic ruleset installation
- TestRulesetInstallationWithLatest - Resolve latest version
- TestRulesetInstallationWithBranch - Install from branch
- TestRulesetInstallationWithPriority - Custom priority
- TestRulesetInstallationToMultipleSinks - Multi-sink install
- TestPromptsetInstallation - Promptset installation
- TestInstallWithPatterns - Include/exclude patterns

### update_test.go
- TestUpdateWithinConstraints - Update respects constraints
- TestUpgradeIgnoresConstraints - Upgrade to latest
- TestUpdatePromptset - Update promptset
- TestUpgradePromptset - Upgrade promptset
- TestManifestFilesUpdated - Manifest/lock file updates

### compile_test.go
- TestCompilationToolFormats - All tool formats
- TestCompilationPromptsets - Promptset compilation
- TestCompilationValidation - Validation errors
- TestCompilationIndexGeneration - Priority index
- TestCompilationHierarchicalLayout - Directory structure
- TestCompilationMultiplePriorities - Priority ordering

### storage_test.go
- TestPackageCachedAfterFirstInstall - Cache creation
- TestCacheReusedOnSecondInstall - Cache reuse
- TestCacheKeyGenerationWithPatterns - Cache keys
- TestCleanCacheWithDefaultAge - Age-based cleanup
- TestCleanCacheWithMaxAge - Custom age cleanup
- TestCleanCacheWithNuke - Complete cache removal
- TestCacheStructure - Directory structure

### version_test.go
- TestVersionResolutionLatest - Latest version
- TestVersionResolutionMajorConstraint - Major constraint
- TestVersionResolutionMinorConstraint - Minor constraint
- TestVersionResolutionExactVersion - Exact version
- TestVersionResolutionBranchToCommit - Branch resolution
- TestVersionResolutionLatestWithNoTags - No tags fallback

### registry_test.go
- TestGitRegistryManagement - Add/remove Git registry
- TestGitRegistryWithBranches - Branch tracking

### sink_test.go
- TestSinkManagement - Add/remove sinks
- TestSinkLayoutModes - Hierarchical vs flat

### manifest_test.go
- TestManifestCreation - Initial manifest
- TestLockFileCreation - Lock file generation
- TestLockFileBranchResolution - Branch to commit
- TestIndexFileCreation - arm-index.json
- TestManifestJSONValidity - Valid JSON structure
- TestManifestPreservesConfiguration - Config persistence

### integrity_test.go
- TestIntegrityVerification_E2E - SHA256 verification
- TestIntegrityVerification_BackwardsCompatibility - No integrity in old locks

### auth_test.go
- TestAuthenticationWithArmrc - Token authentication
- TestArmrcFilePermissions - Security warnings
- TestArmrcSectionMatching - URL matching

### archive_test.go
- TestArchiveTarGz - .tar.gz extraction
- TestArchiveZip - .zip extraction
- TestArchiveMixedWithLooseFiles - Archive + loose files
- TestArchivePrecedenceOverLooseFiles - Archive precedence
- TestArchiveWithIncludeExcludePatterns - Pattern filtering

### multisink_test.go
- TestMultiSinkCrossToolInstallation - Install to multiple tools
- TestMultiSinkSwitching - Change sink configuration
- TestMultiSinkUpdate - Update across sinks

### cleanup_test.go
- TestUninstallCleanup - Empty directory removal
- TestUninstallCleanup/EmptyDirectoriesRemoved - Recursive cleanup
- TestUninstallCleanup/IndexFileRemoved - arm-index.json removal
- TestUninstallCleanup/PriorityIndexRemoved - arm_index.* removal
- TestUninstallCleanup/MultiplePackages - Partial cleanup

### errors_test.go
- TestErrorHandling - Error messages and codes

## Test Helpers

### helpers/arm.go
- NewARMRunner - Create ARM CLI runner
- Run - Execute ARM command
- MustRun - Execute and assert success
- MustFail - Execute and assert failure

### helpers/git.go
- NewGitRepo - Create test Git repository
- WriteFile - Add file to repository
- Commit - Create commit
- Tag - Create tag
- Branch - Create branch
- Checkout - Switch branch

### helpers/assertions.go
- AssertFileExists - Verify file exists
- AssertFileNotExists - Verify file doesn't exist
- AssertFileContains - Verify file content
- AssertDirExists - Verify directory exists
- AssertDirNotExists - Verify directory doesn't exist
- ReadJSON - Parse JSON file
- AssertJSONField - Verify JSON field value
- CountFiles - Count files in directory
- CountFilesRecursive - Count files recursively

## Implementation Mapping

**Source files:**
- `test/e2e/install_test.go` - Installation tests
- `test/e2e/update_test.go` - Update/upgrade tests
- `test/e2e/compile_test.go` - Compilation tests
- `test/e2e/storage_test.go` - Cache tests
- `test/e2e/version_test.go` - Version resolution tests
- `test/e2e/registry_test.go` - Registry tests
- `test/e2e/sink_test.go` - Sink tests
- `test/e2e/manifest_test.go` - Manifest tests
- `test/e2e/integrity_test.go` - Integrity tests
- `test/e2e/auth_test.go` - Authentication tests
- `test/e2e/archive_test.go` - Archive tests
- `test/e2e/multisink_test.go` - Multi-sink tests
- `test/e2e/cleanup_test.go` - Cleanup tests
- `test/e2e/errors_test.go` - Error handling tests
- `test/e2e/helpers/` - Test helper functions

## Examples

### Basic Test Structure
```go
func TestRulesetInstallation(t *testing.T) {
    // Setup: Create isolated environment
    workDir := t.TempDir()
    arm := helpers.NewARMRunner(t, workDir)
    
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
    
    // Test: Install ruleset
    arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "cursor-rules")
    
    // Verify: Check manifest
    armJSON := filepath.Join(workDir, "arm.json")
    manifest := helpers.ReadJSON(t, armJSON)
    helpers.AssertJSONField(t, manifest, "dependencies.test-registry/test-ruleset.type", "ruleset")
    
    // Verify: Check lock file
    lockFile := filepath.Join(workDir, "arm-lock.json")
    helpers.AssertFileExists(t, lockFile)
    
    // Verify: Check compiled files
    sinkDir := filepath.Join(workDir, ".cursor", "rules")
    helpers.AssertDirExists(t, sinkDir)
    fileCount := helpers.CountFilesRecursive(t, sinkDir)
    if fileCount == 0 {
        t.Error("expected compiled files in sink directory")
    }
}
```

### Test Isolation
```go
func TestCacheIsolation(t *testing.T) {
    // Each test gets isolated temporary directory
    workDir := t.TempDir()
    
    // Set environment variables for isolation
    t.Setenv("ARM_HOME", t.TempDir())
    t.Setenv("ARM_CONFIG_PATH", filepath.Join(t.TempDir(), ".armrc"))
    
    // Test runs without polluting user's directories
    // Cleanup is automatic via t.TempDir()
}
```

### Verification Patterns
```go
// Verify file exists
helpers.AssertFileExists(t, filepath.Join(workDir, "arm.json"))

// Verify directory exists
helpers.AssertDirExists(t, filepath.Join(workDir, ".cursor", "rules"))

// Verify JSON field
manifest := helpers.ReadJSON(t, armJSON)
helpers.AssertJSONField(t, manifest, "dependencies.pkg.version", "^1.0.0")

// Verify file content
helpers.AssertFileContains(t, ruleFile, "priority: 100")

// Count files
fileCount := helpers.CountFilesRecursive(t, sinkDir)
if fileCount != 5 {
    t.Errorf("expected 5 files, got %d", fileCount)
}
```
