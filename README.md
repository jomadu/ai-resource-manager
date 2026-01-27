![ARM Header](assets/header.png)

# AI Resource Manager (ARM)

## What is ARM?

ARM is a dependency manager for AI packages, designed to treat them as code dependencies. You can install packages as rulesets, promptsets, or other resource types. ARM introduces semantic versioning, reproducible installs, and straightforward distribution to your AI tools.

Seamlessly connect to Git repositories such as [awesome-cursorrules](https://github.com/PatrickJS/awesome-cursorrules/tree/main) or your team's private collections. Install and manage versioned AI resources across projects, and keep everything in sync with your source of truth.

## Why ARM?

Managing AI resources for coding assistants is cumbersome:

- **Manual duplication**: Copying resources disconnects them from updates and the original source
- **Hidden breaking changes**: Updates may unexpectedly alter your AI's behavior
- **Poor scalability**: Coordinating AI resources across multiple projects becomes chaotic
- **Incompatible formats**: Frequent manual conversions between different tool formats. ([Obligatory we need a new standard xkcd link](https://xkcd.com/927/))

ARM solves these problems with a modern dependency manager approach.

### Key Features of ARM

- **Consistent, versioned installs** using semantic versioning (except for [git based registry](docs/git-registry.md) without semver tags, which gets a little funky)
- **Reliable, reproducible environments** through manifest and lock files (similar to npm's `package.json` and `package-lock.json`)
- **Backwards compatibility** with existing repositories like [awesome-cursorrules](https://github.com/PatrickJS/awesome-cursorrules) - use them immediately without conversion
- **Write once, deploy everywhere** - compile ARM resource schemas to any tool format (Cursor, Copilot, Amazon Q, Markdown)
- **Priority-based rule composition** for layering multiple rulesets with clear conflict resolution (your team's standards > internet best practices)
- **Flexible registry support** for managing AI resources from Git, GitLab, and Cloudsmith
- **Automated update workflow:** easily check for updates and apply them across projects ([nice](https://media1.giphy.com/media/v1.Y2lkPTc5MGI3NjExNTlycTA1ejFwdnZtZHNzOG5tYnVwajF3bDAwYzllcnU1dm5oNWplMCZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/NEvPzZ8bd1V4Y/giphy.gif))

## Installation

### Quick Install

```bash
curl -fsSL https://raw.githubusercontent.com/jomadu/ai-resource-manager/main/scripts/install.sh | bash
```

### Install Specific Version

```bash
curl -fsSL https://raw.githubusercontent.com/jomadu/ai-resource-manager/main/scripts/install.sh | bash -s v1.0.0
```

### Manual Installation

1. Download the latest release from [GitHub](https://github.com/jomadu/ai-resource-manager/releases)
2. Extract and move the binary to your PATH
3. Run `arm help` to verify installation

### Verify Installation

```bash
arm version
arm help
```

> **Upgrading from v2?** See the [migration guide](docs/migration-v2-to-v3.md) for breaking changes and upgrade steps. TL;DR: Sorry, nuke and pave. We made some poor design choices in v1 and v2. Honestly, we've probably made them in v3 too, but hey, better is the enemy of good.

## Uninstall

```bash
curl -fsSL https://raw.githubusercontent.com/jomadu/ai-resource-manager/main/scripts/uninstall.sh | bash
```

## Quick Start

### Resource Types

Packages can be installed as:

- **Rulesets** - Collections of AI rules with priority-based conflict resolution
- **Promptsets** - Collections of AI prompts for reusable templates

### Setup

Add Git registry:
```bash
arm add registry git --url https://github.com/jomadu/ai-rules-manager-sample-git-registry ai-rules
```

Add GitLab registry:
```bash
arm add registry gitlab --url https://gitlab.example.com --project-id 123 my-gitlab
```

Add Cloudsmith registry:
```bash
arm add registry cloudsmith --owner myorg --repo ai-rules my-cloudsmith
```

Configure sinks (output destinations):

```bash
# Cursor
arm add sink --tool cursor cursor-rules .cursor/rules
arm add sink --tool cursor cursor-commands .cursor/commands

# GitHub Copilot
arm add sink --tool copilot copilot-rules .github/copilot

# Amazon Q
arm add sink --tool amazonq q-rules .amazonq/rules
```

Install ruleset:
```bash
arm install ruleset ai-rules/clean-code-ruleset cursor-rules
```

Install promptset:
```bash
arm install promptset ai-rules/code-review-promptset cursor-commands
```

Install to multiple tools:
```bash
arm install ruleset ai-rules/clean-code-ruleset cursor-rules copilot-rules q-rules
```

Install with priority (higher priority rules take precedence):
```bash
# Install team rules with high priority
arm install ruleset --priority 200 ai-rules/team-standards cursor-rules

# Install general rules with default priority (100)
arm install ruleset ai-rules/clean-code-ruleset cursor-rules
```

Install specific files from a ruleset (useful for git-based registries):
```bash
# Only install TypeScript-related rules
arm install ruleset --include "**/typescript-*.yml" ai-rules/language-rules cursor-rules

# Install security rules but exclude experimental ones
arm install ruleset --include "security/**/*.yml" --exclude "**/experimental/**" ai-rules/security-ruleset cursor-rules

# Install only specific prompt files
arm install promptset --include "review/**/*.yml" --include "refactor/**/*.yml" ai-rules/code-review-promptset cursor-commands
```

Compile local ARM resource files to tool-specific formats:
```bash
# Compile a single ruleset file
arm compile ruleset my-rules.yml --tool cursor --output .cursor/rules/

# Compile all YAML files in a directory
arm compile ruleset ./rules/ --tool copilot --output .github/copilot/ --recursive

# Compile with namespace for organization
arm compile ruleset ./rules/ --tool cursor --output .cursor/rules/ --namespace my-team
```

## Documentation

### Getting Started

- **[Concepts](docs/concepts.md)** - Core concepts, file types, and resource definitions
- **[Commands](docs/commands.md)** - Complete command reference and usage examples

### Using ARM

- **[Registries](docs/registries.md)** - Registry management and types
- **[Sinks](docs/sinks.md)** - Sink configuration and compilation
- **[Resource Schemas](docs/resource-schemas.md)** - ARM resource YAML schemas

### Publishing Resources

- **[Publishing Guide](docs/publishing-guide.md)** - How to create and publish rulesets to Git repositories

### Registry Types

- **[Git Registry](docs/git-registry.md)** - GitHub, GitLab, and Git remotes
- **[GitLab Registry](docs/gitlab-registry.md)** - GitLab Package Registry
- **[Cloudsmith Registry](docs/cloudsmith-registry.md)** - Cloudsmith package repository

Buy Me a Coffee: buymeacoffee.com/max.dunn