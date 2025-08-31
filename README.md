![ARM Header](assets/header.png)

# AI Rules Manager (ARM)

A package manager for AI coding rules that syncs rulesets across different AI tools like Cursor and Amazon Q.

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
