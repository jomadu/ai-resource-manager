# Authentication

## Job to be Done
Authenticate with registries requiring tokens (GitLab, Cloudsmith) to access private packages, enabling secure package distribution within organizations.

## Activities
1. **Parse .armrc files** - Read INI-formatted configuration from local and global locations
2. **Resolve tokens** - Look up tokens using hierarchical precedence (local > global)
3. **Expand environment variables** - Substitute `${VAR}` references with environment values
4. **Generate auth keys** - Create unique keys for registry identification
5. **Inject authentication** - Add tokens to HTTP requests (Bearer or Token headers)

## Acceptance Criteria
- [ ] Local .armrc overrides global ~/.armrc for same section
- [ ] Environment variables expanded using `${VAR}` syntax
- [ ] Missing environment variables expand to empty string
- [ ] Section names match registry URLs exactly (protocol-aware)
- [ ] GitLab uses Bearer token authentication
- [ ] Cloudsmith uses Token authentication
- [ ] Missing .armrc files handled gracefully (no auth)
- [ ] Missing sections return error when token required
- [ ] File permissions should be 0600 for security
- [ ] Multiple registry sections supported in single file
- [ ] Config manager accepts working directory and home directory as constructor parameters
- [ ] Default constructor calls os.Getwd() and os.UserHomeDir() for production use
- [ ] Test constructor accepts temporary directory paths for isolation

## Data Structures

### .armrc File Format (INI)
```ini
[registry https://hostname/project/id]
token = your_token_here

[registry https://hostname/group/id]
token = ${ENV_VAR_NAME}
```

**Section naming:**
- GitLab project: `registry https://gitlab.example.com/project/123`
- GitLab group: `registry https://gitlab.example.com/group/456`
- Cloudsmith: `registry https://api.cloudsmith.io/owner/repository`

### Config Manager Interface
```go
type Manager interface {
    GetAllSections(ctx context.Context) (map[string]map[string]string, error)
    GetSection(ctx context.Context, section string) (map[string]string, error)
    GetValue(ctx context.Context, section, key string) (string, error)
}
```

**Methods:**
- `GetAllSections` - Retrieve all sections from both files (local overrides global)
- `GetSection` - Retrieve single section using hierarchical lookup
- `GetValue` - Retrieve single key from section

### File Manager
```go
type FileManager struct {
    workingDir  string  // Current working directory (local .armrc)
    userHomeDir string  // User home directory (global .armrc)
}
```

**Fields:**
- `workingDir` - Path to project directory containing `.armrc`
- `userHomeDir` - Path to user home directory containing `.armrc`

**Construction:**
- `NewFileManager()` - Calls os.Getwd() and os.UserHomeDir() for production
- `NewFileManagerWithPaths(workingDir, userHomeDir)` - Direct path injection for testing

## Algorithm

### 1. Parse .armrc File
Read INI file and extract sections with key-value pairs.

**Pseudocode:**
```
function loadFileIfExists(filePath string) (*ini.File, error):
    cfg = ini.Load(filePath)
    if error is NotExist:
        return nil, nil  // File doesn't exist, not an error
    if error:
        return error("failed to load .armrc")
    return cfg

function getAllSectionsFromFile(filePath string) (map[string]map[string]string, error):
    cfg = loadFileIfExists(filePath)
    if cfg == nil:
        return empty map
    
    sections = {}
    for section in cfg.Sections():
        // Skip empty default section
        if section.Name == "DEFAULT" and len(section.Keys) == 0:
            continue
        
        values = {}
        for key in section.Keys():
            values[key.Name] = expandEnvVars(key.Value)
        
        if len(values) > 0:
            sections[section.Name] = values
    
    return sections
```

**Implementation:** `internal/arm/config/manager.go:getAllSectionsFromFile()`

### 2. Get All Sections (Hierarchical Merge)
Retrieve all sections from both files, with local overriding global.

**Pseudocode:**
```
function GetAllSections(ctx context.Context) (map[string]map[string]string, error):
    result = {}
    
    // Load global file first (base)
    if userHomeDir != "":
        userRcPath = userHomeDir + "/.armrc"
        userSections = getAllSectionsFromFile(userRcPath)
        if error:
            return error
        
        // Copy user sections to result
        for section, values in userSections:
            result[section] = values
    
    // Load local file and override
    if workingDir != "":
        projectRcPath = workingDir + "/.armrc"
        projectSections = getAllSectionsFromFile(projectRcPath)
        if error:
            return error
        
        // Project sections override user sections
        for section, values in projectSections:
            result[section] = values
    
    return result
```

**Implementation:** `internal/arm/config/manager.go:GetAllSections()`

### 3. Get Section (Hierarchical Lookup)
Retrieve single section, checking local first then global.

**Pseudocode:**
```
function GetSection(ctx context.Context, section string) (map[string]string, error):
    // Try local .armrc first
    if workingDir != "":
        projectRcPath = workingDir + "/.armrc"
        values = getSectionFromFile(projectRcPath, section)
        if no error:
            return values
        // If file doesn't exist or section not found, continue to global
    
    // Fallback to global .armrc
    if userHomeDir == "":
        return error("section not found")
    
    userRcPath = userHomeDir + "/.armrc"
    return getSectionFromFile(userRcPath, section)

function getSectionFromFile(filePath string, section string) (map[string]string, error):
    cfg = loadFileIfExists(filePath)
    if cfg == nil:
        return error("section not found")
    
    sec = cfg.GetSection(section)
    if error:
        return error("section not found")
    
    values = {}
    for key in sec.Keys():
        values[key.Name] = expandEnvVars(key.Value)
    
    return values
```

**Implementation:** `internal/arm/config/manager.go:GetSection()`

### 4. Get Value
Retrieve single key from section.

**Pseudocode:**
```
function GetValue(ctx context.Context, section string, key string) (string, error):
    values = GetSection(ctx, section)
    if error:
        return error
    
    value = values[key]
    if not exists:
        return error("key not found in section")
    
    return value
```

**Implementation:** `internal/arm/config/manager.go:GetValue()`

### 5. Expand Environment Variables
Substitute `${VAR}` references with environment values.

**Pseudocode:**
```
function expandEnvVars(value string) string:
    return os.ExpandEnv(value)
```

**Behavior:**
- `${GITLAB_TOKEN}` → value of `GITLAB_TOKEN` environment variable
- `${MISSING_VAR}` → empty string (if variable not set)
- `literal-value` → unchanged

**Implementation:** `internal/arm/config/manager.go:expandEnvVars()`

### 6. Generate Auth Key (GitLab)
Create unique key for registry identification.

**Pseudocode:**
```
function getAuthKey() string:
    baseURL = config.URL
    if config.ProjectID != "":
        return baseURL + "/project/" + config.ProjectID
    return baseURL + "/group/" + config.GroupID
```

**Examples:**
- Project: `https://gitlab.example.com/project/123`
- Group: `https://gitlab.example.com/group/456`

**Implementation:** `internal/arm/registry/gitlab.go:getAuthKey()`

### 7. Generate Auth Key (Cloudsmith)
Create unique key for Cloudsmith registry.

**Pseudocode:**
```
function getAuthKey() string:
    return config.URL + "/" + config.Owner + "/" + config.Repository
```

**Example:** `https://api.cloudsmith.io/myorg/ai-rules`

**Implementation:** `internal/arm/registry/cloudsmith.go:loadToken()`

### 8. Load Token (GitLab)
Load token from .armrc for GitLab registry.

**Pseudocode:**
```
function loadToken(ctx context.Context) error:
    if client.token != "":
        return nil  // Already loaded
    
    if configMgr == nil:
        return nil  // No config manager, no auth
    
    authKey = getAuthKey()
    token = configMgr.GetValue(ctx, "registry " + authKey, "token")
    if error:
        return error("failed to load token from .armrc")
    
    client.token = token
    return nil
```

**Implementation:** `internal/arm/registry/gitlab.go:loadToken()`

### 9. Load Token (Cloudsmith)
Load token from .armrc for Cloudsmith registry.

**Pseudocode:**
```
function loadToken(ctx context.Context) error:
    if client.token != "":
        return nil  // Already loaded
    
    if configMgr == nil:
        return error("no token configured for Cloudsmith registry")
    
    authKey = config.URL + "/" + config.Owner + "/" + config.Repository
    token = configMgr.GetValue(ctx, "registry " + authKey, "token")
    if error:
        return error("failed to load token from .armrc")
    
    client.token = token
    return nil
```

**Implementation:** `internal/arm/registry/cloudsmith.go:loadToken()`

### 10. Inject Authentication (GitLab)
Add Bearer token to HTTP request.

**Pseudocode:**
```
function makeRequest(ctx context.Context, method string, path string) (*http.Response, error):
    req = http.NewRequest(method, baseURL + path, nil)
    
    if token != "":
        req.Header.Set("Authorization", "Bearer " + token)
    
    return httpClient.Do(req)
```

**Implementation:** `internal/arm/registry/gitlab.go` (client methods)

### 11. Inject Authentication (Cloudsmith)
Add Token authentication to HTTP request.

**Pseudocode:**
```
function makeRequest(ctx context.Context, method string, path string) (*http.Response, error):
    req = http.NewRequest(method, baseURL + path, nil)
    
    if token != "":
        req.Header.Set("Authorization", "Token " + token)
    
    return httpClient.Do(req)
```

**Implementation:** `internal/arm/registry/cloudsmith.go:makeRequest()`

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| No .armrc files exist | No authentication (works for public registries) |
| Local .armrc only | Use local configuration |
| Global .armrc only | Use global configuration |
| Both files exist, same section | Local overrides global completely |
| Both files exist, different sections | Merge sections (local + global) |
| Section not found | Return error when token required |
| Key not found in section | Return error |
| Environment variable not set | Expand to empty string |
| Invalid INI format | Return parse error |
| File permissions too open (not 0600) | Warning (security best practice) |
| Expired token | Registry returns 401 (handled by registry) |
| Invalid token format | Registry returns 401 (handled by registry) |
| Missing protocol in section name | Section not found (exact match required) |
| HTTP vs HTTPS mismatch | Section not found (protocol-aware) |
| Multiple keys in section | All keys available |
| Empty section | Section exists but no keys |

## Dependencies

- `gopkg.in/ini.v1` - INI file parsing
- `os` package - Environment variable expansion, file operations
- `path/filepath` - Path manipulation
- `context` package - Context propagation

## Implementation Mapping

**Source files:**
- `internal/arm/config/manager.go` - .armrc parsing and hierarchical lookup (accepts workingDir and homeDir parameters)
- `internal/arm/registry/gitlab.go` - GitLab authentication (Bearer token)
- `internal/arm/registry/cloudsmith.go` - Cloudsmith authentication (Token)

**Related specs:**
- `registry-management.md` - Registry configuration and factory pattern
- `package-installation.md` - Package installation using authenticated registries

## Examples

### Example 1: Local .armrc Overrides Global

**Global ~/.armrc:**
```ini
[registry https://gitlab.example.com/project/123]
token = global-token-123
```

**Local ./.armrc:**
```ini
[registry https://gitlab.example.com/project/123]
token = local-token-123
```

**Lookup:**
```
GetValue(ctx, "registry https://gitlab.example.com/project/123", "token")
→ "local-token-123"
```

**Verification:**
- Local token used
- Global token ignored

### Example 2: Environment Variable Expansion

**.armrc:**
```ini
[registry https://gitlab.example.com/project/456]
token = ${GITLAB_TOKEN}
```

**Environment:**
```bash
export GITLAB_TOKEN=glpat-abc123def456
```

**Lookup:**
```
GetValue(ctx, "registry https://gitlab.example.com/project/456", "token")
→ "glpat-abc123def456"
```

**Verification:**
- Environment variable expanded
- Token value from environment

### Example 3: Multiple Registries

**.armrc:**
```ini
[registry https://gitlab.example.com/project/111]
token = token-111

[registry https://gitlab.example.com/project/222]
token = token-222

[registry https://api.cloudsmith.io/org1/repo1]
token = ckcy_token_1
```

**Lookup:**
```
GetAllSections(ctx)
→ {
    "registry https://gitlab.example.com/project/111": {"token": "token-111"},
    "registry https://gitlab.example.com/project/222": {"token": "token-222"},
    "registry https://api.cloudsmith.io/org1/repo1": {"token": "ckcy_token_1"}
}
```

**Verification:**
- All sections retrieved
- Each section independent

### Example 4: GitLab Authentication Flow

**Registry Configuration:**
```
URL: https://gitlab.example.com
ProjectID: 123
```

**.armrc:**
```ini
[registry https://gitlab.example.com/project/123]
token = glpat-abc123def456
```

**Authentication Flow:**
1. Registry calls `loadToken(ctx)`
2. Generate auth key: `https://gitlab.example.com/project/123`
3. Lookup section: `registry https://gitlab.example.com/project/123`
4. Get token: `glpat-abc123def456`
5. HTTP request: `Authorization: Bearer glpat-abc123def456`

**Verification:**
- Token loaded from .armrc
- Bearer authentication used
- HTTP header set correctly

### Example 5: Cloudsmith Authentication Flow

**Registry Configuration:**
```
URL: https://api.cloudsmith.io
Owner: myorg
Repository: ai-rules
```

**.armrc:**
```ini
[registry https://api.cloudsmith.io/myorg/ai-rules]
token = ckcy_abc123def456
```

**Authentication Flow:**
1. Registry calls `loadToken(ctx)`
2. Generate auth key: `https://api.cloudsmith.io/myorg/ai-rules`
3. Lookup section: `registry https://api.cloudsmith.io/myorg/ai-rules`
4. Get token: `ckcy_abc123def456`
5. HTTP request: `Authorization: Token ckcy_abc123def456`

**Verification:**
- Token loaded from .armrc
- Token authentication used (not Bearer)
- HTTP header set correctly

### Example 6: Missing .armrc (Public Registry)

**No .armrc files exist**

**Authentication Flow:**
1. Registry calls `loadToken(ctx)`
2. configMgr is nil or section not found
3. GitLab: Returns nil (no auth)
4. Cloudsmith: Returns error (auth required)
5. HTTP request: No Authorization header

**Verification:**
- GitLab works without auth (public repos)
- Cloudsmith requires auth (returns error)

### Example 7: Section Not Found

**.armrc:**
```ini
[registry https://gitlab.example.com/project/123]
token = token-123
```

**Lookup:**
```
GetValue(ctx, "registry https://gitlab.example.com/project/999", "token")
→ error: "section not found: registry https://gitlab.example.com/project/999"
```

**Verification:**
- Section name must match exactly
- Protocol-aware (https vs http)
- Returns error when not found

## Notes

### Design Decisions

1. **Hierarchical precedence** - Local .armrc overrides global for same section, enabling project-specific tokens
2. **Complete section override** - Local section replaces entire global section (not key-by-key merge)
3. **Environment variable expansion** - Enables CI/CD integration without hardcoding tokens
4. **Protocol-aware section names** - Prevents ambiguity between HTTP/HTTPS registries
5. **Graceful degradation** - Missing .armrc files don't cause errors (public registries work)
6. **INI format** - Simple, human-readable, widely supported

### Security Considerations

- **File permissions**: .armrc should be 0600 (owner read/write only)
- **Gitignore**: Add .armrc to .gitignore to prevent committing tokens
- **Environment variables**: Safer for CI/CD than hardcoded tokens
- **Token scopes**: Use minimal required scopes (read_api, read_repository)
- **Token rotation**: Regularly rotate tokens for security

### GitLab vs Cloudsmith Differences

| Aspect | GitLab | Cloudsmith |
|--------|--------|------------|
| Auth header | `Authorization: Bearer <token>` | `Authorization: Token <token>` |
| Token format | `glpat-*` or `gldt-*` | `ckcy_*` |
| Auth required | Optional (public repos) | Required |
| Auth key format | `URL/project/ID` or `URL/group/ID` | `URL/owner/repo` |
| Missing token | No error (no auth) | Error (auth required) |

### Testing Considerations

- Test local-only, global-only, and both configurations
- Test section override behavior (local > global)
- Test environment variable expansion (set and unset)
- Test multiple sections in single file
- Test missing .armrc files (graceful degradation)
- Test section not found errors
- Test key not found errors
- Test protocol-aware section matching (http vs https)
- Test file permissions (security warning)
