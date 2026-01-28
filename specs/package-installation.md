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
- [ ] UpdateAll continues on error for partial success (BUG: returns on first error)
- [x] Upgrade ignores constraints and fetches latest
- [ ] UpgradeAll continues on error for partial success (BUG: returns on first error)
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
    "registry/clean-code-ruleset": {
      "type": "ruleset",
      "version": "^1.0.0",
      "priority": 100,
      "sinks": ["cursor-rules", "q-rules"],
      "include": ["**/*.yml"],
      "exclude": ["**/experimental/**"]
    },
    "registry/code-review-promptset": {
      "type": "promptset",
      "version": "^2.0.0",
      "sinks": ["cursor-commands"]
    }
  }
}
```

### Lock Entry (arm-lock.json)
```json
{
  "version": 1,
  "dependencies": {
    "registry/package@v1.0.0": {
      "integrity": "sha256-abc123..."
    }
  }
}
```

**Note:** Lock file uses composite key format: `registry/package@version`

## Data Structures

### Sink Path Structure
```
{sink}/arm/{registry}/{package}/{version}/{file}

Example:
.cursor/rules/arm/ai-rules/clean-code-ruleset/1.0.0/rules/cleanCode_ruleOne.mdc
```

### ARM Index (arm-index.json)
```json
{
  "version": 1,
  "rulesets": {
    "registry/package@v1.0.0": {
      "priority": 100,
      "files": ["arm/registry/package/v1.0.0/rules/rule.mdc"]
    }
  },
  "promptsets": {
    "registry/package@v2.0.0": {
      "files": ["arm/registry/package/v2.0.0/prompts/prompt.md"]
    }
  }
}
```

**Note:** Tracks installed packages per sink using composite keys

## Algorithm

### Install
1. Validate sinks exist in manifest
2. Resolve version from registry (semver, branch, or tag)
3. Fetch package files from registry (cached in ~/.arm/storage/)
4. Calculate SHA256 integrity hash
5. Verify integrity if package already locked (prevents tampering)
6. Update manifest (version constraint, sinks, patterns, priority)
7. Update lock file (resolved version, integrity) using composite key `registry/package@version`
8. For each sink, call `sinkMgr.InstallRuleset()` or `sinkMgr.InstallPromptset()` which:
   - Uninstalls all existing versions of package
   - Parses ARM resource files and compiles to tool format
   - Copies non-resource files directly
   - Writes to hierarchical path: `{sink}/arm/{registry}/{package}/{version}/{file}`
   - Updates arm-index.json with installed files and priority
   - Regenerates arm_index.* priority file (rulesets only)

### Update
1. For each dependency in manifest:
   - Read version constraint from manifest
   - Resolve highest version satisfying constraint
   - Skip if version unchanged
   - If version changed:
     - Uninstall old version from sinks
     - Remove old lock entry
     - Install new version to sinks
     - Update lock file with new version and integrity
   - Manifest constraint unchanged
2. Continue on error for partial success (BUG: UpdateAll returns on first error)

### Upgrade
1. For each dependency in manifest:
   - Fetch latest version from registry (ignore constraint)
   - Skip if version unchanged
   - If version changed:
     - Uninstall old version from sinks
     - Remove old lock entry
     - Install new version to sinks
     - Update lock file with new version and integrity
     - Update manifest constraint to `^{major}.0.0`
2. Continue on error for partial success (BUG: UpgradeAll returns on first error)

### Uninstall
1. For each sink:
   - Load arm-index.json
   - Find all versions matching `registry/package@*` prefix
   - Remove all files from disk
   - Delete entries from index
   - If no packages remain, remove arm-index.json and arm_index.*
   - Otherwise, save index and regenerate arm_index.*
   - Clean up empty directories recursively
2. Remove from lock file
3. Remove from manifest

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Package already installed | Uninstall old versions, then install new |
| Version not found | Error with available versions |
| Integrity mismatch on install | Error and refuse to install (package tampered) |
| No lock file on update | Treat as fresh install |
| Sink doesn't exist | Error (must add sink first) |
| Empty directory after uninstall | Remove directory recursively |
| Nested empty directories | Remove all empty ancestors |
| Sink root directory empty | Keep sink root, only remove subdirs |
| Update/upgrade with errors | Continue processing remaining packages (BUG: UpdateAll/UpgradeAll return on first error) |
| Uninstall package | Removes all versions matching `registry/package@*` |

## Known Bugs

### Bug: UpdateAll/UpgradeAll Don't Continue on Error
**Files:** `internal/arm/service/service.go` (UpdateAll, UpgradeAll)  
**Issue:** Return on first error instead of continuing for partial success  
**Expected:** Continue processing remaining packages, return error only if all fail  
**Note:** UpdatePackages and UpgradePackages correctly implement partial success

## Dependencies

- Version resolution (version-resolution.md)
- Registry access (registry-management.md)
- Sink compilation (sink-compilation.md)
- Pattern filtering (pattern-filtering.md)
- Cache storage (cache-management.md)

## Implementation Mapping

**Source files:**
- `internal/arm/service/service.go` - InstallRuleset, InstallPromptset, UpdatePackages (partial success ✓), UpdateAll (BUG: fail-fast), UpgradePackages (partial success ✓), UpgradeAll (BUG: fail-fast), UninstallPackages
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
