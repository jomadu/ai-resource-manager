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
  "url": "https://github.com/org/repo",
  "type": "git",
  "group_id": "123",
  "project_id": "456",
  "owner": "myorg",
  "repository": "myrepo"
}
```

**Note:** `group_id` and `project_id` are GitLab-specific. `owner` and `repository` are Cloudsmith-specific.

### Package Metadata
Package metadata is the raw package key object serialized to JSON. Fields vary based on the package key:

```json
{
  "name": "clean-code-ruleset",
  "include": ["**/*.yml"],
  "exclude": ["**/experimental/**"],
  "version": "1.0.0"
}
```

### Version Metadata
```json
{
  "version": {
    "major": 1,
    "minor": 0,
    "patch": 0,
    "prerelease": "",
    "build": "",
    "version": "1.0.0",
    "isSemver": true
  },
  "createdAt": "2024-01-15T10:30:00Z",
  "updatedAt": "2024-01-15T10:30:00Z",
  "accessedAt": "2024-01-20T14:45:00Z"
}
```

## Algorithm

### Generate Cache Key
1. JSON marshal the registry key object
2. Hash with SHA256
3. Return full 64-character hex digest

**Note:** Cache key generation relies on JSON marshaling determinism. The same registry configuration will produce the same cache key.

### Store Package
1. Generate cache key for registry
2. Create directory structure: registries/{key}/packages/{package}/{version}/
3. Write metadata.json with timestamps (createdAt, updatedAt, accessedAt)
4. Copy files to files/ subdirectory
5. Update accessedAt timestamp on subsequent reads

### Clean by Age
1. Traverse all version directories
2. Read metadata.json for each version
3. Calculate age from updatedAt timestamp
4. If age > threshold, remove version directory
5. If package has no versions, remove package directory
6. If registry has no packages, remove registry directory

### Clean by Access Time
1. Traverse all version directories
2. Read metadata.json for each version
3. Calculate time since accessedAt timestamp
4. If time > threshold, remove version directory
5. Clean up empty parent directories

### File Locking
1. Create .lock file in target directory using O_CREATE|O_EXCL for atomicity
2. Retry with 10ms sleep if lock file exists
3. Wait with timeout if lock held (default: 10s)
4. Check context cancellation during wait
5. Perform operation
6. Release lock and remove .lock file

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Cache key collision | Extremely unlikely (SHA256), but would share cache |
| Missing metadata.json | Treat as corrupted, skip or remove |
| Concurrent access | File lock prevents corruption, second process waits |
| Lock timeout | Return error after timeout (default: 10s) |
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
registryKey := map[string]interface{}{
    "url":  "https://github.com/org/repo",
    "type": "git",
}
key := GenerateKey(registryKey) // Full 64-char SHA256 hash

// Same URL, different type = different key
registryKey2 := map[string]interface{}{
    "url":  "https://github.com/org/repo",
    "type": "gitlab",
}
key2 := GenerateKey(registryKey2) // Different hash
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
