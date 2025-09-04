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
      "type": "git"
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
        "url": "https://github.com/jomadu/ai-rules-manager-sample-git-registry",
        "type": "git",
        "constraint": "main",
        "resolved": "abc1234adfdafdfda12355434314...",
        "include": ["rules/amazonq/*.md"],
        "exclude": []
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
ai-rules/amazonq-rules@abc1234
```

**arm info (no args)**
```bash
./arm info
ai-rules/amazonq-rules
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
Ruleset: ai-rules/amazonq-rules
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

Git registries determine available versions based on repository content:

### Semver Tags Present
If the repository contains any semver tags (with or without `v` prefix), those tags are the **only** versions available from the registry.

**Example:** Repository has tags `v1.0.0`, `v1.1.0`, `v2.0.0`
- Available versions: `1.0.0`, `1.1.0`, `2.0.0`
- Branch constraints like `main` are **not available**
- Installing `ai-rules/ruleset@main` would fail with helpful error:

```bash
./arm install ai-rules/ruleset@main
Error: Version 'main' not found for ruleset 'ai-rules/ruleset'

This registry uses semver tags for versioning. Available versions:
  1.0.0, 1.1.0, 2.0.0

To install the latest version, use:
  arm install ai-rules/ruleset@2.0.0
  # or
  arm install ai-rules/ruleset  # installs latest (2.0.0)
```

### No Semver Tags
If the repository has no semver tags, only commits on branches configured by the registry are available as versions.

**Example:** Repository has no semver tags, registry configured with branches `["main", "develop"]`
- Available versions: `main`, `develop`
- Branch constraint must match one of the configured branches
- Installing `ai-rules/ruleset@feature-branch` would fail with helpful error:

```bash
./arm install ai-rules/ruleset@feature-branch
Error: Version 'feature-branch' not found for ruleset 'ai-rules/ruleset'

This registry uses branch-based versioning. Available versions:
  main, develop

To install from an available branch, use:
  arm install ai-rules/ruleset@main
  # or
  arm install ai-rules/ruleset@develop
```

### Registry Configuration

**Default Git Registry:**
```bash
arm config registry add ai-rules https://github.com/example/repo --type git
# Creates Git registry with default branches: ["main", "master"]
```

**Custom Git Branches:**
```bash
arm config registry add ai-rules https://github.com/example/repo --type git --branches main,develop,staging
# Creates Git registry with custom branches: ["main", "develop", "staging"]
```

**Other Registry Types:**
```bash
arm config registry add api-rules https://api.example.com --type http
# Creates HTTP registry (no Git-specific configuration)
```

**Resulting Configuration:**
```json
{
  "registries": {
    "ai-rules": {
      "url": "https://github.com/example/repo",
      "type": "git",
      "branches": ["main", "develop", "staging"]
    }
  }
}
```

**Default branches when no `--branches` flag is provided:** `["main", "master"]`

## Implementation Notes

- Use short commit hashes (7-8 characters) for display and directory names
- Store full commit hashes in lock file for precision
- Branch constraints still track latest commits on that branch
- When constraint is a branch, "latest" shows latest commit on that branch
- When constraint is semver, "latest" shows latest semver tag

## Implementation Decisions

### Breaking Changes
- **No migration support**: This is a breaking change that requires users to reinstall rulesets
- **Directory cleanup**: Old branch-named directories will be orphaned and must be manually removed
- **Lock file format**: Existing lock files with branch names will become invalid

### Technical Details
- **Semver detection**: Strict semver parsing only (1.0.0, v1.0.0 format)
- **Constraint parsing**: Use existing `BranchHead` constraint type for branch names
- **Update frequency**: Check for new commits on every ARM operation (no caching)
- **CLI display**: Show normalized constraint forms in tables (e.g., ">=1.0.0 <2.0.0" instead of "^1.0.0")

### Type-First Configuration
- **CLI parsing**: Parse `--type` flag first, then route to type-specific methods
- **Git registries**: Use `AddGitRegistry(name, url, branches)` for Git-specific configuration
- **Other registries**: Use `AddRegistry(name, url, type)` for generic registries
- **Factory validation**: Registry factory validates type and extracts appropriate configuration

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
