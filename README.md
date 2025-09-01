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

Registries are remote sources where rulesets are stored and versioned, similar to npm registries for JavaScript packages. ARM supports Git-based registries that can point to GitHub repositories, GitLab projects, or any Git remote containing rule collections. When you configure a registry like `ai-rules`, you're creating a named connection to a repository that contains multiple rulesets with proper semantic versioning.

### Sinks

Sinks define where installed rules should be placed in your local filesystem and which AI tools should receive them. Each sink targets specific directories (like `.amazonq/rules` or `.cursor/rules`) and can filter rulesets using include/exclude patterns. This allows you to automatically distribute the right rules to the right AI tools without manual file management.

### Rulesets

Rulesets are collections of AI rules packaged as versioned units, identified by names like `ai-rules/amazonq-rules` where `ai-rules` is the registry and `amazonq-rules` is the ruleset name. Each ruleset contains rule files (markdown, text, etc.) along with metadata defining version constraints, file patterns, and compatibility. Rulesets can be installed with specific version constraints and will automatically update according to semantic versioning rules.

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

```bash
# Configure registry
arm config add registry ai-rules https://github.com/jomadu/ai-rules-manager-sample-git-registry --type git

# Configure sinks for different AI tools
arm config add sink q --directories .amazonq/rules --include "ai-rules/amazonq-*"
arm config add sink cursor --directories .cursor/rules --include "ai-rules/cursor-*"

# Install rulesets
arm install ai-rules/amazonq-rules --include "rules/amazonq/*.md"
arm install ai-rules/cursor-rules --include "rules/cursor/*.mdc"
```

## Key Commands

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

## Files

- `.armrc.json` - Registry and sink configuration
- `arm.json` - Project manifest with dependencies
- `arm-lock.json` - Locked versions for reproducible installs

ARM follows npm-like patterns for predictable dependency management across AI development environments.
