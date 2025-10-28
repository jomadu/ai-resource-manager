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
  - [Ruleset Management](#ruleset-management)
    - [arm install ruleset](#arm-install-ruleset)
    - [arm uninstall ruleset](#arm-uninstall-ruleset)
    - [arm set ruleset](#arm-set-ruleset)
    - [arm list ruleset](#arm-list-ruleset)
    - [arm info ruleset](#arm-info-ruleset)
    - [arm update ruleset](#arm-update-ruleset)
    - [arm upgrade ruleset](#arm-upgrade-ruleset)
    - [arm outdated ruleset](#arm-outdated-ruleset)
  - [Promptset Management](#promptset-management)
    - [arm install promptset](#arm-install-promptset)
    - [arm uninstall promptset](#arm-uninstall-promptset)
    - [arm set promptset](#arm-set-promptset)
    - [arm list promptset](#arm-list-promptset)
    - [arm info promptset](#arm-info-promptset)
    - [arm update promptset](#arm-update-promptset)
    - [arm upgrade promptset](#arm-upgrade-promptset)
    - [arm outdated promptset](#arm-outdated-promptset)
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
AI Resource Manager v1.2.3
Build: 2024-01-15T10:30:45Z
Commit: a1b2c3d4
Platform: darwin/arm64
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

Add a new Cloudsmith registry to the ARM configuration. Cloudsmith registries use Cloudsmith's package repository service for single-file artifacts. The URL defaults to `https://api.cloudsmith.io` if not specified.

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

List all configured registries. This command displays all registries that have been added to the ARM configuration, showing their names, types, and basic information in a tabular format.

**Example:**
```bash
$ arm list registry
NAME              TYPE        URL                                          CONFIG
my-org            git         https://github.com/my-org/arm-registry
my-gitlab         gitlab      https://gitlab.com                           group_id=123
my-gitlab-project gitlab      https://gitlab.com                           project_id=101
cloudsmith-registry cloudsmith  https://api.cloudsmith.io                  owner=my-org, repository=my-repo
```

### arm info registry

`arm info registry [NAME]...`

Display detailed information about one or more registries. This command shows comprehensive details about the specified registries, including configuration settings, available packages, and status information. If no names are provided, it shows information for all configured registries.

**Examples:**
```bash
# Show info for all registries
$ arm info registry

# Show info for specific registries
$ arm info registry my-org cloudsmith-registry
```

**Sample output:**
```bash
$ arm info registry my-org
Registry: my-org
Type: git
URL: https://github.com/my-org/arm-registry

$ arm info registry my-gitlab
Registry: my-gitlab
Type: gitlab
URL: https://gitlab.example.com
Group Id: 123

$ arm info registry cloudsmith-registry
Registry: cloudsmith-registry
Type: cloudsmith
URL: https://api.cloudsmith.io
Owner: my-org
Repository: my-repo
```

## Sink Management

### arm add sink

`arm add sink [--type <cursor|copilot|amazonq>] [--layout <hierarchical|flat>] [--compile-to <md|cursor|amazonq|copilot>] [--force] NAME PATH`

Add a new sink to the ARM configuration. A sink defines where and how compiled rulesets and promptsets should be output. The `--type` flag is a shortcut that sets combinations of `--layout` and `--compile-to` (e.g., `--type cursor` sets `--layout hierarchical --compile-to cursor`). You can also specify `--layout` and `--compile-to` individually for custom configurations. The `--force` flag allows overwriting an existing sink with the same name.

**Examples:**
```bash
# Add Cursor rules sink
$ arm add sink --type cursor cursor-rules .cursor/rules

# Add Cursor prompts sink
$ arm add sink --type cursor cursor-commands .cursor/commands

# Add Amazon Q rules sink
$ arm add sink --type amazonq q-rules .amazonq/rules

# Add Amazon Q prompts sink
$ arm add sink --type amazonq q-prompts .amazonq/prompts

# Add GitHub Copilot sink
$ arm add sink --type copilot copilot-rules .github/copilot

# Overwrite an existing sink
$ arm add sink --type cursor --force cursor-rules .cursor/new-rules
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

Set configuration values for a specific sink. This command allows you to configure sink-specific settings. The available configuration keys are `layout` (hierarchical or flat), `directory` (output path), and `compile_target` (md, cursor, amazonq, or copilot).

**Examples:**
```bash
# Change sink layout
$ arm set sink cursor-rules layout flat

# Update sink directory
$ arm set sink cursor-rules directory .cursor/new-rules

# Change compilation target
$ arm set sink cursor-rules compile_target md
```

### arm list sink

`arm list sink`

List all configured sinks. This command displays all sinks that have been added to the ARM configuration, showing their names, types, output directories, and basic configuration in a tabular format.

**Example:**
```bash
$ arm list sink
NAME           LAYOUT        COMPILE_TARGET  DIRECTORY
cursor-rules   hierarchical  cursor          .cursor/rules
q-rules        hierarchical  md              .amazonq/rules
copilot-rules  flat          copilot         .github/copilot
```

### arm info sink

`arm info sink [NAME]...`

Display detailed information about one or more sinks. This command shows comprehensive details about the specified sinks, including configuration settings, output directories, layout preferences, and status information. If no names are provided, it shows information for all configured sinks.

**Examples:**
```bash
# Show info for all sinks
$ arm info sink

# Show info for specific sinks
$ arm info sink cursor-rules q-rules
```

**Sample output:**
```bash
$ arm info sink cursor-rules
Sink: cursor-rules
Layout: hierarchical
Compile Target: cursor
Directory: .cursor/rules
```

## Package Management

### arm install

`arm install`

Install all configured packages (rulesets and promptsets) to their assigned sinks. This command processes the ARM configuration file and installs all rulesets and promptsets that are configured. Packages can be compiled from source files or installed as pre-compiled files from repositories, then placed in the correct output directories for each sink.

**Example:**
```bash
# Install all configured packages
$ arm install
```

### arm outdated

`arm outdated [--output <table|json|list>]`

Check for outdated packages across all configured registries. This command compares the currently installed versions of rulesets and promptsets with the latest available versions in their respective registries. It shows which packages have newer versions available, displaying the constraint, current version, wanted version, and latest version for each outdated package. The output format can be specified as table (default), JSON, or list.

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

Update all installed packages to their latest available versions. This command checks for updates to all currently installed rulesets and promptsets and updates them to the latest versions that satisfy their version constraints. It performs the same installation process as `arm install` but with the updated versions.

**Example:**
```bash
# Update all installed packages
$ arm update
```

### arm upgrade

`arm upgrade`

Upgrade all installed packages to their latest available versions, ignoring version constraints. This command updates all currently installed rulesets and promptsets to their absolute latest versions, even if they would violate the version constraints specified in the configuration. This is useful for testing or when you want to move to the newest versions regardless of compatibility constraints.

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

List all installed packages across all sinks. This command displays all currently installed rulesets and promptsets, showing their names, versions, source registries, and which sinks they are installed to. It provides a comprehensive overview of the current installation state.

**Example:**
```bash
# List all installed packages
$ arm list
```

**Sample output:**
```bash
$ arm list
PACKAGE                         VERSION  SINKS
my-org/clean-code-ruleset       1.1.0    cursor-rules, q-rules
my-org/code-review-promptset    1.1.0    cursor-commands, q-prompts
```

### arm info

`arm info`

Display detailed information about all installed packages. This command shows comprehensive details about all currently installed rulesets and promptsets, including their metadata, configuration, dependencies, and status information. It provides a detailed overview of the entire installation state.

**Example:**
```bash
# Show detailed info for all installed packages
$ arm info
```

**Sample output:**
```bash
$ arm info
Package: my-org/clean-code-ruleset
Type: Ruleset
Version: 1.1.0
Constraint: ^1.0.0
Priority: 100
Sinks: cursor-rules, q-rules
Include: **/*.yml, **/*.yaml

Package: my-org/code-review-promptset
Type: Promptset
Version: 1.1.0
Constraint: ~1.1.0
Sinks: cursor-commands, q-prompts
Include: **/*.yml, **/*.yaml
```

## Ruleset Management

### arm install ruleset

`arm install ruleset [--priority PRIORITY] [--include GLOB...] [--exclude GLOB...] REGISTRY_NAME/RULESET_NAME[@VERSION] SINK_NAME...`

Install a specific ruleset from a registry to one or more sinks. This command allows you to specify priority (default: 100), include/exclude patterns for filtering rules (default include: all .yml and .yaml files), and target specific sinks. The ruleset can be installed from a specific version or the latest version that satisfies the constraint.

**Examples:**
```bash
# Install ruleset to single sink
$ arm install ruleset my-org/clean-code-ruleset cursor-rules

# Install specific version to multiple sinks
$ arm install ruleset my-org/clean-code-ruleset@1.0.0 cursor-rules q-rules

# Install with custom priority
$ arm install ruleset --priority 200 my-org/clean-code-ruleset cursor-rules

# Install with include/exclude patterns
$ arm install ruleset --include "**/*.yml" --exclude "**/README.md" my-org/clean-code-ruleset cursor-rules
```

### arm uninstall ruleset

`arm uninstall ruleset REGISTRY_NAME/RULESET_NAME`

Uninstall a specific ruleset from all sinks. This command removes the specified ruleset from all sinks where it is currently installed, cleaning up all ruleset files. The ruleset is also removed from the ARM configuration. The ruleset can be reinstalled later using `arm install ruleset`.

**Example:**
```bash
$ arm uninstall ruleset my-org/clean-code-ruleset
```

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

### arm list ruleset

`arm list ruleset`

List all installed rulesets. This command displays all currently installed rulesets in list format, showing their names, versions, source registries, priority, and which sinks they are installed to.

**Example:**
```bash
$ arm list ruleset
```

**Sample output:**
```bash
$ arm list ruleset
RULESET                         VERSION  PRIORITY  SINKS
my-org/clean-code-ruleset       1.0.1    100       cursor-rules, q-rules
my-org/security-ruleset         2.1.0    200       cursor-rules, q-rules, copilot-rules
```

### arm info ruleset

`arm info ruleset [REGISTRY_NAME/RULESET_NAME...]`

Display detailed information about one or more rulesets. This command shows comprehensive details about the specified rulesets, including registry, name, version constraint, resolved version, include, exclude, sinks, and directories where it's installed. If no names are provided, it shows information for all installed rulesets.

**Examples:**
```bash
# Show info for all rulesets
$ arm info ruleset

# Show info for specific rulesets
$ arm info ruleset my-org/clean-code-ruleset my-org/security-ruleset
```

**Sample output:**
```bash
$ arm info ruleset my-org/clean-code-ruleset
Package: my-org/clean-code-ruleset
Type: ruleset
Version: 1.0.1
Constraint: ^1.0.0
Priority: 100
Sinks: cursor-rules, q-rules
Include: **/*.yml, **/*.yaml
```

### arm update ruleset

`arm update ruleset [REGISTRY_NAME/RULESET_NAME...]`

Update one or more rulesets to their latest available versions. This command checks for updates to the specified rulesets and updates them to the latest versions that satisfy their version constraints. If no ruleset names are provided, it updates all installed rulesets. It performs the same installation process as `arm install ruleset` but with the updated versions.

**Examples:**
```bash
# Update all rulesets
$ arm update ruleset

# Update specific rulesets
$ arm update ruleset my-org/clean-code-ruleset my-org/security-ruleset
```

### arm upgrade ruleset

`arm upgrade ruleset [REGISTRY_NAME/RULESET_NAME...]`

Upgrade one or more rulesets to their latest available versions, ignoring version constraints. This command updates the specified rulesets to their absolute latest versions, even if they would violate the version constraints specified in the configuration. If no ruleset names are provided, it upgrades all installed rulesets. This is useful for testing or when you want to move to the newest versions regardless of compatibility constraints.

**Examples:**
```bash
# Upgrade all rulesets to latest versions
$ arm upgrade ruleset

# Upgrade specific rulesets to latest versions
$ arm upgrade ruleset my-org/clean-code-ruleset my-org/security-ruleset
```

### arm outdated ruleset

`arm outdated [--output <table|json|list>] ruleset`

Check for outdated rulesets across all configured registries. This command compares the currently installed versions of rulesets with the latest available versions in their respective registries. It shows which rulesets have newer versions available, displaying the constraint, current version, wanted version, and latest version for each outdated ruleset. The output format can be specified as table (default), JSON, or list.

**Examples:**
```bash
# Check for outdated rulesets (table format)
$ arm outdated ruleset

# Check for outdated rulesets in JSON format
$ arm outdated --output json ruleset

# Check for outdated rulesets in list format
$ arm outdated --output list ruleset
```

**Sample output:**
```bash
$ arm outdated ruleset
PACKAGE                         TYPE     CONSTRAINT  CURRENT  WANTED  LATEST
my-org/clean-code-ruleset       ruleset  ^1.0.0      1.0.1    1.0.2   2.0.0
my-org/security-ruleset         ruleset  ~2.1.0      2.1.0    2.1.1   2.1.1

$ arm outdated --output json ruleset
[
  {
    "package": "my-org/clean-code-ruleset",
    "type": "ruleset",
    "constraint": "^1.0.0",
    "current": "1.0.1",
    "wanted": "1.0.2",
    "latest": "2.0.0"
  }
]

$ arm outdated --output list ruleset
my-org/clean-code-ruleset
my-org/security-ruleset
```

## Promptset Management

### arm install promptset

`arm install promptset [--include GLOB...] [--exclude GLOB...] REGISTRY_NAME/PROMPTSET[@VERSION] SINK_NAME...`

Install a specific promptset from a registry to one or more sinks. This command allows you to specify include/exclude patterns for filtering prompts (default include: all .yml and .yaml files), and target specific sinks. The promptset can be installed from a specific version or the latest version that satisfies the constraint.

**Examples:**
```bash
# Install promptset to single sink
$ arm install promptset my-org/code-review-promptset cursor-commands

# Install specific version to multiple sinks
$ arm install promptset my-org/code-review-promptset@1.0.0 cursor-commands q-prompts

# Install with include/exclude patterns
$ arm install promptset --include "**/*.yml" --exclude "**/README.md" my-org/code-review-promptset cursor-commands
```

### arm uninstall promptset

`arm uninstall promptset REGISTRY_NAME/PROMPTSET`

Uninstall a specific promptset from all sinks. This command removes the specified promptset from all sinks where it is currently installed, cleaning up all promptset files. The promptset is also removed from the ARM configuration. The promptset can be reinstalled later using `arm install promptset`.

**Example:**
```bash
$ arm uninstall promptset my-org/code-review-promptset
```

### arm set promptset

`arm set promptset REGISTRY_NAME/PROMPTSET KEY VALUE`

Set configuration values for a specific promptset. This command allows you to configure promptset-specific settings. The available configuration keys are `version`, `sinks`, `include`, and `exclude`.

**Examples:**
```bash
# Update version constraint
$ arm set promptset my-org/code-review-promptset version ^2.0.0

# Update sinks
$ arm set promptset my-org/code-review-promptset sinks cursor-commands,q-prompts

# Update include pattern
$ arm set promptset my-org/code-review-promptset include "**/*.yml,**/*.yaml"
```

### arm list promptset

`arm list promptset`

List all installed promptsets. This command displays all currently installed promptsets in list format, showing their names, versions, source registries, and which sinks they are installed to.

**Example:**
```bash
$ arm list promptset
```

**Sample output:**
```bash
$ arm list promptset
PROMPTSET                       VERSION  SINKS
my-org/code-review-promptset    1.1.0    cursor-commands, q-prompts
my-org/testing-promptset        2.0.1    cursor-commands, q-prompts
```

### arm info promptset

`arm info promptset [REGISTRY_NAME/PROMPTSET...]`

Display detailed information about one or more promptsets. This command shows comprehensive details about the specified promptsets, including registry, name, version constraint, resolved version, include, exclude, sinks, and directories where it's installed. If no names are provided, it shows information for all installed promptsets.

**Examples:**
```bash
# Show info for all promptsets
$ arm info promptset

# Show info for specific promptsets
$ arm info promptset my-org/code-review-promptset my-org/testing-promptset
```

**Sample output:**
```bash
$ arm info promptset my-org/code-review-promptset
Package: my-org/code-review-promptset
Type: promptset
Version: 1.1.0
Constraint: ^1.0.0
Sinks: cursor-commands, q-prompts
Include: **/*.yml, **/*.yaml
Exclude: none
```

### arm update promptset

`arm update promptset [REGISTRY_NAME/PROMPTSET...]`

Update one or more promptsets to their latest available versions. This command checks for updates to the specified promptsets and updates them to the latest versions that satisfy their version constraints. If no promptset names are provided, it updates all installed promptsets. It performs the same installation process as `arm install promptset` but with the updated versions.

**Examples:**
```bash
# Update all promptsets
$ arm update promptset

# Update specific promptsets
$ arm update promptset my-org/code-review-promptset my-org/testing-promptset
```

### arm upgrade promptset

`arm upgrade promptset [REGISTRY_NAME/PROMPTSET...]`

Upgrade one or more promptsets to their latest available versions, ignoring version constraints. This command updates the specified promptsets to their absolute latest versions, even if they would violate the version constraints specified in the configuration. If no promptset names are provided, it upgrades all installed promptsets. This is useful for testing or when you want to move to the newest versions regardless of compatibility constraints.

**Examples:**
```bash
# Upgrade all promptsets to latest versions
$ arm upgrade promptset

# Upgrade specific promptsets to latest versions
$ arm upgrade promptset my-org/code-review-promptset my-org/testing-promptset
```

### arm outdated promptset

`arm outdated [--output <table|json|list>] promptset`

Check for outdated promptsets across all configured registries. This command compares the currently installed versions of promptsets with the latest available versions in their respective registries. It shows which promptsets have newer versions available, displaying the constraint, current version, wanted version, and latest version for each outdated promptset. The output format can be specified as table (default), JSON, or list.

**Examples:**
```bash
# Check for outdated promptsets (table format)
$ arm outdated promptset

# Check for outdated promptsets in JSON format
$ arm outdated --output json promptset

# Check for outdated promptsets in list format
$ arm outdated --output list promptset
```

**Sample output:**
```bash
$ arm outdated promptset
PACKAGE                         TYPE       CONSTRAINT  CURRENT  WANTED  LATEST
my-org/code-review-promptset    promptset  ^1.0.0      1.1.0    1.1.2   2.0.0
my-org/testing-promptset        promptset  ~2.0.0      2.0.1    2.0.2   2.0.2

$ arm outdated --output json promptset
[
  {
    "package": "my-org/code-review-promptset",
    "type": "promptset",
    "constraint": "^1.0.0",
    "current": "1.1.0",
    "wanted": "1.1.2",
    "latest": "2.0.0"
  }
]

$ arm outdated --output list promptset
my-org/code-review-promptset
my-org/testing-promptset
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

`arm compile [--target <md|cursor|amazonq|copilot>] [--force] [--recursive] [--validate-only] [--include GLOB...] [--exclude GLOB...] [--fail-fast] INPUT_PATH... [OUTPUT_PATH]`

Compile rulesets and promptsets from source files. This command compiles source ruleset and promptset files to platform-specific formats. It supports different target platforms (md, cursor, amazonq, copilot), recursive directory processing, validation-only mode, and various filtering and output options. This is useful for development and testing of rulesets and promptsets before publishing to registries.

**Note:** When using `--validate-only`, the OUTPUT_PATH argument is optional and will be ignored if provided. The command will only validate the syntax of the input files without generating any output.

**Examples:**
```bash
# Compile single file to Cursor format
$ arm compile --target cursor ruleset.yml ./output/

# Compile directory recursively to Amazon Q format
$ arm compile --target amazonq --recursive ./src/ ./build/

# Validate only (no output files) - OUTPUT_PATH is optional
$ arm compile --validate-only ruleset.yml

# Validate multiple files without output
$ arm compile --validate-only ./rulesets/*.yml

# Compile with include/exclude patterns
$ arm compile --target cursor --include "**/*.yml" --exclude "**/README.md" ./src/ ./build/

# Compilation with force overwrite
$ arm compile --target copilot --force ruleset.yml ./output/

# Validate and fail fast on first error (useful for CI)
$ arm compile --validate-only --fail-fast ./rulesets/
```
