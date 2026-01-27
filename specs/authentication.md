# Authentication

## Job to be Done
Securely authenticate with private registries using token-based authentication configured in .armrc files.

## Activities
1. Load authentication tokens from .armrc files
2. Support hierarchical .armrc lookup (project > user)
3. Expand environment variables in tokens
4. Apply authentication to registry requests

## Acceptance Criteria
- [x] Load .armrc from project directory (./.armrc)
- [x] Load .armrc from user home directory (~/.armrc)
- [x] Project .armrc overrides user .armrc
- [x] Support ARM_CONFIG_PATH to override .armrc location (bypasses hierarchy)
- [x] Parse INI-style sections with `registry` prefix
- [x] Expand ${ENV_VAR} in token values
- [x] Support `token` field
- [x] GitLab: Apply Bearer token to HTTP requests
- [x] Cloudsmith: Apply Token header to HTTP requests

## Data Structures

### .armrc File Format
```ini
# GitLab with project ID
[registry https://gitlab.example.com/project/123]
token = glpat-xyz123

# GitLab with group ID
[registry https://gitlab.example.com/group/456]
token = ${GITLAB_TOKEN}

# Cloudsmith
[registry https://api.cloudsmith.io/myorg/myrepo]
token = ${CLOUDSMITH_TOKEN}
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

**GitLab:**
1. Construct auth key: `{url}/project/{id}` or `{url}/group/{id}`
2. Prepend `"registry "` to auth key
3. Look up section in .armrc by exact string match
4. Extract `token` field value
5. Expand environment variables
6. Set header: `Authorization: Bearer {token}`

**Cloudsmith:**
1. Construct auth key: `{url}/{owner}/{repo}`
2. Prepend `"registry "` to auth key
3. Look up section in .armrc by exact string match
4. Extract `token` field value
5. Expand environment variables
6. Set header: `Authorization: Token {token}`

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| No .armrc files | No authentication applied |
| ARM_CONFIG_PATH set | Use only that file, ignore ./.armrc and ~/.armrc |
| ARM_CONFIG_PATH not set | Use hierarchical lookup (./.armrc overrides ~/.armrc) |
| Environment variable not set | Expands to empty string |
| Invalid .armrc syntax | Skip malformed lines, continue parsing |
| Section not found | Return error |
| Token field not found | Return error |
| .armrc file permissions | Should be 0600 (owner read/write only) |

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
[registry https://gitlab.example.com/project/123]
token = ${GITLAB_TOKEN}

[registry https://api.cloudsmith.io/myorg/myrepo]
token = ${CLOUDSMITH_TOKEN}
```

### Project .armrc (./.armrc)
```ini
[registry https://gitlab.example.com/project/456]
token = glpat-project-token
```

### Environment Variables
```bash
export GITLAB_TOKEN=glpat-xyz789
export CLOUDSMITH_TOKEN=ckcy-abc123
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
token = ${GITLAB_TOKEN}

# After expansion (GITLAB_TOKEN=glpat-xyz789)
token = glpat-xyz789
```

### Authentication Headers
```ini
# GitLab
[registry https://gitlab.example.com/project/123]
token = glpat-xyz789
# Results in: Authorization: Bearer glpat-xyz789

# Cloudsmith
[registry https://api.cloudsmith.io/myorg/myrepo]
token = ckcy-abc123
# Results in: Authorization: Token ckcy-abc123
```

### Section Name Format

**GitLab with Project ID:**
```ini
[registry https://gitlab.example.com/project/123]
token = glpat-xyz
```

**GitLab with Group ID:**
```ini
[registry https://gitlab.example.com/group/456]
token = glpat-xyz
```

**Cloudsmith:**
```ini
[registry https://api.cloudsmith.io/myorg/myrepo]
token = ckcy-xyz
```

**Important:** Section names must exactly match the constructed auth key:
- GitLab: `registry {url}/project/{id}` or `registry {url}/group/{id}`
- Cloudsmith: `registry {url}/{owner}/{repo}`
- No URL normalization is performed
- Exact string match required

### File Permissions
```bash
# Recommended permissions (owner read/write only)
chmod 600 ~/.armrc
chmod 600 ./.armrc
```
