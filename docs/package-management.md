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

Compile resources from source files. The compile command accepts both individual files and directories as inputs:

**INPUT_PATH** accepts:
- **Files**: Directly processes the specified file(s)
- **Directories**: Discovers files within using `--include`/`--exclude` patterns
- **Mixed**: Can combine files and directories in the same command

**Note:** Shell glob patterns (e.g., `*.yml`) are expanded to individual files by your shell before ARM processes them.

**Examples:**
```bash
# Compile single file
$ arm compile --tool cursor ruleset.yml ./output/

# Compile multiple files
$ arm compile --tool cursor file1.yml file2.yml file3.yml ./output/

# Compile directory (non-recursive by default)
$ arm compile --tool cursor ./rulesets/ ./output/

# Compile directory recursively
$ arm compile --tool amazonq --recursive ./src/ ./build/

# Mix files and directories
$ arm compile --tool cursor specific-file.yml ./more-rulesets/ ./output/

# Use shell glob expansion (expands to individual files)
$ arm compile --tool cursor ./rulesets/*.yml ./output/

# Validate only (no output)
$ arm compile --validate-only ruleset.yml

# Compile with force overwrite
$ arm compile --tool copilot --force ruleset.yml ./output/

# Compile with include/exclude patterns
$ arm compile --tool cursor --include "**/*.yml" --exclude "**/README.md" ./src/ ./build/

# Validate and fail fast on first error (useful for CI)
$ arm compile --validate-only --fail-fast ./rulesets/
```

Available tools: `md`, `cursor`, `amazonq`, `copilot`

Compile options:
- `--recursive` - Process directories recursively
- `--validate-only` - Validate without generating output files
- `--force` - Overwrite existing files
- `--include <pattern>` - Include files matching pattern
- `--exclude <pattern>` - Exclude files matching pattern
- `--fail-fast` - Stop on first error
