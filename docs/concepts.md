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

**Sinks** are local directories where ARM places compiled files, each configured for a specific AI tool (Cursor, Amazon Q, GitHub Copilot).

## Versioning

ARM supports two versioning models:

**Git Registries** - Flexible versioning using Git's native features:
- Semantic version tags (`1.0.0`, `v2.1.0`) for production releases
- Branches (`main`, `develop`) for development and testing (resolves to commit hash)
- Version constraints: `@1` (major), `@1.1` (major.minor), `@1.0.0` (exact)

**GitLab/Cloudsmith Registries** - Semantic versioning only:
- Only semantic versions (`1.0.0`, `2.1.3`)
- No branches or commit hashes
- Version constraints: `@1` (major), `@1.1` (major.minor), `@1.0.0` (exact)

## How to Install Packages

1. Add registries where packages are stored
2. Configure sinks with your desired AI tool and target directory
3. Install packages from registries as rulesets or promptsets to sinks

## How to Publish Packages

**ARM resource files** (YAML) - Cross-platform definitions that ARM compiles to any tool format.

**Raw tool files** - Tool-specific files that are used as-is without compilation.

This approach enables backwards compatibility with existing repositories like awesome-cursorrules.

See [Resource Schemas](resource-schemas.md) for detailed YAML format specifications.