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
Returns empty list. Git registries don't have predefined packages - users define package names as arbitrary labels when installing. Package boundaries are determined by include/exclude patterns, not repository structure.

### Git Registry - List Versions
1. Ensure repository cloned/fetched: `git fetch --all --tags --force`
2. List all tags via `git tag -l`
3. Filter to only semantic version tags (non-semver tags ignored)
4. If branches configured:
   - Get remote branches via `git branch -r`
   - Match configured branches using glob patterns (e.g., `main`, `feature/*`)
   - Preserve config order for matched branches
5. Sort semantic versions descending (highest first)
6. Append matched branches in config order
7. Return version list

**Note:** packageName parameter is ignored (versions are repository-wide)

### Git Registry - Fetch Package

**Cache Strategy:**
- Cache key: `{version, normalized_include_patterns, normalized_exclude_patterns}`
- Package name NOT in cache key (allows pattern reuse across "packages")
- Pattern normalization: trim whitespace, normalize path separators (`\` → `/`), sort alphabetically

**Algorithm:**
1. Generate cache key from version + normalized patterns
2. Check cache first - return if available
3. If cache miss:
   - Ensure repository cloned/fetched
   - List files: `git ls-tree -r --name-only {commit}`
   - Extract content: `git show {commit}:{path}` for each file
   - Extract archives (.zip, .tar.gz) and merge with loose files
   - Apply include/exclude pattern filtering
   - Cache filtered results
4. Calculate SHA256 integrity hash (sorted paths for determinism)
5. Return package with files and metadata

**Note:** No git checkout performed - files extracted directly from commit objects

### GitLab Registry - List Packages
1. Load authentication token from .armrc
2. Determine API endpoint (projectId and groupId are mutually exclusive):
   - If projectId: `/api/v4/projects/{id}/packages`
   - If groupId: `/api/v4/groups/{id}/packages`
3. Paginate through all results (100 per page):
   - Fetch page with `?page={n}&per_page=100`
   - Continue until page returns < 100 items
4. Filter for `package_type == "generic"`
5. Deduplicate by package name (multiple versions → single entry)
6. Return unique package names

### GitLab Registry - List Versions
1. Load authentication token
2. Call API with pagination (same as List Packages)
3. Filter for matching package name and `package_type == "generic"`
4. Extract version from each package
5. Sort by semantic version descending
6. Return version list

### GitLab Registry - Fetch Package
1. Load authentication token
2. Find package by name and version
3. Get package files:
   - Project: `/api/v4/projects/{id}/packages/{pkg_id}/package_files`
   - Group: `/api/v4/groups/{id}/-/packages/{pkg_id}/package_files`
4. Download each file:
   - Project: `/api/v4/projects/{id}/packages/generic/{name}/{version}/{filename}`
   - Group: `/api/v4/groups/{id}/-/packages/generic/{name}/{version}/{filename}`
5. Extract archives and merge with loose files
6. Apply include/exclude pattern filtering
7. Cache filtered results (same cache key strategy as Git)
8. Calculate SHA256 integrity hash
9. Return package with files and metadata

### Cloudsmith Registry - List Packages
1. Load authentication token
2. Call `/v1/packages/{owner}/{repo}/`
3. Paginate using Link header (RFC 5988):
   - Parse `Link: <url>; rel="next"` header
   - Extract path component from full URL
   - Continue until no `rel="next"` link
4. Filter for `format == "raw"` packages
5. Deduplicate by package name
6. Return unique package names

### Cloudsmith Registry - List Versions
1. Load authentication token
2. Query packages: `/v1/packages/{owner}/{repo}/?query={packageName}`
3. Paginate using Link header
4. Filter for `format == "raw"` packages
5. Match by name OR filename prefix:
   - Include if `pkg.Name == packageName`
   - Include if `pkg.Filename` starts with `packageName`
6. Deduplicate versions (multiple files → single version)
7. Sort versions:
   - If both semver: semantic comparison descending
   - If either non-semver: lexicographic ascending
8. Return version list

### Cloudsmith Registry - Fetch Package
1. Load authentication token
2. Query and download matching packages (same matching as List Versions)
3. Extract archives and merge with loose files
4. Apply include/exclude pattern filtering
5. Cache filtered results (same cache key strategy as Git)
6. Calculate SHA256 integrity hash
7. Return package with files and metadata

### Calculate Integrity Hash
1. Collect all file paths from package
2. Sort paths alphabetically (ensures deterministic hash)
3. Create SHA256 hasher
4. For each path in sorted order:
   - Hash the path string (UTF-8 bytes)
   - Hash the file content bytes
5. Return hex digest with `"sha256-"` prefix

**Example:** Files `[{path: "b.yml", ...}, {path: "a.yml", ...}]` → hash order: `a.yml` (path+content), then `b.yml` (path+content) → `"sha256-abc123..."`

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
| Non-semver tags in Git | Filtered out (only semver tags included) |
| Branch glob pattern | Match using `filepath.Match()` (e.g., `feature/*`) |
| Archive extraction fails | Error, don't install |
| Integrity hash mismatch | Error, refuse to install |
| Concurrent git operations | File lock prevents corruption |
| Cache key collision | Multiple packages with same patterns share cache (by design) |
| GitLab pagination >100 | Automatically fetch all pages |
| Cloudsmith Link header | Parse RFC 5988 format for next page |

## Dependencies

- Git operations (clone, fetch, ls-tree, show) - no checkout required
- HTTP client (for GitLab/Cloudsmith APIs)
- Authentication (authentication.md) - see note below on .armrc format
- Archive extraction (pattern-filtering.md)
- Cache storage (cache-management.md)
- File locking (cross-process git operation protection)

**Authentication Note:** GitLab and Cloudsmith registries use registry-specific .armrc sections:
- GitLab: `[registry {url}/project/{id}]` or `[registry {url}/group/{id}]`
- Cloudsmith: `[registry {url}/{owner}/{repo}]`
- Token field: `token` (not `authToken`)
- See authentication.md for details

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

### List Registries
```bash
arm list registry
# Output:
# - my-org
# - my-gitlab
# - cloudsmith-registry
```

### Registry Info
```bash
arm info registry my-org
# Output:
# my-org:
#     type: git
#     url: https://github.com/my-org/arm-registry
```

### List Package Versions
```bash
# List all available versions for a package
arm list versions my-org/clean-code-ruleset
# Output:
# my-org/clean-code-ruleset:
#   - 2.1.0
#   - 2.0.0
#   - 1.5.0
#   - 1.0.0
#   - main (branch)
#   - develop (branch)
```

### List Packages (GitLab/Cloudsmith only)
```bash
# Note: Git registries don't support package listing (returns empty by design)
# GitLab/Cloudsmith registries have explicit packages that can be queried via API
# Currently no CLI command to list packages from a registry
# Packages are discovered during installation
```

### Branch Glob Patterns
```bash
# Match multiple branches with patterns
arm add registry git --url https://github.com/org/repo --branches "main,feature/*,release-*" my-rules

# Matches: main, feature/auth, feature/ui, release-1.0, release-2.0
# Does not match: develop, hotfix/bug
```

### Package Naming Models

**Git Registry (User-Named):**
```bash
# Repository structure:
# - clean-code.yml
# - security.yml

# You choose the package name (arbitrary label):
arm install ruleset my-rules/my-clean-code cursor-rules
arm install ruleset my-rules/team-standards cursor-rules

# Package name is NOT derived from repository structure
# Multiple "packages" with same patterns share cache
```

**GitLab/Cloudsmith (Registry-Named):**
```bash
# Registry has explicit packages:
# - clean-code-ruleset
# - security-ruleset

# Must use exact package name from registry:
arm install ruleset my-gitlab/clean-code-ruleset cursor-rules
arm install ruleset my-gitlab/security-ruleset cursor-rules
```

### Cache Key Behavior
```bash
# These two installs share the same cache (same version + patterns):
arm install ruleset my-rules/package-a@1.0.0 --include "**/*.yml" cursor-rules
arm install ruleset my-rules/package-b@1.0.0 --include "**/*.yml" q-rules

# Cache key: {version: "1.0.0", include: ["**/*.yml"], exclude: []}
# Package names "package-a" and "package-b" are just labels
```

### Authentication Examples

**GitLab (.armrc):**
```ini
[registry https://gitlab.com/project/123]
token = glpat-xyz123

# Or for group:
[registry https://gitlab.com/group/456]
token = glpat-abc789
```

**Cloudsmith (.armrc):**
```ini
[registry https://api.cloudsmith.io/myorg/ai-rules]
token = ${CLOUDSMITH_TOKEN}
```
