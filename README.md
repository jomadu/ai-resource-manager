![ARM Header](assets/header.png)

# AI Resource Manager (ARM)

## What is ARM?

ARM is a package manager for AI resources, designed to treat rulesets and promptsets as code dependencies. It introduces semantic versioning, reproducible installs, and straightforward distribution to your AI tools.

Seamlessly connect to Git repositories such as [awesome-cursorrules](https://github.com/PatrickJS/awesome-cursorrules/tree/main) or your team's private collections. Install and manage versioned resources across projects, and keep everything in sync with your source of truth.

## Why ARM?

Managing rules and prompts for AI coding assistants like Cursor or Amazon Q is cumbersome:

- **Manual duplication**: Copying resources disconnects them from updates and the original source
- **Hidden breaking changes**: Updates may unexpectedly alter your AI's behavior
- **Poor scalability**: Coordinating resources across multiple projects becomes chaotic
- **Incompatible formats**: Frequent manual conversions between different tool formats. ([Obligatory we need a new standard xkcd link](https://xkcd.com/927/))

ARM solves these problems with a modern package manager approach.

### Key Features of ARM

- **Consistent, versioned installs** using semantic versioning (except for [git based registry](docs/registries/git-registry.md) without semver tags, which gets a little funky)
- **Reliable, reproducible environments** through manifest and lock files (similar to npm's `package.json` and `package-lock.json`)
- **[Unified resource definitions](https://xkcd.com/927/)** that compile to formats needed by any AI tool (the audacity! *clutches pearls*)
- **Priority-based rule composition** for layering multiple rulesets with clear conflict resolution (your team's standards > internet best practices)
- **Flexible registry support** for managing resources from Git, GitLab, and Cloudsmith
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

Configure sinks:

**Cursor:**
```bash
arm add sink --type cursor cursor-rules .cursor/rules
arm add sink --type cursor cursor-commands .cursor/commands
```

**GitHub Copilot:**
```bash
arm add sink --type copilot copilot-rules .github/copilot
```

**Amazon Q:**
```bash
arm add sink --type amazonq q-rules .amazonq/rules
arm add sink --type amazonq q-prompts .amazonq/prompts
```

Install ruleset:
```bash
arm install ruleset ai-rules/clean-code-ruleset cursor-rules
```

Install promptset:
```bash
arm install promptset ai-rules/code-review-promptset cursor-commands
```

Install to multiple sinks:
```bash
arm install ruleset ai-rules/clean-code-ruleset cursor-rules q-rules
arm install promptset ai-rules/code-review-promptset cursor-commands q-prompts
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

## Documentation

- **[Concepts](docs/concepts.md)** - Core concepts and files
- **[Rulesets](docs/rulesets.md)** - Working with rulesets
- **[Promptsets](docs/promptsets.md)** - Working with promptsets
- **[Registries](docs/registries.md)** - Registry management and types
- **[Sinks](docs/sinks.md)** - Sink configuration and compilation
- **[Package Management](docs/package-management.md)** - Installing, updating, and managing packages

### Registry Types

- **[Git Registry](docs/registries/git-registry.md)** - GitHub, GitLab, and Git remotes
- **[GitLab Registry](docs/registries/gitlab-registry.md)** - GitLab Package Registry
- **[Cloudsmith Registry](docs/registries/cloudsmith-registry.md)** - Cloudsmith package repository
