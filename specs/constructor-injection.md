# Test Isolation via Constructor Injection and Environment Variables

## Job to be Done
Enable test isolation through environment variables and constructor injection, preventing tests from polluting user's home directory. Additionally, ensure lock file is always colocated with manifest file.

## Activities
1. **Add environment variable support** - ARM_MANIFEST_PATH, ARM_CONFIG_PATH, ARM_HOME
2. **Accept directory paths as constructor parameters** - Components receive paths instead of calling os.UserHomeDir()
3. **Provide default constructors** - Production code checks env vars, falls back to os.UserHomeDir()
4. **Provide test constructors** - Tests pass t.TempDir() directly as string parameters
5. **Derive lock file path from manifest path** - Ensure arm-lock.json lives next to arm.json

## Acceptance Criteria
- [ ] ARM_MANIFEST_PATH controls manifest and lock file location (already exists)
- [ ] ARM_CONFIG_PATH overrides .armrc location (single file, bypasses hierarchy)
- [ ] ARM_HOME overrides home directory for .arm/ directory only (not .armrc)
- [ ] Lock file always colocated with manifest file (same directory)
- [ ] Lock path derived from manifest path (arm.json → arm-lock.json)
- [ ] Components accept home directory path as constructor parameter
- [ ] Default constructors check ARM_HOME before calling os.UserHomeDir()
- [ ] Test constructors accept directory paths directly (no OS calls)
- [ ] No direct os.UserHomeDir() calls in component methods
- [ ] Tests can use env vars or direct path injection
- [ ] Tests don't pollute user's actual home directory

## Pattern

## Environment Variables

ARM supports three environment variables for controlling file locations:

### ARM_MANIFEST_PATH (Already Exists)
Controls the location of `arm.json` and `arm-lock.json`.

```bash
ARM_MANIFEST_PATH=/tmp/test/arm.json
# Results in:
# - /tmp/test/arm.json (manifest)
# - /tmp/test/arm-lock.json (lock file, colocated)
```

**Default:** `./arm.json` and `./arm-lock.json` in current working directory

### ARM_CONFIG_PATH (New)
Overrides the `.armrc` configuration file location. When set, this is the ONLY config file used (no hierarchical lookup).

```bash
ARM_CONFIG_PATH=/tmp/test/.armrc
# Results in:
# - Only reads /tmp/test/.armrc
# - Ignores both ./.armrc and ~/.armrc
```

**Default:** Hierarchical lookup (`./.armrc` overrides `~/.armrc`)

### ARM_HOME (New)
Overrides the home directory for the `.arm/` directory (storage, cache, etc.). Does NOT affect `.armrc` location.

```bash
ARM_HOME=/tmp/test
# Results in:
# - /tmp/test/.arm/storage/registries/... (package cache)
# Does NOT affect .armrc location
# Future: /tmp/test/.arm/logs/, /tmp/test/.arm/cache/, etc.
```

**Default:** User's home directory from `os.UserHomeDir()`

### Priority Order

**For .armrc lookup:**
1. **ARM_CONFIG_PATH** - If set, use this exact file (bypasses hierarchy)
2. **workingDir/.armrc** - Project config (highest priority in hierarchy)
3. **userHomeDir/.armrc** - User config (fallback in hierarchy)

**For .arm/storage/ lookup:**
1. **ARM_HOME/.arm/storage/** - If ARM_HOME is set
2. **~/.arm/storage/** - Default

### Implementation Pattern

```go
// Get config path - ARM_CONFIG_PATH bypasses hierarchy
func getConfigPath() string {
    // Explicit override - bypasses hierarchical lookup
    if path := os.Getenv("ARM_CONFIG_PATH"); path != "" {
        return path
    }
    
    // Normal hierarchical lookup (handled by FileManager)
    // - workingDir/.armrc (project)
    // - userHomeDir/.armrc (user)
    return "" // FileManager handles hierarchy
}

// Get storage path - ARM_HOME only affects .arm/ directory
func getStoragePath() string {
    armHome := os.Getenv("ARM_HOME")
    if armHome == "" {
        armHome, _ = os.UserHomeDir()
    }
    
    return filepath.Join(armHome, ".arm", "storage")
}
```

### Test Usage

**Option 1: Environment Variables**
```go
func TestWithEnvVars(t *testing.T) {
    testDir := t.TempDir()
    t.Setenv("ARM_MANIFEST_PATH", filepath.Join(testDir, "arm.json"))
    t.Setenv("ARM_CONFIG_PATH", filepath.Join(testDir, ".armrc"))
    t.Setenv("ARM_HOME", testDir)
    
    // All ARM operations now isolated to testDir
}
```

**Option 2: Direct Path Injection**
```go
func TestWithDirectPaths(t *testing.T) {
    testDir := t.TempDir()
    registry, err := storage.NewRegistryWithHomeDir(registryKey, testDir)
    // Uses testDir/.arm/storage/ instead of ~/.arm/storage/
}
```

## Constructor Injection Pattern

### Before (Direct OS Call in Methods)
```go
type Registry struct {
    registryKey interface{}
    registryDir string
}

func NewRegistry(registryKey interface{}) (*Registry, error) {
    homeDir, err := os.UserHomeDir()  // ❌ Called every time
    if err != nil {
        return nil, err
    }
    baseDir := filepath.Join(homeDir, ".arm")
    return NewRegistryWithPath(baseDir, registryKey)
}
```

**Problem:** Tests create registries in user's actual ~/.arm/ directory

### After (Constructor Injection + Environment Variables)
```go
type Registry struct {
    registryKey interface{}
    registryDir string
}

// Production constructor - checks env vars, then calls os.UserHomeDir()
func NewRegistry(registryKey interface{}) (*Registry, error) {
    armHome := os.Getenv("ARM_HOME")
    if armHome == "" {
        var err error
        armHome, err = os.UserHomeDir()
        if err != nil {
            return nil, err
        }
    }
    return NewRegistryWithHomeDir(registryKey, armHome)
}

// Test constructor - accepts directory path directly
func NewRegistryWithHomeDir(registryKey interface{}, homeDir string) (*Registry, error) {
    baseDir := filepath.Join(homeDir, ".arm")
    return NewRegistryWithPath(baseDir, registryKey)
}
```

**Benefits:** 
- Production code can use ARM_HOME environment variable
- Tests can use environment variables OR direct path injection
- No OS calls in component methods

## Components to Update

### 1. storage.Registry
**Current:**
```go
func NewRegistry(registryKey interface{}) (*Registry, error) {
    homeDir, err := os.UserHomeDir()
    // ...
}
```

**Updated:**
```go
func NewRegistry(registryKey interface{}) (*Registry, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return nil, err
    }
    return NewRegistryWithHomeDir(registryKey, homeDir)
}

func NewRegistryWithHomeDir(registryKey interface{}, homeDir string) (*Registry, error) {
    baseDir := filepath.Join(homeDir, ".arm")
    return NewRegistryWithPath(baseDir, registryKey)
}
```

### 2. config.FileManager
**Current:**
```go
func NewFileManager() *FileManager {
    workingDir, _ := os.Getwd()
    userHomeDir, _ := os.UserHomeDir()
    return &FileManager{
        workingDir:  workingDir,
        userHomeDir: userHomeDir,
    }
}
```

**No changes needed** - ARM_HOME does not affect .armrc location

**Note:** ARM_CONFIG_PATH should be handled at a higher level to bypass the hierarchical lookup entirely

### 3. service.ArmService (Cache Methods)
**Current:**
```go
func (s *ArmService) CleanCacheByAge(ctx context.Context, maxAge time.Duration) error {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return err
    }
    storageDir := filepath.Join(homeDir, ".arm", "storage")
    return s.cleanCacheByAgeWithPath(ctx, maxAge, storageDir)
}
```

**Updated:**
```go
func (s *ArmService) CleanCacheByAge(ctx context.Context, maxAge time.Duration) error {
    return s.CleanCacheByAgeWithHomeDir(ctx, maxAge, "")
}

func (s *ArmService) CleanCacheByAgeWithHomeDir(ctx context.Context, maxAge time.Duration, homeDir string) error {
    if homeDir == "" {
        // Check ARM_HOME environment variable
        homeDir = os.Getenv("ARM_HOME")
        if homeDir == "" {
            var err error
            homeDir, err = os.UserHomeDir()
            if err != nil {
                return err
            }
        }
    }
    storageDir := filepath.Join(homeDir, ".arm", "storage")
    return s.cleanCacheByAgeWithPath(ctx, maxAge, storageDir)
}
```

**Apply same pattern to:**
- `CleanCacheByTimeSinceLastAccess` → `CleanCacheByTimeSinceLastAccessWithHomeDir`
- `NukeCache` → `NukeCacheWithHomeDir`

## Test Usage

### Production
```go
// Uses actual user home directory (or ARM_HOME if set)
registry, err := storage.NewRegistry(registryKey)
```

### Testing with Environment Variables
```go
func TestRegistry(t *testing.T) {
    testHome := t.TempDir()
    t.Setenv("ARM_HOME", testHome)
    
    // Uses testHome/.arm/storage/ instead of ~/.arm/storage/
    registry, err := storage.NewRegistry(registryKey)
}
```

### Testing with Direct Path Injection
```go
func TestRegistry(t *testing.T) {
    testHome := t.TempDir()
    
    // Uses testHome/.arm/storage/ instead of ~/.arm/storage/
    registry, err := storage.NewRegistryWithHomeDir(registryKey, testHome)
}
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| os.UserHomeDir() returns error | Production constructor propagates error |
| Empty string passed to test constructor | Creates directories relative to empty path (caller's responsibility) |
| Non-existent directory passed | Component creates directories as needed |

## Benefits

### Flexibility
- Production can use environment variables for custom locations
- Tests can use environment variables OR direct path injection
- No interfaces or mocking frameworks needed

### Test Reliability
- Tests use isolated temporary directories
- No pollution of user's ~/.arm/ or ~/.armrc
- Parallel test execution is safe
- Automatic cleanup via t.TempDir()

### Production Use Cases
- **CI/CD:** Set ARM_HOME to build-specific cache directory
- **Docker:** Mount volumes and point ARM_HOME to mounted path
- **Multi-user systems:** Separate ARM cache directories per user/project
- **Network storage:** Point ARM_HOME to shared network drive for team caches
- **Custom config:** Use ARM_CONFIG_PATH for non-standard .armrc locations

### Backward Compatibility
- Existing production code unchanged (uses default constructors)
- Existing tests can migrate incrementally
- No breaking changes to public API
- Environment variables are optional

## Implementation Mapping

**Environment variables:**
- `ARM_MANIFEST_PATH` - Already exists, controls arm.json location
- `ARM_CONFIG_PATH` - New, overrides .armrc location (bypasses hierarchy)
- `ARM_HOME` - New, overrides home directory for .arm/ directory only

**Lock file colocation:**
- `cmd/arm/main.go` - All command handlers (24 locations) must derive lock path from manifest path

**Components to update:**
- `internal/arm/storage/registry.go` - Add NewRegistryWithHomeDir(), check ARM_HOME in NewRegistry()
- `internal/arm/service/service.go` - Add *WithHomeDir() variants for cache methods, check ARM_HOME
- Registry creation code - Handle ARM_CONFIG_PATH to bypass hierarchical .armrc lookup

**No changes needed:**
- `internal/arm/config/manager.go` - ARM_HOME does not affect .armrc location
- `internal/arm/manifest/manager.go` - Uses relative paths (no home dir needed)
- `internal/arm/packagelockfile/manager.go` - Uses relative paths (no home dir needed)

## Lock File Colocation

### Problem

When `ARM_MANIFEST_PATH` is set to a custom location, the lock file is NOT colocated with the manifest file.

**Current Behavior:**
```bash
ARM_MANIFEST_PATH=/tmp/test/arm.json
# Results in:
# - /tmp/test/arm.json (manifest) ✅
# - ./arm-lock.json (lock) ❌ NOT colocated!
```

**Root Cause:**
```go
manifestPath := os.Getenv("ARM_MANIFEST_PATH")
if manifestPath == "" {
    manifestPath = "arm.json"
}
manifestMgr := manifest.NewFileManagerWithPath(manifestPath)
lockfileMgr := packagelockfile.NewFileManager()  // ❌ Always uses "./arm-lock.json"
```

### Solution

Derive lock path from manifest path:

```go
manifestPath := os.Getenv("ARM_MANIFEST_PATH")
if manifestPath == "" {
    manifestPath = "arm.json"
}

// Derive lock path from manifest path
lockPath := strings.TrimSuffix(manifestPath, ".json") + "-lock.json"
// Examples:
// - arm.json → arm-lock.json
// - /tmp/test/arm.json → /tmp/test/arm-lock.json
// - /custom/path/manifest.json → /custom/path/manifest-lock.json

manifestMgr := manifest.NewFileManagerWithPath(manifestPath)
lockfileMgr := packagelockfile.NewFileManagerWithPath(lockPath)
```

### Files to Update

**cmd/arm/main.go** - All command handlers that create lockfileMgr:
- handleAddGitRegistry
- handleAddGitLabRegistry
- handleAddCloudsmithRegistry
- handleAddSink
- handleInstallRuleset
- handleInstallPromptset
- handleInstallAll
- handleUninstall
- handleUpdate
- handleUpgrade
- handleSetRuleset
- handleSetPromptset
- handleOutdated
- handleListDependencies
- handleInfoDependency
- handleCleanSinks
- (All other handlers that use lockfileMgr)

### Benefits

- Lock file always in same directory as manifest file
- Consistent with npm (package-lock.json next to package.json)
- Tests can isolate both files with single `ARM_MANIFEST_PATH`
- No new environment variables needed

## Related Specs
- `cache-management.md` - Storage components accept homeDir parameter
- `authentication.md` - Config manager accepts workingDir and homeDir parameters
- `e2e-testing.md` - Tests use direct path injection
