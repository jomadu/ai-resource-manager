# Authentication

## Job to be Done
Securely authenticate with private registries using token-based authentication configured in .armrc files.

## Activities
1. Load authentication tokens from .armrc files
2. Support hierarchical .armrc lookup (project > user)
3. Expand environment variables in tokens
4. Apply authentication to registry requests
5. Support Bearer and Token authentication schemes

## Acceptance Criteria
- [x] Load .armrc from project directory (./.armrc)
- [x] Load .armrc from user home directory (~/.armrc)
- [x] Project .armrc overrides user .armrc
- [x] Support ARM_CONFIG_PATH to override .armrc location (bypasses hierarchy)
- [x] Parse INI-style sections [registry-url]
- [x] Expand ${ENV_VAR} in token values
- [x] Support authToken field
- [x] Apply Bearer token to HTTP requests
- [x] Apply Token header to HTTP requests
- [x] Match registry URL to .armrc section

## Data Structures

### .armrc File Format
```ini
[https://github.com/org/private-repo]
authToken = ${GITHUB_TOKEN}

[https://gitlab.com]
authToken = glpat-xyz123

[https://api.cloudsmith.io]
authToken = Bearer ${CLOUDSMITH_TOKEN}
```

## Algorithm

### Load .armrc
1. Check if ARM_CONFIG_PATH is set:
   - If set, load only that file (bypass hierarchy)
   - If not set, use hierarchical lookup:
     - Load ~/.armrc (user config)
     - Load ./.armrc (project config)
     - Project config overrides user config
2. Parse INI format
3. Extract sections and key-value pairs
4. Return configuration map

### Expand Environment Variables
1. Find ${VAR} patterns in token value
2. Look up VAR in environment
3. Replace ${VAR} with environment value
4. Return expanded token

### Apply Authentication
1. Normalize registry URL (strip protocol, trailing slashes)
2. Find matching section in .armrc
3. Extract authToken
4. Expand environment variables
5. Detect authentication scheme:
   - If token starts with "Bearer ", use as-is
   - Otherwise, prepend "Bearer "
6. Add Authorization header to request

### Section Matching
1. Normalize both URLs (strip protocol, trailing slashes)
2. Compare normalized URLs
3. Support exact match and prefix match
4. Return matching section or nil

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| No .armrc files | No authentication applied |
| ARM_CONFIG_PATH set | Use only that file, ignore ./.armrc and ~/.armrc |
| ARM_CONFIG_PATH not set | Use hierarchical lookup (./.armrc overrides ~/.armrc) |
| Environment variable not set | Leave ${VAR} unexpanded (will fail auth) |
| Invalid .armrc syntax | Skip malformed lines, continue parsing |
| Multiple matching sections | Use first match |
| Token without Bearer prefix | Prepend "Bearer " automatically |
| Token with Bearer prefix | Use as-is |
| .armrc file permissions | Warn if world-readable (security risk) |

## Dependencies

- Environment variable expansion
- INI file parsing
- HTTP client

## Implementation Mapping

**Source files:**
- `internal/arm/config/manager.go` - LoadConfig, GetSection, expandEnvVars
- `internal/arm/registry/git.go` - loadToken (Git registry)
- `internal/arm/registry/gitlab.go` - loadToken, makeRequest (GitLab registry)
- `internal/arm/registry/cloudsmith.go` - loadToken, makeRequest (Cloudsmith registry)
- `test/e2e/auth_test.go` - E2E authentication tests

## Examples

### User .armrc (~/.armrc)
```ini
[https://github.com]
authToken = ${GITHUB_TOKEN}

[https://gitlab.com]
authToken = glpat-abc123
```

### Project .armrc (./.armrc)
```ini
[https://github.com/myorg/private-repo]
authToken = ${PROJECT_GITHUB_TOKEN}
```

### Environment Variables
```bash
export GITHUB_TOKEN=ghp_xyz789
export PROJECT_GITHUB_TOKEN=ghp_abc123
export CLOUDSMITH_TOKEN=cs_def456
```

### Hierarchical Lookup
```bash
# Without ARM_CONFIG_PATH (hierarchical)
# 1. Load ~/.armrc (user config)
# 2. Load ./.armrc (project config)
# 3. Project overrides user for matching sections

# With ARM_CONFIG_PATH (bypass hierarchy)
export ARM_CONFIG_PATH=/custom/path/.armrc
# Only loads /custom/path/.armrc
# Ignores both ~/.armrc and ./.armrc
```

### Token Expansion
```ini
# Before expansion
authToken = ${GITHUB_TOKEN}

# After expansion (GITHUB_TOKEN=ghp_xyz789)
authToken = ghp_xyz789
```

### Authentication Schemes
```ini
# Automatic Bearer prefix
[https://github.com]
authToken = ghp_xyz789
# Results in: Authorization: Bearer ghp_xyz789

# Explicit Bearer prefix
[https://api.cloudsmith.io]
authToken = Bearer cs_abc123
# Results in: Authorization: Bearer cs_abc123

# Token scheme (alternative)
[https://gitlab.com]
authToken = Token glpat-xyz789
# Results in: Authorization: Token glpat-xyz789
```

### Section Matching
```ini
[https://github.com/myorg/repo]
authToken = ${GITHUB_TOKEN}
```

```bash
# Matches:
# - https://github.com/myorg/repo
# - https://github.com/myorg/repo.git
# - github.com/myorg/repo

# Does not match:
# - https://github.com/otherorg/repo
# - https://gitlab.com/myorg/repo
```

### File Permissions Warning
```bash
# If .armrc is world-readable (chmod 644)
$ arm install ruleset private-registry/rules cursor-rules
Warning: .armrc file is world-readable (permissions: 644)
Consider running: chmod 600 ~/.armrc

# Recommended permissions (chmod 600)
$ chmod 600 ~/.armrc
$ chmod 600 ./.armrc
```
