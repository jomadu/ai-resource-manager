# Registry Management

## Job to be Done
Configure and manage remote package registries (Git, GitLab, Cloudsmith) for discovering and fetching versioned AI resource packages.

## Activities
1. Add Git registry (GitHub, GitLab, any Git remote)
2. Add GitLab Package Registry
3. Add Cloudsmith registry
4. List packages and versions from registries
5. Fetch package files with integrity verification

## Acceptance Criteria
- [x] Add Git registry with URL and optional branches
- [x] Add GitLab registry with project ID or group ID
- [x] Add Cloudsmith registry with owner and repository
- [x] List packages from registry
- [x] List versions for a package
- [x] Fetch package files by version
- [x] Calculate SHA256 integrity hash for packages
- [x] Support authentication via .armrc
- [x] Cache registry metadata and packages

## Data Structures

### Git Registry Config
```json
{
  "type": "git",
  "url": "https://github.com/org/repo",
  "branches": ["main", "develop"]
}
```

### GitLab Registry Config
```json
{
  "type": "gitlab",
  "url": "https://gitlab.com",
  "projectId": "123",
  "groupId": "",
  "apiVersion": "v4"
}
```

### Cloudsmith Registry Config
```json
{
  "type": "cloudsmith",
  "url": "https://api.cloudsmith.io",
  "owner": "myorg",
  "repository": "ai-rules"
}
```

## Algorithm

### Git Registry - List Packages
1. Clone or fetch repository to cache
2. List all .yml and .yaml files
3. Extract package names from file paths
4. Return unique package names

### Git Registry - List Versions
1. List all tags (semantic versions)
2. List configured branches
3. Sort versions (semver tags first, then branches)
4. Return version list

### Git Registry - Fetch Package
1. Resolve version to commit hash
2. Checkout commit in cached repository
3. Collect files matching package name and patterns
4. Extract archives if present
5. Calculate SHA256 integrity hash
6. Return files and metadata

### GitLab Registry - List Packages
1. Call GitLab API: GET /projects/{id}/packages
2. Filter for generic packages
3. Return package names

### GitLab Registry - List Versions
1. Call GitLab API: GET /projects/{id}/packages?package_name={name}
2. Extract version from each package
3. Sort by semantic version
4. Return version list

### GitLab Registry - Fetch Package
1. Call GitLab API: GET /projects/{id}/packages/{package_id}/package_files
2. Download each file
3. Extract archives if present
4. Calculate SHA256 integrity hash
5. Return files and metadata

### Cloudsmith Registry - Similar to GitLab
1. Use Cloudsmith API endpoints
2. Support pagination for large package lists
3. Download raw package files
4. Calculate integrity hash

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Registry URL unreachable | Error with connection details |
| Authentication required | Load token from .armrc |
| Invalid token | Error with 401/403 status |
| Package not found | Error with available packages |
| Version not found | Error with available versions |
| Git repository empty | Return empty package list |
| No semantic version tags | Use branches only |
| Archive extraction fails | Error, don't install |
| Integrity hash mismatch | Error, refuse to install |

## Dependencies

- Git operations (clone, fetch, checkout)
- HTTP client (for GitLab/Cloudsmith APIs)
- Authentication (authentication.md)
- Archive extraction (pattern-filtering.md)
- Cache storage (cache-management.md)

## Implementation Mapping

**Source files:**
- `internal/arm/registry/git.go` - Git registry implementation
- `internal/arm/registry/gitlab.go` - GitLab registry implementation
- `internal/arm/registry/cloudsmith.go` - Cloudsmith registry implementation
- `internal/arm/registry/factory.go` - Registry factory for creating instances
- `internal/arm/registry/integrity.go` - SHA256 integrity calculation
- `internal/arm/service/service.go` - AddGitRegistry, AddGitLabRegistry, AddCloudsmithRegistry
- `test/e2e/registry_test.go` - E2E registry tests

## Examples

### Add Git Registry
```bash
# GitHub repository
arm add registry git --url https://github.com/org/ai-rules my-rules

# With branch tracking
arm add registry git --url https://github.com/org/ai-rules --branches main,develop my-rules

# GitLab project
arm add registry git --url https://gitlab.com/org/ai-rules my-rules
```

### Add GitLab Package Registry
```bash
# With project ID
arm add registry gitlab --project-id 123 my-gitlab

# With group ID
arm add registry gitlab --group-id 456 my-gitlab-group

# Self-hosted GitLab
arm add registry gitlab --url https://gitlab.example.com --project-id 123 my-gitlab
```

### Add Cloudsmith Registry
```bash
# Public repository
arm add registry cloudsmith --owner myorg --repo ai-rules my-cloudsmith

# With custom URL
arm add registry cloudsmith --url https://api.cloudsmith.io --owner myorg --repo ai-rules my-cloudsmith
```

### List Packages
```bash
arm list registry my-rules
# Output:
# my-rules:
#   - clean-code-ruleset
#   - security-ruleset
#   - code-review-promptset
```

### List Versions
```bash
arm info dependency my-rules/clean-code-ruleset
# Output:
# my-rules/clean-code-ruleset:
#     type: ruleset
#     version: 2.0.0
#     constraint: ^2.0.0
#     priority: 100
#     sinks:
#         - cursor-rules
# 
# Available versions:
#     - 2.0.0
#     - 1.1.0
#     - 1.0.0
#     - main (branch)
#     - develop (branch)
```

### Package Naming Models

**Git Registry (User-Named):**
```bash
# Repository structure:
# - clean-code.yml
# - security.yml

# You choose the package name:
arm install ruleset my-rules/my-clean-code cursor-rules
arm install ruleset my-rules/team-standards cursor-rules
# Package name is just a label
```

**GitLab/Cloudsmith (Registry-Named):**
```bash
# Registry has explicit packages:
# - clean-code-ruleset
# - security-ruleset

# Must use exact package name:
arm install ruleset my-gitlab/clean-code-ruleset cursor-rules
arm install ruleset my-gitlab/security-ruleset cursor-rules
# Package name must match registry
```
