# Lock File Simplification Design

## Overview

Simplify `arm-lock.json` by removing redundant registry and ruleset configuration that already exists in `arm.json`. The lock file should only contain resolved versions and integrity checksums.

## Current State

### arm.json
```json
{
    "registries": {
        "ai-rules": {
            "url": "https://github.com/my-user/ai-rules",
            "type": "git"
        }
    },
    "rulesets": {
        "ai-rules": {
            "amazonq-rules": {
                "version": "^2.0.0",
                "include": ["rules/amazonq/*.md"]
            },
            "dev-rules": {
                "version": "main",
                "include": ["rules/dev/*.md"]
            }
        }
    }
}
```

### arm-lock.json (Current)
```json
{
    "rulesets": {
        "ai-rules": {
            "amazonq-rules": {
                "url": "https://github.com/my-user/ai-rules",
                "type": "git",
                "constraint": "^2.0.0",
                "resolved": "2.1.0",
                "include": ["rules/amazonq/*.md"]
            }
        }
    }
}
```

## Target State

### arm-lock.json (New)
```json
{
    "rulesets": {
        "ai-rules": {
            "amazonq-rules": {
                "resolved": "2.1.0",
                "checksum": "sha256:abc123def456789..."
            },
            "dev-rules": {
                "resolved": "abc123f",
                "checksum": "sha256:def789abc123456..."
            }
        }
    }
}
```

## Design Decisions

### Lock File Dependencies
- Lock file depends on `arm.json` for registry and constraint information
- Both files must exist together for reproducible installs
- Reduces lock file size and eliminates duplication

### Integrity Verification
- Use SHA-256 checksums of resolved ruleset files
- Checksum covers only the files that get installed locally
- Provides integrity verification across all registry types

### Structure Consistency
- Lock file mirrors `arm.json` nested structure
- Maintains registry/ruleset hierarchy for easy cross-referencing
- Consistent with npm-style dependency management

## Implementation Requirements

### Checksum Generation
- Calculate SHA-256 of concatenated file contents in deterministic order
- Include file paths in checksum calculation for structure integrity
- Store as `sha256:` prefixed hex string

### Resolution Process
1. Read constraints from `arm.json`
2. Resolve versions from registries
3. Download and checksum resolved files
4. Write minimal lock entry with resolved version and checksum

### Verification Process
1. Read resolved versions from lock file
2. Look up registry configuration from `arm.json`
3. Download resolved files and verify checksums
4. Install if checksums match, error if mismatch

## Breaking Changes

### Migration Strategy
- This is a breaking change requiring manual migration
- Users must regenerate lock files with new ARM version
- Clear error message when old lock file format detected

### Error Handling
```
Error: Incompatible lock file format detected.
Please run 'arm install' to regenerate arm-lock.json with the new format.
```

## Benefits

- **Reduced Duplication**: Eliminates redundant configuration storage
- **Single Source of Truth**: Registry and constraint config only in `arm.json`
- **Smaller Lock Files**: Significant size reduction for projects with many rulesets
- **Integrity Verification**: SHA-256 checksums ensure file integrity
- **Registry Agnostic**: Checksum approach works with any registry type

## File Size Impact

For a project with 10 rulesets:
- **Current**: ~2KB per ruleset (registry URL, type, constraint, include patterns)
- **New**: ~100 bytes per ruleset (version + checksum)
- **Savings**: ~95% reduction in lock file size
