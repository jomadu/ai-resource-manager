![ARM Header](assets/header.png)

# AI Rules Manager (ARM)

## What is ARM?

A package manager for AI rules that treats rulesets like code dependencies - with semantic versioning, reproducible installs, and automatic distribution to your AI tools.

Connect to Git repositories like awesome-cursorrules or your team's rule collections, install versioned rulesets across projects, and keep them automatically synced with their source of truth.

## Why ARM?

AI coding assistants like Cursor and Amazon Q rely on rules to guide their behavior, but managing these rules is broken:

- **Manual copying** severs the connection to the source of truth - once copied, rules are orphaned with no way to get updates
- **Breaking changes blindness** - when you pull latest rules, you have no idea if they'll break your AI's behavior
- **Doesn't scale** - managing rules across even 3 projects becomes unmanageable overhead

ARM treats AI rules like code dependencies - versioned, distributable packages that stay in sync across your entire development environment.

## Concepts

### Registries

Registries are remote sources where rulesets are stored and versioned, similar to npm registries for JavaScript packages. These are shared across your team and stored in `arm.json`. ARM supports Git-based registries that can point to GitHub repositories, GitLab projects, or any Git remote containing rule collections. When you configure a registry like `ai-rules`, you're creating a named connection to a repository that contains multiple rulesets with proper semantic versioning.

### Rulesets

Rulesets are collections of AI rules packaged as versioned units, identified by names like `ai-rules/amazonq-rules` where `ai-rules` is the registry and `amazonq-rules` is the ruleset name. These are shared across your team and tracked in `arm.json`. Each ruleset contains rule files (markdown, text, etc.) along with metadata defining version constraints, file patterns, and compatibility. Rulesets can be installed with specific version constraints and will automatically update according to semantic versioning rules.

### Sinks

Sinks define where installed rules should be placed in your local filesystem and which AI tools should receive them. These are personal configuration settings stored in `.armrc.json` on each developer's machine. Each sink targets specific directories (like `.amazonq/rules` or `.cursor/rules`) and can filter rulesets using include/exclude patterns. Sinks support two layout modes:

- **Hierarchical Layout** (default): Preserves directory structure from rulesets
- **Flat Layout**: Places all files in a single directory with hash-prefixed names for tools that require flat file structures

This allows you to automatically distribute the right rules to the right AI tools without manual file management.

## Installation

### Quick Install

```bash
curl -fsSL https://raw.githubusercontent.com/jomadu/ai-rules-manager/main/scripts/install.sh | bash
```

### Install Specific Version

```bash
curl -fsSL https://raw.githubusercontent.com/jomadu/ai-rules-manager/main/scripts/install.sh | bash -s v1.0.0
```

### Manual Installation

1. Download the latest release from [GitHub](https://github.com/jomadu/ai-rules-manager/releases)
2. Extract and move the binary to your PATH
3. Run `arm help` to verify installation

## Quick Start

Configure registry:
```bash
arm config registry add ai-rules https://github.com/jomadu/ai-rules-manager-sample-git-registry --type git
```

Configure sinks for different AI tools:
```bash
arm config sink add q --directories .amazonq/rules
```

```bash
arm config sink add cursor --directories .cursor/rules
```

```bash
arm config sink add copilot --directories .github/copilot --layout flat
```

Install rulesets:
```bash
arm install ai-rules/rules
```

## Key Commands

- `arm config registry` - Manage registries
- `arm config sink` - Manage sinks
- `arm install <ruleset>[@version]` - Install rulesets with semantic versioning
- `arm update [ruleset]` - Update to latest compatible versions
- `arm uninstall <ruleset>` - Remove rulesets
- `arm list` - Show installed rulesets
- `arm outdated` - Check for updates
- `arm info [ruleset]` - Show detailed information
- `arm cache clean` - Remove old cached versions

## Version Constraints

- `arm install ai-rules/rules@2.1.0` - Exact version (=2.1.0)
- `arm install ai-rules/rules@2.1` - Minor updates (~2.1.0)
- `arm install ai-rules/rules@2` - Major updates (^2.0.0)
- `arm install ai-rules/rules@main` - Track branch

## Layout Modes

### Hierarchical Layout (Default)

Preserves the original directory structure from rulesets. Files are organized as `sink-dir/arm/registry/ruleset/version/original-path`:

```
.cursor/rules/
└── arm/
    └── ai-rules/
        └── rules/
            └── 1.0.0/
                └── rules/
                    ├── grug-brained-dev.md
                    └── clean-code.md
```

### Flat Layout

Places all files in a single directory with hash-prefixed names. Each filename starts with an 8-character hash (derived from registry/ruleset@version:filepath) followed by the original path with slashes replaced by underscores:

```
.github/copilot/
├── 183791a9_rules_clean-code.md
├── 3554667c_rules_grug-brained-dev.md
└── arm-index.json
```

The `arm-index.json` file maps hashed filenames back to their original paths:

```json
{
  "183791a9_rules_clean-code.md": {
    "registry": "ai-rules",
    "ruleset": "rules",
    "version": "1.0.0",
    "filePath": "rules/clean-code.md"
  },
  "3554667c_rules_grug-brained-dev.md": {
    "registry": "ai-rules",
    "ruleset": "rules",
    "version": "1.0.0",
    "filePath": "rules/grug-brained-dev.md"
  }
}
```

Configure via CLI flags:
```bash
arm config sink add copilot --directories .github/copilot --layout flat
```

Or configure in `.armrc.json`:
```json
{
  "sinks": {
    "cursor": {
      "directories": [".cursor/rules"],
      "layout": "hierarchical"
    },
    "copilot": {
      "directories": [".github/copilot"],
      "layout": "flat"
    }
  }
}
```

## Files

- `.armrc.json` - Personal sink configuration (where rules are installed on your machine)
- `arm.json` - Team-shared project manifest with registries and dependencies
- `arm-lock.json` - Team-shared locked versions for reproducible installs
- `arm-index.json` - Local flat layout index (maps hashes to original file paths)

ARM follows npm-like patterns for predictable dependency management across AI development environments.
