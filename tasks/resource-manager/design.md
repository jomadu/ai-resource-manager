# AI Resource Manager (ARM) - Resource Extension Design

## Overview

Extend AI Rules Manager to AI Resource Manager, supporting both rulesets and promptsets as first-class resources while maintaining the same package management approach.

## Command Changes

### Current vs New Commands

| Current | New |
|---------|-----|
| `arm install <ruleset>` | `arm install ruleset <ruleset>` |
| N/A | `arm install promptset <promptset>` |
| `arm info [ruleset]` | `arm info [ruleset\|promptset <name>]` |
| N/A | `arm info promptset <promptset>` |
| `arm list` | `arm list` (shows both, separate sections) |
| `arm outdated` | `arm outdated` (shows both, separate sections) |
| `arm uninstall <ruleset>` | `arm uninstall ruleset <ruleset>` |
| N/A | `arm uninstall promptset <promptset>` |
| `arm update [ruleset]` | `arm update [ruleset\|promptset <name>]` |

### Command Specifications

#### Install (Breaking Change)
```bash
arm install ruleset my-reg/my-ruleset --sinks cursor-rules --include "**/*.yml"
arm install promptset my-reg/my-promptset --sinks cursor-prompts --include "**/*.yml"
```

#### Info (Enhanced)
```bash
arm info                            # shows all installed resources
arm info ruleset my-reg/my-ruleset  # shows specific ruleset
arm info promptset my-reg/my-promptset # shows specific promptset
```

#### List (Enhanced)
```bash
arm list
```
Displays all installed resources organized into separate sections for rulesets and promptsets. Each resource shows its registry/name and installed version. May include additional metadata like installation date, sink assignments, or status indicators.

#### Outdated (Enhanced)
```bash
arm outdated
```
Shows resources with available updates, grouped by type. Displays current and available versions with clear upgrade paths. May include release notes summaries or breaking change warnings.

When called without arguments, displays summary information for all installed resources grouped by type. When called with a specific resource, provides detailed information including metadata, version history, dependencies, sink assignments, file contents, and registry information. May show compilation status for URF resources.

#### Uninstall (Breaking Change)
```bash
arm uninstall ruleset my-ruleset
arm uninstall promptset my-promptset
```

#### Update (Enhanced)
```bash
arm update                          # updates all outdated resources
arm update ruleset my-ruleset       # updates specific ruleset
arm update promptset my-promptset   # updates specific promptset
```

## File Structure Changes

### Configuration Files

#### arm.json (Enhanced)
```json
{
  "version": "3.0.0",
  "registries": { ... },
  "sinks": { ... },
  "rulesets": {
    "my-reg/clean-code": {
      "version": "^1.0.0",
      "sinks": ["cursor-rules"],
      "priority": 100
    }
  },
  "promptsets": {
    "my-reg/code-review": {
      "version": "^1.2.0",
      "sinks": ["cursor-prompts"]
    }
  }
}
```

#### arm-lock.json (Enhanced)
```json
{
  "rulesets": {
    "my-reg": {
      "clean-code": {
        "version": "v1.0.0",
        "display": "v1.0.0",
        "checksum": "sha256:..."
      }
    }
  },
  "promptsets": {
    "my-reg": {
      "code-review": {
        "version": "v1.2.0",
        "display": "v1.2.0",
        "checksum": "sha256:..."
      }
    }
  }
}
```

### Directory Structure

```
.cursor/
├── rules/           # Rulesets only
│   └── arm/
│       ├── my-reg/
│       ├── arm_index.mdc
│       └── arm-index.json
└── commands/        # Promptsets only
    └── arm/
        ├── my-reg/
        └── arm-index.json  # No priority ordering
```

## URF Extension for Prompts

### Prompt URF Format
```yaml
version: "1.0"
metadata:
  id: "code-review-prompts"
  name: "Code Review Prompt Set"
  description: "Prompts for code review tasks"
prompts:
  security-review:
    name: "Security Review"
    description: "Review code for security vulnerabilities"
    body: |
      Review this code for security vulnerabilities.
      Focus on input validation, authentication, and data exposure.
  performance-review:
    name: "Performance Review"
    description: "Review code for performance issues"
    body: |
      Analyze this code for performance bottlenecks.
      Look for inefficient algorithms and resource usage.
```

### Compilation Targets

#### Cursor (.md)
```markdown
---
namespace: my-reg
metadata:
  id: code-review-prompts
  name: Code Review Prompt Set
  version: 1.0.0
  prompts:
    - security-review
    - performance-review
prompt:
  id: security-review
  name: Security Review
---

# Security Review

Review this code for security vulnerabilities.
Focus on input validation, authentication, and data exposure.
```

#### Amazon Q (.md)
```markdown
---
namespace: my-reg
metadata:
  id: code-review-prompts
  name: Code Review Prompt Set
  version: 1.0.0
  prompts:
    - security-review
    - performance-review
prompt:
  id: security-review
  name: Security Review
---

# Security Review

Review this code for security vulnerabilities.
Focus on input validation, authentication, and data exposure.
```

## Implementation Plan

### Phase 1: Ruleset Migration
1. Update command parsing to require `ruleset` keyword
2. Update configuration file structure
3. Maintain backward compatibility warnings
4. Update documentation

### Phase 2: Promptset Support
1. Add promptset command handling
2. Extend URF compiler for prompts
3. Add promptset installation logic
4. Update list/outdated commands for dual display
5. Update metadata block compilation to use `metadata` section header instead of `ruleset`
6. Refactor cache system for unified resource support:
   - Rename `rulesets/` directory to `resources/` with type subdirectories
   - Generalize `RegistryIndex` to handle resource types (rulesets, promptsets)
   - Create generic `RegistryResourceCache` interface replacing `RegistryRulesetCache`
   - Update cache key generation to include resource type

### Phase 3: Polish
1. Enhanced error messages
2. Migration tooling
3. Updated examples and documentation

## Breaking Changes

- All install/info/uninstall commands now require resource type
- Configuration file format version bump to 3.0.0
- Old command syntax will show migration hints

## Backward Compatibility

- Provide clear error messages with migration examples
- Include migration guide in documentation
- Consider temporary alias support during transition period
