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

GitLab registries require explicit token authentication in `.armrc`:

```ini
[registry gitlab.example.com/project/123]
token = ${GITLAB_TOKEN}
```

Or for group-level packages:

```ini
[registry gitlab.example.com/group/456]
token = ${GITLAB_TOKEN}
```

Set your GitLab token:
```bash
export GITLAB_TOKEN=glpat-xxxxxxxxxxxxxxxxxxxx
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
└── build/                      # Pre-compiled rules
    ├── cursor/
    │   ├── clean-code.mdc
    │   └── security.mdc
    └── amazonq/
        ├── clean-code.md
        └── security.md
```

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
