# Compiling Resources

Compile resources from source files. The compile command accepts both individual files and directories as inputs:

**INPUT_PATH** accepts:
- **Files**: Directly processes the specified file(s)
- **Directories**: Discovers files within using `--include`/`--exclude` patterns
- **Mixed**: Can combine files and directories in the same command

**Note:** Shell glob patterns (e.g., `*.yml`) are expanded to individual files by your shell before ARM processes them.

## Examples

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

## Available Tools

- `md` - Markdown output
- `cursor` - Cursor IDE format
- `amazonq` - Amazon Q format
- `copilot` - GitHub Copilot format

## Options

- `--recursive` - Process directories recursively
- `--validate-only` - Validate without generating output files
- `--force` - Overwrite existing files
- `--include <pattern>` - Include files matching pattern
- `--exclude <pattern>` - Exclude files matching pattern
- `--fail-fast` - Stop on first error