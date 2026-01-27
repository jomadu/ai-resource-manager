# Package Installation

## Job to be Done
Install, update, upgrade, and uninstall AI resource packages (rulesets and promptsets) from registries to local sinks with version tracking and integrity verification.

## Activities
1. Install package from registry to sink(s)
2. Update package within version constraints
3. Upgrade package to latest version (ignoring constraints)
4. Uninstall package from sink(s) and clean up empty directories

## Acceptance Criteria
- [x] Install ruleset/promptset from registry to one or more sinks
- [x] Track installation in arm.json (manifest) and arm-lock.json (lock file)
- [x] Verify package integrity using SHA256 hash on install
- [x] Update respects version constraints from manifest
- [x] Upgrade ignores constraints and fetches latest
- [x] Uninstall removes files from sinks and cleans empty directories
- [x] Uninstall removes arm-index.json when all packages removed
- [x] Uninstall removes arm_index.* when all rulesets removed
- [x] Support --priority flag for rulesets (default: 100)
- [x] Support --include and --exclude patterns for file filtering

## Data Structures

### Manifest Entry (arm.json)
```json
{
  "dependencies": {
    "registry/package": {
      "type": "ruleset",
      "version": "^1.0.0",
      "priority": 100,
      "sinks": ["cursor-rules", "q-rules"],
      "include": ["**/*.yml"],
      "exclude": ["**/experimental/**"]
    }
  }
}
```

### Lock Entry (arm-lock.json)
```json
{
  "dependencies": {
    "registry/package@v1.0.0": {
      "version": "v1.0.0",
      "resolved": "https://github.com/org/repo",
      "integrity": "sha256-abc123...",
      "commit": "abc123def456"
    }
  }
}
```

## Algorithm

### Install
1. Parse package key (registry/package@version)
2. Resolve version from registry (semver, branch, or tag)
3. Fetch package files from registry
4. Calculate SHA256 integrity hash
5. Store in cache (~/.arm/storage/)
6. Compile to tool-specific format for each sink
7. Write to sink directories
8. Update arm.json and arm-lock.json

### Update
1. Read manifest for version constraint
2. Resolve highest version satisfying constraint
3. If newer than locked version, install new version
4. Uninstall old version from sinks
5. Update lock file

### Upgrade
1. Fetch latest version from registry (ignore constraint)
2. Install new version
3. Uninstall old version from sinks
4. Update manifest constraint and lock file

### Uninstall
1. Remove files from each sink
2. Remove package from arm-index.json
3. If no packages remain, remove arm-index.json
4. If no rulesets remain, remove arm_index.* priority files
5. Clean up empty directories recursively
6. Update arm.json and arm-lock.json

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Package already installed | Reinstall (replace files) |
| Version not found | Error with available versions |
| Integrity mismatch | Error and refuse to install |
| No lock file on update | Treat as fresh install |
| Sink doesn't exist | Error (must add sink first) |
| Empty directory after uninstall | Remove directory recursively |
| Nested empty directories | Remove all empty ancestors |
| Sink root directory empty | Keep sink root, only remove subdirs |

## Dependencies

- Version resolution (version-resolution.md)
- Registry access (registry-management.md)
- Sink compilation (sink-compilation.md)
- Pattern filtering (pattern-filtering.md)
- Cache storage (cache-management.md)

## Implementation Mapping

**Source files:**
- `internal/arm/service/service.go` - InstallRuleset, InstallPromptset, UpdatePackages, UpgradePackages, UninstallPackages
- `internal/arm/service/dependency_test.go` - Unit tests for install workflows
- `internal/arm/sink/manager.go` - InstallRuleset, InstallPromptset, Uninstall, CleanupEmptyDirectories
- `cmd/arm/main.go` - CLI handlers for install, update, upgrade, uninstall
- `test/e2e/install_test.go` - E2E tests for installation
- `test/e2e/update_test.go` - E2E tests for update/upgrade
- `test/e2e/cleanup_test.go` - E2E tests for uninstall cleanup

## Examples

### Install with Priority
```bash
arm install ruleset --priority 200 ai-rules/team-standards cursor-rules
```

### Install with Patterns
```bash
arm install ruleset --include "security/**/*.yml" --exclude "**/experimental/**" ai-rules/security cursor-rules
```

### Install to Multiple Sinks
```bash
arm install ruleset ai-rules/clean-code cursor-rules q-rules copilot-rules
```

### Update Within Constraints
```bash
# Manifest has ^1.0.0, currently on 1.0.0
arm update
# Updates to 1.1.0 if available, but not 2.0.0
```

### Upgrade Ignoring Constraints
```bash
# Manifest has ^1.0.0, currently on 1.0.0
arm upgrade
# Upgrades to 2.0.0 if available, updates manifest to ^2.0.0
```

### Uninstall
```bash
arm uninstall ai-rules/clean-code
# Removes from all sinks, cleans empty directories, updates manifest/lock
```
