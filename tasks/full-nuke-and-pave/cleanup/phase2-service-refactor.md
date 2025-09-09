# Phase 2: Service Layer Refactoring - READY TO START

## Current State Analysis

### ⚠️ Critical Issues in `internal/arm/service.go` (600+ lines)

#### `InstallRuleset` Method (120+ lines) - NEEDS IMMEDIATE REFACTORING
**Current Problems:**
- Does 6 different things: validate, resolve, download, update manifest, update lockfile, install to sinks
- Hard to test individual steps
- Error handling scattered throughout
- Mixed abstraction levels
- Hard-coded dependency creation (`registry.NewRegistry`)

#### Other Large Methods:
- `GetOutdatedRulesets` (60+ lines) - nested loops, complex version comparison
- `SyncSink` (50+ lines) - complex installation/removal logic
- `installFromLockfile` (30+ lines) - duplicate installation logic

## Task 2.1: Break Down Large Methods - PRIORITY 1

**Target Refactored Structure:**
```go
func (s *ArmService) InstallRuleset(ctx context.Context, registry, ruleset, version string, include, exclude []string) error {
    // Create install request
    req := InstallRequest{
        Registry: registry,
        Ruleset:  ruleset,
        Version:  version,
        Include:  include,
        Exclude:  exclude,
    }

    // Validate input
    if err := s.validateInstallRequest(ctx, req); err != nil {
        return fmt.Errorf("invalid install request: %w", err)
    }

    // Resolve version
    resolved, err := s.resolveVersion(ctx, req)
    if err != nil {
        return fmt.Errorf("failed to resolve version: %w", err)
    }

    // Download content
    files, err := s.downloadContent(ctx, req, resolved)
    if err != nil {
        return fmt.Errorf("failed to download content: %w", err)
    }

    // Update tracking files
    if err := s.updateTrackingFiles(ctx, req, resolved, files); err != nil {
        return fmt.Errorf("failed to update tracking files: %w", err)
    }

    // Install to sinks
    return s.installToSinks(ctx, req, resolved, files)
}
```

#### `Outdated` Method (60 lines)
**Problems:**
- Nested loops with complex logic
- Registry client creation mixed with version checking
- Hard to test version comparison logic

**Refactoring Plan:**
```go
func (s *ArmService) Outdated(ctx context.Context) ([]OutdatedRuleset, error) {
    installed, err := s.getInstalledRulesets(ctx)
    if err != nil {
        return nil, err
    }

    registryClients, err := s.createRegistryClients(ctx)
    if err != nil {
        return nil, err
    }

    return s.versionService.FindOutdatedRulesets(ctx, installed, registryClients)
}
```

### Implementation Steps

### IMMEDIATE ACTION PLAN

#### Step 2.1.1: Extract InstallRuleset Components (2-3 hours)

**1. Create helper types in `internal/arm/types.go`:**
```go
type InstallRequest struct {
    Registry string
    Ruleset  string
    Version  string
    Include  []string
    Exclude  []string
}

func (r InstallRequest) Validate() error {
    if r.Registry == "" {
        return errors.New("registry is required")
    }
    if r.Ruleset == "" {
        return errors.New("ruleset is required")
    }
    return nil
}

type ResolvedInstall struct {
    Request InstallRequest
    Version types.Version
    Files   []types.File
}
```

**2. Extract validation method:**
```go
func (s *ArmService) validateInstallRequest(ctx context.Context, req InstallRequest) error {
    if err := req.Validate(); err != nil {
        return err
    }

    // Check registry exists in manifest
    registries, err := s.manifestManager.GetRawRegistries(ctx)
    if err != nil {
        return fmt.Errorf("failed to get registries: %w", err)
    }
    if _, exists := registries[req.Registry]; !exists {
        return fmt.Errorf("registry %s not configured", req.Registry)
    }

    return nil
}
```

**3. Extract version resolution:**
```go
func (s *ArmService) resolveVersion(ctx context.Context, req InstallRequest) (types.Version, error) {
    registries, err := s.manifestManager.GetRawRegistries(ctx)
    if err != nil {
        return types.Version{}, fmt.Errorf("failed to get registries: %w", err)
    }

    registryConfig := registries[req.Registry]
    registryClient, err := registry.NewRegistry(req.Registry, registryConfig)
    if err != nil {
        return types.Version{}, fmt.Errorf("failed to create registry: %w", err)
    }

    version := req.Version
    if version == "" {
        version = "latest"
    }

    resolved, err := registryClient.ResolveVersion(ctx, expandVersionShorthand(version))
    if err != nil {
        return types.Version{}, err
    }

    return resolved.Version, nil
}
```

**4. Extract content download:**
```go
func (s *ArmService) downloadContent(ctx context.Context, req InstallRequest, version types.Version) ([]types.File, error) {
    registries, err := s.manifestManager.GetRawRegistries(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to get registries: %w", err)
    }

    registryConfig := registries[req.Registry]
    registryClient, err := registry.NewRegistry(req.Registry, registryConfig)
    if err != nil {
        return nil, fmt.Errorf("failed to create registry: %w", err)
    }

    selector := types.ContentSelector{
        Include: req.Include,
        Exclude: req.Exclude,
    }

    return registryClient.GetContent(ctx, version, selector)
}
```

**5. Extract tracking file updates:**
```go
func (s *ArmService) updateTrackingFiles(ctx context.Context, req InstallRequest, version types.Version, files []types.File) error {
    // Update manifest
    manifestEntry := manifest.Entry{
        Version: req.Version,
        Include: req.Include,
        Exclude: req.Exclude,
    }
    if err := s.manifestManager.CreateEntry(ctx, req.Registry, req.Ruleset, manifestEntry); err != nil {
        if err := s.manifestManager.UpdateEntry(ctx, req.Registry, req.Ruleset, manifestEntry); err != nil {
            return fmt.Errorf("failed to update manifest: %w", err)
        }
    }

    // Update lockfile
    checksum := lockfile.GenerateChecksum(files)
    lockEntry := &lockfile.Entry{
        Version:  version.Version,
        Display:  version.Display,
        Checksum: checksum,
    }
    if err := s.lockFileManager.CreateEntry(ctx, req.Registry, req.Ruleset, lockEntry); err != nil {
        if err := s.lockFileManager.UpdateEntry(ctx, req.Registry, req.Ruleset, lockEntry); err != nil {
            return fmt.Errorf("failed to update lockfile: %w", err)
        }
    }

    return nil
}
```

**6. Extract sink installation:**
```go
func (s *ArmService) installToSinks(ctx context.Context, req InstallRequest, version types.Version, files []types.File) error {
    sinks, err := s.configManager.GetSinks(ctx)
    if err != nil {
        return fmt.Errorf("failed to get sinks: %w", err)
    }

    rulesetKey := req.Registry + "/" + req.Ruleset
    var matchingSinkNames []string

    for sinkName, sink := range sinks {
        if s.matchesSink(rulesetKey, &sink) {
            matchingSinkNames = append(matchingSinkNames, sinkName)
            installer := installer.NewInstaller(&sink)
            for _, dir := range sink.Directories {
                if err := installer.Install(ctx, dir, req.Registry, req.Ruleset, version.Display, files); err != nil {
                    return fmt.Errorf("failed to install to directory %s: %w", dir, err)
                }
            }
        }
    }

    if len(matchingSinkNames) == 0 {
        slog.WarnContext(ctx, "Ruleset not targeted by any sinks", "registry", req.Registry, "ruleset", req.Ruleset)
    } else {
        slog.InfoContext(ctx, "Ruleset installed to sinks", "registry", req.Registry, "ruleset", req.Ruleset, "sinks", matchingSinkNames)
    }

    return nil
}
```

## Task 2.2: Extract Business Logic Components - PRIORITY 2

### AFTER Task 2.1 - Create Service Components (4-5 hours)

**Target Architecture:**
```
internal/arm/
├── service.go              # Main orchestration (150 lines max)
├── installer_service.go    # Installation/uninstallation logic
├── version_service.go      # Version resolution and comparison
├── content_service.go      # Content download/validation
├── tracking_service.go     # Manifest/lockfile management
├── sync_service.go         # Sink synchronization
└── types.go               # Service-specific types
```

#### Step 2.2.1: Create Version Service
```go
type VersionService interface {
    ResolveVersion(ctx context.Context, registryName, constraint string) (types.Version, error)
    FindOutdatedRulesets(ctx context.Context) ([]OutdatedRuleset, error)
    ExpandVersionShorthand(constraint string) string
}
```

#### Step 2.2.2: Create Installer Service
```go
type InstallerService interface {
    InstallToSinks(ctx context.Context, req InstallRequest, version types.Version, files []types.File) error
    UninstallFromSinks(ctx context.Context, registry, ruleset string) error
    SyncSink(ctx context.Context, sinkName string, sink *config.SinkConfig) error
}
```

#### Step 2.2.3: Create Content Service
```go
type ContentService interface {
    DownloadContent(ctx context.Context, registryName string, version types.Version, selector types.ContentSelector) ([]types.File, error)
    ValidateContent(files []types.File, checksum string) error
}
```

#### Step 2.2.4: Create Tracking Service
```go
type TrackingService interface {
    UpdateManifest(ctx context.Context, req InstallRequest) error
    UpdateLockfile(ctx context.Context, req InstallRequest, version types.Version, files []types.File) error
    RemoveEntry(ctx context.Context, registry, ruleset string) error
}
```

## Task 2.3: Improve Dependency Injection - PRIORITY 3

### Current Issues in `NewArmService()`
- Hard-coded dependencies (`config.NewFileManager()`, `manifest.NewFileManager()`, `lockfile.NewFileManager()`)
- No interfaces for dependencies
- Difficult to test with mocks
- Single constructor approach

### Target State (After Service Extraction)
```go
type ServiceDependencies struct {
    ConfigManager   config.Manager
    ManifestManager manifest.Manager
    LockfileManager lockfile.Manager
}

func NewArmService(deps ServiceDependencies) *ArmService {
    return &ArmService{
        configManager:   deps.ConfigManager,
        manifestManager: deps.ManifestManager,
        lockFileManager: deps.LockfileManager,
    }
}

// For backward compatibility and CLI usage
func NewArmServiceWithDefaults() *ArmService {
    return NewArmService(ServiceDependencies{
        ConfigManager:   config.NewFileManager(),
        ManifestManager: manifest.NewFileManager(),
        LockfileManager: lockfile.NewFileManager(),
    })
}
```

### Implementation Steps
1. **Update constructor** to accept dependencies
2. **Update CLI** to use `NewArmServiceWithDefaults()`
3. **Create test constructor** for unit tests with mocks

## SUMMARY: Phase 2 Action Plan

### 🔥 IMMEDIATE PRIORITY (Next 2-3 hours)
**Task 2.1.1: Break down `InstallRuleset` method**
1. Extract `validateInstallRequest` method
2. Extract `resolveVersion` method
3. Extract `downloadContent` method
4. Extract `updateTrackingFiles` method
5. Extract `installToSinks` method
6. Update main `InstallRuleset` to orchestrate these methods

**Expected Outcome:**
- `InstallRuleset` reduced from 120+ lines to ~30 lines
- Each extracted method <25 lines
- Clear separation of concerns
- Easier to test individual steps

### 🔄 FOLLOW-UP PRIORITY (Next 4-5 hours)
**Task 2.1.2: Break down other large methods**
- `GetOutdatedRulesets` (60+ lines)
- `SyncSink` (50+ lines)
- `installFromLockfile` (30+ lines)

### 🏢 FUTURE WORK (After method extraction)
**Task 2.2: Extract service components**
**Task 2.3: Improve dependency injection**

### ✅ SUCCESS CRITERIA
- No methods >50 lines
- `service.go` file <400 lines
- Clear single responsibility per method
- Improved testability context.Context) (map[string]config.SinkConfig, error)
    GetSink(ctx context.Context, name string) (*config.SinkConfig, error)
    // ... other methods
}

type ManifestManager interface {
    GetEntries(ctx context.Context) (map[string]map[string]manifest.Entry, error)
    CreateEntry(ctx context.Context, registry, ruleset string, entry manifest.Entry) error
    // ... other methods
}
```

#### Step 2.3.2: Update Service Constructor
- Accept interfaces instead of concrete types
- Provide default constructor for CLI usage
- Enable dependency injection for tests

#### Step 2.3.3: Update CLI Integration
```go
// In cmd/arm/main.go
func main() {
    armService = arm.NewArmServiceWithDefaults()
    // ... rest of main
}
```

## Acceptance Criteria

### Task 2.1 Complete When:
- [ ] No method in service layer >50 lines
- [ ] Each method has single responsibility
- [ ] Complex operations broken into testable components
- [ ] Clear error handling at each step

### Task 2.2 Complete When:
- [ ] Business logic separated into focused services
- [ ] Each service has clear interface and responsibility
- [ ] Main service orchestrates without business logic
- [ ] Services are independently testable

### Task 2.3 Complete When:
- [ ] All dependencies injected through constructor
- [ ] Interfaces defined for all external dependencies
- [ ] Default constructor available for CLI usage
- [ ] Easy to create test doubles for all dependencies

## Testing Strategy

### Unit Tests for Each Service
- Mock all external dependencies
- Test happy path and error scenarios
- Test edge cases and validation

### Integration Tests
- Test service composition
- Test end-to-end workflows
- Test error propagation between services

## Time Estimate: 10-12 hours total
- Task 2.1: 5-6 hours
- Task 2.2: 4-5 hours
- Task 2.3: 1-2 hours
