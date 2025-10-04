# Commands

## Core

`arm version`

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

`arm help`

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
$ arm help add registry
$ arm add registry --help
```

## Registry Management

`arm add registry --type <git|gitlab|cloudsmith> [--gitlab-group-id ID] [--gitlab-project-id ID] NAME URL`

Add a new registry to the ARM configuration. This command supports different registry types (git, gitlab, cloudsmith) and allows specifying additional parameters like GitLab group and project IDs for more precise targeting.

**Examples:**
```bash
# Add a Git registry
$ arm add registry --type git my-org https://github.com/my-org/arm-registry

# Add a GitLab registry with group ID
$ arm add registry --type gitlab --gitlab-group-id 123 my-gitlab https://gitlab.com

# Add a GitLab registry with project ID
$ arm add registry --type gitlab --gitlab-project-id 456 my-gitlab-project https://gitlab.com

# Add a Cloudsmith registry
$ arm add registry --type cloudsmith sample-registry https://app.cloudsmith.com/sample-org/arm-registry
```

`arm remove registry NAME`

Remove a registry from the ARM configuration by name. This command removes the specified registry and all its associated configuration from the ARM setup. After removal, the registry will no longer be available for installing rulesets or promptsets.

**Example:**
```bash
$ arm remove registry my-org
```

`arm config registry set NAME KEY VALUE`

Set configuration values for a specific registry. This command allows you to configure registry-specific settings such as URL endpoints or other registry-specific parameters. The exact configuration keys available depend on the registry type.

**Examples:**
```bash
# Update registry URL
$ arm config registry set my-org url https://github.com/my-org/new-arm-registry

# Set GitLab group ID
$ arm config registry set my-gitlab gitlab_group_id 789

# Set GitLab project ID
$ arm config registry set my-gitlab-project gitlab_project_id 101
```

`arm list registry`

List all configured registries. This command displays all registries that have been added to the ARM configuration, showing their names, types, and basic information in a tabular format.

**Example:**
```bash
$ arm list registry
NAME              TYPE        URL                                          CONFIG
my-org            git         https://github.com/my-org/arm-registry
my-gitlab         gitlab      https://gitlab.com                          gitlab_group_id=123
my-gitlab-project gitlab      https://gitlab.com                          gitlab_project_id=101
sample-registry   cloudsmith  https://app.cloudsmith.com/sample-org/arm-registry
```

`arm info registry [NAME]...`

Display detailed information about one or more registries. This command shows comprehensive details about the specified registries, including configuration settings, available packages, and status information. If no names are provided, it shows information for all configured registries.

**Examples:**
```bash
# Show info for all registries
$ arm info registry

# Show info for specific registries
$ arm info registry my-org sample-registry
```

**Sample output:**
```bash
$ arm info registry my-org
Registry: my-org
Type: git
URL: https://github.com/my-org/arm-registry
Status: connected
Available packages:
  - clean-code-ruleset (versions: 1.0.0, 1.0.1, 1.1.0)
  - security-ruleset (versions: 2.0.0, 2.1.0)
  - code-review-promptset (versions: 1.0.0, 1.1.0)
```

## Sink Management

`arm add sink [--type <cursor|copilot|amazonq>] [--layout <hierarchical|flat>] [--compile-to <md|cursor|amazonq|copilot>] NAME PATH`

Add a new sink to the ARM configuration. A sink defines where and how compiled rulesets and promptsets should be output. The `--type` flag is a shortcut that sets combinations of `--layout` and `--compile-to` (e.g., `--type cursor` sets `--layout hierarchical --compile-to cursor`). You can also specify `--layout` and `--compile-to` individually for custom configurations.

`arm remove sink NAME`

Remove a sink from the ARM configuration by name. This command removes the specified sink and all its associated configuration from the ARM setup. After removal, the sink will no longer be available for installing rulesets or promptsets.

`arm config sink set NAME KEY VALUE`
Set configuration values for a specific sink. This command allows you to configure sink-specific settings. The available configuration keys are `layout` (hierarchical or flat), `directory` (output path), and `compile_target` (md, cursor, amazonq, or copilot).

`arm list sink`

List all configured sinks. This command displays all sinks that have been added to the ARM configuration, showing their names, types, output directories, and basic configuration in a tabular format.

`arm info sink [NAME]...`

Display detailed information about one or more sinks. This command shows comprehensive details about the specified sinks, including configuration settings, output directories, layout preferences, and status information. If no names are provided, it shows information for all configured sinks.

## Resource Management

`arm install`

Install all configured resources (rulesets and promptsets) to their assigned sinks. This command processes the ARM configuration file and installs all rulesets and promptsets that are configured. Resources can be compiled from source files or installed as pre-compiled files from repositories, then placed in the correct output directories for each sink.

`arm outdated [--output <table|json|list>]`

Check for outdated resources across all configured registries. This command compares the currently installed versions of rulesets and promptsets with the latest available versions in their respective registries. It shows which resources have newer versions available, displaying the constraint, current version, wanted version, and latest version for each outdated package. The output format can be specified as table (default), JSON, or list.

`arm update`

Update all installed resources to their latest available versions. This command checks for updates to all currently installed rulesets and promptsets and updates them to the latest versions that satisfy their version constraints. It performs the same installation process as `arm install` but with the updated versions.

`arm upgrade`

Upgrade all installed resources to their latest available versions, ignoring version constraints. This command updates all currently installed rulesets and promptsets to their absolute latest versions, even if they would violate the version constraints specified in the configuration. This is useful for testing or when you want to move to the newest versions regardless of compatibility constraints.

`arm uninstall`

Uninstall all configured resources from their assigned sinks. This command removes all currently installed rulesets and promptsets from their output directories, cleaning up the sink directories while preserving the ARM configuration. The resources can be reinstalled later using `arm install`.

`arm list`

List all installed resources across all sinks. This command displays all currently installed rulesets and promptsets, showing their names, versions, source registries, and which sinks they are installed to. It provides a comprehensive overview of the current installation state.

`arm info`

Display detailed information about all installed resources. This command shows comprehensive details about all currently installed rulesets and promptsets, including their metadata, configuration, dependencies, and status information. It provides a detailed overview of the entire installation state.

### Ruleset Management

`arm install ruleset [--priority PRIORITY] [--include GLOB...] [--exclude GLOB...] REGISTRY_NAME/RULESET_NAME[@VERSION] SINK_NAME...`

Install a specific ruleset from a registry to one or more sinks. This command allows you to specify priority (default: 100), include/exclude patterns for filtering rules (default include: all .yml and .yaml files), and target specific sinks. The ruleset can be installed from a specific version or the latest version that satisfies the constraint.

`arm uninstall ruleset REGISTRY_NAME/RULESET_NAME`

Uninstall a specific ruleset from all sinks. This command removes the specified ruleset from all sinks where it is currently installed, cleaning up all ruleset files. The ruleset is also removed from the ARM configuration. The ruleset can be reinstalled later using `arm install ruleset`.

`arm config ruleset set REGISTRY_NAME/RULESET_NAME KEY VALUE`

Set configuration values for a specific ruleset. This command allows you to configure ruleset-specific settings. The available configuration keys are `version`, `priority`, `sinks`, `includes`, and `excludes`.

`arm list ruleset`

List all installed rulesets. This command displays all currently installed rulesets in list format, showing their names, versions, source registries, priority, and which sinks they are installed to.

`arm info ruleset [REGISTRY_NAME/RULESET_NAME...]`

Display detailed information about one or more rulesets. This command shows comprehensive details about the specified rulesets, including registry, name, version constraint, resolved version, includes, excludes, sinks, and directories where it's installed. If no names are provided, it shows information for all installed rulesets.

`arm update ruleset [REGISTRY_NAME/RULESET_NAME...]`

Update one or more rulesets to their latest available versions. This command checks for updates to the specified rulesets and updates them to the latest versions that satisfy their version constraints. If no ruleset names are provided, it updates all installed rulesets. It performs the same installation process as `arm install ruleset` but with the updated versions.

`arm upgrade ruleset [REGISTRY_NAME/RULESET_NAME...]`

Upgrade one or more rulesets to their latest available versions, ignoring version constraints. This command updates the specified rulesets to their absolute latest versions, even if they would violate the version constraints specified in the configuration. If no ruleset names are provided, it upgrades all installed rulesets. This is useful for testing or when you want to move to the newest versions regardless of compatibility constraints.

`arm outdated [--output <table|json|list>] ruleset`

Check for outdated rulesets across all configured registries. This command compares the currently installed versions of rulesets with the latest available versions in their respective registries. It shows which rulesets have newer versions available, displaying the constraint, current version, wanted version, and latest version for each outdated ruleset. The output format can be specified as table (default), JSON, or list.

### Promptset Management

`arm install promptset [--include GLOB...] [--exclude GLOB...] REGISTRY_NAME/PROMPTSET[@VERSION] SINK_NAME...`

Install a specific promptset from a registry to one or more sinks. This command allows you to specify include/exclude patterns for filtering prompts (default include: all .yml and .yaml files), and target specific sinks. The promptset can be installed from a specific version or the latest version that satisfies the constraint.

`arm uninstall promptset REGISTRY_NAME/PROMPTSET`

Uninstall a specific promptset from all sinks. This command removes the specified promptset from all sinks where it is currently installed, cleaning up all promptset files. The promptset is also removed from the ARM configuration. The promptset can be reinstalled later using `arm install promptset`.

`arm config promptset set REGISTRY_NAME/PROMPTSET KEY VALUE`

Set configuration values for a specific promptset. This command allows you to configure promptset-specific settings. The available configuration keys are `version`, `sinks`, `includes`, and `excludes`.

`arm list promptset`

List all installed promptsets. This command displays all currently installed promptsets in list format, showing their names, versions, source registries, and which sinks they are installed to.

`arm info promptset [REGISTRY_NAME/PROMPTSET...]`

Display detailed information about one or more promptsets. This command shows comprehensive details about the specified promptsets, including registry, name, version constraint, resolved version, includes, excludes, sinks, and directories where it's installed. If no names are provided, it shows information for all installed promptsets.

`arm update promptset [REGISTRY_NAME/PROMPTSET...]`

Update one or more promptsets to their latest available versions. This command checks for updates to the specified promptsets and updates them to the latest versions that satisfy their version constraints. If no promptset names are provided, it updates all installed promptsets. It performs the same installation process as `arm install promptset` but with the updated versions.

`arm upgrade promptset [REGISTRY_NAME/PROMPTSET...]`

Upgrade one or more promptsets to their latest available versions, ignoring version constraints. This command updates the specified promptsets to their absolute latest versions, even if they would violate the version constraints specified in the configuration. If no promptset names are provided, it upgrades all installed promptsets. This is useful for testing or when you want to move to the newest versions regardless of compatibility constraints.

`arm outdated [--output <table|json|list>] promptset`

Check for outdated promptsets across all configured registries. This command compares the currently installed versions of promptsets with the latest available versions in their respective registries. It shows which promptsets have newer versions available, displaying the constraint, current version, wanted version, and latest version for each outdated promptset. The output format can be specified as table (default), JSON, or list.

## Utilities

`arm clean cache [--nuke]`

Clean the local cache directory. This command removes cached registry data and downloaded packages from the local cache. The `--nuke` flag performs a more aggressive cleanup, removing all cached data including registry indexes and package archives. Without the flag, it performs a standard cleanup of outdated or corrupted cache entries.

`arm clean sinks [--nuke]`

Clean sink directories based on the ARM index. This command removes files from sink directories that shouldn't be there according to the arm-index.json file. The `--nuke` flag performs a more aggressive cleanup, clearing out the entire ARM directory entirely. Without the flag, it performs a selective cleanup based on the index.

`arm compile [--target <md|cursor|amazonq|copilot>] [--force] [--recursive] [--verbose] [--validate-only] [--include GLOB...] [--output GLOB...] [--fail-fast] INPUT_PATH... OUTPUT_PATH`
Compile rulesets and promptsets from source files. This command compiles source ruleset and promptset files to platform-specific formats. It supports different target platforms (md, cursor, amazonq, copilot), recursive directory processing, validation-only mode, and various filtering and output options. This is useful for development and testing of rulesets and promptsets before publishing to registries.
