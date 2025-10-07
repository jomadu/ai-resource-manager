![ARM Header](assets/header.png)

# AI Rules Manager (ARM)

## What is ARM?

A package manager for AI rules that treats rulesets and promptsets like code dependencies - with semantic versioning, reproducible installs, and automatic distribution to your AI tools.

Connect to Git repositories like awesome-cursorrules or your team's rule collections, install versioned resources across projects, and keep them automatically synced with their source of truth.

## Why ARM?

AI coding assistants like Cursor and Amazon Q rely on rules and prompts to guide their behavior, but managing these resources is broken:

- **Manual copying** severs the connection to the source of truth - once copied, resources are orphaned with no way to get updates
- **Breaking changes blindness** - when you pull latest resources, you have no idea if they'll break your AI's behavior
- **Doesn't scale** - managing resources across even 3 projects becomes unmanageable overhead
- **Format fragmentation** - each AI tool uses different formats, requiring manual conversion

ARM solves these problems with a **package manager approach** - semantic versioning, reproducible installs, and automatic distribution keep resources connected to their source of truth. ARM resource definitions provide a unified format that compiles to platform-specific outputs.

## Installation

### Quick Install

```bash
curl -fsSL https://raw.githubusercontent.com/jomadu/ai-rules-manager/main/scripts/install.sh | bash
```

### Install Specific Version

```bash
curl -fsSL https://raw.githubusercontent.com/jomadu/ai-rules-manager/main/scripts/install.sh | bash -s v1.0.0
```

### Manual Installation

1. Download the latest release from [GitHub](https://github.com/jomadu/ai-rules-manager/releases)
2. Extract and move the binary to your PATH
3. Run `arm help` to verify installation

### Verify Installation

```bash
arm version
arm help
```

> **Upgrading from v2?** See the [migration guide](docs/migration-v2-to-v3.md) for breaking changes and upgrade steps.

## Uninstall

```bash
curl -fsSL https://raw.githubusercontent.com/jomadu/ai-rules-manager/main/scripts/uninstall.sh | bash
```

## Quick Start

Add Git registry:
```bash
arm add registry --type git ai-rules https://github.com/jomadu/ai-rules-manager-sample-git-registry
```

Add GitLab registry:
```bash
arm add registry --type gitlab --project-id 123 my-gitlab https://gitlab.example.com
```

Add Cloudsmith registry:
```bash
arm add registry --type cloudsmith --owner myorg --repo ai-rules my-cloudsmith https://app.cloudsmith.com
```

Configure sinks:

**Cursor:**
```bash
arm add sink --type cursor cursor-rules .cursor/rules
arm add sink --type cursor cursor-prompts .cursor/prompts
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
arm install promptset ai-rules/code-review-promptset cursor-prompts
```

Install to multiple sinks:
```bash
arm install ruleset ai-rules/clean-code-ruleset cursor-rules q-rules
arm install promptset ai-rules/code-review-promptset cursor-prompts q-prompts
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
