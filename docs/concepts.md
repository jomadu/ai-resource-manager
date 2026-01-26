# Concepts

## Core Files

- `arm.json` - Project manifest containing registries, packages, and sinks
- `arm-lock.json` - Locked versions for reproducible installs
- `arm-index.json` - Local inventory tracking installed packages and file paths
- `arm_index.*` - Generated priority rules that help AI tools resolve conflicts

## Core Components

**Registries** are remote sources where packages are stored and versioned (Git repositories, GitLab Package Registry, Cloudsmith).

**Packages** are versioned collections stored in registries that can be installed as different resource types:
- **Rulesets** - Collections of AI rules with priority-based conflict resolution
- **Promptsets** - Collections of AI prompts for reusable templates

**Important:** Packages are independent units. They do not have dependencies on each other. When we refer to "dependencies" in ARM, we mean packages installed in your project (similar to npm's package.json), not dependencies between packages.

**Sinks** are output destinations where ARM compiles and places files. Each sink specifies a target directory and compilation format for a specific AI tool (Cursor, Amazon Q, GitHub Copilot).

## File Patterns

ARM uses glob patterns to filter which files are installed from packages.

**Default behavior**: Include all `.yml` and `.yaml` files

**Pattern matching**:
- Archives (`.tar.gz`, `.zip`) are extracted first, then patterns filter the contents
- Patterns use standard glob syntax: `**/*.yml`, `security/**/*.md`, `**/typescript-*`
- Multiple `--include` patterns are OR'd together (match any)
- `--exclude` patterns override includes

**Examples**:
```bash
# Install only TypeScript rules
arm install ruleset --include "**/typescript-*.yml" ai-rules/language-rules cursor-rules

# Exclude experimental files
arm install ruleset --exclude "**/experimental/**" ai-rules/security-ruleset cursor-rules

# Multiple includes
arm install promptset --include "review/**/*.yml" --include "refactor/**/*.yml" ai-rules/prompts cursor-commands
```

## Versioning

ARM supports two versioning models:

**Git Registries** - Flexible versioning using Git's native features:
- Semantic version tags (`1.0.0`, `v2.1.0`) for production releases (always prioritized)
- Branches (`main`, `develop`) for development and testing (resolves to commit hash)
- Version resolution: Semver tags always take precedence over branches
- Branch priority: Determined by order in registry `--branches` configuration
- Version constraints: `@1` (major), `@1.1` (major.minor), `@1.0.0` (exact)

**GitLab/Cloudsmith Registries** - Semantic versioning only:
- Only semantic versions (`1.0.0`, `2.1.3`)
- No branches or commit hashes
- Version constraints: `@1` (major), `@1.1` (major.minor), `@1.0.0` (exact)

## How to Install Packages

1. Add registries where packages are stored
2. Configure sinks with your desired AI tool and target directory
3. Install packages from registries as rulesets or promptsets to sinks

## Environment Variables

ARM supports environment variables to customize file locations for testing, CI/CD, and multi-user environments.

### ARM_MANIFEST_PATH

Controls the location of `arm.json` and `arm-lock.json`.

```bash
ARM_MANIFEST_PATH=/custom/path/arm.json
# Results in:
# - /custom/path/arm.json (manifest)
# - /custom/path/arm-lock.json (lock file, colocated)
```

**Default:** `./arm.json` and `./arm-lock.json` in current working directory

**Use cases:**
- Testing with isolated manifests
- CI/CD with build-specific configurations
- Managing multiple ARM projects in same directory

### ARM_CONFIG_PATH

Overrides the `.armrc` configuration file location. When set, this is the ONLY config file used (no hierarchical lookup).

```bash
ARM_CONFIG_PATH=/custom/path/.armrc
# Results in:
# - Only reads /custom/path/.armrc
# - Ignores both ./.armrc and ~/.armrc
```

**Default:** Hierarchical lookup (`./.armrc` overrides `~/.armrc`)

**Use cases:**
- CI/CD with centralized authentication
- Testing with isolated credentials
- Custom config locations in containerized environments

### ARM_HOME

Overrides the home directory for the `.arm/` directory (storage, cache). Does NOT affect `.armrc` location.

```bash
ARM_HOME=/custom/home
# Results in:
# - /custom/home/.arm/storage/ (package cache)
```

**Default:** User's home directory from `os.UserHomeDir()`

**Use cases:**
- Docker with mounted volumes
- Multi-user systems with separate caches
- Network storage for shared team caches
- CI/CD with build-specific cache directories

### Priority Order

**For .armrc lookup:**
1. `ARM_CONFIG_PATH` - If set, use this exact file (bypasses hierarchy)
2. `./armrc` - Project config (highest priority in hierarchy)
3. `~/.armrc` - User config (fallback in hierarchy)

**For .arm/storage/ lookup:**
1. `$ARM_HOME/.arm/storage/` - If ARM_HOME is set
2. `~/.arm/storage/` - Default

## How to Publish Packages

**ARM resource files** (YAML) - Cross-platform definitions that ARM compiles to any tool format.

**Raw tool files** - Tool-specific files that are used as-is without compilation.

This approach enables backwards compatibility with existing repositories like awesome-cursorrules.

See [Resource Schemas](resource-schemas.md) for detailed YAML format specifications.