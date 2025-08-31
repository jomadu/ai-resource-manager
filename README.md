![ARM Header](assets/header.png)

# AI Rules Manager (ARM)

A package manager for AI coding rules that syncs rulesets across different AI tools like Cursor and Amazon Q.

## Quick Start

```bash
# Configure registry
arm config add registry ai-rules https://github.com/my-user/ai-rules --type git

# Configure sinks for different AI tools
arm config add sink q --directories .amazonq/rules --include ai-rules/amazonq-*
arm config add sink cursor --directories .cursor/rules --include ai-rules/cursor-*

# Install rulesets
arm install ai-rules/amazonq-rules --include rules/amazonq/*.md
arm install ai-rules/cursor-rules --include rules/cursor/*.mdc
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
