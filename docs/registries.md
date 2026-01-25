# Registries

Registries are remote sources where packages are stored and versioned, similar to npm registries. ARM supports:

- **Git registries**: GitHub repositories, GitLab projects, or any Git remote
- **GitLab Package registries**: GitLab's Generic Package Registry for versioned packages
- **Cloudsmith registries**: Cloudsmith's package repository service for single-file artifacts

## Commands

For detailed command usage and examples, see [Registry Management](commands.md#registry-management) in the commands reference.

## Registry Structure Models

ARM supports two fundamentally different package naming models:

### Git Registries (User-Named Packages)

**How it works**: The repository is a source of files. You choose the package name when installing.

**Repository structure:**
```
awesome-cursorrules/         # Repository (not the package name)
├── clean-code-ruleset.yml   # ARM resource definitions
├── security-ruleset.yml
├── code-review-promptset.yml
├── rules.tar.gz             # Optional: archived rules
└── build/                   # Pre-compiled rules (optional)
    ├── cursor/
    │   ├── clean-code.mdc
    │   └── security.mdc
    └── amazonq/
        ├── clean-code.md
        └── security.md
```

**Install examples**: 
```bash
# You choose the package name - can be anything
arm install ruleset git-registry/my-rules cursor-rules
arm install ruleset git-registry/team-standards q-rules
arm install ruleset git-registry/awesome-cursorrules cursor-rules

# Use --include/--exclude to filter files (default: *.yml, *.yaml)
arm install ruleset --include "security-*.yml" git-registry/security-only cursor-rules
```

**Key point**: The package name (after the `/`) is whatever you want - it's just a label for your installation.

### Non-Git Registries (Registry-Named Packages)

**How it works**: Packages have explicit names in the registry. You must use the exact package name when installing.

**Package structure:**
```
Package: clean-code-ruleset (defined in registry)
├── clean-code-ruleset.yml   # ARM resource definition
└── build/                   # Optional pre-compiled rules
    ├── cursor/
    │   └── clean-code.mdc
    └── amazonq/
        └── clean-code.md

Package: security-ruleset (defined in registry)
├── security-ruleset.yml     # ARM resource definition
└── build/
    ├── cursor/
    │   └── security.mdc
    └── amazonq/
        └── security.md
```

**Install examples**:
```bash
# Must use exact package names from registry
arm install ruleset gitlab-registry/clean-code-ruleset cursor-rules
arm install ruleset gitlab-registry/security-ruleset q-rules

# Use --include/--exclude to filter files (default: *.yml, *.yaml)
arm install ruleset --include "**/*.yml" --exclude "**/experimental/**" gitlab-registry/clean-code-ruleset cursor-rules
```

**Key point**: The package name (after the `/`) must match the exact name published in the registry.

## Archive Support

ARM automatically extracts and processes **zip** and **tar.gz** archives during installation:

- **Supported formats**: `.zip` and `.tar.gz` files
- **Automatic extraction**: Archives are detected by extension and extracted transparently
- **Merge behavior**: Extracted files are merged with loose files, with archives taking precedence in case of path conflicts
- **Security**: Path sanitization prevents directory traversal attacks
- **Pattern filtering**: `--include` patterns are applied to the merged content after extraction

## Registry Types

### Git Registry
Uses Git repositories with tags and branches for versioning. See [Git Registry](./git-registry.md) for details.

### GitLab Registry
Uses GitLab's Generic Package Registry for versioned packages. See [GitLab Registry](./gitlab-registry.md) for details.

### Cloudsmith Registry
Uses Cloudsmith's package repository service for single-file artifacts. See [Cloudsmith Registry](./cloudsmith-registry.md) for details.

## Examples

**Community registries:**
- [PatrickJS/awesome-cursorrules](https://github.com/PatrickJS/awesome-cursorrules) - Community collection of Cursor rules
- [snarktank/ai-dev-tasks](https://github.com/snarktank/ai-dev-tasks) - AI development task templates
- [steipete/agent-rules](https://github.com/steipete/agent-rules) - Agent configuration rules
