# Cache Management

## Job to be Done
Cache downloaded packages locally to avoid redundant network requests, enabling faster installs and offline access to previously downloaded versions.

## Activities
1. **Generate Cache Keys** - Create unique hashes for registries and packages
2. **Store Package Versions** - Save downloaded files with metadata and timestamps
3. **Retrieve Cached Versions** - Load previously downloaded packages from disk
4. **Update Access Times** - Track when versions are accessed for cleanup decisions
5. **Clean Old Versions** - Remove cached versions based on age or last access time
6. **Nuke Cache** - Delete entire cache directory

## Acceptance Criteria
- [ ] Registry key generated from registry configuration (URL, type, group_id, project_id, owner, repository)
- [ ] Package key generated from package identification (name for non-git, includes/excludes patterns)
- [ ] Package versions stored in `~/.arm/storage/registries/{registry-key}/packages/{package-key}/{version}/files/`
- [ ] Metadata stored at three levels: registry, package, version
- [ ] Version metadata includes createdAt, updatedAt, accessedAt timestamps
- [ ] Version metadata optionally includes integrity hash for verification
- [ ] Access time updated on every GetPackageVersion call
- [ ] Integrity verification (optional) detects corrupted cached packages
- [ ] Clean by age removes versions where updatedAt < cutoff
- [ ] Clean by last access removes versions where accessedAt < cutoff
- [ ] Empty package directories removed after version cleanup
- [ ] Nuke removes entire storage directory
- [ ] Cross-process file locking prevents concurrent access corruption
- [ ] Git repositories cached in `repo/` subdirectory (git registries only)
- [ ] Storage components accept home directory as constructor parameter for test isolation
- [ ] Default constructors call os.UserHomeDir() for production use
- [ ] Test constructors accept temporary directory paths for isolation

## Data Structures

### Storage Directory Structure
```
~/.arm/storage/
    registries/
        {registry-key}/
            metadata.json           # Registry metadata
            repo/                   # Git clone (git registries only)
                .git/
                # Repository files
            packages/
                {package-key}/
                    metadata.json   # Package metadata
                    v1.0.0/
                        metadata.json # Version metadata + timestamps
                        files/
                            # Extracted package files
                    v1.1.0/
                        metadata.json
                        files/
                            # Extracted package files
```

### Registry Metadata
```json
{
  "url": "https://github.com/example/repo",
  "type": "git",
  "group_id": "123",
  "project_id": "456",
  "owner": "sample-org",
  "repository": "arm-registry"
}
```

**Fields:**
- `url` - Registry URL
- `type` - Registry type (git, gitlab, cloudsmith)
- `group_id` - GitLab group ID (optional)
- `project_id` - GitLab project ID (optional)
- `owner` - Cloudsmith owner (optional)
- `repository` - Cloudsmith repository (optional)

### Package Metadata
```json
{
  "name": "clean-code-ruleset",
  "includes": ["**/*.yml"],
  "excludes": ["**/test/**"]
}
```

**Fields:**
- `name` - Package name (omitted for git registries)
- `includes` - Include patterns (normalized)
- `excludes` - Exclude patterns (normalized)

### Version Metadata
```json
{
  "version": {
    "major": 1,
    "minor": 0,
    "patch": 0
  },
  "integrity": "sha256-abc123def456...",
  "createdAt": "2025-01-08T23:10:43.984784Z",
  "updatedAt": "2025-01-08T23:10:43.984784Z",
  "accessedAt": "2025-01-08T23:10:43.984784Z"
}
```

**Fields:**
- `version` - Semantic version (major, minor, patch)
- `integrity` - SHA256 hash of package contents for verification (optional, recommended)
- `createdAt` - When version was first cached
- `updatedAt` - When version was last updated (currently same as createdAt)
- `accessedAt` - When version was last accessed (updated on every read)

### File Lock
```go
type FileLock struct {
    lockPath string // {directory}.lock
}
```

**Purpose:**
- Prevents concurrent access to same package/registry
- Cross-process locking using file system
- Lock acquired before read/write operations
- Lock released after operation completes
- Prevents concurrent access to same package/registry
- Cross-process locking using file system
- Lock acquired before read/write operations
- Lock released after operation completes

## Algorithm

### Generate Registry Key

1. **Marshal registry config** to JSON
2. **Hash JSON** using SHA256
3. **Return hex-encoded hash**

**Pseudocode:**
```
function GenerateRegistryKey(config):
    keyObj = {
        url: config.url,
        type: config.type,
        group_id: config.group_id,
        project_id: config.project_id,
        owner: config.owner,
        repository: config.repository
    }
    json = Marshal(keyObj)
    hash = SHA256(json)
    return HexEncode(hash)
```

### Generate Package Key

1. **Build key object** (name for non-git, includes/excludes for all)
2. **Marshal to JSON**
3. **Hash JSON** using SHA256
4. **Return hex-encoded hash**

**Pseudocode:**
```
function GeneratePackageKey(name, includes, excludes, isGit):
    if isGit:
        keyObj = {
            includes: Normalize(includes),
            excludes: Normalize(excludes)
        }
    else:
        keyObj = {
            name: name,
            includes: Normalize(includes),
            excludes: Normalize(excludes)
        }
    json = Marshal(keyObj)
    hash = SHA256(json)
    return HexEncode(hash)
```

### Store Package Version

1. **Acquire package lock** (cross-process)
2. **Generate package key** from package metadata
3. **Create directory structure** (packages/{key}/v{version}/files/)
4. **Write files** to files/ directory
5. **Write package metadata** (packages/{key}/metadata.json)
6. **Write version metadata** with timestamps (packages/{key}/v{version}/metadata.json)
7. **Release lock**

**Pseudocode:**
```
function SetPackageVersion(packageKey, version, files):
    lock = GetPackageLock(packageKey)
    lock.Acquire()
    defer lock.Release()
    
    hashedKey = GenerateKey(packageKey)
    versionDir = packagesDir + "/" + hashedKey + "/v" + version
    filesDir = versionDir + "/files"
    
    CreateDirectories(filesDir)
    
    for file in files:
        filePath = filesDir + "/" + file.path
        WriteFile(filePath, file.content)
    
    now = Now()
    packageMetadata = packageKey
    WriteJSON(packagesDir + "/" + hashedKey + "/metadata.json", packageMetadata)
    
    versionMetadata = {
        version: version,
        integrity: integrity,  // Optional: SHA256 hash for verification
        createdAt: now,
        updatedAt: now,
        accessedAt: now
    }
    WriteJSON(versionDir + "/metadata.json", versionMetadata)
```

### Retrieve Package Version

1. **Acquire package lock** (cross-process)
2. **Generate package key** from package metadata
3. **Check if version directory exists** (return error if not)
4. **Update accessedAt timestamp** in version metadata
5. **Read all files** from files/ directory
6. **(Optional) Verify integrity** if stored in metadata
7. **Release lock**
8. **Return files**

**Pseudocode:**
```
function GetPackageVersion(packageKey, version, expectedIntegrity):
    lock = GetPackageLock(packageKey)
    lock.Acquire()
    defer lock.Release()
    
    hashedKey = GenerateKey(packageKey)
    versionDir = packagesDir + "/" + hashedKey + "/v" + version
    filesDir = versionDir + "/files"
    
    if not Exists(filesDir):
        return error("package version not found")
    
    // Update access time
    metadata = ReadJSON(versionDir + "/metadata.json")
    metadata.accessedAt = Now()
    WriteJSON(versionDir + "/metadata.json", metadata)
    
    // Read files
    files = []
    Walk(filesDir, func(path, info):
        if not info.IsDir():
            relPath = RelativePath(filesDir, path)
            content = ReadFile(path)
            files.append({path: relPath, content: content})
    )
    
    // Optional: Verify cached integrity against expected
    if expectedIntegrity and metadata.integrity:
        if metadata.integrity != expectedIntegrity:
            // Cache corrupted - delete and return error to trigger re-fetch
            DeleteDirectory(versionDir)
            return error("cached package integrity mismatch")
    
    return files
```

### Clean Cache by Age

1. **List all registry directories**
2. **For each registry:**
   - Get packages directory
   - Create PackageCache
   - Call RemoveOldVersions(maxAge)
3. **RemoveOldVersions:**
   - List all package directories
   - For each package:
     - List version directories
     - For each version:
       - Read version metadata
       - If updatedAt < cutoff: delete version directory
     - If no versions remain: delete package directory

**Pseudocode:**
```
function CleanCacheByAge(maxAge):
    cutoff = Now() - maxAge
    registries = ListDirectories(storageDir + "/registries")
    
    for registry in registries:
        packagesDir = registry + "/packages"
        if not Exists(packagesDir):
            continue
        
        packages = ListDirectories(packagesDir)
        for package in packages:
            versions = ListDirectories(package)
            hasVersions = false
            
            for version in versions:
                metadata = ReadJSON(version + "/metadata.json")
                if metadata.updatedAt < cutoff:
                    DeleteDirectory(version)
                else:
                    hasVersions = true
            
            if not hasVersions:
                DeleteDirectory(package)
```

### Clean Cache by Last Access

1. **List all registry directories**
2. **For each registry:**
   - Get packages directory
   - Create PackageCache
   - Call RemoveUnusedVersions(maxTimeSinceLastAccess)
3. **RemoveUnusedVersions:**
   - List all package directories
   - For each package:
     - List version directories
     - For each version:
       - Read version metadata
       - If accessedAt < cutoff: delete version directory
     - If no versions remain: delete package directory

**Pseudocode:**
```
function CleanCacheByLastAccess(maxTimeSinceLastAccess):
    cutoff = Now() - maxTimeSinceLastAccess
    registries = ListDirectories(storageDir + "/registries")
    
    for registry in registries:
        packagesDir = registry + "/packages"
        if not Exists(packagesDir):
            continue
        
        packages = ListDirectories(packagesDir)
        for package in packages:
            versions = ListDirectories(package)
            hasVersions = false
            
            for version in versions:
                metadata = ReadJSON(version + "/metadata.json")
                if metadata.accessedAt < cutoff:
                    DeleteDirectory(version)
                else:
                    hasVersions = true
            
            if not hasVersions:
                DeleteDirectory(package)
```

### Nuke Cache

1. **Delete entire storage directory**

**Pseudocode:**
```
function NukeCache():
    DeleteDirectory(storageDir)
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Corrupted metadata.json | Treat as missing. Skip or delete directory. |
| Missing version metadata | Skip version during cleanup. Don't delete. |
| Cached package integrity mismatch | Delete corrupted cache. Re-fetch from registry. |
| Missing integrity field in cache metadata | Skip integrity verification (backwards compatibility). |
| Disk full during write | Return error. Partial writes may leave corrupted cache. |
| Concurrent access (same process) | File lock prevents corruption. Second caller waits. |
| Concurrent access (different processes) | File lock prevents corruption. Second process waits. |
| Lock file stale (process crashed) | Lock implementation should handle stale locks (timeout or force). |
| Cache directory deleted manually | Recreated on next install. No data loss (source is registry). |
| Version directory exists but empty | Treated as valid. Not cleaned up unless metadata indicates old. |
| Invalid version directory name | Skipped during cleanup. Not treated as version. |
| maxAge = 0 | Deletes all versions (cutoff = now). |
| maxAge negative | Undefined. Should validate and reject. |
| Nuke during active install | File lock prevents corruption. Nuke waits or fails. |
| Git repo/ directory corrupted | Registry operations fail. User must clean cache manually. |
| Symlinks in cache | Followed during read. Potential security risk if user-created. |

## Dependencies

- **File System** - Directory and file operations
- **JSON Marshaling** - Metadata serialization
- **SHA256 Hashing** - Key generation
- **File Locking** - Cross-process concurrency control
- **Time** - Timestamp tracking and comparison
- **SHA256 Hashing** - Integrity calculation (optional, for verification)

## Implementation Mapping

**Source files:**
- `internal/arm/storage/storage.go` - GenerateKey, directory structure documentation
- `internal/arm/storage/registry.go` - Registry metadata and directory management (accepts homeDir parameter)
- `internal/arm/storage/package.go` - PackageCache operations (Set, Get, Remove, Cleanup)
- `internal/arm/storage/repo.go` - Git repository management
- `internal/arm/storage/lock.go` - FileLock implementation
- `internal/arm/service/service.go` - CleanCacheByAge, CleanCacheByTimeSinceLastAccess, NukeCache (accepts homeDir parameter)
- `internal/arm/registry/integrity.go` - Integrity calculation (used by cache for verification)
- `internal/arm/config/manager.go` - Config file manager (accepts workingDir and homeDir parameters)

**Related specs:**
- `package-installation.md` - Uses cache to store/retrieve packages, defines integrity verification requirements
- `registry-management.md` - Registry configuration used for key generation
- `pattern-filtering.md` - Include/exclude patterns used in package keys

## Examples

### Example 1: Store and Retrieve Package

**Input:**
```go
packageKey := map[string]interface{}{
    "name": "clean-code-ruleset",
    "includes": []string{"**/*.yml"},
    "excludes": []string{},
}
version := core.Version{Major: 1, Minor: 0, Patch: 0}
files := []*core.File{
    {Path: "rules/rule1.yml", Content: []byte("...")},
    {Path: "rules/rule2.yml", Content: []byte("...")},
}

cache.SetPackageVersion(ctx, packageKey, &version, files)
```

**Expected Output:**

Directory structure:
```
~/.arm/storage/registries/{registry-key}/packages/{package-key}/
    metadata.json
    v1.0.0/
        metadata.json
        files/
            rules/
                rule1.yml
                rule2.yml
```

metadata.json (package):
```json
{
  "name": "clean-code-ruleset",
  "includes": ["**/*.yml"],
  "excludes": []
}
```

metadata.json (version):
```json
{
  "version": {"major": 1, "minor": 0, "patch": 0},
  "createdAt": "2025-01-25T14:00:00Z",
  "updatedAt": "2025-01-25T14:00:00Z",
  "accessedAt": "2025-01-25T14:00:00Z"
}
```

**Verification:**
- Files written to correct paths
- Metadata files created
- Timestamps set to current time

### Example 2: Integrity Verification (Optional)

**Input:**
```go
// Store package with integrity
packageKey := map[string]interface{}{
    "name": "clean-code-ruleset",
}
version := core.Version{Major: 1, Minor: 0, Patch: 0}
files := []*core.File{
    {Path: "rule1.yml", Content: []byte("content1")},
}
integrity := "sha256-abc123..."

cache.SetPackageVersion(ctx, packageKey, &version, files, integrity)

// Later: Retrieve and verify
expectedIntegrity := "sha256-abc123..."
files, err := cache.GetPackageVersion(ctx, packageKey, &version, expectedIntegrity)
```

**Expected Output:**

metadata.json (version):
```json
{
  "version": {"major": 1, "minor": 0, "patch": 0},
  "integrity": "sha256-abc123...",
  "createdAt": "2025-01-25T14:00:00Z",
  "updatedAt": "2025-01-25T14:00:00Z",
  "accessedAt": "2025-01-25T14:00:00Z"
}
```

**Verification:**
- Integrity stored in version metadata
- Retrieval succeeds when integrity matches
- Retrieval fails and deletes cache if integrity mismatches

### Example 3: Update Access Time

**Input:**
```go
cache.GetPackageVersion(ctx, packageKey, &version, "")
// Wait 1 hour
cache.GetPackageVersion(ctx, packageKey, &version, "")
```

**Expected Output:**

metadata.json (version) after first access:
```json
{
  "accessedAt": "2025-01-25T14:00:00Z"
}
```

metadata.json (version) after second access:
```json
{
  "accessedAt": "2025-01-25T15:00:00Z"
}
```

**Verification:**
- accessedAt updated on each GetPackageVersion call
- createdAt and updatedAt unchanged

### Example 4: Clean by Age

**Input:**
```bash
arm clean cache --max-age 7d
```

**State:**
- Package A v1.0.0: updatedAt = 10 days ago
- Package A v1.1.0: updatedAt = 5 days ago
- Package B v2.0.0: updatedAt = 8 days ago

**Expected Output:**
- Package A v1.0.0: DELETED (older than 7 days)
- Package A v1.1.0: KEPT (newer than 7 days)
- Package B v2.0.0: DELETED (older than 7 days)
- Package B directory: DELETED (no versions remain)

**Verification:**
- Only versions older than cutoff deleted
- Empty package directories removed
- Package A directory kept (has v1.1.0)

### Example 5: Clean by Last Access

**Input:**
```bash
arm clean cache --max-age 7d
```

**State:**
- Package A v1.0.0: accessedAt = 10 days ago
- Package A v1.1.0: accessedAt = 5 days ago
- Package B v2.0.0: accessedAt = 8 days ago

**Expected Output:**
- Package A v1.0.0: DELETED (not accessed in 7 days)
- Package A v1.1.0: KEPT (accessed within 7 days)
- Package B v2.0.0: DELETED (not accessed in 7 days)
- Package B directory: DELETED (no versions remain)

**Verification:**
- Only versions not accessed within cutoff deleted
- Empty package directories removed
- Package A directory kept (has v1.1.0)

### Example 6: Nuke Cache

**Input:**
```bash
arm clean cache --nuke
```

**Expected Output:**
- Entire `~/.arm/storage/` directory deleted
- All registries, packages, versions removed
- Git repositories removed

**Verification:**
- Storage directory does not exist
- Next install recreates directory structure

### Example 7: Concurrent Access

**Input:**
```go
// Process 1
go cache.SetPackageVersion(ctx, packageKey, &version1, files1, integrity1)

// Process 2 (simultaneously)
go cache.SetPackageVersion(ctx, packageKey, &version2, files2, integrity2)
```

**Expected Output:**
- One process acquires lock first
- Second process waits for lock
- Both operations complete successfully
- No corrupted metadata or files

**Verification:**
- File lock prevents concurrent writes
- Both versions stored correctly
- No race conditions or data corruption

## Notes

### Why Store Integrity in Cache?

**Benefits:**
- Detect cache corruption without re-fetching from registry
- Enable offline integrity verification
- Faster verification (no need to recalculate from files)

**Trade-offs:**
- Optional feature (not required for basic caching)
- Adds small overhead to metadata storage
- Requires integrity calculation before caching

### Why Three-Level Metadata?

**Registry metadata:**
- Identifies registry configuration
- Enables cache inspection without reading manifest

**Package metadata:**
- Identifies package by name and patterns
- Enables listing packages without reading versions

**Version metadata:**
- Tracks timestamps for cleanup decisions
- Enables age-based and access-based cleanup

### Why Update Access Time?

Access time tracking enables intelligent cleanup:
- Frequently used versions kept (even if old)
- Unused versions removed (even if recent)
- Balances disk space with performance

### Why Hash-Based Keys?

Hash-based keys provide:
- Unique identification (collision-resistant)
- Filesystem-safe names (no special characters)
- Consistent length (no path length issues)
- Privacy (hides registry URLs from directory names)

### Why Per-Package Locking?

Per-package locking provides:
- Fine-grained concurrency (multiple packages simultaneously)
- Prevents deadlocks (no lock ordering issues)
- Simpler implementation (no global lock coordination)

### Why Not Database?

File-based storage provides:
- Simplicity (no database dependency)
- Portability (works on all platforms)
- Transparency (users can inspect cache manually)
- Reliability (no database corruption issues)

### Design Decisions

**Why store files in files/ subdirectory?**
- Separates files from metadata
- Prevents name collisions (metadata.json vs package files)
- Cleaner directory structure

**Why use updatedAt for age-based cleanup?**
- Represents when version was cached
- Independent of access patterns
- Predictable cleanup behavior

**Why delete empty package directories?**
- Prevents orphaned directories
- Cleaner cache structure
- Easier cache inspection

**Why not compress cached files?**
- Simplicity (no compression/decompression overhead)
- Performance (faster reads)
- Transparency (users can inspect files directly)

**Why not deduplicate files across versions?**
- Simplicity (no content-addressable storage)
- Reliability (no broken references)
- Isolation (version deletion doesn't affect others)

### Testing Considerations

Tests should verify:
- Registry key generation (different configs produce different keys)
- Package key generation (name and patterns affect key)
- Store and retrieve round-trip (files match)
- Access time updated on read
- Clean by age removes old versions
- Clean by last access removes unused versions
- Empty package directories removed
- Nuke removes entire storage
- Concurrent access prevented by locks
- Corrupted metadata handled gracefully
- Missing directories created automatically
