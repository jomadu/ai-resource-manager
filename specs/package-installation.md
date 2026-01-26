# Package Installation

## Job to be Done
Install, update, upgrade, and uninstall packages from registries to sinks, maintaining manifest and lock files for reproducible builds.

**Note:** Packages are independent units with no inter-package dependencies. The term "dependencies" in this spec refers to packages installed in a project (similar to npm's package.json), not dependencies between packages.

## Activities
1. Install packages (rulesets/promptsets) from registries to sinks
2. Update packages to newer versions within constraint
3. Upgrade packages to latest versions
4. Uninstall packages from sinks
5. Maintain manifest (arm.json) with user constraints
6. Maintain lock file (arm-lock.json) with resolved versions

## Acceptance Criteria
- [ ] Install resolves version constraint and fetches package from registry
- [ ] Install validates all specified sinks exist before proceeding
- [ ] Install calculates integrity hash (SHA256) of fetched package files
- [ ] Install verifies integrity hash matches locked hash (if package previously installed)
- [ ] Install fails if integrity verification fails with clear error message
- [ ] Install updates manifest with package configuration (version, sinks, patterns, priority)
- [ ] Install updates lock file with resolved version and integrity hash
- [ ] Install compiles and writes files to all specified sinks
- [ ] Reinstall to different sinks removes package from old sinks (within same sink, old versions are replaced)
- [ ] Update resolves new version within existing constraint
- [ ] Update only proceeds if newer version available
- [ ] Update calculates and verifies integrity hash of new version
- [ ] Update updates lock file with new resolved version and integrity hash
- [ ] Update recompiles to all configured sinks
- [ ] Upgrade changes constraint to "latest" and resolves highest version
- [ ] Upgrade calculates and verifies integrity hash of new version
- [ ] Upgrade updates both manifest and lock file
- [ ] Uninstall removes files from all configured sinks
- [ ] Uninstall removes entries from manifest and lock file
- [ ] Install fails if registry doesn't exist
- [ ] Install fails if sink doesn't exist
- [ ] Install fails if no versions satisfy constraint
- [ ] Concurrent installs to same sink are safe (no corruption)

## Data Structures

### Manifest Dependency Config (arm.json)
```json
{
  "dependencies": {
    "registry-name/package-name": {
      "version": "^1.0.0",
      "sinks": ["cursor-rules", "q-rules"],
      "include": ["**/*.yml"],
      "exclude": ["**/experimental/**"],
      "priority": 200
    }
  }
}
```

**Fields:**
- `version` - Version constraint (e.g., "1.0.0", "^1.0.0", "~1.2.0", "latest")
- `sinks` - Array of sink names where package is installed
- `include` - Optional glob patterns to include files
- `exclude` - Optional glob patterns to exclude files
- `priority` - Optional priority for rulesets (default 100, higher wins)

### Lock File Entry (arm-lock.json)
```json
{
  "dependencies": {
    "registry-name/package-name@v1.2.3": {
      "integrity": "sha256-abc123..."
    }
  }
}
```

**Fields:**
- Key format: `registry-name/package-name@resolved-version`
- `integrity` - SHA256 hash of package contents for verification

### Package Metadata
```json
{
  "RegistryName": "test-registry",
  "Name": "clean-code-ruleset",
  "Version": "v1.2.3",
  "Integrity": "sha256-abc123..."
}
```

## Algorithm

### Install Ruleset/Promptset
1. Validate registry exists in manifest
2. Validate all specified sinks exist in manifest
3. Create registry instance from config
4. List available versions from registry
5. Resolve version using constraint (see version-resolution.md)
6. Fetch package from registry with include/exclude patterns
7. Calculate integrity hash (SHA256) of fetched files
8. Check if package already installed (exists in lock file)
9. If already installed, verify integrity hash matches locked hash
10. If integrity mismatch, return error "integrity verification failed"
11. Get old sinks from manifest (if package already installed)
12. Remove package from all old sinks
13. Update manifest with dependency config (version, sinks, patterns, priority)
14. Update lock file with resolved version and integrity hash
15. For each new sink:
   - Create sink manager
   - Compile and install package to sink
16. Return success

**Pseudocode:**
```
function InstallRuleset(registryName, packageName, version, priority, include, exclude, sinks):
    // Validate
    registryConfig = manifest.GetRegistryConfig(registryName)
    if not registryConfig:
        return error "registry not found"
    
    allSinks = manifest.GetAllSinks()
    for sinkName in sinks:
        if sinkName not in allSinks:
            return error "sink not found: " + sinkName
    
    // Resolve and fetch
    registry = createRegistry(registryName, registryConfig)
    availableVersions = registry.ListPackageVersions(packageName)
    resolvedVersion = ResolveVersion(version, availableVersions)
    package = registry.GetPackage(packageName, resolvedVersion, include, exclude)
    
    // Calculate and verify integrity
    calculatedIntegrity = CalculateIntegrity(package.files)
    
    // Check if already installed and verify integrity
    lockInfo = lockfile.GetDependencyLock(registryName, packageName)
    if lockInfo exists:
        if lockInfo.integrity != calculatedIntegrity:
            return error "integrity verification failed: expected " + lockInfo.integrity + ", got " + calculatedIntegrity
    
    // Remove from old sinks (if package already installed)
    oldDepConfig = manifest.GetRulesetDependency(registryName, packageName)
    if oldDepConfig exists:
        for oldSinkName in oldDepConfig.sinks:
            if oldSinkName in allSinks:
                oldSinkConfig = allSinks[oldSinkName]
                oldSinkManager = NewSinkManager(oldSinkConfig.Directory, oldSinkConfig.Tool)
                oldSinkManager.Uninstall(registryName, packageName)
    
    // Update manifest
    depConfig = {
        version: version,
        sinks: sinks,
        include: include,
        exclude: exclude,
        priority: priority
    }
    manifest.UpsertRulesetDependency(registryName, packageName, depConfig)
    
    // Update lock file with verified integrity
    lockConfig = {
        integrity: calculatedIntegrity
    }
    lockfile.UpsertDependencyLock(registryName, packageName, resolvedVersion, lockConfig)
    
    // Install to new sinks
    for sinkName in sinks:
        sinkConfig = allSinks[sinkName]
        sinkManager = NewSinkManager(sinkConfig.Directory, sinkConfig.Tool)
        sinkManager.InstallRuleset(package, priority)
    
    return success
```

### Update Packages
1. Get dependency config from manifest (version constraint, sinks, patterns)
2. Get current locked version from lock file
3. Create registry instance
4. List available versions from registry
5. Resolve version using existing constraint
6. If resolved version == current version, skip (already up to date)
7. Fetch package from registry
8. Calculate integrity hash (SHA256) of fetched files
9. Update lock file with new resolved version and integrity hash
10. For each configured sink:
   - Uninstall old version
   - Install new version
11. Return success

**Pseudocode:**
```
function UpdatePackage(registryName, packageName):
    // Get current state
    depConfig = manifest.GetDependency(registryName, packageName)
    lockInfo = lockfile.GetDependencyLock(registryName, packageName)
    currentVersion = lockInfo.version
    
    // Resolve new version
    registry = createRegistry(registryName, registryConfig)
    availableVersions = registry.ListPackageVersions(packageName)
    newVersion = ResolveVersion(depConfig.version, availableVersions)
    
    if newVersion == currentVersion:
        return "already up to date"
    
    // Fetch and calculate integrity
    package = registry.GetPackage(packageName, newVersion, depConfig.include, depConfig.exclude)
    calculatedIntegrity = CalculateIntegrity(package.files)
    
    // Update lock file with new integrity
    lockfile.UpsertDependencyLock(registryName, packageName, newVersion, {integrity: calculatedIntegrity})
    
    // Reinstall to all sinks
    for sinkName in depConfig.sinks:
        sinkConfig = allSinks[sinkName]
        sinkManager = NewSinkManager(sinkConfig.Directory, sinkConfig.Tool)
        sinkManager.Uninstall(registryName, packageName)
        sinkManager.InstallRuleset(package, depConfig.priority)
    
    return success
```

### Upgrade Packages
1. Get dependency config from manifest
2. Update manifest constraint to "latest"
3. Create registry instance
4. List available versions from registry
5. Resolve version using "latest" constraint (highest semver)
6. Fetch package from registry
7. Calculate integrity hash (SHA256) of fetched files
8. Update manifest with new constraint
9. Update lock file with new resolved version and integrity hash
10. For each configured sink:
   - Uninstall old version
   - Install new version
11. Return success

**Pseudocode:**
```
function UpgradePackage(registryName, packageName):
    // Get current state
    depConfig = manifest.GetDependency(registryName, packageName)
    
    // Resolve latest version
    registry = createRegistry(registryName, registryConfig)
    availableVersions = registry.ListPackageVersions(packageName)
    latestVersion = ResolveVersion("latest", availableVersions)
    
    // Fetch and calculate integrity
    package = registry.GetPackage(packageName, latestVersion, depConfig.include, depConfig.exclude)
    calculatedIntegrity = CalculateIntegrity(package.files)
    
    // Update manifest constraint to latest
    depConfig.version = "latest"
    manifest.UpsertDependency(registryName, packageName, depConfig)
    
    // Update lock file with new integrity
    lockfile.UpsertDependencyLock(registryName, packageName, latestVersion, {integrity: calculatedIntegrity})
    
    // Reinstall to all sinks
    for sinkName in depConfig.sinks:
        sinkConfig = allSinks[sinkName]
        sinkManager = NewSinkManager(sinkConfig.Directory, sinkConfig.Tool)
        sinkManager.Uninstall(registryName, packageName)
        sinkManager.InstallRuleset(package, depConfig.priority)
    
    return success
```

### Uninstall Packages
1. Get dependency config from manifest (to find sinks)
2. For each configured sink:
   - Create sink manager
   - Uninstall package (removes all files)
3. Remove entry from lock file
4. Remove entry from manifest
5. Return success

**Pseudocode:**
```
function UninstallPackage(registryName, packageName):
    // Get current state
    depConfig = manifest.GetDependency(registryName, packageName)
    allSinks = manifest.GetAllSinks()
    
    // Remove from all sinks
    for sinkName in depConfig.sinks:
        if sinkName in allSinks:
            sinkConfig = allSinks[sinkName]
            sinkManager = NewSinkManager(sinkConfig.Directory, sinkConfig.Tool)
            sinkManager.Uninstall(registryName, packageName)
    
    // Remove from manifest and lock file
    lockfile.RemoveDependencyLock(registryName, packageName)
    manifest.RemoveDependency(registryName, packageName)
    
    return success
```

### Install All (Reproducible Install)
1. Get all dependencies from manifest
2. For each dependency:
   - Get locked version from lock file
   - If no lock, resolve version from constraint
   - Fetch package from registry
   - Calculate integrity hash (SHA256) of fetched files
   - If lock exists, verify integrity hash matches locked hash
   - If integrity mismatch, return error "integrity verification failed"
   - Install to all configured sinks
3. Return success

**Pseudocode:**
```
function InstallAll():
    allDeps = manifest.GetAllDependencies()
    
    for (registryName, packageName, depConfig) in allDeps:
        // Try to use locked version
        lockInfo = lockfile.GetDependencyLock(registryName, packageName)
        if lockInfo:
            version = lockInfo.version
        else:
            // No lock, resolve from constraint
            registry = createRegistry(registryName, registryConfig)
            availableVersions = registry.ListPackageVersions(packageName)
            version = ResolveVersion(depConfig.version, availableVersions)
        
        // Fetch and calculate integrity
        package = registry.GetPackage(packageName, version, depConfig.include, depConfig.exclude)
        calculatedIntegrity = CalculateIntegrity(package.files)
        
        // Verify integrity if locked
        if lockInfo:
            if lockInfo.integrity != calculatedIntegrity:
                return error "integrity verification failed for " + registryName + "/" + packageName + ": expected " + lockInfo.integrity + ", got " + calculatedIntegrity
        
        for sinkName in depConfig.sinks:
            sinkConfig = allSinks[sinkName]
            sinkManager = NewSinkManager(sinkConfig.Directory, sinkConfig.Tool)
            sinkManager.Uninstall(registryName, packageName)
            sinkManager.InstallRuleset(package, depConfig.priority)
    
    return success
```

### Calculate Integrity Hash

1. **Create SHA256 hasher**
2. **Extract file paths** from package files
3. **Sort paths alphabetically** for deterministic ordering
4. **For each sorted path:**
   - Hash the file path
   - Hash the file content
5. **Return** "sha256-" + hex-encoded hash

**Pseudocode:**
```
function CalculateIntegrity(files):
    hasher = SHA256.New()
    
    // Extract and sort paths
    paths = []
    fileMap = {}
    for file in files:
        paths.append(file.path)
        fileMap[file.path] = file.content
    
    sort(paths)  // Alphabetical sort for deterministic ordering
    
    // Hash paths and contents in sorted order
    for path in paths:
        hasher.Write(path)
        hasher.Write(fileMap[path])
    
    return "sha256-" + HexEncode(hasher.Sum())
```

**Properties:**
- **Deterministic**: Same files always produce same hash (due to sorting)
- **Content-sensitive**: Any change to file content changes hash
- **Path-sensitive**: Renaming files changes hash
- **Collision-resistant**: SHA256 provides strong cryptographic guarantees

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Registry doesn't exist | Error: "registry not found" |
| Sink doesn't exist | Error: "sink not found: {name}" |
| No versions satisfy constraint | Error: "no version satisfies constraint" |
| Integrity verification fails | Error: "integrity verification failed: expected {expected}, got {actual}" |
| Package already installed (same version, same sinks) | Reinstall (idempotent), verify integrity |
| Package already installed (same version, different sinks) | Remove from old sinks, install to new sinks, verify integrity |
| Package already installed (different version) | Remove from old sinks, install new version to new sinks, calculate new integrity |
| Reinstall to same sink | Old version removed, new version installed, verify integrity |
| Update with no newer version | Skip, report "already up to date" |
| Upgrade to same version | Idempotent, updates manifest constraint to "latest", verify integrity |
| Uninstall non-existent package | Error: "dependency not found" |
| Uninstall with missing sink | Skip missing sink, continue with others |
| Network failure during fetch | Error propagated to user |
| Corrupted package (integrity mismatch) | Error: "integrity verification failed: expected {expected}, got {actual}" |
| Concurrent installs to same sink | Last write wins (no locking currently) |
| Install with empty sinks array | Error: "at least one sink required" |
| Manifest file missing | Create new manifest |
| Lock file missing | Resolve versions from constraints, no integrity verification |
| Lock file exists but no integrity field | Skip integrity verification (backwards compatibility) |

## Dependencies

- Version resolution (see version-resolution.md)
- Registry management (see registry-management.md)
- Sink compilation (see sink-compilation.md)
- Pattern filtering (see pattern-filtering.md)
- Cache management (see cache-management.md)

## Implementation Mapping

**Source files:**
- `internal/arm/service/service.go` - Install, update, upgrade, uninstall operations
- `internal/arm/manifest/` - Manifest file management (arm.json)
- `internal/arm/packagelockfile/` - Lock file management (arm-lock.json)
- `internal/arm/sink/manager.go` - Sink-level install/uninstall operations
- `internal/arm/registry/` - Registry implementations for fetching packages

**Related specs:**
- `version-resolution.md` - How versions are resolved
- `registry-management.md` - How registries are configured
- `sink-compilation.md` - How packages are compiled to sinks
- `pattern-filtering.md` - How include/exclude patterns work

## Examples

### Example 1: Install Ruleset

**Input:**
```bash
arm install ruleset test-registry/clean-code@^1.0.0 --priority 200 cursor-rules q-rules
```

**Expected Behavior:**
1. Resolve version: ^1.0.0 → v1.2.3 (highest matching)
2. Fetch package from test-registry
3. Update arm.json:
```json
{
  "dependencies": {
    "test-registry/clean-code": {
      "version": "^1.0.0",
      "sinks": ["cursor-rules", "q-rules"],
      "priority": 200
    }
  }
}
```
4. Update arm-lock.json:
```json
{
  "dependencies": {
    "test-registry/clean-code@v1.2.3": {
      "integrity": "sha256-abc123..."
    }
  }
}
```
5. Install to cursor-rules sink (.cursor/rules/)
6. Install to q-rules sink (.amazonq/rules/)

**Verification:**
- arm.json contains dependency with constraint ^1.0.0
- arm-lock.json contains resolved version v1.2.3
- Files exist in .cursor/rules/
- Files exist in .amazonq/rules/

### Example 2: Update Package

**Input:**
```bash
arm update test-registry/clean-code
```

**Current State:**
- Manifest constraint: ^1.0.0
- Locked version: v1.2.3
- Available versions: v1.2.3, v1.3.0, v2.0.0

**Expected Behavior:**
1. Resolve version: ^1.0.0 → v1.3.0 (highest matching, newer than v1.2.3)
2. Fetch package from test-registry
3. Update arm-lock.json with v1.3.0
4. Reinstall to all configured sinks

**Verification:**
- arm.json constraint unchanged (^1.0.0)
- arm-lock.json updated to v1.3.0
- Files updated in all sinks

### Example 3: Upgrade Package

**Input:**
```bash
arm upgrade test-registry/clean-code
```

**Current State:**
- Manifest constraint: ^1.0.0
- Locked version: v1.3.0
- Available versions: v1.3.0, v2.0.0, v2.1.0

**Expected Behavior:**
1. Resolve version: latest → v2.1.0 (highest semver)
2. Fetch package from test-registry
3. Update arm.json constraint to "latest"
4. Update arm-lock.json with v2.1.0
5. Reinstall to all configured sinks

**Verification:**
- arm.json constraint changed to "latest"
- arm-lock.json updated to v2.1.0
- Files updated in all sinks

### Example 4: Uninstall Package

**Input:**
```bash
arm uninstall test-registry/clean-code
```

**Expected Behavior:**
1. Get configured sinks from manifest
2. Remove files from cursor-rules sink
3. Remove files from q-rules sink
4. Remove entry from arm-lock.json
5. Remove entry from arm.json

**Verification:**
- Files removed from .cursor/rules/
- Files removed from .amazonq/rules/
- arm.json no longer contains dependency
- arm-lock.json no longer contains lock entry

### Example 5: Reinstall to Different Sink

**Input:**
```bash
# Initial install
arm install ruleset test-registry/clean-code@1.0.0 sink-a

# Reinstall to different sink
arm install ruleset test-registry/clean-code@1.0.0 sink-b
```

**Expected Behavior:**
1. First install: Files written to sink-a, manifest shows sinks: ["sink-a"]
2. Second install:
   - Read manifest, see package currently in sink-a
   - Remove files from sink-a
   - Install files to sink-b
   - Update manifest to show sinks: ["sink-b"]

**Verification:**
- arm.json shows sinks: ["sink-b"]
- Files exist in sink-b
- Files removed from sink-a (automatic cleanup)

### Example 6: Reinstall with Partial Overlap

**Input:**
```bash
# Initial install to two sinks
arm install ruleset test-registry/clean-code@1.0.0 sink-a sink-b

# Reinstall to different sinks with partial overlap
arm install ruleset test-registry/clean-code@1.0.0 sink-b sink-c
```

**Expected Behavior:**
1. First install: Files written to sink-a and sink-b
2. Second install:
   - Read manifest, see package currently in sink-a and sink-b
   - Remove files from sink-a (not in new list)
   - Remove files from sink-b (will be reinstalled)
   - Install files to sink-b (fresh install)
   - Install files to sink-c (new sink)
   - Update manifest to show sinks: ["sink-b", "sink-c"]

**Verification:**
- arm.json shows sinks: ["sink-b", "sink-c"]
- Files removed from sink-a
- Files exist in sink-b (reinstalled)
- Files exist in sink-c (new)

### Example 7: Install All (Reproducible)

**Input:**
```bash
# In a fresh clone with arm.json and arm-lock.json
arm install
```

**Expected Behavior:**
1. Read all dependencies from arm.json
2. For each dependency:
   - Read locked version from arm-lock.json
   - Fetch exact version from registry
   - Install to configured sinks
3. Reproduce exact environment

**Verification:**
- All dependencies installed with locked versions
- All sinks populated with correct files
- No version resolution (uses locked versions)

## Notes

**Idempotency**: Install operations are idempotent. Installing the same package version multiple times produces the same result.

**Reinstall Behavior**: When reinstalling a package to different sinks, ARM automatically removes files from all old sinks (as listed in the manifest) before installing to the new sinks. The command specifies the complete intent - the new sink list replaces the old sink list entirely.

**Within-Sink Replacement**: When installing a package to a sink that already has a different version of that package, the old version is automatically removed before installing the new version.

**Lock File Format**: The lock file uses a composite key format `registry/package@version` to support multiple versions of the same package (though currently only one version per package is supported).

**Integrity Verification**: The integrity hash (SHA256) is calculated from sorted file paths and contents, stored in the lock file, and verified during install operations. This ensures packages haven't been corrupted or tampered with. Verification is skipped if the lock file doesn't contain an integrity field (backwards compatibility).

**Integrity Calculation**: The integrity hash is calculated by:
1. Sorting all file paths alphabetically
2. For each file in sorted order: hash(path + content)
3. Return "sha256-" + hex-encoded hash

**Concurrent Safety**: Currently, ARM does not implement file locking for concurrent operations. Concurrent installs to the same sink may result in race conditions. This is a known limitation.

**Error Handling**: All operations are atomic at the file level but not at the operation level. If an install fails partway through, some sinks may be updated while others are not. Users should re-run the install to complete the operation.
