# GitLab Registry Support - Technical Specification

## Architecture Overview

The GitLab registry implementation extends ARM's existing registry interface to support GitLab's Generic Package Registry API. It follows ARM's established patterns for caching, version resolution, and content retrieval while adding GitLab-specific authentication and API integration.

## Component Design

### Registry Interface Implementation

```go
// GitLabRegistry implements the Registry interface for GitLab package registries
type GitLabRegistry struct {
    cache    cache.RegistryRulesetCache
    config   GitLabRegistryConfig
    resolver resolver.ConstraintResolver
    client   *GitLabClient
}

// GitLabRegistryConfig extends RegistryConfig with GitLab-specific settings
type GitLabRegistryConfig struct {
    registry.RegistryConfig
    ProjectID  string `json:"project_id,omitempty"`
    GroupID    string `json:"group_id,omitempty"`
    APIVersion string `json:"api_version"`
}

// GitLabClient handles HTTP communication with GitLab API
type GitLabClient struct {
    baseURL    string
    apiVersion string
    httpClient *http.Client
    token      string
}
```

### Authentication Implementation

#### Token Storage
- **Only**: Project `.armrc` file with environment variable expansion
- **Never**: Plain text in shared configuration files (arm.json)

#### RC File Service
```go
// RCFileService handles .armrc file operations
type RCFileService struct {
    filePath string
}

func NewRCFileService() *RCFileService {
    return &RCFileService{
        filePath: ".armrc",
    }
}

func (r *RCFileService) GetValue(section, key string) (string, error)
func (r *RCFileService) expandEnvVars(value string) string
```

#### Authentication Interface
```go
func (g *GitLabRegistry) loadToken(rcService *RCFileService) (string, error)
```

#### .armrc File Format
```ini
# Registry authentication
[registry my-gitlab]
    token = ${GITLAB_TOKEN}

[registry ci-registry]
    token = ${CI_JOB_TOKEN}
```

### GitLab API Integration

#### Factory Integration
```go
// Add "gitlab" case to registry/factory.go
case "gitlab":
    return newGitLabRegistry(name, rawConfig)

func newGitLabRegistry(name string, rawConfig map[string]interface{}) (*GitLabRegistry, error) {
    // Parse raw config into GitLabRegistryConfig
    configBytes, err := json.Marshal(rawConfig)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal config: %w", err)
    }

    var gitlabConfig GitLabRegistryConfig
    if err := json.Unmarshal(configBytes, &gitlabConfig); err != nil {
        return nil, fmt.Errorf("failed to parse gitlab registry config: %w", err)
    }

    // Build registry key object for cache uniqueness
    registryKeyObj := map[string]string{
        "url":  gitlabConfig.URL,
        "type": gitlabConfig.Type,
    }
    if gitlabConfig.ProjectID != "" {
        registryKeyObj["project_id"] = gitlabConfig.ProjectID
    }
    if gitlabConfig.GroupID != "" {
        registryKeyObj["group_id"] = gitlabConfig.GroupID
    }

    rulesetCache, err := cache.NewRegistryRulesetCache(registryKeyObj)
    if err != nil {
        return nil, err
    }

    return NewGitLabRegistry(gitlabConfig, rulesetCache), nil
}
```

#### API Endpoints
```go
// Endpoint templates - version will be injected at runtime
const (
    // Project-level endpoints
    ProjectPackageListTemplate     = "/api/%s/projects/%s/packages"
    ProjectPackageFilesTemplate    = "/api/%s/projects/%s/packages/%d/package_files"
    ProjectPackageDownloadTemplate = "/api/%s/projects/%s/packages/generic/%s/%s/%s"

    // Group-level endpoints
    GroupPackageListTemplate     = "/api/%s/groups/%s/packages"
    GroupPackageFilesTemplate    = "/api/%s/groups/%s/-/packages/%d/package_files"
    GroupPackageDownloadTemplate = "/api/%s/groups/%s/-/packages/generic/%s/%s/%s"
)

// Client methods build URLs with API version
func (c *GitLabClient) buildProjectPackageListURL(projectID string) string {
    return fmt.Sprintf(c.baseURL+ProjectPackageListTemplate, c.apiVersion, projectID)
}

func (c *GitLabClient) buildProjectPackageFilesURL(projectID string, packageID int) string {
    return fmt.Sprintf(c.baseURL+ProjectPackageFilesTemplate, c.apiVersion, projectID, packageID)
}
```

### Registry Interface Implementation

```go
// GitLabRegistry implements the Registry interface
func (g *GitLabRegistry) ListVersions(ctx context.Context) ([]types.Version, error)
func (g *GitLabRegistry) ResolveVersion(ctx context.Context, constraint string) (*resolver.ResolvedVersion, error)
func (g *GitLabRegistry) GetContent(ctx context.Context, version types.Version, selector types.ContentSelector) ([]types.File, error)
```

#### GitLab Package Types
```go
type GitLabPackage struct {
    ID          int                  `json:"id"`
    Name        string               `json:"name"`
    Version     string               `json:"version"`
    PackageType string               `json:"package_type"`
    CreatedAt   time.Time           `json:"created_at"`
    Files       []GitLabPackageFile `json:"package_files"`
}

type GitLabPackageFile struct {
    ID       int    `json:"id"`
    FileName string `json:"file_name"`
    Size     int64  `json:"size"`
}
```

#### GitLab Client Interface
```go
func (c *GitLabClient) ListProjectPackages(ctx context.Context, projectID string) ([]GitLabPackage, error)
func (c *GitLabClient) ListGroupPackages(ctx context.Context, groupID string) ([]GitLabPackage, error)
func (c *GitLabClient) DownloadProjectPackage(ctx context.Context, projectID, packageName, version string, selector types.ContentSelector) ([]types.File, error)
func (c *GitLabClient) DownloadGroupPackage(ctx context.Context, groupID, packageName, version string, selector types.ContentSelector) ([]types.File, error)
```

#### Package Download Workflow

GitLab Generic Package Registry requires downloading individual files rather than entire packages. The client implements a two-step process:

1. **List Package Files**: Use the package files API to get all files in a package
2. **Download Each File**: Download individual files using the generic package download endpoint

#### Package Structure Optimization

GitLab packages are structured to optimize for URF files:
- URF files (*.yml, *.yaml) are placed at package root for easy targeting
- Pre-compiled rules are organized under `build/{tool}/` directories
- Command-level defaults should target URF files for GitLab registries

```go
// Package files endpoint for listing all files in a package
const (
    ProjectPackageFilesTemplate = "/api/%s/projects/%s/packages/%d/package_files"
    GroupPackageFilesTemplate   = "/api/%s/groups/%s/-/packages/%d/package_files"
)

// Download workflow implementation
func (c *GitLabClient) downloadPackageFiles(ctx context.Context, pkg GitLabPackage, selector types.ContentSelector) ([]types.File, error) {
    var files []types.File

    // Filter package files based on selector
    for _, pkgFile := range pkg.Files {
        if selector.Matches(pkgFile.FileName) {
            // Download individual file using generic package endpoint
            fileContent, err := c.downloadSingleFile(ctx, pkg, pkgFile.FileName)
            if err != nil {
                return nil, fmt.Errorf("failed to download %s: %w", pkgFile.FileName, err)
            }

            files = append(files, types.File{
                Path:    pkgFile.FileName, // Preserves directory structure (e.g., "rules/amazonq/a-q-rule.md")
                Content: fileContent,
            })
        }
    }

    return files, nil
}

func (c *GitLabClient) downloadSingleFile(ctx context.Context, pkg GitLabPackage, fileName string) ([]byte, error) {
    // Use generic package download endpoint: /api/v4/projects/:id/packages/generic/:package_name/:package_version/:file_name
    url := fmt.Sprintf("%s/api/%s/projects/%s/packages/generic/%s/%s/%s",
        c.baseURL, c.apiVersion, pkg.ProjectID, pkg.Name, pkg.Version, fileName)

    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }

    resp, err := c.makeRequest(ctx, req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    return io.ReadAll(resp.Body)
}
```

### Configuration Integration

#### Manifest Configuration

**Project Registry:**
```json
{
  "registries": {
    "my-gitlab-project": {
      "type": "gitlab",
      "url": "https://gitlab.example.com",
      "project_id": "123",
      "api_version": "v4"
    }
  }
}
```

**Group Registry:**
```json
{
  "registries": {
    "my-gitlab-group": {
      "type": "gitlab",
      "url": "https://gitlab.example.com",
      "group_id": "456",
      "api_version": "v4"
    }
  }
}
``` {
      "type": "gitlab",
      "url": "https://gitlab.example.com",
      "group_id": "456",
      "api_version": "v4"
    }
  }
}
```

#### Key Differences from Git Registry
- **Package-level targeting**: Install from standardized `ai-rules` packages
- **No branches**: GitLab uses semantic versioning only
- **Simplified configuration**: Project ID or Group ID instead of complex Git settings
- **Registry scope**: Project registries for single-project packages, Group registries for multi-project packages
- **URF-optimized structure**: URF files at root level for easy command-level targeting
- **Pre-compiled rule support**: Packages can contain both URF source and compiled rules
- **Hierarchical file support**: Maintains directory structure from packages: Preserves directory structure from packages (e.g., `rules/amazonq/file.md`)
- **Mixed format support**: Handles URF files (`.yml`/`.yaml`), tool-specific files (`.mdc`, `.md`), and nested directories

#### CLI Commands
```bash
# Add GitLab project registry
arm config registry add my-gitlab https://gitlab.com/group/project --type gitlab

# Add with explicit project ID
arm config registry add my-gitlab https://gitlab.com --type gitlab --project-id 12345

# Add GitLab group registry
arm config registry add my-gitlab-group https://gitlab.com/group --type gitlab --group-id 456

# Install from GitLab registry with file selection
arm install my-gitlab/cursor-rules --include "*.mdc" --sinks cursor
arm install my-gitlab/mixed-rules --include "rules/amazonq/*.md" --sinks amazonq
```

### Key Differences from Git Registry
- **Package-level targeting**: Install packages by name, then use include/exclude patterns for file selection
- **No branches**: GitLab uses semantic versioning only
- **Simplified configuration**: Project ID or Group ID instead of complex Git settings
- **Registry scope**: Project registries for single-project packages, Group registries for multi-project packages

### Implementation Files
```
internal/registry/
├── config.go              # Add GitLabRegistryConfig
├── factory.go             # Add gitlab case
├── git_registry.go        # Existing
├── gitlab_registry.go     # New GitLab implementation
├── gitlab_client.go       # New GitLab API client
└── types.go               # Existing Registry interface

internal/resolver/
├── constraint.go          # Existing
└── semantic.go            # New SemanticConstraintResolver

internal/rcfile/
└── service.go             # New RC file service
```go
// GitLabRegistry uses the same cache.RegistryRulesetCache as GitRegistry
// The factory creates the registry key object with URL, type, project_id/group_id
// which ensures cache uniqueness across different GitLab registries
func NewGitLabRegistry(config GitLabRegistryConfig, rulesetCache cache.RegistryRulesetCache) *GitLabRegistry {
    return &GitLabRegistry{
        cache:    rulesetCache,
        config:   config,
        resolver: resolver.NewSemanticConstraintResolver(), // Use semantic-only resolver
        client:   NewGitLabClient(config.URL, config.APIVersion),
    }
}
```

#### Rate Limiting
```go
// Simple rate limiting for GitLab API
func (c *GitLabClient) makeRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
    // Add Authorization header
    req.Header.Set("Authorization", "Bearer "+c.auth.Token)
    req.Header.Set("Content-Type", "application/json")

    return c.httpClient.Do(req.WithContext(ctx))
}
```

### Testing Strategy

#### Unit Tests
```go
// Follow existing test patterns in internal/registry/
func TestGitLabRegistry_ListVersions(t *testing.T) {
    // Mock GitLab API responses
    // Test version parsing and sorting
}

func TestGitLabRegistry_ResolveVersion(t *testing.T) {
    // Test constraint resolution
    // Test error handling
}

func TestGitLabAuth_LoadFromArmRC(t *testing.T) {
    // Test .armrc parsing
    // Test environment variable expansion
}
```

#### Integration Tests
- Real GitLab API interactions (with test tokens)
- Package download workflows
- Cache integration validation

### Implementation Plan

#### Phase 1: Core Infrastructure
1. Add `GitLabRegistryConfig` to `internal/registry/config.go`
2. Implement `GitLabRegistry` struct with Registry interface
3. Add GitLab case to `registry/factory.go`
4. Implement `.armrc` authentication loading

#### Phase 2: GitLab API Client
1. Implement `GitLabClient` with HTTP operations
2. Add package listing and downloading
3. Integrate with existing cache infrastructure
4. Add rate limiting and error handling

#### Phase 3: CLI Integration
1. Update CLI commands to support `--type gitlab`
2. Add GitLab-specific configuration options
3. Update documentation and examples
4. Comprehensive testing

### File Structure
```
internal/registry/
├── config.go              # Add GitLabRegistryConfig
├── factory.go             # Add gitlab case
├── git_registry.go        # Existing
├── gitlab_registry.go     # New GitLab implementation
├── gitlab_client.go       # New GitLab API client
└── types.go               # Existing Registry interface
```
