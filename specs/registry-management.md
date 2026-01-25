# Registry Management

## Job to be Done
Configure and manage registries as remote sources for AI packages, enabling users to add, list, and remove registries with proper authentication and validation.

## Activities
1. Add registry (Git, GitLab, Cloudsmith) with configuration validation
2. List all configured registries with their types and URLs
3. Remove registry by name with dependency checking
4. Validate registry configuration (URL format, required fields)
5. Store registry configuration in manifest file
6. Generate authentication keys for GitLab and Cloudsmith registries

## Acceptance Criteria
- [ ] Git registry can be added with URL and optional branches
- [ ] GitLab registry can be added with URL and optional project-id, group-id, api-version
- [ ] Cloudsmith registry can be added with URL, owner, and repository
- [ ] Registry names must be unique (error if duplicate without --force flag)
- [ ] --force flag allows overwriting existing registry configuration
- [ ] Registry configuration is persisted to arm.json manifest file
- [ ] List registries displays name, type, and URL for all configured registries
- [ ] Remove registry deletes configuration from manifest
- [ ] Remove registry fails if packages depend on it (unless --force used)
- [ ] Invalid URLs are rejected with clear error messages
- [ ] Missing required fields (owner, repository) are rejected with clear error messages
- [ ] Registry factory creates correct registry type from configuration
- [ ] Authentication keys are generated for GitLab and Cloudsmith registries
- [ ] Registry configuration supports arbitrary JSON fields for extensibility

## Data Structures

### Manifest Registry Storage
```json
{
  "registries": {
    "registry-name": {
      "type": "git|gitlab|cloudsmith",
      "url": "https://...",
      "branches": ["main", "develop"],
      "projectId": "123",
      "groupId": "456",
      "apiVersion": "v4",
      "owner": "myorg",
      "repository": "ai-rules"
    }
  }
}
```

**Fields:**
- `type` - Registry type: "git", "gitlab", or "cloudsmith" (required)
- `url` - Registry URL (required)
- `branches` - Git branches to include as versions (optional, Git only)
- `projectId` - GitLab project ID (optional, GitLab only)
- `groupId` - GitLab group ID (optional, GitLab only)
- `apiVersion` - GitLab API version (optional, GitLab only, default: "v4")
- `owner` - Cloudsmith organization/owner (required, Cloudsmith only)
- `repository` - Cloudsmith repository name (required, Cloudsmith only)

### Registry Interface
```go
type Registry interface {
    ListPackages(ctx context.Context) ([]*core.PackageMetadata, error)
    ListPackageVersions(ctx context.Context, packageName string) ([]core.Version, error)
    GetPackage(ctx context.Context, packageName string, version *core.Version, include []string, exclude []string) (*core.Package, error)
}
```

**Methods:**
- `ListPackages` - Returns available packages (empty for Git registries)
- `ListPackageVersions` - Returns available versions for a package
- `GetPackage` - Downloads and returns package content with optional filtering

### Registry Config Types
```go
type RegistryConfig struct {
    URL  string `json:"url"`
    Type string `json:"type"`
}

type GitRegistryConfig struct {
    RegistryConfig
    Branches []string `json:"branches,omitempty"`
}

type GitLabRegistryConfig struct {
    RegistryConfig
    ProjectID  string `json:"projectId,omitempty"`
    GroupID    string `json:"groupId,omitempty"`
    APIVersion string `json:"apiVersion,omitempty"`
}

type CloudsmithRegistryConfig struct {
    RegistryConfig
    Owner      string `json:"owner"`
    Repository string `json:"repository"`
}
```

## Algorithm

### Add Registry
1. Parse command-line arguments (type, name, flags)
2. Validate required fields based on registry type:
   - Git: url, name
   - GitLab: url, name
   - Cloudsmith: url, owner, repository, name
3. Load existing manifest from arm.json
4. Check if registry name already exists:
   - If exists and --force not set → error "registry already exists"
   - If exists and --force set → overwrite configuration
5. Create registry configuration map with type and fields
6. Store configuration in manifest.Registries[name]
7. Save manifest to arm.json
8. Generate authentication key if needed (GitLab, Cloudsmith)
9. Print success message

**Pseudocode:**
```
function AddRegistry(type, name, config, force):
    validate_required_fields(type, config)
    
    manifest = load_manifest()
    
    if manifest.Registries[name] exists and not force:
        return error "registry already exists, use --force to overwrite"
    
    registry_config = {
        "type": type,
        "url": config.url,
        ...additional_fields
    }
    
    manifest.Registries[name] = registry_config
    save_manifest(manifest)
    
    if type in ["gitlab", "cloudsmith"]:
        generate_auth_key(name)
    
    return success
```

### List Registries
1. Load manifest from arm.json
2. If no registries configured → print "No registries configured"
3. For each registry in manifest.Registries:
   - Extract name, type, url
   - Format output: "name (type): url"
4. Print formatted list

**Pseudocode:**
```
function ListRegistries():
    manifest = load_manifest()
    
    if len(manifest.Registries) == 0:
        print "No registries configured"
        return
    
    for name, config in manifest.Registries:
        print f"{name} ({config.type}): {config.url}"
```

### Remove Registry
1. Parse registry name from arguments
2. Load manifest from arm.json
3. Check if registry exists → error if not found
4. If not --force, check dependencies:
   - Scan manifest.Dependencies for packages using this registry
   - If found → error "registry in use by packages: [list]"
5. Delete manifest.Registries[name]
6. Save manifest to arm.json
7. Print success message

**Pseudocode:**
```
function RemoveRegistry(name, force):
    manifest = load_manifest()
    
    if name not in manifest.Registries:
        return error "registry not found"
    
    if not force:
        dependent_packages = find_packages_using_registry(manifest, name)
        if len(dependent_packages) > 0:
            return error f"registry in use by: {dependent_packages}"
    
    delete manifest.Registries[name]
    save_manifest(manifest)
    
    return success
```

### Registry Factory
1. Receive registry name and configuration map
2. Extract "type" field from configuration
3. Convert configuration map to typed struct based on type
4. Create registry instance:
   - Git: NewGitRegistry(name, GitRegistryConfig)
   - GitLab: NewGitLabRegistry(name, GitLabRegistryConfig, configManager)
   - Cloudsmith: NewCloudsmithRegistry(name, CloudsmithRegistryConfig, configManager)
5. Return registry instance

**Pseudocode:**
```
function CreateRegistry(name, config_map):
    type = config_map["type"]
    
    switch type:
        case "git":
            git_config = convert_to_struct(config_map, GitRegistryConfig)
            return NewGitRegistry(name, git_config)
        case "gitlab":
            gitlab_config = convert_to_struct(config_map, GitLabRegistryConfig)
            return NewGitLabRegistry(name, gitlab_config, config_manager)
        case "cloudsmith":
            cloudsmith_config = convert_to_struct(config_map, CloudsmithRegistryConfig)
            return NewCloudsmithRegistry(name, cloudsmith_config, config_manager)
        default:
            return error "unsupported registry type"
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Duplicate registry name without --force | Error: "registry 'name' already exists, use --force to overwrite" |
| Duplicate registry name with --force | Overwrite existing configuration, print warning |
| Invalid URL format | Error: "invalid URL format" |
| Missing required field (owner, repository) | Error: "field 'X' is required for cloudsmith registries" |
| Remove registry with dependent packages | Error: "registry in use by packages: [list], use --force to remove" |
| Remove registry with --force | Remove registry and orphan dependent packages |
| Remove non-existent registry | Error: "registry 'name' not found" |
| List registries when none configured | Print: "No registries configured" |
| Empty registry name | Error: "registry name cannot be empty" |
| Registry name with special characters | Accept any non-empty string (no validation) |
| Network failure during validation | Defer validation to first use (add succeeds) |
| Corrupted manifest file | Error: "failed to parse manifest: [details]" |

## Dependencies

- **Manifest management** - Registry configuration stored in arm.json
- **Config management** - Authentication tokens stored in .armrc
- **Storage layer** - Registry cache and package storage
- **Version resolution** - Registry provides versions for constraint matching
- **Package installation** - Registry provides package content

## Implementation Mapping

**Source files:**
- `cmd/arm/main.go` - Command handlers: handleAddRegistry, handleAddGitRegistry, handleAddGitLabRegistry, handleAddCloudsmithRegistry, handleListRegistries, handleRemoveRegistry
- `internal/arm/service/service.go` - Service methods: AddGitRegistry, AddGitLabRegistry, AddCloudsmithRegistry, ListRegistries, RemoveRegistry
- `internal/arm/manifest/manager.go` - Manifest persistence: AddRegistry, GetRegistry, RemoveRegistry, ListRegistries
- `internal/arm/registry/factory.go` - Registry factory: CreateRegistry, convertMapToStruct
- `internal/arm/registry/git.go` - Git registry implementation: NewGitRegistry, GitRegistryConfig
- `internal/arm/registry/gitlab.go` - GitLab registry implementation: NewGitLabRegistry, GitLabRegistryConfig
- `internal/arm/registry/cloudsmith.go` - Cloudsmith registry implementation: NewCloudsmithRegistry, CloudsmithRegistryConfig
- `internal/arm/registry/registry.go` - Registry interface definition
- `internal/arm/config/manager.go` - Authentication key generation and storage

**Related specs:**
- `version-resolution.md` - Registries provide versions for resolution
- `package-installation.md` - Registries provide package content for installation
- `authentication.md` - Registries use .armrc for authentication tokens

## Examples

### Example 1: Add Git Registry

**Input:**
```bash
arm add registry git --url https://github.com/PatrickJS/awesome-cursorrules --branches main,develop awesome-rules
```

**Expected Output:**
```
Added git registry 'awesome-rules'
```

**Manifest (arm.json):**
```json
{
  "version": 1,
  "registries": {
    "awesome-rules": {
      "type": "git",
      "url": "https://github.com/PatrickJS/awesome-cursorrules",
      "branches": ["main", "develop"]
    }
  }
}
```

**Verification:**
- arm.json contains registry configuration
- Registry name is "awesome-rules"
- Registry type is "git"
- Branches array contains "main" and "develop"

### Example 2: Add GitLab Registry

**Input:**
```bash
arm add registry gitlab --url https://gitlab.example.com --project-id 123 my-gitlab
```

**Expected Output:**
```
Added gitlab registry 'my-gitlab'
Generated authentication key: gitlab.my-gitlab
```

**Manifest (arm.json):**
```json
{
  "version": 1,
  "registries": {
    "my-gitlab": {
      "type": "gitlab",
      "url": "https://gitlab.example.com",
      "projectId": "123",
      "apiVersion": "v4"
    }
  }
}
```

**Verification:**
- arm.json contains registry configuration
- .armrc contains authentication key entry
- Default apiVersion is "v4"

### Example 3: Add Cloudsmith Registry

**Input:**
```bash
arm add registry cloudsmith --url https://dl.cloudsmith.io --owner myorg --repo ai-rules my-cloudsmith
```

**Expected Output:**
```
Added cloudsmith registry 'my-cloudsmith'
Generated authentication key: cloudsmith.my-cloudsmith
```

**Manifest (arm.json):**
```json
{
  "version": 1,
  "registries": {
    "my-cloudsmith": {
      "type": "cloudsmith",
      "url": "https://dl.cloudsmith.io",
      "owner": "myorg",
      "repository": "ai-rules"
    }
  }
}
```

**Verification:**
- arm.json contains registry configuration
- .armrc contains authentication key entry
- Owner and repository fields are present

### Example 4: Duplicate Registry Error

**Input:**
```bash
arm add registry git --url https://github.com/example/repo existing-registry
```

**Expected Output:**
```
Error: registry 'existing-registry' already exists, use --force to overwrite
```

**Verification:**
- Command exits with non-zero status
- Manifest is unchanged
- Error message mentions --force flag

### Example 5: List Registries

**Input:**
```bash
arm list registries
```

**Expected Output:**
```
awesome-rules (git): https://github.com/PatrickJS/awesome-cursorrules
my-gitlab (gitlab): https://gitlab.example.com
my-cloudsmith (cloudsmith): https://dl.cloudsmith.io
```

**Verification:**
- All configured registries are listed
- Format: "name (type): url"
- One registry per line

### Example 6: Remove Registry with Dependencies

**Input:**
```bash
arm remove registry awesome-rules
```

**Expected Output:**
```
Error: registry 'awesome-rules' is in use by packages: awesome-rules/clean-code, awesome-rules/security
Use --force to remove anyway
```

**Verification:**
- Command exits with non-zero status
- Manifest is unchanged
- Error lists dependent packages

### Example 7: Remove Registry with Force

**Input:**
```bash
arm remove registry awesome-rules --force
```

**Expected Output:**
```
Removed registry 'awesome-rules'
Warning: 2 packages now reference a missing registry
```

**Verification:**
- Registry removed from manifest
- Dependent packages remain but are orphaned
- Warning message indicates orphaned packages

## Notes

### Registry Naming Model Differences

ARM supports two fundamentally different package naming models:

1. **Git Registries (User-Named Packages)**: The repository is a source of files. Users choose the package name when installing. The package name is just a label for the installation, not tied to repository structure.

2. **Non-Git Registries (Registry-Named Packages)**: Packages have explicit names in the registry. Users must use the exact package name when installing.

This distinction affects how packages are referenced:
- Git: `registry-name/any-name-you-want`
- GitLab/Cloudsmith: `registry-name/exact-package-name`

### Configuration Extensibility

Registry configuration uses `map[string]interface{}` to support arbitrary JSON fields. This allows:
- Future registry types without code changes
- Custom fields for specific registry implementations
- Backward compatibility when adding new fields

The factory pattern converts the generic map to typed structs for type safety in implementation.

### Authentication Integration

GitLab and Cloudsmith registries require authentication tokens stored in .armrc. The key generation process:
1. Generate unique key: `{type}.{name}` (e.g., "gitlab.my-gitlab")
2. Store key in .armrc with placeholder value
3. User manually updates .armrc with actual token
4. Registry reads token from .armrc during operations

See `authentication.md` for .armrc format and token resolution.

### Force Flag Behavior

The --force flag has different semantics for add vs remove:
- **Add**: Overwrite existing registry configuration (safe operation)
- **Remove**: Remove registry even if packages depend on it (dangerous operation)

This asymmetry reflects the risk level: overwriting configuration is recoverable, but orphaning packages may break installations.
