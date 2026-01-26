# Test Isolation via Constructor Injection

## Job to be Done
Enable test isolation by injecting directory paths directly into component constructors, preventing tests from polluting user's home directory.

## Activities
1. **Accept directory paths as constructor parameters** - Components receive paths instead of calling os.UserHomeDir()
2. **Provide default constructors** - Production code calls os.UserHomeDir() once in constructor
3. **Provide test constructors** - Tests pass t.TempDir() directly as string parameters

## Acceptance Criteria
- [ ] Components accept home directory path as constructor parameter
- [ ] Default constructors call os.UserHomeDir() internally for production use
- [ ] Test constructors accept directory paths directly (no OS calls)
- [ ] No direct os.UserHomeDir() calls in component methods
- [ ] Tests pass t.TempDir() to test constructors
- [ ] Tests don't pollute user's actual home directory

## Pattern

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

### After (Constructor Injection)
```go
type Registry struct {
    registryKey interface{}
    registryDir string
}

// Production constructor - calls os.UserHomeDir() once
func NewRegistry(registryKey interface{}) (*Registry, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return nil, err
    }
    return NewRegistryWithHomeDir(registryKey, homeDir)
}

// Test constructor - accepts directory path directly
func NewRegistryWithHomeDir(registryKey interface{}, homeDir string) (*Registry, error) {
    baseDir := filepath.Join(homeDir, ".arm")
    return NewRegistryWithPath(baseDir, registryKey)
}
```

**Benefits:** Tests pass t.TempDir() to NewRegistryWithHomeDir()

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

**Already correct!** - Already has NewFileManagerWithPaths() for testing

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
        var err error
        homeDir, err = os.UserHomeDir()
        if err != nil {
            return err
        }
    }
    storageDir := filepath.Join(homeDir, ".arm", "storage")
    return s.cleanCacheByAgeWithPath(ctx, maxAge, storageDir)
}
```

## Test Usage

### Production
```go
// Uses actual user home directory
registry, err := storage.NewRegistry(registryKey)
```

### Testing
```go
func TestRegistry(t *testing.T) {
    testHome := t.TempDir()
    registry, err := storage.NewRegistryWithHomeDir(registryKey, testHome)
    // Uses testHome instead of actual user home
}
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| os.UserHomeDir() returns error | Production constructor propagates error |
| Empty string passed to test constructor | Creates directories relative to empty path (caller's responsibility) |
| Non-existent directory passed | Component creates directories as needed |

## Benefits

### Simplicity
- No interfaces or abstractions needed
- Just string parameters
- Minimal code changes

### Test Reliability
- Tests use isolated temporary directories
- No pollution of user's ~/.arm/ or ~/.armrc
- Parallel test execution is safe
- Automatic cleanup via t.TempDir()

### Backward Compatibility
- Existing production code unchanged (uses default constructors)
- Existing tests can migrate incrementally
- No breaking changes to public API

## Implementation Mapping

**Components to update:**
- `internal/arm/storage/registry.go` - Add NewRegistryWithHomeDir()
- `internal/arm/service/service.go` - Add *WithHomeDir() variants for cache methods

**Already correct:**
- `internal/arm/config/manager.go` - Already has NewFileManagerWithPaths()
- `internal/arm/manifest/manager.go` - Uses relative paths (no home dir needed)
- `internal/arm/packagelockfile/manager.go` - Uses relative paths (no home dir needed)

## Related Specs
- `cache-management.md` - Storage components accept homeDir parameter
- `authentication.md` - Config manager accepts workingDir and homeDir parameters
- `e2e-testing.md` - Tests use direct path injection
