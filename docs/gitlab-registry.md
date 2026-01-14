# GitLab Registry

GitLab registries use GitLab's Generic Package Registry to store and distribute packages as versioned packages.

## Authentication

GitLab registries require explicit token authentication in `.armrc`. See the [.armrc documentation](armrc.md) for complete details.

**Quick Setup:**

1. Create `.armrc` file:
   ```ini
   [registry https://gitlab.example.com/project/123]
   token = glpat-xxxxxxxxxxxxxxxxxxxx
   ```

2. Set file permissions:
   ```bash
   chmod 600 .armrc
   ```

3. Test installation:
   ```bash
   arm install ruleset my-gitlab/clean-code-ruleset cursor-rules
   ```

**Note**: Unlike Git registries which use Git's built-in authentication, GitLab registries require explicit token configuration because they use GitLab's Package Registry API.

## Package Structure

**Key Concept**: GitLab packages have explicit names defined in the registry. You must use the exact package name when installing.

```
Package: clean-code-ruleset (exact name in registry)
├── clean-code-ruleset.yml      # ARM resource definition
├── rules.tar.gz                # Optional: archived rules
└── build/                      # Pre-compiled rules (optional)
    ├── cursor/
    │   └── clean-code.mdc
    └── amazonq/
        └── clean-code.md

Package: security-ruleset (exact name in registry)
├── security-ruleset.yml        # ARM resource definition
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

### Archive Support

GitLab packages automatically extract and process archives during installation:

- **Supported formats**: `.zip` and `.tar.gz` files
- **Automatic processing**: Archives are detected and extracted transparently
- **Merge behavior**: Extracted files are merged with loose files, with archives taking precedence
- **Security**: Path sanitization prevents directory traversal attacks

## Version Resolution

GitLab registries support semantic versioning:

- **Semantic versions**: `1.0.0`, `^1.0.0`, `~1.1.0`
- **Latest**: `latest` (resolves to highest semantic version)
- Versions are sorted by semantic version in descending order

## Publishing Packages

Use GitLab CI/CD to publish packages:

```yaml
publish:
  script:
    - |
      curl --header "JOB-TOKEN: $CI_JOB_TOKEN" \
           --upload-file clean-code-ruleset.yml \
           "${CI_API_V4_URL}/projects/${CI_PROJECT_ID}/packages/generic/clean-code-ruleset/${SEMVER_STRING}/clean-code-ruleset.yml"
```