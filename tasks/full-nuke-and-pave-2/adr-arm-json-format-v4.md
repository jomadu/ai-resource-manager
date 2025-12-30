# ADR: ARM Manifest Format v4

## Status
Proposed

## Context
The current arm.json format has several usability and clarity issues:
- `packages` section mixes rulesets and promptsets, making them hard to distinguish
- `resourceType` field is redundant when structure could indicate type
- Sink configuration exposes too much internal complexity (`compileTarget`, `layout`)
- No support for selective package installs (include/exclude patterns)
- Lock file format needs simplification and consistency

## Decision
Restructure both manifest and lock file formats to improve clarity and hide complexity:

### 1. Manifest Format (arm.json)

**Complete new format:**
```json
{
  "version": 1,
  "registries": {
    "sample-registry": {
      "type": "cloudsmith",
      "url": "https://api.cloudsmith.io",
      "owner": "sample-org",
      "repository": "arm-registry"
    }
  },
  "sinks": {
    "cursor-rules": {
      "directory": ".cursor/rules",
      "tool": "cursor"
    },
    "amazonq-rules": {
      "directory": ".amazonq/rules",
      "tool": "amazonq"
    }
  },
  "dependencies": {
    "rulesets": {
      "sample-registry/clean-code": {
        "version": "^1.0.0",
        "priority": 100,
        "sinks": ["cursor-rules", "amazonq-rules"],
        "include": ["**/*.yml"],
        "exclude": ["**/experimental/**"]
      }
    },
    "promptsets": {
      "sample-registry/prompts": {
        "version": "^2.0.0",
        "sinks": ["cursor-commands"]
      }
    }
  }
}
```

**Key changes:**
- `packages` → `dependencies` with separate `rulesets` and `promptsets` sections
- Remove `resourceType` field (structure indicates type)
- Simplify sinks: `compileTarget` + `layout` → just `tool` field
- Add optional `include`/`exclude` patterns for selective installs
- Integer `version` field for migration support

### 2. Lock File Format (arm-lock.json)

**New lock file format:**
```json
{
  "version": 1,
  "rulesets": {
    "sample-registry/clean-code@1.2.3": {
      "integrity": "sha256-abc123..."
    },
    "git-registry/experimental@abc123def": {
      "integrity": "sha256-def456..."
    }
  },
  "promptsets": {
    "sample-registry/prompts@2.1.0": {
      "integrity": "sha256-789xyz..."
    }
  }
}
```

**Lock file design principles:**
- Use `registry/package@version` as key (semver or commit hash)
- Separate sections for rulesets and promptsets (matches manifest)
- Store only essential data: integrity hash for verification
- No `resolved` URLs (can reconstruct from manifest + registry config)
- No file lists (ARM manages file tracking internally)
- Integer `version` field for format evolution

### 3. Version Strategy

**Manifest version:** Start at 1 (new v4 format)
**Lock file version:** Start at 1 (new simplified format)
**Format:** Integer versions for easy comparison and migration scripts

## Rationale

### Manifest Changes
1. **Clear separation**: Rulesets and promptsets in separate sections eliminates confusion
2. **Simplified sinks**: Developers specify intent (`tool`), ARM handles implementation details
3. **Selective installs**: Include/exclude patterns support partial package installs from git registries
4. **Reduced redundancy**: Structure indicates resource type, no need for explicit field

### Lock File Simplification
1. **Minimal data**: Only store what's essential for reproducible installs
2. **Consistent structure**: Matches manifest organization (rulesets/promptsets)
3. **Simple keys**: `registry/package@version` format is clear and parseable
4. **No URL duplication**: Registry config in manifest provides source information
5. **Integrity focus**: Hash verification without unnecessary metadata

### Comparison with npm
- **Similar**: Dependencies structure, semver constraints, separate lock file
- **Simpler**: No nested dependency trees, shorter lock entries
- **ARM-specific**: Registry configuration, sink targeting, resource-specific metadata

## Consequences

### Positive
- Clear separation between rulesets and promptsets
- Simpler sink configuration - developers specify intent, not implementation
- Support for partial package installs
- Reduced redundancy and file size
- More intuitive structure following npm patterns
- Simplified lock file reduces complexity
- Integer versions enable easy migration scripting

### Negative
- Breaking change requiring migration from v3
- Need to update all tooling and documentation
- Existing arm.json and arm-lock.json files become invalid

## Implementation Notes
- ARM will auto-detect layout and compile target from `tool` field
- Registry `type` field remains unchanged (different purpose than sink `tool`)
- Include/exclude patterns use glob syntax
- Lock file keys use exact resolved versions (semver or git commit hash)
- Migration tool needed for existing projects
- `arm migrate` command should handle v3 → v1 conversion