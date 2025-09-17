# Git Registry

Git registries use Git repositories (GitHub, GitLab, or any Git remote) to store and distribute AI rules using Git tags and branches.

## Configuration

Add a Git registry:

```bash
arm config registry add ai-rules https://github.com/jomadu/ai-rules-manager-sample-git-registry --type git
```

With specific branches:
```bash
arm config registry add ai-rules https://github.com/myorg/rules --type git --branches main,develop
```

## Authentication

Git registries use Git's built-in authentication system:

- **SSH**: Use SSH keys configured with `ssh-agent` or `~/.ssh/config`
- **HTTPS**: Use Git credential helpers or environment variables
- **GitHub CLI**: Use `gh auth login` for automatic GitHub authentication
- **GitLab CLI**: Use `glab auth login` for automatic GitLab authentication

No additional configuration in `.armrc` is required - ARM uses the same authentication as your Git commands.

## Repository Structure

Git repositories should follow the recommended structure:

```
Repository Root:
├── clean-code.yml              # URF rulesets
├── security.yml
├── performance.yml
└── build/                      # Pre-compiled rules (optional)
    ├── cursor/
    │   ├── clean-code.mdc
    │   └── security.mdc
    └── amazonq/
        ├── clean-code.md
        └── security.md
```

## Installing Rules

Install from latest tag:
```bash
arm install ai-rules/rules --sinks cursor,amazonq
```

Install from specific tag:
```bash
arm install ai-rules/rules@v1.0.0 --sinks cursor
```

Install from branch:
```bash
arm install ai-rules/rules@main --sinks cursor
```

Install specific files:
```bash
arm install ai-rules/rules --include "security.yml" --sinks cursor
```

Install pre-compiled rules:
```bash
arm install ai-rules/rules --include "build/cursor/**" --sinks cursor
```

## Version Resolution

Git registries support multiple version types:

### Tags (Semantic Versions)
- **Semantic versions**: `v1.0.0`, `1.2.3`
- **Version constraints**: `^1.0.0`, `~1.1.0`
- **Latest**: `latest` (resolves to highest semantic version tag)

### Branches
- **Branch names**: `main`, `develop`, `feature/new-rules`
- Resolves to the HEAD commit of the branch
- Displayed as short commit hash (first 7 characters)

### Priority Order
1. Semantic version tags (sorted descending)
2. Non-semantic tags
3. Branch commits (in configuration order)

## Publishing Rules

Simply push tags to your Git repository:

```bash
git tag v1.0.0
git push origin v1.0.0
```

ARM automatically discovers new tags and makes them available for installation.
