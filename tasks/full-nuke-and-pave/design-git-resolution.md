# Git Branch Resolution Design

## Problem Statement

Currently, ARM resolves Git branches to their branch names (e.g., "main") rather than specific commit hashes. This creates several issues:

1. **No reproducible builds** - "main" can point to different commits over time
2. **No change detection** - Can't detect when upstream has new commits
3. **Unclear versioning** - Directory structure uses branch names instead of commit hashes

## Current Behavior Analysis

When installing `ai-rules/amazonq-rules@main`, ARM currently:
- Creates directory: `.amazonq/rules/arm/ai-rules/amazonq-rules/main/`
- Stores `"resolved": "main"` in lock file
- Shows `main` in all CLI outputs

### Current Test Output

```bash
./arm install ai-rules/amazonq-rules@main --include "rules/amazonq/*.md"
# Creates: .amazonq/rules/arm/ai-rules/amazonq-rules/main/
```

**Current arm-lock.json:**
```json
{
  "rulesets": {
    "ai-rules": {
      "amazonq-rules": {
        "constraint": "main",
        "resolved": "main"
      }
    }
  }
}
```

## Desired Behavior

ARM should resolve branches to specific commit hashes for reproducible, trackable installations.

### Directory Structure

**Current:**
```
.amazonq/rules/arm/ai-rules/amazonq-rules/main/
```

**Desired:**
```
.amazonq/rules/arm/ai-rules/amazonq-rules/abc1234/
```

### Configuration Files

**arm.json** (unchanged)
```json
{
  "registries": {
    "ai-rules": {
      "url": "https://github.com/jomadu/ai-rules-manager-sample-git-registry",
      "type": "git",
      "branches" [
        "main"
      ]
    }
  },
  "rulesets": {
    "ai-rules": {
      "amazonq-rules": {
        "version": "main",
        "include": ["rules/amazonq/*.md"]
      }
    }
  }
}
```

**arm-lock.json** (key changes)
```json
{
  "rulesets": {
    "ai-rules": {
      "amazonq-rules": {
        "resolved": "abc1234adfdafdfda12355434314...",
        "checksum": "sha256:bbbccdd4412566234..."
      }
    }
  }
}
```

### CLI Output Changes

**arm outdated**
```bash
./arm outdated
┌──────────┬───────────────┬────────────┬─────────────┬──────────────┬──────────────┐
│ REGISTRY │    RULESET    │ CONSTRAINT │   CURRENT   │    WANTED    │    LATEST    │
├──────────┼───────────────┼────────────┼─────────────┼──────────────┼──────────────┤
│ ai-rules │ amazonq-rules │ main       │ abc1234     │ ddedf22      │ ddedf22      │
└──────────┴───────────────┴────────────┴─────────────┴──────────────┴──────────────┘
```
*Note: Current behavior shows current=main, wanted=main, and latest=v2.1.0*

**arm list**
```bash
./arm list
ai-rules/amazonq-rules@v1.2.0 (^1.0.0)
ai-rules/amazonq-rules@abc1234 (main)
```

**arm info (no args)**
```bash
./arm info
ai-rules/amazonq-rules@abc1234 (main)
  Registry: ()
  include:
    - rules/amazonq/*.md
  Installed:
    - .amazonq/rules/arm/ai-rules/amazonq-rules/abc1234
  Sinks:
    - q
  Constraint: main | Resolved: abc1234
```

**arm info (specific ruleset)**
```bash
./arm info ai-rules/amazonq-rules
Ruleset: ai-rules/amazonq-rules@abc1234 (main)
Registry: https://github.com/jomadu/ai-rules-manager-sample-git-registry (git)
include:
  - rules/amazonq/*.md
Installed:
  - .amazonq/rules/arm/ai-rules/amazonq-rules/abc1234
Sinks:
  - q
Constraint: main
Resolved: abc1234
```

## Version Listing Behavior

Git registries provide both semver tags and branch commits as available versions, with tags taking priority.

### Mixed Version Sources
Registries return all available versions in priority order:
1. **Semver tags** (sorted by semantic version, latest first)
2. **Branch HEAD commits** (in registry configuration order)

**Example:** Repository has tags `v1.0.0`, `v1.1.0`, `v2.0.0` and configured branches `["main", "develop"]`
- Available versions: `v2.0.0`, `v1.1.0`, `v1.0.0`, `abc1234` (main), `def5678` (develop)
- Installing `ai-rules/ruleset` (no version) selects `v2.0.0` (highest priority)
- Both tag and branch constraints are valid

### Version Resolution Priority
When no version is specified, ARM selects the first version from the priority-ordered list:
- **Tags present:** Latest semver tag is selected
- **No tags:** Latest commit from first configured branch is selected

### Error Messages
When a version isn't found, ARM shows categorized available versions:

```bash
./arm install ai-rules/ruleset@feature-branch
Error: Version 'feature-branch' not found for ruleset 'ai-rules/ruleset'

Available versions:
  Tags: v2.0.0, v1.1.0, v1.0.0
  Branches: main, develop

To install a specific version:
  arm install ai-rules/ruleset@v2.0.0
  arm install ai-rules/ruleset@main
```

## Implementation Notes

- Use short commit hashes (8 characters) for display and directory names
- Store full commit hashes in lock file for precision
- Branch constraints still track latest commits on that branch
- When constraint is a branch, "latest" shows latest commit on that branch
- When constraint is semver, "latest" shows latest semver tag

## Implementation Decisions

### Resolution Strategy
- **Eager resolution**: Resolve branch names to commit hashes immediately during `install`
- **Always update**: `arm update` on branch constraints fetches latest commit from tracked branch
- **Mixed version support**: Both semver tags and branch commits available simultaneously
- **Priority ordering**: Tags first (semver sorted), then branch commits (config order)
- **Fail fast**: All network operations fail immediately on connection issues

### Display Format
- **Consistent constraint display**: Always show `@abc123 (main)` format in all CLI output
- **8-character commit hashes**: Use for display and directory names
- **Full commit hashes**: Store in lock file for precision

### Breaking Changes
- **No migration support**: This is a breaking change that requires users to reinstall rulesets
- **Directory cleanup**: Old branch-named directories will be orphaned (no documentation)
- **Lock file validation**: Fail with generic error "Invalid lock file format. Please reinstall rulesets."

### Error Handling
- **Git error parsing**: Parse network, authentication, and missing ref errors into friendly messages
- **No concurrent protection**: Let filesystem/Git handle concurrent access

### Technical Details
- **Version listing**: Return all semver tags plus HEAD commits from configured branches
- **Priority ordering**: Tags first (semver sorted desc), then branches (config order)
- **Session caching**: Cache version lists within single ARM command execution
- **Constraint storage**: Store original user constraint in arm.json only
- **Lock file format**: Contains only resolved version and checksum
- **CLI display**: Show original constraint in parentheses (e.g., "@v1.2.0 (^1.0.0)", "@abc1234 (main)")

### Simplified Data Structures
- **ResolvedVersion**: Flattened from nested `VersionRef` to simple `Version` string field
- **Configuration storage**: Git-specific fields (like `branches`) stored directly in base config, not nested

## Architecture Notes

- **Keep registry-specific logic isolated**: Avoid polluting high-level structs like `service.go` with Git registry implementation details
- Git-specific version resolution, branch handling, and semver tag detection should be encapsulated within the Git registry implementation
- High-level services should work with abstract version concepts, letting each registry type handle its own versioning semantics

### ResolvedVersion Struct

Introduce a `ResolvedVersion` struct that combines constraint and resolved version information for clean API design:

```go
type ResolvedVersion struct {
    Constraint Constraint // Original constraint struct
    Version    string     // Resolved version (e.g., "abc1234", "1.2.0")
}
```

**Usage:**
- `ResolveVersion()` returns `ResolvedVersion` with simplified structure
- Encapsulates both the original user intent (constraint) and the concrete resolution (version string)
- Clean, flat structure without unnecessary nesting
