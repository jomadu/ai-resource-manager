# Git Registry

Git registries use Git repositories (GitHub, GitLab, or any Git remote) to store and distribute AI rulesets and promptsets using Git tags and branches.

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