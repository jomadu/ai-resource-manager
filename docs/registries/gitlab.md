# GitLab Registry

GitLab registries use GitLab's Generic Package Registry to store and distribute AI rules as versioned packages.

## Configuration

Add a GitLab registry:

```bash
arm config registry add my-gitlab https://gitlab.example.com --type gitlab --project-id 123
```

Or for group-level packages:

```bash
arm config registry add my-gitlab https://gitlab.example.com --type gitlab --group-id 456
```

## Authentication

GitLab registries require explicit token authentication in `.armrc`. See the [.armrc documentation](../armrc.md) for complete details.

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
   arm install my-gitlab/ai-rules --sinks cursor
   ```

**Note**: Unlike Git registries which use Git's built-in authentication, GitLab registries require explicit token configuration because they use GitLab's Package Registry API.

## Package Structure

GitLab packages use the standardized "ai-rules" name with semantic versioning:

```
Package Name: ai-rules
Version: 1.0.0, 1.1.0, 2.0.0, etc.
Files:
├── clean-code.yml              # URF rulesets
├── security.yml
├── rules.tar.gz                # Optional: archived rules
├── legacy-rules.zip            # Optional: archived rules
└── build/                      # Pre-compiled rules
    ├── cursor/
    │   ├── clean-code.mdc
    │   └── security.mdc
    └── amazonq/
        ├── clean-code.md
        └── security.md
```

### Archive Support

GitLab packages automatically extract and process archives during installation:

- **Supported formats**: `.zip` and `.tar.gz` files
- **Automatic processing**: Archives are detected and extracted transparently
- **Merge behavior**: Extracted files are merged with loose files, with archives taking precedence
- **Security**: Path sanitization prevents directory traversal attacks

## Installing Rules

Install URF rulesets (default):
```bash
arm install my-gitlab/ai-rules --sinks cursor,amazonq
```

Install specific version:
```bash
arm install my-gitlab/ai-rules@1.0.0 --sinks cursor
```

Install pre-compiled rules:
```bash
arm install my-gitlab/ai-rules --include "build/cursor/**" --sinks cursor
```

Install from archives:
```bash
# Install all files from archives and loose files
arm install my-gitlab/ai-rules --sinks cursor

# Install only YAML files (from both archives and loose files)
arm install my-gitlab/ai-rules --include "**/*.yml" --sinks cursor

# Install specific archive
arm install my-gitlab/ai-rules --include "rules.tar.gz" --sinks cursor
```

## Publishing Packages

Use GitLab CI/CD to publish packages:

```yaml
publish:
  script:
    - |
      curl --header "JOB-TOKEN: $CI_JOB_TOKEN" \
           --upload-file clean-code.yml \
           "${CI_API_V4_URL}/projects/${CI_PROJECT_ID}/packages/generic/ai-rules/1.0.0/clean-code.yml"
```

## Version Resolution

- **Semantic versions**: `1.0.0`, `^1.0.0`, `~1.1.0`
- **Latest**: `latest` (resolves to highest semantic version)
- Versions are sorted by semantic version in descending order
