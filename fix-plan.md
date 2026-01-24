# Version Constraint Bug - Fix Summary

## Bug Summary

When installing a package with a version constraint that requires a downgrade, ARM updates the manifest and lock file but **does not replace the installed files**. The old version remains in the sink directory despite the constraint specifying a different version range.

## Root Cause

1. **Sink Manager (`internal/arm/sink/manager.go`):**
   - `InstallRuleset()` and `InstallPromptset()` only checked if the **exact same version** was already installed
   - Did not remove **other versions** of the same package before installing

2. **Lock File Manager (`internal/arm/packagelockfile/manager.go`):**
   - `UpsertDependencyLock()` did not remove old version entries before adding new ones
   - Multiple version entries accumulated in the lock file

## Solution Implemented

### 1. Sink Manager Changes

**File:** `internal/arm/sink/manager.go`

- **Simplified `Uninstall(registryName, packageName)`:**
  - Changed signature from `Uninstall(metadata *core.PackageMetadata)` to `Uninstall(registryName, packageName string)`
  - Now removes **all versions** of a package (by registry + name)
  - Iterates through index and removes all matching entries

- **Updated `InstallRuleset()` and `InstallPromptset()`:**
  - Now simply call `Uninstall(pkg.Metadata.RegistryName, pkg.Metadata.Name)` before installing
  - No need to list and iterate - `Uninstall()` handles it

### 2. Lock File Manager Changes

**File:** `internal/arm/packagelockfile/manager.go`

- **Enhanced `UpsertDependencyLock()`:**
  - Now calls `RemoveDependencyLock(registry, packageName)` internally before upserting
  - Automatically cleans up old versions

- **Simplified `RemoveDependencyLock(registry, packageName)`:**
  - Changed signature from `RemoveDependencyLock(registry, packageName, version)` to `RemoveDependencyLock(registry, packageName)`
  - Removes **all versions** of a package
  - Uses prefix matching to find and delete all version entries

### 3. Service Layer Changes

**File:** `internal/arm/service/service.go`

- **Simplified `InstallRuleset()` and `InstallPromptset()`:**
  - Removed explicit `RemoveDependencyLock()` calls
  - Just call `UpsertDependencyLock()` - it handles cleanup automatically

- **Updated all `Uninstall()` calls:**
  - Changed from `sinkMgr.Uninstall(&core.PackageMetadata{...})` to `sinkMgr.Uninstall(registryName, packageName)`
  - Removed unnecessary version parameter

### 4. Interface Updates

**File:** `internal/arm/packagelockfile/interface.go`

- Updated `RemoveDependencyLock` signature to `(ctx, registry, packageName)`
- Removed separate `RemoveDependencyLocks` method

### 5. Test Updates

- Updated mock implementations to match new signatures
- Updated test cases to use simplified API
- All tests pass

## Benefits

1. **Automatic cleanup:** Both managers handle version cleanup internally
2. **Simpler API:** Service layer just calls install/upsert - no manual cleanup needed
3. **Proper encapsulation:** Complexity is in the right layers
4. **Cleaner code:** Service layer is much more readable and maintainable
5. **No duplicate logic:** Single implementation of cleanup in each manager

## Verification

✅ Lock file only contains the new version entry  
✅ Old version files removed from sink directories  
✅ New version files correctly installed  
✅ All tests pass  
✅ Manual testing confirms correct behavior

## Test Case

```bash
# Install latest (v2.1.0)
arm install ruleset ai-rules/grug-brained-dev --include "rulesets/grug-brained-dev.yml" q-rules

# Install with @1 constraint (installs v1.1.0, removes v2.1.0)
arm install ruleset ai-rules/grug-brained-dev@1 --include "rules/amazonq/*.md" q-rules
```

**Result:** Only v1.1.0 files exist, lock file contains only v1.1.0 entry.
