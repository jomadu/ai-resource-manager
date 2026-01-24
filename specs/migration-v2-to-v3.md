# Migration Guide: ARM v2 to v3

## TL;DR

**Recommendation: Nuke and Pave**

Due to fundamental changes in file formats, command structure, and internal data models, we recommend a clean migration:

1. Document your current v2 configuration
2. Uninstall v2 completely
3. Install v3
4. Reconfigure from scratch using new commands

**Why?** File format incompatibilities between v2 and v3 make in-place upgrades unreliable. Starting fresh ensures a clean, working installation.

---

## What Changed in v3?

### 1. Command Structure Overhaul

The command hierarchy was completely restructured for better usability and consistency.

#### Registry Management

**v2:**
```bash
arm config registry add <name> <url> --type git
arm config registry remove <name>
```

**v3:**
```bash
arm add registry git --url <url> <name>
arm add registry gitlab --url <url> --project-id <id> <name>
arm add registry cloudsmith --owner <owner> --repo <repo> <name>
arm remove registry <name>
```

#### Sink Management

**v2:**
```bash
arm config sink add <name> --directories <path> --include <pattern>
arm config sink remove <name>
```

**v3:**
```bash
arm add sink --tool <tool> <name> <path>
arm remove sink <name>
```

#### Listing Configuration

**v2:**
```bash
arm config list
```

**v3:**
```bash
arm list                    # List all entities
arm list registry           # List registries only
arm list sink              # List sinks only
arm list dependency        # List dependencies only
```

#### Information Display

**v2:**
```bash
# No dedicated info command
```

**v3:**
```bash
arm info registry <name>
arm info sink <name>
arm info dependency <name>
```

### 2. Installation Command Changes

The install command now requires explicit resource types and sink targets.

**v2:**
```bash
# Implicit sink from config
arm install awesome-cursorrules/python --include "rules-new/python.mdc"
```

**v3:**
```bash
# Explicit resource type and sink
arm install ruleset ai-rules/clean-code-ruleset cursor-rules
arm install promptset ai-rules/code-review-promptset cursor-commands

# Multiple sinks
arm install ruleset ai-rules/clean-code-ruleset cursor-rules q-rules

# With priority
arm install ruleset --priority 200 ai-rules/team-standards cursor-rules

# With file filtering
arm install ruleset --include "**/typescript-*.yml" ai-rules/language-rules cursor-rules
```

### 3. Terminology Changes

**v2: AI Rules Manager**
- Focused exclusively on "rules" and "rulesets"
- Single resource type

**v3: AI Resource Manager**
- Expanded to multiple resource types:
  - **Rulesets** - Collections of AI rules with priority-based conflict resolution
  - **Promptsets** - Collections of AI prompts for reusable templates
- More flexible and extensible architecture

### 4. Resource Format Changes

**v2: URF (Universal Rule Format)**
```yaml
version: "1.0"
metadata:
  name: example-rule
  description: Example rule
  author: user
rules:
  - content: Rule content
```

**v3: ARM Resource Format**
```yaml
# Version field removed (managed by registry)
metadata:
  name: example-rule
  description: Example rule
  author: user
  # More optional fields
rules:
  - content: Rule content
    priority: 100
```

**Key Changes:**
- JSON keys changed from `snake_case` to `camelCase`
- `version` field removed from resource schema (managed externally)
- Additional optional fields for metadata
- Enhanced priority system for conflict resolution

### 5. Registry Types Expanded

**v2:**
- Git registries only

**v3:**
- **Git Registry** - GitHub, GitLab, and Git remotes
- **GitLab Registry** - GitLab Package Registry integration
- **Cloudsmith Registry** - Cloudsmith package repository

### 6. Sink Configuration

**v2:**
- Generic sink configuration
- Include/exclude patterns at sink level

**v3:**
- Tool-specific sinks (`--tool cursor`, `--tool amazonq`, `--tool copilot`, `--tool markdown`)
- Include/exclude patterns at install time
- Better separation of concerns

### 7. New Commands

v3 introduces several new commands:

```bash
arm compile              # Compile rulesets/promptsets to tool formats
arm outdated            # Check for outdated packages
arm upgrade             # Upgrade to latest versions
arm set registry        # Configure registry settings
arm set sink           # Configure sink settings
arm set ruleset        # Configure ruleset settings
arm set promptset      # Configure promptset settings
arm clean cache        # Clean cache with --nuke or --max-age
arm clean sinks        # Clean sinks with --nuke
```

### 8. Internal Changes

**File Formats:**
- Manifest format changed (`arm.json`)
- Lock file format changed (`arm-lock.json`)
- Index format changed (`arm-index.json`)
- Cache structure updated

**Architecture:**
- Index manager redesigned
- Cache system overhauled
- Lockfile manager updated
- Manifest manager refactored
- Service layer restructured

---

## Migration Steps

### Step 1: Document Current Configuration

Before uninstalling v2, document your setup:

```bash
# Save your current configuration
arm config list > arm-v2-config.txt

# Note your installed rulesets
arm list > arm-v2-installed.txt
```

### Step 2: Uninstall v2

```bash
curl -fsSL https://raw.githubusercontent.com/jomadu/ai-resource-manager/main/scripts/uninstall.sh | bash
```

Or manually:
```bash
rm -f $(which arm)
rm -rf ~/.arm
```

### Step 3: Clean Project Files

In each project using ARM:

```bash
# Remove v2 files
rm -f arm.json arm-lock.json arm-index.json
rm -rf .arm/

# Optionally clean sink directories
rm -rf .cursor/rules .amazonq/rules .github/copilot
```

### Step 4: Install v3

```bash
curl -fsSL https://raw.githubusercontent.com/jomadu/ai-resource-manager/main/scripts/install.sh | bash
```

Or install specific version:
```bash
curl -fsSL https://raw.githubusercontent.com/jomadu/ai-resource-manager/main/scripts/install.sh | bash -s v3.1.2
```

### Step 5: Reconfigure Registries

Using your saved v2 configuration, recreate registries with new commands:

**v2 Config:**
```
awesome-cursorrules: https://github.com/PatrickJS/awesome-cursorrules (git)
```

**v3 Commands:**
```bash
arm add registry git --url https://github.com/PatrickJS/awesome-cursorrules awesome-cursorrules
```

For GitLab or Cloudsmith registries (new in v3):
```bash
arm add registry gitlab --url https://gitlab.example.com --project-id 123 my-gitlab
arm add registry cloudsmith --owner myorg --repo ai-rules my-cloudsmith
```

### Step 6: Reconfigure Sinks

**v2 Config:**
```
cursor:
  directories: [.cursor/rules]
  include: [awesome-cursorrules/python]
```

**v3 Commands:**
```bash
arm add sink --tool cursor cursor-rules .cursor/rules
arm add sink --tool cursor cursor-commands .cursor/commands
arm add sink --tool amazonq q-rules .amazonq/rules
arm add sink --tool amazonq q-prompts .amazonq/prompts
```

### Step 7: Reinstall Dependencies

**v2 Command:**
```bash
arm install awesome-cursorrules/python --include "rules-new/python.mdc"
```

**v3 Command:**
```bash
# Determine if it's a ruleset or promptset
arm install ruleset awesome-cursorrules/python cursor-rules --include "rules-new/python.mdc"

# Or for prompts
arm install promptset awesome-cursorrules/python-prompts cursor-commands
```

### Step 8: Verify Installation

```bash
arm version
arm list
arm list registry
arm list sink
arm list dependency
```

---

## Common Migration Scenarios

### Scenario 1: Simple Git Registry with Cursor

**v2:**
```bash
arm config registry add ai-rules https://github.com/user/rules --type git
arm config sink add cursor --directories .cursor/rules
arm install ai-rules/python-rules
```

**v3:**
```bash
arm add registry git --url https://github.com/user/rules ai-rules
arm add sink --tool cursor cursor-rules .cursor/rules
arm install ruleset ai-rules/python-rules cursor-rules
```

### Scenario 2: Multiple Sinks with Filtering

**v2:**
```bash
arm config sink add cursor --directories .cursor/rules --include "ai-rules/cursor-*"
arm config sink add q --directories .amazonq/rules --include "ai-rules/amazonq-*"
arm install ai-rules/cursor-python
arm install ai-rules/amazonq-python
```

**v3:**
```bash
arm add sink --tool cursor cursor-rules .cursor/rules
arm add sink --tool amazonq q-rules .amazonq/rules
arm install ruleset ai-rules/cursor-python cursor-rules
arm install ruleset ai-rules/amazonq-python q-rules

# Or install to both sinks at once
arm install ruleset ai-rules/python-rules cursor-rules q-rules
```

### Scenario 3: Team Standards with Priority

**v2:**
```bash
# No priority support in v2
arm install team-rules/standards
arm install community-rules/best-practices
```

**v3:**
```bash
# Team rules take precedence with higher priority
arm install ruleset --priority 200 team-rules/standards cursor-rules
arm install ruleset --priority 100 community-rules/best-practices cursor-rules
```

---

## Why "Nuke and Pave"?

### File Format Incompatibilities

1. **Manifest Format** - `arm.json` structure changed fundamentally
2. **Lock File Format** - `arm-lock.json` uses different schema
3. **Index Format** - `arm-index.json` tracking mechanism redesigned
4. **Cache Structure** - Internal cache organization completely different

### Data Model Changes

1. **JSON Keys** - Changed from `snake_case` to `camelCase`
2. **Resource Schema** - Version field removed, optional fields added
3. **Dependency Tracking** - New relationship model between resources
4. **Sink Configuration** - Tool-specific configuration added

### Command Structure

1. **Command Hierarchy** - Completely restructured
2. **Flag Names** - Many flags renamed or removed
3. **Argument Order** - Different positional argument requirements

### Risk of Corruption

Attempting to upgrade in-place could result in:
- Corrupted manifest files
- Broken dependency resolution
- Inconsistent cache state
- Failed installations
- Unpredictable behavior

**Clean installation ensures:**
- ✅ Consistent file formats
- ✅ Proper dependency tracking
- ✅ Reliable cache state
- ✅ Predictable behavior

---

## Troubleshooting

### Issue: "Registry not found"

**Cause:** Registry names or URLs changed

**Solution:**
```bash
arm list registry
arm add registry git --url <correct-url> <name>
```

### Issue: "Sink not configured"

**Cause:** Sinks require explicit tool specification in v3

**Solution:**
```bash
arm add sink --tool cursor cursor-rules .cursor/rules
```

### Issue: "Invalid resource type"

**Cause:** Must specify `ruleset` or `promptset` in v3

**Solution:**
```bash
# Determine resource type from registry
arm info dependency <name>

# Install with correct type
arm install ruleset <registry>/<name> <sink>
arm install promptset <registry>/<name> <sink>
```

### Issue: "Command not found"

**Cause:** Command structure changed

**Solution:** Refer to command mapping in this guide or run:
```bash
arm help
arm help install
arm help add
```

---

## Getting Help

- **Documentation:** See `specs/` directory for detailed specifications
- **Commands:** Run `arm help` or `arm help <command>`
- **Issues:** Report bugs at https://github.com/jomadu/ai-resource-manager/issues

---

## Summary

v3 represents a significant evolution of ARM with:
- ✅ Better command structure
- ✅ Multiple resource types (rulesets + promptsets)
- ✅ More registry options (Git, GitLab, Cloudsmith)
- ✅ Enhanced priority system
- ✅ Tool-specific sink configuration
- ✅ Improved dependency management

While migration requires starting fresh, the improved architecture and features make v3 worth the effort.
