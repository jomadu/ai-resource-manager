# Commands

## Table of Contents

- [Commands](#commands)
  - [Table of Contents](#table-of-contents)
  - [Core](#core)
    - [arm version](#arm-version)
    - [arm help](#arm-help)
  - [Registry Management](#registry-management)
    - [arm add registry git](#arm-add-registry-git)
    - [arm add registry gitlab](#arm-add-registry-gitlab)
    - [arm add registry cloudsmith](#arm-add-registry-cloudsmith)
    - [arm remove registry](#arm-remove-registry)
    - [arm set registry](#arm-set-registry)
    - [arm list registry](#arm-list-registry)
    - [arm info registry](#arm-info-registry)
  - [Sink Management](#sink-management)
    - [arm add sink](#arm-add-sink)
    - [arm remove sink](#arm-remove-sink)
    - [arm set sink](#arm-set-sink)
    - [arm list sink](#arm-list-sink)
    - [arm info sink](#arm-info-sink)
  - [Package Management](#package-management)
    - [arm install](#arm-install)
    - [arm outdated](#arm-outdated)
    - [arm update](#arm-update)
    - [arm upgrade](#arm-upgrade)
    - [arm uninstall](#arm-uninstall)
    - [arm list](#arm-list)
    - [arm info](#arm-info)
  - [Resource Management](#resource-management)
    - [arm install ruleset](#arm-install-ruleset)
    - [arm set ruleset](#arm-set-ruleset)
    - [arm install promptset](#arm-install-promptset)
    - [arm set promptset](#arm-set-promptset)
  - [Utilities](#utilities)
    - [arm clean cache](#arm-clean-cache)
    - [arm clean sinks](#arm-clean-sinks)
    - [arm compile](#arm-compile)

## Core

### arm version

Display the current version, build information, and build datetime of the AI Rules Manager tool. This command shows detailed version information including version number (e.g., v1.2.3), build identifier/hash, and build timestamp showing when the binary was compiled.

This information is useful for:
- Verifying which version is installed
- Debugging compatibility issues with specific builds
- Checking if updates are available
- Reporting issues with precise version context
- Understanding when the binary was built (useful for troubleshooting time-sensitive issues)

**Example:**
```bash
$ arm version
version: v1.2.3
build-id: a1b2c3d4
build-timestamp: 2024-01-15T10:30:45Z
build-platform: darwin/arm64
```

### arm help

Display comprehensive help information and usage instructions for the AI Resource Manager. This command shows available commands and their syntax, command descriptions and usage examples, global flags and options, quick reference for common operations, and links to documentation and examples.

The help system provides contextual information based on what the user is trying to do, making it easy to discover commands and understand their proper usage.

**Examples:**
```bash
# Show main help
$ arm help
$ arm --help

# Show help for a specific command
$ arm help install
$ arm install --help

# Show help for subcommands
$ arm help add registry git
$ arm add registry git --help
```

## Registry Management

### arm add registry git

`arm add registry git --url URL [--branches BRANCH...] [--force] NAME`

Add a new Git registry to the ARM configuration. Git registries use Git repositories (GitHub, GitLab, or any Git remote) to store and distribute rulesets and promptsets using Git tags and branches for versioning.

**Examples:**
```bash
# Add a Git registry
$ arm add registry git --url https://github.com/my-org/arm-registry my-org

# Add a Git registry with specific branches
$ arm add registry git --url https://github.com/my-org/arm-registry --branches main,develop my-org

# Overwrite an existing registry
$ arm add registry git --url https://github.com/my-org/new-arm-registry --force my-org
```

### arm add registry gitlab

`arm add registry gitlab [--url URL] [--group-id ID] [--project-id ID] [--api-version VERSION] [--force] NAME`

Add a new GitLab registry to the ARM configuration. GitLab registries use GitLab's Generic Package Registry for versioned packages. The URL defaults to `https://gitlab.com` if not specified. You must specify either `--group-id` or `--project-id`.

**Examples:**
```bash
# Add a GitLab registry with project ID (using default gitlab.com)
$ arm add registry gitlab --project-id 456 my-gitlab-project

# Add a self-hosted GitLab registry with group ID
$ arm add registry gitlab --url https://gitlab.example.com --group-id 123 my-gitlab

# Add with custom API version
$ arm add registry gitlab --url https://gitlab.example.com --project-id 456 --api-version v4 my-gitlab-project

# Overwrite an existing registry
$ arm add registry gitlab --url https://gitlab.example.com --group-id 123 --force my-gitlab
```

### arm add registry cloudsmith

`arm add registry cloudsmith [--url URL] --owner OWNER --repo REPO [--force] NAME`

Add a new Cloudsmith registry to the ARM configuration. Cloudsmith registries use Cloudsmith's raw package repository service for single-file artifacts. The URL defaults to `https://api.cloudsmith.io` if not specified.

**Examples:**
```bash
# Add a Cloudsmith registry (using default API URL)
$ arm add registry cloudsmith --owner my-org --repo my-repo cloudsmith-registry

# Add a self-hosted Cloudsmith instance
$ arm add registry cloudsmith --url https://cloudsmith.mycompany.com --owner my-org --repo my-repo private-registry

# Overwrite an existing registry
$ arm add registry cloudsmith --owner my-org --repo my-repo --force cloudsmith-registry
```

### arm remove registry

`arm remove registry NAME`

Remove a registry from the ARM configuration by name. This command removes the specified registry and all its associated configuration from the ARM setup. After removal, the registry will no longer be available for installing rulesets or promptsets.

**Example:**
```bash
$ arm remove registry my-org
```

### arm set registry

`arm set registry NAME KEY VALUE`

Set configuration values for a specific registry. This command allows you to configure registry-specific settings such as URL endpoints or other registry-specific parameters. The exact configuration keys available depend on the registry type.

**Examples:**
```bash
# Update registry URL
$ arm set registry my-org url https://github.com/my-org/new-arm-registry

# Set GitLab group ID
$ arm set registry my-gitlab group_id 789

# Set GitLab project ID
$ arm set registry my-gitlab-project project_id 101

# Set Cloudsmith owner
$ arm set registry cloudsmith-registry owner new-org

# Set Cloudsmith repository
$ arm set registry cloudsmith-registry repository new-repo
```

### arm list registry

`arm list registry`

List all configured registries. This command displays all registries that have been added to the ARM configuration as a simple list.

**Example:**

```bash
$ arm list registry
- my-org
- my-gitlab
- my-gitlab-project
- cloudsmith-registry
```

### arm info registry

`arm info registry [NAME]...`

Display detailed information about one or more registries. This command shows comprehensive details about the specified registries, including configuration settings. If no names are provided, it shows information for all configured registries.

**Examples:**

```bash
# Show info for all registries
$ arm info registry
my-org:
    type: git
    url: https://github.com/my-org/arm-registry
my-gitlab:
    type: gitlab
    url: https://gitlab.example.com
    group_id: 123
cloudsmith-registry:
    type: cloudsmith
    url: https://api.cloudsmith.io
    owner: my-org
    repository: my-repo

# Show info for specific registries
$ arm info registry my-org
my-org:
    type: git
    url: https://github.com/my-org/arm-registry
```

## Sink Management

### arm add sink

`arm add sink --tool <cursor|copilot|amazonq|markdown> [--force] NAME PATH`

Add a new sink to the ARM configuration. A sink defines where resources should be output for a specific use case. The `--force` flag allows overwriting an existing sink with the same name.

**Examples:**
```bash
# Add Cursor rules sink
$ arm add sink --tool cursor cursor-rules .cursor/rules

# Add Cursor prompts sink
$ arm add sink --tool cursor cursor-commands .cursor/commands

# Add Amazon Q rules sink
$ arm add sink --tool amazonq q-rules .amazonq/rules

# Add Amazon Q prompts sink
$ arm add sink --tool amazonq q-prompts .amazonq/prompts

# Add GitHub Copilot sink
$ arm add sink --tool copilot copilot-rules .github/copilot

# Overwrite an existing sink
$ arm add sink --tool cursor --force cursor-rules .cursor/new-rules
```

### arm remove sink

`arm remove sink NAME`

Remove a sink from the ARM configuration by name. This command removes the specified sink and all its associated configuration from the ARM setup. After removal, the sink will no longer be available for installing rulesets or promptsets.

**Example:**
```bash
$ arm remove sink cursor-rules
```

### arm set sink

`arm set sink NAME KEY VALUE`

Set configuration values for a specific sink. This command allows you to configure sink-specific settings. The available configuration keys are `tool` (cursor, amazonq, copilot, or markdown) and `directory` (output path).

**Examples:**
```bash
# Change sink tool
$ arm set sink cursor-rules tool amazonq

# Change sink tool to markdown
$ arm set sink cursor-rules tool markdown

# Update sink directory
$ arm set sink cursor-rules directory .cursor/new-rules
```

### arm list sink

`arm list sink`

List all configured sinks. This command displays all sinks that have been added to the ARM configuration as a simple list.

**Example:**

```bash
$ arm list sink
- cursor-rules
- q-rules
- copilot-rules
```

### arm info sink

`arm info sink [NAME]...`

Display detailed information about one or more sinks. This command shows comprehensive details about the specified sinks, including configuration settings and output directories. If no names are provided, it shows information for all configured sinks.

**Examples:**

```bash
# Show info for all sinks
$ arm info sink
cursor-rules:
    directory: .cursor/rules
    tool: cursor
q-rules:
    directory: .amazonq/rules
    tool: amazonq
copilot-rules:
    directory: .github/copilot
    tool: copilot

# Show info for specific sinks
$ arm info sink cursor-rules
cursor-rules:
    directory: .cursor/rules
    tool: cursor
```

## Package Management

### arm install

`arm install`

Install all configured dependencies to their assigned sinks.

**Example:**
```bash
# Install all configured packages
$ arm install
```

### arm outdated

`arm outdated [--output <table|json|list>]`

Check for outdated dependencies across all configured registries. This command compares the currently installed versions of rulesets and promptsets with the latest available versions in their respective registries. It shows which packages have newer versions available, displaying the constraint, current version, wanted version, and latest version for each outdated package. The output format can be specified as table (default), JSON, or list.

**Examples:**
```bash
# Check for outdated packages (table format)
$ arm outdated

# Check for outdated packages in JSON format
$ arm outdated --output json

# Check for outdated packages in list format
$ arm outdated --output list
```

**Sample output:**
```bash
$ arm outdated
PACKAGE                         TYPE       CONSTRAINT  CURRENT  WANTED  LATEST
my-org/clean-code-ruleset       ruleset    ^1.0.0      1.0.1    1.1.0   2.0.0
my-org/code-review-promptset    promptset  ^1.0.0      1.0.1    1.1.0   2.0.0

$ arm outdated --output json
[
  {
    "package": "my-org/clean-code-ruleset",
    "type": "ruleset",
    "constraint": "^1.0.0",
    "current": "1.0.1",
    "wanted": "1.1.0",
    "latest": "2.0.0"
  }
]

$ arm outdated --output list
my-org/clean-code-ruleset
my-org/code-review-promptset
```

### arm update

`arm update`

Update all installed packages to their latest available versions. This command checks for updates to all currently installed dependencies and updates them to the latest versions that satisfy their version constraints. It performs the same installation process as `arm install` but with the updated versions.

**Example:**
```bash
# Update all installed packages
$ arm update
```

### arm upgrade

`arm upgrade`

Upgrade all installed packages to their latest available versions, ignoring version constraints. This command updates all currently installed rulesets and promptsets to their absolute latest versions, even if they would violate the version constraints specified in the configuration. By default, the upgrade command also modifies the version constraint to use a major constraint (^X.0.0) based on the newly installed version, allowing future updates within the same major version.

**Example:**
```bash
# Upgrade all packages to latest versions
$ arm upgrade
```

### arm uninstall

`arm uninstall`

Uninstall all configured packages from their assigned sinks. This command removes all currently installed rulesets and promptsets from their output directories, cleaning up the sink directories while preserving the ARM configuration. The packages can be reinstalled later using `arm install`.

**Example:**
```bash
# Uninstall all configured packages
$ arm uninstall
```

### arm list

`arm list`

List all configured entities in the ARM environment. This command provides a comprehensive overview showing all registries, sinks, and installed packages (rulesets and promptsets), grouped by category.

**Example:**

```bash
$ arm list
registries:
    - my-org
    - cloudsmith-registry
sinks:
    - cursor-rules
    - q-rules
    - cursor-commands
    - q-prompts
dependencies:
    - my-org/clean-code-ruleset@1.1.0
    - my-org/security-ruleset@2.1.0
    - my-org/code-review-promptset@1.1.0
    - my-org/testing-promptset@2.0.1
```

### arm info

`arm info`

Display detailed information about all configured entities in the ARM environment. This command shows comprehensive details about all registries, sinks, and installed packages (rulesets and promptsets), including their metadata, configuration, dependencies, and status information. It provides a complete hierarchical view of the entire ARM environment.

**Example:**

```bash
$ arm info
registries:
    sample-registry:
        type: cloudsmith
        url: https://api.cloudsmith.io
        owner: sample-org
        repository: arm-registry
sinks:
    cursor-rules:
        directory: .cursor/rules
        tool: cursor
    amazonq-rules:
        directory: .amazonq/rules
        tool: amazonq
    cursor-commands:
        directory: .cursor/commands
        tool: cursor
    amazonq-prompts:
        directory: .amazonq/prompts
        tool: amazonq
    copilot-instructions:
        directory: .github/instructions
        tool: copilot
dependencies:
    sample-registry/clean-code-ruleset:
        type: ruleset
        version: 1.0.0
        constraint: ^1.0.0
        priority: 100
        sinks:
            - cursor-rules
            - amazonq-rules
            - copilot-instructions
        include:
            - "**/*.yml"
        exclude:
            - "**/experimental/**"
    sample-registry/code-review-promptset:
        type: promptset
        version: 1.0.0
        constraint: ^1.0.0
        sinks:
            - cursor-commands
            - amazonq-prompts
        include:
            - "review/**/*.yml"
```

## Resource Management

### arm install ruleset

`arm install ruleset [--priority PRIORITY] [--include GLOB...] [--exclude GLOB...] REGISTRY_NAME/RULESET_NAME[@VERSION] SINK_NAME...`

Install a specific ruleset from a registry to one or more sinks. This command allows you to specify priority (default: 100), include/exclude patterns for filtering rules (default include: all .yml and .yaml files), and target specific sinks. The ruleset can be installed from a specific version or the latest version that satisfies the constraint.

**Important:** When installing a ruleset to specific sinks, ARM will automatically uninstall the ruleset from any previous sinks that are not in the new sink list. This ensures a clean state and prevents orphaned files across your sinks.

**Examples:**
```bash
# Install ruleset to single sink
$ arm install ruleset my-org/clean-code-ruleset cursor-rules

# Install specific version to multiple sinks
$ arm install ruleset my-org/clean-code-ruleset@1.0.0 cursor-rules q-rules

# Reinstall to different sinks (removes from previous sinks)
# If previously installed to cursor-rules, this will remove it from cursor-rules
$ arm install ruleset my-org/clean-code-ruleset q-rules copilot-rules

# Install with custom priority
$ arm install ruleset --priority 200 my-org/clean-code-ruleset cursor-rules

# Install with include/exclude patterns
$ arm install ruleset --include "**/*.yml" --exclude "**/README.md" my-org/clean-code-ruleset cursor-rules
```

**Reinstall Behavior Example:**
- Initial install to sinks A and B: `arm install ruleset repo/pkg A B`
- Later install to only sink C: `arm install ruleset repo/pkg C`
- Result: Package is removed from A and B, only exists in C

### arm set ruleset

`arm set ruleset REGISTRY_NAME/RULESET_NAME KEY VALUE`

Set configuration values for a specific ruleset. This command allows you to configure ruleset-specific settings. The available configuration keys are `version`, `priority`, `sinks`, `include`, and `exclude`.

**Examples:**
```bash
# Update version constraint
$ arm set ruleset my-org/clean-code-ruleset version ^2.0.0

# Change priority
$ arm set ruleset my-org/clean-code-ruleset priority 150

# Update sinks
$ arm set ruleset my-org/clean-code-ruleset sinks cursor-rules,q-rules,copilot-rules
```

### arm install promptset

`arm install promptset [--include GLOB...] [--exclude GLOB...] REGISTRY_NAME/PROMPTSET[@VERSION] SINK_NAME...`

Install a specific promptset from a registry to one or more sinks. This command allows you to specify include/exclude patterns for filtering prompts (default include: all .yml and .yaml files), and target specific sinks. The promptset can be installed from a specific version or the latest version that satisfies the constraint.

**Important:** When installing a promptset to specific sinks, ARM will automatically uninstall the promptset from any previous sinks that are not in the new sink list. This ensures a clean state and prevents orphaned files across your sinks.

**Examples:**
```bash
# Install promptset to single sink
$ arm install promptset my-org/code-review-promptset cursor-commands

# Install specific version to multiple sinks
$ arm install promptset my-org/code-review-promptset@1.0.0 cursor-commands q-prompts

# Reinstall to different sinks (removes from previous sinks)
$ arm install promptset my-org/code-review-promptset q-prompts

# Install with include/exclude patterns
$ arm install promptset --include "**/*.yml" --exclude "**/README.md" my-org/code-review-promptset cursor-commands
```

### arm set promptset

`arm set promptset REGISTRY_NAME/PROMPTSET_NAME KEY VALUE`

Set configuration values for a specific promptset. This command allows you to configure promptset-specific settings. The available configuration keys are `version`, `sinks`, `include`, and `exclude`.

**Examples:**
```bash
# Update version constraint
$ arm set promptset my-org/code-review-promptset version ^2.0.0

# Update sinks
$ arm set promptset my-org/code-review-promptset sinks cursor-commands,q-prompts

# Update include patterns
$ arm set promptset my-org/code-review-promptset include "review/**/*.yml","refactor/**/*.yml"
```

## Utilities

### arm clean cache

`arm clean cache [--nuke | --max-age DURATION]`

Clean the local cache directory. This command removes cached registry data and downloaded packages from the local cache. The `--nuke` flag performs a more aggressive cleanup, removing all cached data including registry indexes and package archives. The `--max-age` flag allows you to specify how old cached data should be before it's removed. Without any flags, it performs a standard cleanup of outdated or corrupted cache entries (default: 7 days).

**Flags:**
- `--nuke`: Aggressive cleanup (remove all cached data)
- `--max-age`: Remove cached data older than specified duration (e.g., "30m", "2h", "7d")

**Duration Format:**
The `--max-age` flag supports duration strings with units:
- **Minutes**: `30m`, `60m`
- **Hours**: `2h`, `24h`
- **Days**: `1d`, `7d`
- **Combined**: `1h30m`, `2d12h`

**Examples:**
```bash
# Standard cache cleanup (removes data older than 7 days)
$ arm clean cache

# Remove data older than 30 minutes
$ arm clean cache --max-age 30m

# Remove data older than 2 hours
$ arm clean cache --max-age 2h

# Remove data older than 1 day
$ arm clean cache --max-age 1d

# Remove data older than 1 hour and 30 minutes
$ arm clean cache --max-age 1h30m

# Aggressive cleanup (remove all cached data)
$ arm clean cache --nuke
```

**Note:** The `--nuke` and `--max-age` flags are mutually exclusive and cannot be used together.

### arm clean sinks

`arm clean sinks [--nuke]`

Clean sink directories based on the ARM index. This command removes files from sink directories that shouldn't be there according to the arm-index.json file. The `--nuke` flag performs a more aggressive cleanup, clearing out the entire ARM directory entirely. Without the flag, it performs a selective cleanup based on the index.

**Examples:**
```bash
# Selective cleanup based on ARM index
$ arm clean sinks

# Complete cleanup (remove entire ARM directory)
$ arm clean sinks --nuke
```

### arm compile

`arm compile [--tool <markdown|cursor|amazonq|copilot>] [--force] [--recursive] [--validate-only] [--include GLOB...] [--exclude GLOB...] [--fail-fast] INPUT_PATH... [OUTPUT_PATH]`

Compile rulesets and promptsets from source files. This command compiles source ruleset and promptset files to platform-specific formats. It supports different tool platforms (markdown, cursor, amazonq, copilot), recursive directory processing, validation-only mode, and various filtering and output options. This is useful for development and testing of rulesets and promptsets before publishing to registries.

**INPUT_PATH** accepts both files and directories:
- **Files**: Directly processes the specified file(s)
- **Directories**: Discovers files within using `--include`/`--exclude` patterns (non-recursive by default)
- **Mixed**: Can combine files and directories in the same command

**Note:** Shell glob patterns (e.g., `*.yml`) are expanded to individual file paths by your shell before ARM processes them. When using `--validate-only`, the OUTPUT_PATH argument is optional and will be ignored if provided.

**Examples:**
```bash
# Compile single file to Cursor format
$ arm compile --tool cursor ruleset.yml ./output/

# Compile multiple files
$ arm compile --tool cursor file1.yml file2.yml ./output/

# Compile with shell glob expansion (expands to individual files)
$ arm compile --tool cursor rulesets/*.yml ./output/

# Compile directory recursively to Amazon Q format
$ arm compile --tool amazonq --recursive ./src/ ./build/

# Compile directory non-recursively (default)
$ arm compile --tool cursor ./rulesets/ ./output/

# Mix files and directories
$ arm compile --tool cursor specific.yml ./more-rulesets/ ./output/

# Validate only (no output files) - OUTPUT_PATH is optional
$ arm compile --validate-only ruleset.yml

# Validate multiple files without output
$ arm compile --validate-only ./rulesets/*.yml

# Compile with include/exclude patterns
$ arm compile --tool cursor --include "**/*.yml" --exclude "**/README.md" ./src/ ./build/

# Compilation with force overwrite
$ arm compile --tool copilot --force ruleset.yml ./output/

# Validate and fail fast on first error (useful for CI)
$ arm compile --validate-only --fail-fast ./rulesets/
```
