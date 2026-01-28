# Constructor Injection

## Job to be Done
Enable test isolation by allowing components to accept custom paths via constructors and environment variables, preventing tests from polluting user directories.

## Activities
1. Support ARM_HOME for custom .arm/storage/ location
2. Support ARM_CONFIG_PATH for custom .armrc location (bypasses hierarchy)
3. Support ARM_MANIFEST_PATH for custom arm.json location
4. Provide *WithPath constructors for test injection
5. Colocate lock file with manifest file

## Acceptance Criteria
- [x] ARM_HOME overrides home directory for .arm/storage/
- [x] ARM_CONFIG_PATH overrides .armrc location (bypasses ./.armrc and ~/.armrc)
- [x] ARM_MANIFEST_PATH resolved at CLI level (not in components)
- [x] Lock file always colocated with manifest file
- [x] NewRegistryWithHomeDir() accepts custom home directory
- [x] *WithHomeDir() variants for cache methods
- [x] NewFileManagerWithPath() accepts custom manifest path
- [x] Tests use t.TempDir() for isolation
- [x] No direct os.UserHomeDir() calls in component methods

## Data Structures

### Environment Variables
```bash
ARM_HOME=/custom/home
ARM_CONFIG_PATH=/custom/path/.armrc
ARM_MANIFEST_PATH=/custom/path/arm.json
```

### Constructor Variants
```go
// Default constructors (check env vars)
NewRegistry(registryKey string) (*Registry, error)
NewFileManager() (*FileManager, error)

// Test constructors (accept paths)
NewRegistryWithHomeDir(registryKey, homeDir string) (*Registry, error)
NewFileManagerWithPath(manifestPath string) (*FileManager, error)
```

## Algorithm

### ARM_HOME Resolution
1. Check if ARM_HOME environment variable is set
2. If set, use ARM_HOME/.arm/storage/
3. If not set, use os.UserHomeDir()/.arm/storage/
4. Return storage path

### ARM_CONFIG_PATH Resolution
1. Check if ARM_CONFIG_PATH environment variable is set
2. If set, load only that file (bypass hierarchy)
3. If not set, use hierarchical lookup:
   - Load ~/.armrc (user config)
   - Load ./.armrc (project config)
   - Project overrides user
4. Return configuration

### ARM_MANIFEST_PATH Resolution (CLI Level)
1. CLI handlers check ARM_MANIFEST_PATH environment variable
2. If set, use that path
3. If not set, use "arm.json"
4. Pass resolved path to manifest.NewFileManagerWithPath()

**Note:** The manifest manager does not check ARM_MANIFEST_PATH. Path resolution is a CLI concern.

### Lock File Colocation
1. Take manifest path as input
2. Strip `.json` suffix: `strings.TrimSuffix(manifestPath, ".json")`
3. Append `-lock.json` suffix
4. Return lock file path

**Examples:**
- `arm.json` → `arm-lock.json`
- `./config/arm.json` → `./config/arm-lock.json`
- `my-manifest.json` → `my-manifest-lock.json`
- `arm` (no extension) → `arm-lock.json`

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| ARM_HOME not set | Use os.UserHomeDir() |
| ARM_HOME set to relative path | Resolve to absolute path |
| ARM_CONFIG_PATH not set | Use hierarchical lookup |
| ARM_CONFIG_PATH set | Bypass hierarchy, use only that file |
| ARM_MANIFEST_PATH not set | Use ./arm.json |
| ARM_MANIFEST_PATH in subdirectory | Lock file in same subdirectory |
| Multiple tests in parallel | Each uses isolated t.TempDir() |
| Test cleanup | Automatic via t.TempDir() |

## Dependencies

- Environment variable access
- File system operations
- Path manipulation

## Implementation Mapping

**Source files:**
- `internal/arm/storage/registry.go` - NewRegistry, NewRegistryWithHomeDir, NewRegistryWithPath
- `internal/arm/service/service.go` - *WithHomeDir cache methods
- `internal/arm/config/manager.go` - NewFileManager, NewFileManagerWithPaths, NewFileManagerWithConfigPath
- `internal/arm/manifest/manager.go` - NewFileManager, NewFileManagerWithPath
- `internal/arm/packagelockfile/manager.go` - NewFileManager, NewFileManagerWithPath
- `cmd/arm/main.go` - deriveLockPath helper, ARM_MANIFEST_PATH resolution
- `test/e2e/storage_test.go` - E2E cache isolation tests

## Examples

### Default Behavior (No Env Vars)
```bash
# Storage: ~/.arm/storage/
# Config: ./.armrc (overrides ~/.armrc)
# Manifest: ./arm.json
# Lock: ./arm-lock.json
```

### Custom Storage Location
```bash
export ARM_HOME=/tmp/test
# Storage: /tmp/test/.arm/storage/
# Config: ./.armrc (overrides ~/.armrc)
# Manifest: ./arm.json
# Lock: ./arm-lock.json
```

### Custom Config Location (Bypass Hierarchy)
```bash
export ARM_CONFIG_PATH=/tmp/test/.armrc
# Storage: ~/.arm/storage/
# Config: /tmp/test/.armrc (ONLY this file, no hierarchy)
# Manifest: ./arm.json
# Lock: ./arm-lock.json
```

### Custom Manifest Location
```bash
export ARM_MANIFEST_PATH=/tmp/test/my-manifest.json
# Storage: ~/.arm/storage/
# Config: ./.armrc (overrides ~/.armrc)
# Manifest: /tmp/test/my-manifest.json
# Lock: /tmp/test/my-manifest-lock.json (colocated)
```

### All Custom Paths
```bash
export ARM_HOME=/tmp/test/home
export ARM_CONFIG_PATH=/tmp/test/config/.armrc
export ARM_MANIFEST_PATH=/tmp/test/project/arm.json
# Storage: /tmp/test/home/.arm/storage/
# Config: /tmp/test/config/.armrc (ONLY this file)
# Manifest: /tmp/test/project/arm.json
# Lock: /tmp/test/project/arm-lock.json (colocated)
```

### Test Isolation (Go)
```go
func TestSomething(t *testing.T) {
    // Option 1: Environment variables
    t.Setenv("ARM_HOME", t.TempDir())
    t.Setenv("ARM_CONFIG_PATH", filepath.Join(t.TempDir(), ".armrc"))
    t.Setenv("ARM_MANIFEST_PATH", filepath.Join(t.TempDir(), "arm.json"))
    
    // Option 2: Direct constructor injection
    homeDir := t.TempDir()
    registry := storage.NewRegistryWithHomeDir("test-key", homeDir)
    
    manifestPath := filepath.Join(t.TempDir(), "arm.json")
    manifestMgr := manifest.NewFileManagerWithPath(manifestPath)
    
    // Test runs in isolation, no pollution of user directories
}
```

### Lock File Colocation Examples
```bash
# Default
ARM_MANIFEST_PATH=./arm.json
# Lock: ./arm-lock.json

# Subdirectory
ARM_MANIFEST_PATH=./config/arm.json
# Lock: ./config/arm-lock.json

# Absolute path
ARM_MANIFEST_PATH=/tmp/test/arm.json
# Lock: /tmp/test/arm-lock.json

# Custom name
ARM_MANIFEST_PATH=/tmp/test/my-manifest.json
# Lock: /tmp/test/my-manifest-lock.json

# No extension
ARM_MANIFEST_PATH=arm
# Lock: arm-lock.json
```

### Priority Order

**For .armrc lookup:**
1. ARM_CONFIG_PATH - If set, use this exact file (bypasses hierarchy)
2. ./.armrc - Project config (highest priority in hierarchy)
3. ~/.armrc - User config (fallback in hierarchy)

**For .arm/storage/ lookup:**
1. $ARM_HOME/.arm/storage/ - If ARM_HOME is set
2. ~/.arm/storage/ - Default

**For arm.json lookup:**
1. $ARM_MANIFEST_PATH - If set (resolved at CLI level)
2. ./arm.json - Default

**For arm-lock.json lookup:**
1. Always colocated with manifest file
2. Same directory as manifest
3. Filename: `strings.TrimSuffix(manifestPath, ".json") + "-lock.json"`
