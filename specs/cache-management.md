# Cache Management

## Job to be Done
Store downloaded packages locally to avoid redundant downloads, enable offline usage, and provide cleanup mechanisms for old cached data.

## Activities
1. Store packages in ~/.arm/storage/ organized by registry and version
2. Generate cache keys from registry configuration
3. Track access times for cache cleanup
4. Clean cache by age or last access time
5. Lock cache during concurrent operations

## Acceptance Criteria
- [x] Store packages in ~/.arm/storage/registries/{key}/packages/{package}/{version}/
- [x] Generate consistent cache keys from registry URL and configuration
- [x] Track creation time and last access time in metadata.json
- [x] Clean cache by age (remove versions older than threshold)
- [x] Clean cache by access time (remove versions not accessed recently)
- [x] Nuke entire cache (remove all cached data)
- [x] File locking prevents concurrent access corruption
- [x] Support ARM_HOME environment variable for custom cache location

## Data Structures

### Registry Metadata
```json
{
  "key": "abc123def456",
  "url": "https://github.com/org/repo",
  "created": "2024-01-15T10:30:00Z"
}
```

### Package Metadata
```json
{
  "name": "clean-code-ruleset",
  "created": "2024-01-15T10:30:00Z"
}
```

### Version Metadata
```json
{
  "version": "v1.0.0",
  "created": "2024-01-15T10:30:00Z",
  "accessed": "2024-01-20T14:45:00Z",
  "commit": "abc123def456",
  "integrity": "sha256-xyz789..."
}
```

## Algorithm

### Generate Cache Key
1. Extract registry URL and configuration
2. Normalize URL (strip protocol, trailing slashes)
3. Sort configuration keys alphabetically
4. Concatenate URL + sorted config
5. Hash with SHA256
6. Return first 16 characters of hex digest

### Store Package
1. Generate cache key for registry
2. Create directory structure: registries/{key}/packages/{package}/{version}/
3. Write metadata.json with timestamps
4. Copy files to files/ subdirectory
5. Update accessed timestamp

### Clean by Age
1. Traverse all version directories
2. Read metadata.json for each version
3. Calculate age from created timestamp
4. If age > threshold, remove version directory
5. If package has no versions, remove package directory
6. If registry has no packages, remove registry directory

### Clean by Access Time
1. Traverse all version directories
2. Read metadata.json for each version
3. Calculate time since last access
4. If time > threshold, remove version directory
5. Clean up empty parent directories

### File Locking
1. Create .lock file in target directory
2. Use flock (Unix) or LockFileEx (Windows)
3. Wait with timeout if lock held
4. Perform operation
5. Release lock and remove .lock file

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Cache key collision | Extremely unlikely (SHA256), but would share cache |
| Missing metadata.json | Treat as corrupted, skip or remove |
| Concurrent access | File lock prevents corruption, second process waits |
| Lock timeout | Return error after timeout (default: 30s) |
| ARM_HOME not set | Use os.UserHomeDir() |
| ARM_HOME set | Use ARM_HOME/.arm/storage/ |
| Partial download | Atomic write or cleanup on error |
| Disk full | Return error, don't corrupt cache |

## Dependencies

- Registry configuration (manifest)
- File system operations
- SHA256 hashing

## Implementation Mapping

**Source files:**
- `internal/arm/storage/registry.go` - NewRegistry, GetRegistryDir, GetPackagesDir
- `internal/arm/storage/package.go` - SetPackageVersion, GetPackageVersion, RemoveOldVersions
- `internal/arm/storage/storage.go` - GenerateKey
- `internal/arm/storage/lock.go` - FileLock, Lock, Unlock
- `internal/arm/service/service.go` - CleanCacheByAge, CleanCacheByTimeSinceLastAccess, NukeCache
- `test/e2e/storage_test.go` - E2E cache tests

## Examples

### Cache Structure
```
~/.arm/storage/
└── registries/
    └── abc123def456/          # Registry cache key
        ├── metadata.json
        └── packages/
            └── clean-code-ruleset/
                ├── metadata.json
                └── v1.0.0/
                    ├── metadata.json
                    └── files/
                        ├── rules/
                        │   └── clean-code.yml
                        └── build/
                            └── cursor/
                                └── clean-code.mdc
```

### Cache Key Generation
```go
// Git registry
url := "https://github.com/org/repo"
branches := []string{"main", "develop"}
key := GenerateKey(url, branches) // "abc123def456"

// Same URL, different branches = different key
branches2 := []string{"main"}
key2 := GenerateKey(url, branches2) // "def789abc123"
```

### Clean by Age
```bash
# Remove versions older than 30 days
arm clean cache --max-age 30d

# Remove versions older than 7 days
arm clean cache --max-age 7d
```

### Clean by Access Time
```bash
# Remove versions not accessed in 30 days
arm clean cache --max-age 30d --by-access

# Remove versions not accessed in 7 days
arm clean cache --max-age 7d --by-access
```

### Nuke Cache
```bash
# Remove all cached data
arm clean cache --nuke
```

### Custom Cache Location
```bash
# Use custom cache directory
export ARM_HOME=/custom/path
arm install ruleset ai-rules/clean-code cursor-rules
# Caches to /custom/path/.arm/storage/
```
