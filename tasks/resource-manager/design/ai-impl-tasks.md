# AI Implementation Tasks for Resource Manager Refactoring

## Overview

This document provides a detailed investigation of the codebase to identify specific tasks needed to implement the resource manager design that supports both rulesets and promptsets. The refactoring transforms ARM from a ruleset-only system to a unified resource manager.

## Current State Analysis

### Current Architecture
- **URF (Universal Rule Format)**: Currently only supports rulesets with `apiVersion: v1`, `kind: Ruleset`
- **Commands**: Ruleset-focused commands like `arm install`, `arm uninstall`, etc.
- **Data Structures**: All focused on rulesets (e.g., `LockFile.Rulesets`, `Manifest.Rulesets`)
- **Cache**: `RegistryRulesetCache` interface and `file_registry_ruleset_cache.go`
- **Index**: `IndexManager` only handles rulesets for `arm_index.*` generation

### Target Architecture
- **Resources**: Support both `kind: Ruleset` and `kind: Promptset`
- **Commands**: Unified package management with resource-specific subcommands
- **Data Structures**: Generic `packages` with `rulesets` and `promptsets` sections
- **Cache**: Generic package cache supporting both resource types
- **Index**: Rulesets impact `arm_index.*`, promptsets do not

## Detailed Implementation Tasks

### 1. Command Structure Refactoring (`cmd/`) ✅ COMPLETED

#### 1.1 Root Command Updates
**File**: `cmd/arm/main.go`
- [x] Update root command description from "AI rule rulesets" to "AI resources (rulesets and promptsets)"
- [x] Add new command structure supporting both unified and resource-specific commands

#### 1.2 New Command Structure
**Files**: All command files need major restructuring

**Current Commands** → **New Commands**:
- `arm install` → `arm install` (unified) + `arm install ruleset` + `arm install promptset`
- `arm uninstall` → `arm uninstall` (unified) + `arm uninstall ruleset` + `arm uninstall promptset`
- `arm update` → `arm update` (unified) + `arm update ruleset` + `arm update promptset`
- `arm upgrade` → `arm upgrade` (unified) + `arm upgrade ruleset` + `arm upgrade promptset`
- `arm outdated` → `arm outdated` (unified) + `arm outdated ruleset` + `arm outdated promptset`
- `arm list` → `arm list` (unified) + `arm list ruleset` + `arm list promptset`
- `arm info` → `arm info` (unified) + `arm info ruleset` + `arm info promptset`

#### 1.3 Command Implementation Tasks
**Files**: `cmd/arm/install.go`, `cmd/arm/uninstall.go`, `cmd/arm/update.go`, `cmd/arm/outdated.go`, `cmd/arm/list.go`, `cmd/arm/info.go`

- [x] **Install Commands**:
  - [x] Create unified `arm install` command that installs all configured packages
  - [x] Create `arm install ruleset REGISTRY/RULESET[@VERSION] SINK...` with priority, include/exclude flags
  - [x] Create `arm install promptset REGISTRY/PROMPTSET[@VERSION] SINK...` with include/exclude flags
  - [x] Update argument parsing to handle resource type detection

- [x] **Uninstall Commands**:
  - [x] Create unified `arm uninstall` command
  - [x] Create `arm uninstall ruleset REGISTRY/RULESET`
  - [x] Create `arm uninstall promptset REGISTRY/PROMPTSET`

- [x] **Update/Upgrade Commands**:
  - [x] Create unified `arm update` and `arm upgrade` commands
  - [x] Create resource-specific update/upgrade commands
  - [x] Handle version constraint logic for both resource types

- [x] **List/Info Commands**:
  - [x] Create unified `arm list` and `arm info` commands
  - [x] Create resource-specific list/info commands
  - [x] Update output formatting to distinguish resource types

- [x] **Outdated Commands**:
  - [x] Create unified `arm outdated` command
  - [x] Create resource-specific outdated commands
  - [x] Support `--output` format options (table, json, list)

#### 1.4 Configuration Commands
**File**: `cmd/arm/config.go`

- [x] **Registry Commands**: Already support the new design
- [x] **Sink Commands**: Already support the new design
- [x] **Resource Config Commands**:
  - [x] Create `arm config ruleset set REGISTRY/RULESET KEY VALUE`
  - [x] Create `arm config promptset set REGISTRY/PROMPTSET KEY VALUE`
  - [x] Support keys: `version`, `priority` (rulesets only), `sinks`, `includes`, `excludes`

#### 1.5 Utility Commands
**Files**: `cmd/arm/cache.go`, `cmd/arm/compile.go`

- [x] **Cache Commands**: Update to handle both resource types
- [x] **Compile Commands**: Update to support both rulesets and promptsets
- [x] Add `arm clean sinks` command for sink cleanup

### 2. Internal Package Refactoring

#### 2.1 URF → Resource Package Rename
**Directory**: `internal/urf/` → `internal/resource/`

- [ ] **Rename Package**: `urf` → `resource`
- [ ] **Update Types** (`internal/resource/types.go`):
  - [ ] Add `Promptset` struct with `apiVersion`, `kind: Promptset`, `metadata`, `spec.prompts`
  - [ ] Add `Prompt` struct with `name`, `description`, `body`
  - [ ] Update `Parser` interface to support both `Ruleset` and `Promptset`
  - [ ] Update `Compiler` interface to handle both resource types
  - [ ] Add `PromptGenerator` interface for prompt compilation
  - [ ] Update `CompileTarget` constants if needed

- [ ] **Update Parser** (`internal/resource/parser.go`):
  - [ ] Modify `IsURF` → `IsResource` to detect both rulesets and promptsets
  - [ ] Update `Parse` to return interface{} or union type
  - [ ] Add validation for promptset structure

- [ ] **Update Compiler** (`internal/resource/compiler.go`):
  - [ ] Support compilation of both rulesets and promptsets
  - [ ] Add prompt compilation logic (simpler than rules - no metadata)
  - [ ] Update `Compile` method signature

- [ ] **Update Generators**:
  - [ ] Create prompt generators for each target (cursor, amazonq, copilot, md)
  - [ ] Update rule generators if needed
  - [ ] Ensure prompt compilation produces content-only files (no frontmatter)

#### 2.2 Service Interface Updates
**File**: `internal/arm/service_interface.go`

- [ ] **Add Promptset Operations**:
  - [ ] `InstallPromptset(ctx context.Context, req *InstallPromptsetRequest) error`
  - [ ] `UninstallPromptset(ctx context.Context, registry, promptset string) error`
  - [ ] `UpdatePromptset(ctx context.Context, registry, promptset string) error`
  - [ ] `UpdatePromptsetConfig(ctx context.Context, registry, promptset, field, value string) error`

- [ ] **Add Unified Operations**:
  - [ ] `InstallAll(ctx context.Context) error` (replaces `InstallManifest`)
  - [ ] `UpdateAll(ctx context.Context) error`
  - [ ] `UpgradeAll(ctx context.Context) error`
  - [ ] `UninstallAll(ctx context.Context) error`

- [ ] **Update Info Operations**:
  - [ ] `ShowPromptsetInfo(ctx context.Context, promptsets []string) error`
  - [ ] Update `ShowOutdated` to handle both resource types
  - [ ] Update `ShowList` to handle both resource types

#### 2.3 Manifest Manager Updates
**File**: `internal/manifest/types.go`

- [ ] **Update Manifest Structure**:
  ```go
  type Manifest struct {
      Registries map[string]map[string]interface{} `json:"registries,omitempty"`
      Packages   struct {
          Rulesets   map[string]map[string]Entry `json:"rulesets,omitempty"`
          Promptsets map[string]map[string]Entry `json:"promptsets,omitempty"`
      } `json:"packages"`
      Sinks      map[string]SinkConfig `json:"sinks,omitempty"`
  }
  ```

- [ ] **Update Entry Structure**:
  - [ ] Make `Priority` optional (only for rulesets)
  - [ ] Add validation to ensure promptsets don't have priority

**File**: `internal/manifest/manager.go`
- [ ] Update all methods to handle both `rulesets` and `promptsets` sections
- [ ] Add validation for resource-specific constraints
- [ ] Update serialization/deserialization logic

#### 2.4 Lockfile Manager Updates
**File**: `internal/lockfile/types.go`

- [ ] **Update LockFile Structure**:
  ```go
  type LockFile struct {
      Rulesets   map[string]map[string]Entry `json:"rulesets"`
      Promptsets map[string]map[string]Entry `json:"promptsets"`
  }
  ```

**File**: `internal/lockfile/manager.go`
- [ ] Update all methods to handle both resource types
- [ ] Add methods for promptset lock management
- [ ] Update serialization/deserialization logic

#### 2.5 Cache System Updates
**File**: `internal/cache/registry_ruleset_cache.go` → `internal/cache/registry_package_cache.go`

- [ ] **Rename Interface**: `RegistryRulesetCache` → `RegistryPackageCache`
- [ ] **Update Methods**:
  - [ ] `GetRulesetVersion` → `GetPackageVersion` (with resource type parameter)
  - [ ] `SetRulesetVersion` → `SetPackageVersion` (with resource type parameter)
  - [ ] Add resource type handling throughout

**File**: `internal/cache/file_registry_ruleset_cache.go` → `internal/cache/file_registry_package_cache.go`
- [ ] Rename file and update implementation
- [ ] Update cache directory structure to support both resource types
- [ ] Update file naming conventions

**File**: `internal/cache/file_git_repo_cache.go`
- [ ] Update to store cloned repo as `repository` (remove name requirement)
- [ ] Simplify structure per design notes

**File**: `internal/cache/manager.go`
- [ ] Update to work with new package cache interface
- [ ] Handle both resource types in cleanup operations

#### 2.6 Index Manager Updates
**File**: `internal/index/manager.go`

- [ ] **Update Interface**: Support both rulesets and promptsets
- [ ] **Ruleset Operations**: Continue to impact `arm_index.*` generation
- [ ] **Promptset Operations**: Do NOT impact `arm_index.*` generation
- [ ] **Update Methods**:
  - [ ] `Create` method should accept resource type parameter
  - [ ] `Read` method should return resource type information
  - [ ] `Delete` method should handle both resource types

**File**: `internal/index/generator.go`
- [ ] Update to only generate index for rulesets
- [ ] Ignore promptsets in index generation

#### 2.7 Installer Updates
**File**: `internal/installer/types.go`

- [ ] **Add Promptset Type**:
  ```go
  type Promptset struct {
      Registry   string
      Promptset  string
      Version    string
      Path       string
      FilePaths  []string
  }
  ```

- [ ] **Update Installer Interface**:
  - [ ] Add methods for promptset installation
  - [ ] Update existing methods to handle both resource types

**File**: `internal/installer/installer.go`
- [ ] Update interface to support both resource types
- [ ] Add resource type parameter to installation methods

**File**: `internal/installer/urf_processor.go` → `internal/installer/resource_processor.go`
- [ ] Rename file and update implementation
- [ ] Support processing both rulesets and promptsets
- [ ] Update compilation logic

#### 2.8 Registry Updates
**File**: `internal/registry/types.go`

- [ ] **Update Registry Interface**:
  - [ ] `ListVersions` should work for both rulesets and promptsets
  - [ ] `ResolveVersion` should work for both resource types
  - [ ] `GetContent` should work for both resource types

**Files**: All registry implementations (`git_registry.go`, `gitlab_registry.go`, `cloudsmith_registry.go`)
- [ ] Update to handle both resource types
- [ ] Ensure they can discover and serve both rulesets and promptsets

#### 2.9 UI Updates
**File**: `internal/ui/ui.go`

- [ ] **Align with Design Outputs**: Review design examples and align UI output
- [ ] **Update Display Methods**: Support both resource types in all display methods
- [ ] **Add Resource Type Indicators**: Show whether output is for ruleset or promptset
- [ ] **Update Table Formats**: Support unified and resource-specific table layouts

### 3. Scripts and Workflows Updates

#### 3.1 Workflow Scripts
**Directory**: `scripts/workflows/`

- [ ] **Update All Workflow Scripts**:
  - [ ] `git/sample-git-workflow.sh`: Update to use new command structure
  - [ ] `gitlab/sample-gitlab-workflow.sh`: Update to use new command structure
  - [ ] `cloudsmith/sample-cloudsmith-workflow.sh`: Update to use new command structure
  - [ ] `compile/sample-compile-workflow.sh`: Update to support both resource types

- [ ] **Update Command Examples**:
  - [ ] Replace `arm install REGISTRY/RULESET` with `arm install ruleset REGISTRY/RULESET`
  - [ ] Add examples for `arm install promptset REGISTRY/PROMPTSET`
  - [ ] Update configuration examples to use new structure
  - [ ] Add examples for unified commands (`arm install`, `arm update`, etc.)

- [ ] **Update Sandbox Configurations**:
  - [ ] Update `arm.json` files to use new `packages` structure
  - [ ] Update `arm-lock.json` files to include both `rulesets` and `promptsets`
  - [ ] Add example promptset configurations

#### 3.2 Example Files
**Directory**: `scripts/workflows/compile/example-rulesets/`

- [ ] **Add Example Promptsets**:
  - [ ] Create example promptset YAML files
  - [ ] Add compilation examples for promptsets
  - [ ] Update workflow scripts to demonstrate both resource types

### 4. Testing Updates

#### 4.1 Unit Tests
- [ ] **Update All Test Files**: Rename `urf` package references to `resource`
- [ ] **Add Promptset Tests**: Create tests for promptset parsing, compilation, and installation
- [ ] **Update Existing Tests**: Modify ruleset tests to work with new structure
- [ ] **Add Integration Tests**: Test unified commands and resource-specific commands

#### 4.2 Integration Tests
- [ ] **Update Service Tests**: Test both ruleset and promptset operations
- [ ] **Update Registry Tests**: Test discovery and serving of both resource types
- [ ] **Update Cache Tests**: Test caching of both resource types
- [ ] **Update Installer Tests**: Test installation of both resource types

### 5. Documentation Updates

#### 5.1 Command Documentation
- [ ] **Update Help Text**: All commands need updated help text
- [ ] **Update Examples**: All examples need to reflect new command structure
- [ ] **Add Resource-Specific Docs**: Document differences between rulesets and promptsets

#### 5.2 API Documentation
- [ ] **Update Service Interface Docs**: Document new methods and parameters
- [ ] **Update Type Documentation**: Document new types and structures
- [ ] **Add Migration Guide**: Document breaking changes from v2 to v3

## Implementation Decisions

### 1. Resource Type Detection
**Decision**: Parse the `kind` field in the YAML to detect resource type.
- YAML files will be parsed to check the `kind` field (`Ruleset` or `Promptset`)
- This provides clear, explicit resource type identification
- Consistent with Kubernetes-style resource definitions

### 2. Backward Compatibility
**Decision**: Clean break with no backward compatibility.
- No automatic migration of old `arm.json` files
- No support for old command syntax with deprecation warnings
- This is acceptable since it's early in development with expected churn

### 3. Priority Handling
**Decision**: Each resource should be validated against their specific schema using common Go modules for YAML schema validation.
- Implement schema validation for both rulesets and promptsets
- Promptset schema will not include priority field
- Validation errors will be returned if priority is specified for promptsets

### 4. Index Generation
**Decision**: Promptsets should be tracked in arm-index.json files, but not in arm_index.* files.
- Promptsets will be tracked in the local inventory (arm-index.json)
- Promptsets will NOT impact arm_index.* file generation (which only handles ruleset conflict resolution)
- Index manager will handle both resource types but only generate index files for rulesets

### 5. Compilation Differences
**Decision**: All promptset compilation targets produce identical output.
- Promptsets compile to content-only files with no metadata or frontmatter
- No target-specific formatting for promptsets
- Simpler compilation process compared to rulesets

### 6. Command Aliases
**Decision**: No aliases will be supported.
- Keep command structure clean and explicit
- Avoid confusion with multiple ways to invoke the same functionality

### 7. Error Handling
**Decision**: Detailed error reporting with resource type context.
- Mixed resource operations should report specific failures for each resource type
- Clear indication of which rulesets/promptsets succeeded or failed
- Partial failures should be clearly communicated to the user

### 8. Performance Considerations
**Decision**: No specific performance improvements needed at this time.
- Current architecture should handle both resource types adequately
- Performance optimizations can be addressed in future iterations if needed

## Implementation Priority

### Phase 1: Core Infrastructure
1. Rename `urf` package to `resource`
2. Add promptset types and parsing
3. Update service interface
4. Update manifest and lockfile structures

### Phase 2: Command Structure
1. Implement new command hierarchy
2. Add resource-specific commands
3. Update argument parsing and validation
4. Update help and documentation

### Phase 3: Backend Services
1. Update cache system
2. Update registry implementations
3. Update installer system
4. Update index manager

### Phase 4: Testing and Polish
1. Update all tests
2. Update workflow scripts
3. Add integration tests
4. Performance testing

### Phase 5: Documentation and Migration
1. Update all documentation
2. Create migration guide
3. Update examples
4. Final testing and validation

## Risk Assessment

### High Risk
- **Breaking Changes**: Complete command structure change
- **Data Migration**: Existing configurations and lock files
- **Testing Coverage**: Large surface area of changes

### Medium Risk
- **Performance**: New resource type handling
- **Cache Management**: Dual resource type caching
- **Registry Compatibility**: Ensuring all registries work with both types

### Low Risk
- **UI Updates**: Mostly cosmetic changes
- **Documentation**: Well-defined target state
- **Script Updates**: Straightforward command replacements

## Success Criteria

1. **Functional**: All new commands work as specified in design
2. **Compatible**: Existing ruleset functionality preserved
3. **Performance**: No significant performance degradation
4. **Usability**: Clear distinction between rulesets and promptsets
5. **Maintainable**: Clean separation of concerns between resource types
6. **Testable**: Comprehensive test coverage for both resource types
