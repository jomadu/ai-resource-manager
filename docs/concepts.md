# Concepts

## Core Components

**Registries** are remote sources where rulesets and promptsets are stored and versioned (Git repositories, GitLab Package Registry, Cloudsmith).

**Packages** are versioned collections of AI rules or prompts:
- **Rulesets** - Collections of AI rules with priority-based conflict resolution
- **Promptsets** - Collections of AI prompts for reusable templates

**Sinks** are local directories where ARM places compiled files, each configured for a specific AI tool (Cursor, Amazon Q, GitHub Copilot).

## How It Works

1. Add registries where packages are stored
2. Configure sinks with your desired AI tool and target directory
3. Install packages from registries to sinks
4. ARM compiles resource files (when needed) and places the correct files for your tool in the project

## Package Author Options

**ARM resource files** (YAML) - Cross-platform definitions that ARM compiles to any tool format.

**Raw tool files** - Tool-specific files that are used as-is without compilation.

This approach enables backwards compatibility with existing repositories like awesome-cursorrules.

## ARM Project Files

- `arm.json` - Project manifest containing registries, packages, and sinks
- `arm-lock.json` - Locked versions for reproducible installs
- `arm-index.json` - Local inventory tracking installed packages and file paths
- `arm_index.*` - Generated priority rules that help AI tools resolve conflicts

See [Resource Schemas](resource-schemas.md) for detailed YAML format specifications.