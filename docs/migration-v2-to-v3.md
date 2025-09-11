# Migration Guide: ARM v2 to v3

## Overview

ARM v3 introduces **explicit sink targeting**, replacing automatic ruleset distribution with precise control over which AI tools receive which rulesets. This is a breaking change that requires manual migration.

## What Changed

### Before (v2): Automatic Distribution
```bash
# Configure sinks with include/exclude patterns
arm config sink add cursor --directories .cursor/rules --include "ai-rules/cursor-*"
arm config sink add q --directories .amazonq/rules --include "ai-rules/amazonq-*"

# Install automatically distributes to matching sinks
arm install ai-rules/cursor-rules
```

### After (v3): Explicit Targeting
```bash
# Configure sinks with single directory
arm config sink add cursor .cursor/rules
arm config sink add q .amazonq/rules

# Explicitly specify target sinks
arm install ai-rules/cursor-rules --sinks cursor
arm install ai-rules/shared-rules --sinks cursor,q
```

## Key Changes

1. **Required `--sinks` flag**: All new installations must specify target sinks
2. **Simplified sink config**: Single directory per sink, no include/exclude patterns
3. **Clean reinstallation**: Reinstalling removes files from previous sink locations
4. **Automatic cleanup**: Removing a sink cleans it from all ruleset configurations

## Migration Steps

### 1. Backup Current Setup
```bash
# Backup your current configuration
cp arm.json arm.json.v2.backup
cp arm-lock.json arm-lock.json.v2.backup

# Note your current installed rulesets
arm list > installed-rulesets.txt
```

### 2. Clean Slate Approach
```bash
# Remove old ARM installation
rm -rf .cursor/rules/arm .amazonq/rules/arm .github/copilot/arm

# Remove old configuration (keep backups)
rm arm.json arm-lock.json
```

### 3. Upgrade ARM Binary
```bash
# Install ARM v3
curl -fsSL https://raw.githubusercontent.com/jomadu/ai-rules-manager/main/scripts/install.sh | bash
```

### 4. Reconfigure Registries
```bash
# Re-add your registries (same as before)
arm config registry add ai-rules https://github.com/your-org/rules-repo --type git
```

### 5. Reconfigure Sinks (New Syntax)
```bash
# Old v2 syntax (don't use):
# arm config sink add cursor --directories .cursor/rules --include "ai-rules/cursor-*"

# New v3 syntax:
arm config sink add cursor .cursor/rules --layout hierarchical
arm config sink add q .amazonq/rules --layout hierarchical
arm config sink add copilot .github/copilot --layout flat
```

### 6. Reinstall Rulesets with Explicit Targeting
```bash
# Review your backup to see what was installed
cat installed-rulesets.txt

# Reinstall with explicit sink targeting
arm install ai-rules/cursor-rules --sinks cursor
arm install ai-rules/amazonq-rules --sinks q
arm install ai-rules/shared-rules --sinks cursor,q
```

### 7. Verify Installation
```bash
# Check that rulesets are installed to correct sinks
arm list

# Verify files are in expected locations
ls -la .cursor/rules/arm/
ls -la .amazonq/rules/arm/
```

## Configuration Format Changes

### v2 arm.json
```json
{
  "registries": { ... },
  "rulesets": {
    "ai-rules": {
      "cursor-rules": {
        "version": "latest",
        "include": ["**/*"],
        "exclude": []
      }
    }
  },
  "sinks": {
    "cursor": {
      "directories": [".cursor/rules"],
      "include": ["ai-rules/cursor-*"],
      "exclude": [],
      "layout": "hierarchical"
    }
  }
}
```

### v3 arm.json
```json
{
  "registries": { ... },
  "rulesets": {
    "ai-rules": {
      "cursor-rules": {
        "version": "latest",
        "include": ["**/*"],
        "exclude": [],
        "sinks": ["cursor"]
      }
    }
  },
  "sinks": {
    "cursor": {
      "directory": ".cursor/rules",
      "layout": "hierarchical"
    }
  }
}
```

## Common Migration Scenarios

### Scenario 1: Tool-Specific Rules
**v2 Setup:**
```bash
arm config sink add cursor --directories .cursor/rules --include "*/cursor-*"
arm config sink add q --directories .amazonq/rules --include "*/amazonq-*"
```

**v3 Migration:**
```bash
arm config sink add cursor .cursor/rules
arm config sink add q .amazonq/rules
arm install ai-rules/cursor-rules --sinks cursor
arm install ai-rules/amazonq-rules --sinks q
```

### Scenario 2: Shared Rules Across Tools
**v2 Setup:**
```bash
arm config sink add cursor --directories .cursor/rules --include "*/shared-*,*/cursor-*"
arm config sink add q --directories .amazonq/rules --include "*/shared-*,*/amazonq-*"
```

**v3 Migration:**
```bash
arm config sink add cursor .cursor/rules
arm config sink add q .amazonq/rules
arm install ai-rules/shared-rules --sinks cursor,q
arm install ai-rules/cursor-rules --sinks cursor
arm install ai-rules/amazonq-rules --sinks q
```

### Scenario 3: Multiple Directories (No Longer Supported)
**v2 Setup:**
```bash
arm config sink add cursor --directories .cursor/rules,.cursor/backup
```

**v3 Migration:**
```bash
# Choose primary directory
arm config sink add cursor .cursor/rules

# Create separate sink for backup if needed
arm config sink add cursor-backup .cursor/backup
arm install ai-rules/rules --sinks cursor,cursor-backup
```

## Benefits of v3

- **Precise Control**: Explicitly choose which tools get which rules
- **Clean State**: Reinstallation removes stale files from previous locations
- **Simplified Config**: No complex include/exclude pattern matching
- **Multi-Tool Support**: Easy management of different AI tool environments
- **Explicit Intent**: Clear understanding of which tools get which rules

## Troubleshooting

### Error: `--sinks is required for installing rulesets`
This is expected in v3. All new installations must specify target sinks:
```bash
arm install ai-rules/rules --sinks cursor,q
```

### Error: `sink X not configured`
Configure the sink first:
```bash
arm config sink add X /path/to/directory
```

### Files in wrong locations after migration
Remove and reinstall with correct sinks:
```bash
arm uninstall ai-rules/rules
arm install ai-rules/rules --sinks correct-sink
```

### Old v2 configuration not working
v3 cannot automatically migrate v2 configurations due to the fundamental change from pattern-based to explicit targeting. Follow the manual migration steps above.

## Need Help?

- Check `arm help` for updated command syntax
- Use `arm config list` to verify your sink configuration
- Use `arm list` to see which sinks each ruleset is installed to
- Refer to the updated README for current examples
