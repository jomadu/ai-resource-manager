# Query Operations

## Job to be Done
Query installed packages, check for outdated dependencies, view detailed package information, and list available versions from registries to maintain awareness of project dependencies and discover packages.

## Activities
1. List all installed packages with versions and sinks
2. Check for outdated packages with available updates
3. View detailed information about specific packages
4. List available versions for a package from its registry
5. Display information in multiple formats (table, JSON, list)

## Acceptance Criteria
- [x] List all installed rulesets and promptsets
- [x] Show current version, constraint, and target sinks for each package
- [x] Check for outdated packages comparing current vs latest versions
- [x] Display outdated packages in table or JSON format
- [x] View detailed info for specific package (registry/package format)
- [ ] List available versions for a package from registry
- [x] Handle missing manifest gracefully
- [x] Handle missing lock file gracefully
- [x] Sort output alphabetically for deterministic results

## Data Structures

### DependencyInfo
```go
type DependencyInfo struct {
    Name         string   // registry/package
    Type         string   // "ruleset" or "promptset"
    Version      string   // Current version from lock file
    Constraint   string   // Version constraint from manifest
    Sinks        []string // Target sink names
}
```

### OutdatedDependency
```go
type OutdatedDependency struct {
    Name         string   // registry/package
    Type         string   // "ruleset" or "promptset"
    Current      string   // Current version from lock file
    Latest       string   // Latest available version from registry
    Constraint   string   // Version constraint from manifest
}
```

## Algorithm

### List All Dependencies
1. Read manifest file (arm.json)
2. Read lock file (arm-lock.json)
3. For each dependency in manifest:
   - Extract name, type, constraint, sinks
   - Find matching lock entry by composite key (registry/package@version)
   - Combine manifest and lock data
4. Sort results alphabetically by name
5. Return list of DependencyInfo

### Check Outdated
1. Read manifest and lock files
2. For each dependency:
   - Get current version from lock file
   - Query registry for latest version
   - Compare current vs latest
   - If different, add to outdated list
3. Sort results alphabetically by name
4. Return list of OutdatedDependency

### Get Dependency Info
1. Parse package key (registry/package)
2. Read manifest for dependency config
3. Read lock file for installed version
4. Return combined DependencyInfo

### List Available Versions
1. Parse package key (registry/package)
2. Extract registry name
3. Create registry instance from manifest config
4. Call registry.ListPackageVersions(packageName)
5. Sort versions (semver descending, branches in config order)
6. Return version list

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| No manifest file | Error: "manifest file not found" |
| No lock file | Show dependencies with "not installed" status |
| Empty dependencies | Return empty list (not an error) |
| Registry unavailable | Error: "failed to fetch latest version" |
| Package not in manifest | Error: "package not found in manifest" |
| Malformed package key | Error: "invalid package format (expected registry/package)" |
| Registry not configured | Error: "registry not found: {name}" |
| Package not found in registry | Error: "package not found in registry" |

## Dependencies

- Manifest Manager - Read arm.json
- Lock File Manager - Read arm-lock.json
- Registry Factory - Create registry instances
- Registry Interface - Query latest versions

## Implementation Mapping

**Source files:**
- `internal/arm/service/service.go` - ListAll, ListOutdated, GetDependencyInfo methods
- `internal/arm/manifest/manager.go` - GetAllDependenciesConfig, GetDependencyConfig
- `internal/arm/packagelockfile/manager.go` - GetLockFile, GetDependencyLock
- `cmd/arm/main.go` - handleList, handleInfo, handleOutdated commands

**Related specs:**
- `package-installation.md` - Manifest and lock file structure
- `version-resolution.md` - Version comparison logic

## Examples

### Example 1: List All Dependencies

**Command:**
```bash
arm list dependency
```

**Expected Output:**
```
- test-registry/clean-code-ruleset@1.0.0
- test-registry/code-review@2.1.0
```

**Verification:**
- All dependencies from arm.json are listed
- Versions match arm-lock.json
- Output is sorted alphabetically
- Format: dash-prefixed list with `@version` suffix

### Example 2: Check Outdated Dependencies

**Command:**
```bash
arm outdated
```

**Expected Output (Table - default):**
```
PACKAGE                           TYPE       CONSTRAINT  CURRENT  WANTED  LATEST
test-registry/clean-code-ruleset  ruleset    ^1.0.0      1.0.0    1.1.0   1.2.0
```

**Expected Output (JSON):**
```json
[
  {
    "package": "test-registry/clean-code-ruleset",
    "type": "ruleset",
    "constraint": "^1.0.0",
    "current": "1.0.0",
    "wanted": "1.1.0",
    "latest": "1.2.0"
  }
]
```

**Expected Output (List):**
```
test-registry/clean-code-ruleset
```

**Verification:**
- Only outdated packages are shown
- Table format includes WANTED column (highest version satisfying constraint)
- JSON uses lowercase keys: "package", "type", "constraint", "current", "wanted", "latest"
- List format shows package names only (no dashes)
- Versions shown without 'v' prefix

### Example 3: View Dependency Info

**Command:**
```bash
arm info dependency test-registry/clean-code-ruleset
```

**Expected Output:**
```
test-registry/clean-code-ruleset:
    type: ruleset
    version: 1.0.0
    constraint: ^1.0.0
    priority: 100
    sinks:
        - cursor-rules
        - q-rules
    include:
        - "**/*.yml"
    exclude:
        - "**/experimental/**"
```

**Verification:**
- All configuration details are shown
- Data comes from both manifest and lock file
- 4-space indentation for nested levels
- Arrays use dash-prefixed items
- Include/exclude patterns displayed with quotes
- Version shown without 'v' prefix

### Example 4: List Available Versions

**Command:**
```bash
arm list versions test-registry/clean-code-ruleset
```

**Expected Output:**
```
test-registry/clean-code-ruleset:
  - 2.1.0
  - 2.0.0
  - 1.5.0
  - 1.0.0
  - main (branch)
  - develop (branch)
```

**Verification:**
- All available versions from registry are listed
- Semantic versions sorted descending (highest first)
- Branches listed after semver versions in config order
- 2-space indentation for version list
- Branch versions labeled with "(branch)"

### Example 5: No Dependencies

**Command:**
```bash
arm list dependency
```

**Expected Output:**
```
(empty output or no dependencies message)
```

**Verification:**
- Graceful handling of empty manifest
- No error, just empty list or informational message

## Notes

- Query operations are read-only and never modify manifest or lock files
- Outdated check requires network access to registries
- List versions requires network access to query registry
- List and info commands work offline using cached manifest/lock data
- Output is always sorted alphabetically for deterministic results (important for testing)
- JSON output format uses lowercase keys: "package", "type", "constraint", "current", "wanted", "latest"
- Table format includes WANTED column showing highest version satisfying constraint
- List format shows package names only (no dashes, no version)
- Versions displayed without 'v' prefix in all outputs
