# ARM Integration: Universal Rule Format Support

## Overview

This document describes how ARM integrates Universal Rule Format (URF) files, including UX changes, command line modifications, and the new ruleset prioritization system.

## Detection and Processing

### Auto-Detection

ARM automatically detects URF files during registry scanning:

```
registry/
├── ruleset.yaml          # URF format - detected by .yaml extension + version field
├── rules/
│   ├── rule1.md          # Legacy format - individual markdown files
│   └── rule2.md
└── README.md
```

**Detection Logic**:
1. Look for `*.yaml` or `*.yml` files in ruleset root
2. Parse and validate `version` field matches URF spec
3. Fall back to legacy format if no valid URF files found

### URF Format Example

ARM processes URF files with this structure:

```yaml
version: "1.0"
metadata:
  id: "ruleset-id"
  name: "Ruleset Name"
  version: "1.0.0"
  description: "Description"
rules:
  - id: "critical-security-check"
    name: "Critical Security Check"
    description: "Validate all user inputs to prevent security vulnerabilities"
    priority: 100
    enforcement: "must"
    scope:
      - files: ["**/*.ext"]
    body: |
      Rule content in markdown format.
  - id: "recommended-best-practice"
    name: "Recommended Best Practice"
    description: "Follow established coding conventions for maintainability"
    priority: 80
    enforcement: "should"
    scope:
      - files: ["**/*.ext"]
    body: |
      Rule content in markdown format.
  - id: "optional-optimization"
    name: "Optional Optimization"
    description: "Consider performance improvements when feasible"
    priority: 60
    enforcement: "may"
    scope:
      - files: ["**/*.ext"]
    body: |
      Rule content in markdown format.
```

### Validation

ARM validates URF files during installation:

- YAML syntax validation
- Required field presence (`version`, `metadata`, `rules`)
- Basic structure validation
- Rule ID uniqueness within ruleset

## Command Line Changes

### Enhanced Install Command

```bash
# Install with ruleset priority
arm install ai-rules/ruleset --sinks cursor --priority 100

# Install multiple rulesets with different priorities
arm install ai-rules/security-rules --sinks cursor,q --priority 200
arm install ai-rules/style-rules --sinks cursor --priority 50
```

**New `--priority` Flag**:
- **Type**: Integer (1-1000+)
- **Behavior**: Overrides existing priority if ruleset already installed
- **Default**: 100 if not specified
- **Validation**: Must be positive integer

### New Config Commands

```bash
# Update ruleset configuration (triggers reinstall)
arm config ruleset update ai-rules/ruleset priority 150
arm config ruleset update ai-rules/ruleset version 1.1.0
arm config ruleset update ai-rules/ruleset sinks cursor,q
arm config ruleset update ai-rules/ruleset include "**/*.ext,**/*.py"
arm config ruleset update ai-rules/ruleset exclude "**/test/**,**/node_modules/**"

# Sink configuration with compilation target
arm config sink add cursor .cursor/rules --compile-to cursor
arm config sink add q .amazonq/rules --compile-to amazonq
```

**New Subcommands**:
- `arm config ruleset update <name> <key> <value>` - Update ruleset config (triggers reinstall)

### Enhanced List Command

```bash
# Show installed rulesets (default: alphanumeric sort)
arm list

# Show with priorities (alphanumeric sort)
arm list --show-priority

# Sort by priority (highest first)
arm list --show-priority --sort-priority

# Example output (alphanumeric):
# ai-rules/ruleset@1.0.0 (priority: 100)
# ai-rules/security-rules@1.5.0 (priority: 200)

# Example output (priority sort):
# ai-rules/security-rules@1.5.0 (priority: 200)
# ai-rules/ruleset@1.0.0 (priority: 100)
```

**New Flags**:
- `--show-priority` - Display ruleset installation priorities
- `--sort-priority` - Sort by priority (highest first) instead of alphanumeric

## Configuration Changes

### arm.json Schema Updates

```json
{
  "registries": { ... },
  "rulesets": {
    "ai-rules": {
      "python-rules": {
        "version": "2.1.0",
        "priority": 100,
        "include": ["**/*"],
        "exclude": [],
        "sinks": ["cursor", "q"]
      },
      "security-rules": {
        "version": "1.5.0",
        "priority": 200,
        "sinks": ["cursor", "q"]
      }
    }
  },
  "sinks": {
    "cursor": {
      "directory": ".cursor/rules",
      "layout": "hierarchical",
      "compileTarget": "cursor"
    },
    "q": {
      "directory": ".amazonq/rules",
      "layout": "hierarchical",
      "compileTarget": "amazonq"
    }
  }
}
```

**New Fields**:
- `rulesets[].priority` - Ruleset installation priority
- `sinks[].compileTarget` - Target format for compilation

## Prioritization System

### Priority Resolution

Rulesets are processed in priority order (higher number = higher priority):

1. **Ruleset Priority** (higher number = higher priority)
2. **Installation Order** (for rulesets with same priority)

### Example Priority Resolution

Given:
- `security-rules` (priority: 200)
- `python-rules` (priority: 100)

**Priority Order**:
```
1. security-rules (priority: 200)
2. python-rules (priority: 100)
```

### Priority Management UX

```bash
# Set priority during install
arm install ai-rules/ruleset --priority 200

# Change priority after install (triggers reinstall)
arm config ruleset update ai-rules/ruleset priority 250

# View current priorities
arm list --show-priority
```

## Compilation Process

### Compilation Triggers

ARM compiles URF files to tool-specific formats when:

1. **Installation**: New rulesets are installed
2. **Update**: Existing rulesets are updated to new versions
3. **Priority Change**: Ruleset priorities are modified

### Compilation Output

**Cursor Format** (`<sink-dir>/arm/<registry>/<ruleset>/<version>/ruleset-id_critical-security-check.mdc`):
```markdown
---
description: "Validate all user inputs to prevent security vulnerabilities"
globs: ["**/*.ext"]
alwaysApply: true
---

---
namespace: <registry>/<ruleset>@<version>
ruleset:
  id: ruleset-id
  name: Ruleset Name
  version: 1.0.0
  rules:
    - critical-security-check
    - recommended-best-practice
    - optional-optimization
rule:
  id: critical-security-check
  name: Critical Security Check
  enforcement: MUST
  priority: 100
  scope:
    - files: "**/*.ext"
---

# Critical Security Check (MUST)

Rule content in markdown format.
```

**Amazon Q Format** (`<sink-dir>/arm/<registry>/<ruleset>/<version>/ruleset-id_critical-security-check.md`):
```markdown
---
namespace: <registry>/<ruleset>@<version>
ruleset:
  id: ruleset-id
  name: Ruleset Name
  version: 1.0.0
  rules:
    - critical-security-check
    - recommended-best-practice
    - optional-optimization
rule:
  id: critical-security-check
  name: Critical Security Check
  enforcement: MUST
  priority: 100
  scope:
    - files: "**/*.ext"
---

# Critical Security Check (MUST)

Rule content in markdown format.
```

## Error Handling

### Validation Errors

```bash
# Invalid URF file
$ arm install ai-rules/broken-rules
Error: Invalid URF format in ai-rules/broken-rules/ruleset.yaml
  - missing required field 'metadata'
  - missing required field 'rules'
```

### Compilation Errors

```bash
# Missing compilation target
$ arm install ai-rules/ruleset --sinks cursor
Error: Sink 'cursor' missing compileTarget configuration
Available targets: cursor, amazonq
Run: arm config sink cursor --compile-to cursor
```

## Migration Path

### Legacy Format Support

ARM continues to support legacy markdown files:

```bash
# Mixed format support
registry/
├── ruleset.yaml          # URF format
├── legacy-rules/
│   ├── rule1.md          # Legacy format
│   └── rule2.md
```

**Behavior**:
- URF rules are compiled to tool-specific formats
- Legacy rules are copied as-is to target directories
- Both formats can coexist in the same project

## Performance Considerations

### Source Caching

- ARM caches URF source files from registries
- Compilation happens at install time for each sink's target format
- Cache invalidation occurs when ruleset versions change

### Incremental Updates

- Ruleset changes trigger recompilation for affected sinks
- Priority changes trigger reinstall and recompilation
- Each sink compiles independently based on its target format
