# Git Registry

Git registries use Git repositories (GitHub, GitLab, or any Git remote) to store and distribute packages using Git tags and branches.

## Authentication

Git registries use Git's built-in authentication system:

- **SSH**: Use SSH keys configured with `ssh-agent` or `~/.ssh/config`
- **HTTPS**: Use Git credential helpers or environment variables
- **GitHub CLI**: Use `gh auth login` for automatic GitHub authentication
- **GitLab CLI**: Use `glab auth login` for automatic GitLab authentication

No additional configuration in `.armrc` is required - ARM uses the same authentication as your Git commands.

## Repository Structure

**Key Concept**: In Git registries, you choose the package name when installing. The repository is just the source of files.

```
Repository (source of files):
├── clean-code-ruleset.yml      # ARM resource definitions
├── security-ruleset.yml
├── code-review-promptset.yml   # ARM resource definitions
├── rules.tar.gz                # Optional: archived rules
├── legacy-rules.zip            # Optional: archived rules
└── build/                      # Pre-compiled rules (optional)
    ├── cursor/
    │   ├── clean-code.mdc
    │   └── security.mdc
    └── amazonq/
        ├── clean-code.md
        └── security.md
```

**Install examples**: 
```bash
# You choose any package name
arm install ruleset git-registry/my-rules cursor-rules
arm install ruleset git-registry/team-standards q-rules

# Use --include/--exclude to filter files (default: *.yml, *.yaml)
arm install ruleset --include "security-*.yml" git-registry/security-only cursor-rules
```

### Archive Support

Git registries automatically extract and process archives during installation:

- **Supported formats**: `.zip` and `.tar.gz` files
- **Automatic processing**: Archives are detected and extracted transparently
- **Merge behavior**: Extracted files are merged with loose files, with archives taking precedence
- **Security**: Path sanitization prevents directory traversal attacks

## Version Resolution

Git registries support two versioning approaches:

### Semantic Version Tags

Preferred for production releases:

```bash
# Install specific semantic version
arm install ruleset git-registry/my-rules@1.0.0 cursor-rules

# Install with version constraints
arm install ruleset git-registry/my-rules@1 cursor-rules    # >= 1.0.0, < 2.0.0
arm install ruleset git-registry/my-rules@1.1 cursor-rules  # >= 1.1.0, < 1.2.0

# Install latest semantic version
arm install ruleset git-registry/my-rules@latest cursor-rules
```

**In arm-lock.json**: `"version": "1.0.0"`

### Branches

Useful for development and testing:

```bash
# Install from branch
arm install ruleset git-registry/my-rules@main cursor-rules
arm install ruleset git-registry/my-rules@develop cursor-rules
arm install ruleset git-registry/my-rules@feature/new-rules cursor-rules
```

**In arm-lock.json**: `"version": "a1b2c3d"` (7-character commit hash)

### Version Resolution Priority

When you specify `@latest` or no version, ARM resolves in this order:

1. **Semantic version tags** (highest version, e.g., `2.1.0` > `1.9.0`)
2. **Branch HEAD commits** (shown as short commit hash)

### Examples

**Repository with semantic tags**:
```bash
# Tags: v1.0.0, v1.1.0, v2.0.0
arm install ruleset git-registry/my-rules cursor-rules
# Resolves to: 2.0.0
```

**Repository with only branches**:
```bash
# Branches: main, develop
arm install ruleset git-registry/my-rules cursor-rules
# Resolves to: a1b2c3d (main branch HEAD)
```

**Repository with tags and branches**:
```bash
# Tags: v1.0.0, v2.0.0
# Branches: main
arm install ruleset git-registry/my-rules cursor-rules
# Resolves to: 2.0.0 (semantic tag wins)
```