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
┌──────────┬───────────────┬─────────────┬──────────────┬──────────────┐
│ REGISTRY │    RULESET    │   CURRENT   │    WANTED    │    LATEST    │
├──────────┼───────────────┼─────────────┼──────────────┼──────────────┤
│ ai-rules │ amazonq-rules │ main:abc1234│ main:ddedf22 │ main:ddedf22 │
└──────────┴───────────────┴─────────────┴──────────────┴──────────────┘
```
*Note: Current behavior shows current=main, wanted=main, and latest=v2.1.0*

**arm list**
```bash
./arm list
ai-rules/amazonq-rules@main:abc1234
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
  Constraint: main | Resolved: main:abc1234
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
Resolved: main:abc1234
```

## Implementation Notes

- Use short commit hashes (7-8 characters) for display and directory names
- Store full commit hashes in lock file for precision
- Branch constraints still track latest commits on that branch
- When constraint is a branch, "latest" shows latest commit on that branch
- When constraint is semver, "latest" shows latest semver tag
