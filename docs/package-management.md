# Package Management

## Installation

Install all configured packages:
```bash
arm install
```

Install specific rulesets:
```bash
arm install ruleset ai-rules/clean-code-ruleset cursor-rules
arm install ruleset ai-rules/security-ruleset cursor-rules q-rules
```

Install specific promptsets:
```bash
arm install promptset ai-rules/code-review-promptset cursor-commands
arm install promptset ai-rules/testing-promptset cursor-commands q-prompts
```

## Updates

Update all packages:
```bash
arm update
```

Update all rulesets:
```bash
arm update ruleset
```

Update specific rulesets:
```bash
arm update ruleset ai-rules/clean-code-ruleset ai-rules/security-ruleset
```

Update all promptsets:
```bash
arm update promptset
```

Update specific promptsets:
```bash
arm update promptset ai-rules/code-review-promptset
```

## Upgrades

Upgrade all packages (ignoring version constraints):
```bash
arm upgrade
```

Upgrade specific rulesets:
```bash
arm upgrade ruleset ai-rules/clean-code-ruleset
```

Upgrade specific promptsets:
```bash
arm upgrade promptset ai-rules/code-review-promptset
```

## Outdated Packages

Check for outdated packages:
```bash
arm outdated
arm outdated --output json
arm outdated --output list
```

Check for outdated rulesets:
```bash
arm outdated ruleset
arm outdated --output json ruleset
```

Check for outdated promptsets:
```bash
arm outdated promptset
arm outdated --output list promptset
```

## Listing Packages

List all installed packages:
```bash
arm list
```

List installed rulesets:
```bash
arm list ruleset
```

List installed promptsets:
```bash
arm list promptset
```

## Package Information

Show detailed information:
```bash
arm info
arm info ruleset
arm info promptset
```

Show specific package information:
```bash
arm info ruleset ai-rules/clean-code-ruleset
arm info promptset ai-rules/code-review-promptset
```

## Utilities

### Clean Cache

Clean the local cache directory:
```bash
arm clean cache
```

Aggressive cleanup (remove all cached data):
```bash
arm clean cache --nuke
```

### Clean Sinks

Clean sink directories based on ARM index:
```bash
arm clean sinks
```

Complete cleanup (remove entire ARM directory):
```bash
arm clean sinks --nuke
```

### Compile Resources

Compile resources from source files:
```bash
arm compile --target cursor ruleset.yml ./output/
arm compile --target amazonq --recursive ./src/ ./build/
arm compile --validate-only ruleset.yml
arm compile --target copilot --force ruleset.yml ./output/
```

Available targets: `md`, `cursor`, `amazonq`, `copilot`

Compile options:
- `--recursive` - Process directories recursively
- `--validate-only` - Validate without generating output files
- `--force` - Overwrite existing files
- `--include <pattern>` - Include files matching pattern
- `--fail-fast` - Stop on first error
