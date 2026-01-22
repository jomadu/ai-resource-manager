# Explicit Sink Targeting

## Overview

Replace ARM's automatic ruleset distribution with explicit sink targeting, requiring users to specify which AI tools should receive each ruleset during installation.

## Problem

Current ARM automatically distributes rulesets to all configured sinks, leading to:
- Unwanted rules appearing in AI tools that shouldn't have them
- No control over which tools get which rulesets
- Difficulty managing multi-tool environments

## Solution

Make sink targeting explicit and required during installation:

```bash
# Configure sinks for different AI tools
arm config sink add cursor .cursor/rules --layout hierarchical
arm config sink add q .amazonq/rules --layout hierarchical

# Explicitly target specific sinks during installation
arm install awesome-cursorrules/python --include "*.mdc" --sinks cursor,q
```

## Key Features

### Required Sink Targeting
- `--sinks [sink...]` becomes mandatory for `arm install`
- No automatic distribution to all sinks
- Explicit control over ruleset placement

### Simplified Sink Configuration
- Single directory per sink (not array)
- Remove include/exclude patterns from sink config
- File filtering moves to ruleset level via `--include`/`--exclude`

### Clean State Management
- Reinstalling a ruleset removes it from previous sink locations
- Only installs to newly specified sinks
- Prevents stale rule accumulation
- Removing a sink automatically removes it from all ruleset configurations

## CLI Changes

### Sink Configuration
```bash
# Before: Multiple directories, include/exclude patterns
arm config sink add cursor --directories .cursor/rules,.cursor/backup --include "*.md"

# After: Single directory, layout only
arm config sink add cursor .cursor/rules --layout hierarchical
```

### Installation
```bash
# Before: Automatic distribution to all sinks
arm install awesome-cursorrules/python

# After: Explicit sink targeting required
arm install awesome-cursorrules/python --sinks cursor,q
```

## Configuration Changes

### SinkConfig Simplification
```go
// Before
type SinkConfig struct {
    Directories []string `json:"directories"`
    Include     []string `json:"include"`
    Exclude     []string `json:"exclude"`
    Layout      string   `json:"layout,omitempty"`
}

// After
type SinkConfig struct {
    Directory string `json:"directory"`
    Layout    string `json:"layout,omitempty"`
}
```

### Ruleset Entry Enhancement
```go
type Entry struct {
    Version string   `json:"version"`
    Include []string `json:"include"`
    Exclude []string `json:"exclude"`
    Sinks   []string `json:"sinks"`  // New required field
}
```

### Example arm.json
```json
{
  "registries": {
    "awesome-cursorrules": {
      "url": "https://github.com/PatrickJS/awesome-cursorrules",
      "type": "git"
    }
  },
  "rulesets": {
    "awesome-cursorrules": {
      "python": {
        "version": "latest",
        "include": ["rules-new/python.mdc"],
        "exclude": [],
        "sinks": ["cursor"]
      }
    }
  },
  "sinks": {
    "cursor": {
      "directory": ".cursor/rules",
      "layout": "hierarchical"
    },
    "q": {
      "directory": ".amazonq/rules",
      "layout": "hierarchical"
    }
  }
}
```

## Benefits

- **Precise Control**: Install specific rulesets to specific AI tools
- **Clean State**: Reinstallation removes stale files from previous locations
- **Simplified Config**: Single directory per sink, no complex filtering
- **Multi-Tool Support**: Easy management of different AI tool environments
- **Explicit Intent**: Clear understanding of which tools get which rules

## Sink Removal Behavior

When a sink is removed, it's automatically cleaned from all ruleset configurations:

```bash
# Current state: python ruleset installed to both sinks
arm list
# awesome-cursorrules/python@latest (sinks: cursor,q)

# Remove cursor sink
arm config sink remove cursor

# Automatically updates all rulesets
arm list
# awesome-cursorrules/python@latest (sinks: q)
```

The removal process:
1. Removes files from the sink's directory
2. Updates all ruleset entries to remove the sink from their `sinks` arrays
3. Maintains consistency between sink configuration and ruleset targeting

## Migration Strategy

Due to the complexity of automatically migrating existing configurations, users will need to manually migrate their ARM setups. A migration guide will be provided with:

1. **Backup existing configuration**: Save current `arm.json` and installed files
2. **Clean slate approach**: Remove existing ARM installation and start fresh
3. **Reconfigure sinks**: Use new simplified sink configuration syntax
4. **Reinstall rulesets**: Explicitly specify target sinks for each ruleset
5. **Verify installation**: Confirm rules are in correct AI tool directories

This manual approach ensures users understand the new explicit targeting model and can make informed decisions about which rulesets go to which AI tools.

## Migration Impact

This is a breaking change requiring:
1. Update to manifest Entry type (add Sinks field)
2. Simplify SinkConfig type (remove directories array, include/exclude)
3. Make --sinks required for install command
4. Implement clean reinstallation behavior
5. Implement sink removal cleanup across all rulesets
6. Update CLI to use positional directory argument for sink add
7. Create comprehensive migration guide for manual user migration
