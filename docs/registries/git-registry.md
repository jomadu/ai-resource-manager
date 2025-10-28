# Git Registry

Git registries use Git repositories (GitHub, GitLab, or any Git remote) to store and distribute AI rules using Git tags and branches.

## Configuration

Add a Git registry:

```bash
arm add registry git --url https://github.com/jomadu/ai-rules-manager-sample-git-registry ai-rules
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

### Archive Support

Git registries automatically extract and process archives during installation:

- **Supported formats**: `.zip` and `.tar.gz` files
- **Automatic processing**: Archives are detected and extracted transparently
- **Merge behavior**: Extracted files are merged with loose files, with archives taking precedence
- **Security**: Path sanitization prevents directory traversal attacks

## Installing Rulesets

Install from latest version:
```bash
arm install ruleset ai-rules/clean-code-ruleset cursor-rules
```

Install from specific version:
```bash
arm install ruleset ai-rules/clean-code-ruleset@v1.0.0 cursor-rules
```

Install to multiple sinks:
```bash
arm install ruleset ai-rules/clean-code-ruleset cursor-rules q-rules
```

Install with custom priority:
```bash
arm install ruleset --priority 200 ai-rules/clean-code-ruleset cursor-rules
```

Install specific files:
```bash
arm install ruleset --include "security.yml" ai-rules/clean-code-ruleset cursor-rules
```

Install pre-compiled rules:
```bash
arm install ruleset --include "build/cursor/**" ai-rules/clean-code-ruleset cursor-rules
```

Install from archives:
```bash
# Install all files from archives and loose files
arm install ruleset ai-rules/clean-code-ruleset cursor-rules

# Install only ruleset files (from both archives and loose files)
arm install ruleset --include "**/*.yml" ai-rules/clean-code-ruleset cursor-rules

# Install specific archive
arm install ruleset --include "rules.tar.gz" ai-rules/clean-code-ruleset cursor-rules
```

## Installing Promptsets

Install from latest version:
```bash
arm install promptset ai-rules/code-review-promptset cursor-commands
```

Install from specific version:
```bash
arm install promptset ai-rules/code-review-promptset@v1.0.0 cursor-commands
```

Install to multiple sinks:
```bash
arm install promptset ai-rules/code-review-promptset cursor-commands q-prompts
```

Install specific files:
```bash
arm install promptset --include "review.yml" ai-rules/code-review-promptset cursor-commands
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

## Management Commands

List all registries:
```bash
arm list registry
```

Show registry information:
```bash
arm info registry ai-rules
```

Update registry configuration:
```bash
arm set registry ai-rules url https://github.com/myorg/new-rules-repo
```

Remove registry:
```bash
arm remove registry ai-rules
```
