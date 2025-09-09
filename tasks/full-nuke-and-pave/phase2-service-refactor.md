# Phase 2: Service Layer Refactoring - Detailed Tasks

## Task 2.1: Break Down Large Methods

### Current Issues Analysis

#### `InstallRuleset` Method (120 lines)
**Problems:**
- Does 6 different things: validate, resolve, download, update manifest, update lockfile, install to sinks
- Hard to test individual steps
- Error handling scattered throughout
- Mixed abstraction levels

**Refactoring Plan:**
```go
func (s *ArmService) InstallRuleset(ctx context.Context, req InstallRequest) error {
    // Validate input
    if err := s.validateInstallRequest(req); err != nil {
        return fmt.Errorf("invalid install request: %w", err)
    }

    // Resolve version
    resolved, err := s.versionService.ResolveVersion(ctx, req)
    if err != nil {
        return fmt.Errorf("failed to resolve version: %w", err)
    }

    // Download content
    content, err := s.contentService.DownloadContent(ctx, resolved)
    if err != nil {
        return fmt.Errorf("failed to download content: %w", err)
    }

    // Update tracking files
    if err := s.updateTrackingFiles(ctx, req, resolved, content); err != nil {
        return fmt.Errorf("failed to update tracking files: %w", err)
    }

    // Install to sinks
    return s.installerService.InstallToSinks(ctx, req, resolved, content)
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

#### Step 2.1.1: Extract InstallRuleset Components

**Create `InstallRequest` type:**
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
```

**Extract validation:**
```go
func (s *ArmService) validateInstallRequest(req InstallRequest) error {
    if err := req.Validate(); err != nil {
        return err
    }

    // Check registry exists
    if !s.registryExists(req.Registry) {
        return fmt.Errorf("registry %s not configured", req.Registry)
    }

    return nil
}
```

**Extract version resolution:**
```go
type ResolvedVersion struct {
    Request InstallRequest
    Version types.Version
    Content types.ContentSelector
}

func (s *versionService) ResolveVersion(ctx context.Context, req InstallRequest) (*ResolvedVersion, error) {
    registry, err := s.getRegistry(req.Registry)
    if err != nil {
        return nil, err
    }

    version := req.Version
    if version == "" {
        version = "latest"
    }

    resolved, err := registry.ResolveVersion(ctx, expandVersionShorthand(version))
    if err != nil {
        return nil, err
    }

    return &ResolvedVersion{
        Request: req,
        Version: resolved.Version,
        Content: types.ContentSelector{
            Include: req.Include,
            Exclude: req.Exclude,
        },
    }, nil
}
```

#### Step 2.1.2: Extract Content Management

**Create content service:**
```go
type ContentService interface {
    DownloadContent(ctx context.Context, resolved *ResolvedVersion) ([]types.File, error)
    ValidateContent(files []types.File, checksum string) error
}

type contentService struct {
    registryFactory registry.Factory
}

func (s *contentService) DownloadContent(ctx context.Context, resolved *ResolvedVersion) ([]types.File, error) {
    registry, err := s.registryFactory.GetRegistry(resolved.Request.Registry)
    if err != nil {
        return nil, err
    }

    return registry.GetContent(ctx, resolved.Version, resolved.Content)
}
```

#### Step 2.1.3: Extract Tracking File Updates

**Create tracking service:**
```go
type TrackingService interface {
    UpdateManifest(ctx context.Context, req InstallRequest) error
    UpdateLockfile(ctx context.Context, req InstallRequest, resolved *ResolvedVersion, files []types.File) error
}

func (s *trackingService) UpdateTrackingFiles(ctx context.Context, req InstallRequest, resolved *ResolvedVersion, files []types.File) error {
    if err := s.UpdateManifest(ctx, req); err != nil {
        return fmt.Errorf("failed to update manifest: %w", err)
    }

    if err := s.UpdateLockfile(ctx, req, resolved, files); err != nil {
        return fmt.Errorf("failed to update lockfile: %w", err)
    }

    return nil
}
```

## Task 2.2: Extract Business Logic Components

### Target Architecture
```
internal/arm/
├── service.go              # Main orchestration service
├── installer_service.go    # Installation logic
├── version_service.go      # Version resolution
├── content_service.go      # Content download/validation
├── tracking_service.go     # Manifest/lockfile management
├── sync_service.go         # Sink synchronization
└── types.go               # Service-specific types
```

### Implementation Steps

#### Step 2.2.1: Create Version Service
```go
type VersionService interface {
    ResolveVersion(ctx context.Context, req InstallRequest) (*ResolvedVersion, error)
    FindOutdatedRulesets(ctx context.Context, installed []InstalledRuleset, clients map[string]registry.Registry) ([]OutdatedRuleset, error)
    ExpandVersionShorthand(constraint string) string
}

type versionService struct {
    registryFactory registry.Factory
}

func NewVersionService(registryFactory registry.Factory) VersionService {
    return &versionService{
        registryFactory: registryFactory,
    }
}
```

#### Step 2.2.2: Create Installer Service
```go
type InstallerService interface {
    InstallToSinks(ctx context.Context, req InstallRequest, resolved *ResolvedVersion, files []types.File) error
    UninstallFromSinks(ctx context.Context, registry, ruleset string) error
    ListInstalled(ctx context.Context) ([]Installation, error)
}

type installerService struct {
    configManager config.Manager
    sinkMatcher   SinkMatcher
}

func (s *installerService) InstallToSinks(ctx context.Context, req InstallRequest, resolved *ResolvedVersion, files []types.File) error {
    sinks, err := s.configManager.GetSinks(ctx)
    if err != nil {
        return err
    }

    rulesetKey := req.Registry + "/" + req.Ruleset
    matchingSinks := s.sinkMatcher.FindMatchingSinks(rulesetKey, sinks)

    if len(matchingSinks) == 0 {
        slog.WarnContext(ctx, "No sinks match ruleset", "ruleset", rulesetKey)
        return nil
    }

    return s.installToMatchingSinks(ctx, matchingSinks, req, resolved, files)
}
```

#### Step 2.2.3: Create Sync Service
```go
type SyncService interface {
    SyncSink(ctx context.Context, sinkName string, sink *config.SinkConfig) error
    SyncRemovedSink(ctx context.Context, removedSink *config.SinkConfig) error
    SyncAllSinks(ctx context.Context) error
}

type syncService struct {
    manifestManager manifest.Manager
    lockfileManager lockfile.Manager
    installerService InstallerService
    sinkMatcher     SinkMatcher
}
```

## Task 2.3: Improve Dependency Injection

### Current Issues
- Hard-coded dependencies (`config.NewFileManager()`)
- No interfaces for dependencies
- Difficult to test with mocks

### Target State
```go
type ServiceDependencies struct {
    ConfigManager   config.Manager
    ManifestManager manifest.Manager
    LockfileManager lockfile.Manager
    RegistryFactory registry.Factory
}

func NewArmService(deps ServiceDependencies) *ArmService {
    versionService := NewVersionService(deps.RegistryFactory)
    contentService := NewContentService(deps.RegistryFactory)
    trackingService := NewTrackingService(deps.ManifestManager, deps.LockfileManager)
    installerService := NewInstallerService(deps.ConfigManager)
    syncService := NewSyncService(deps.ManifestManager, deps.LockfileManager, installerService)

    return &ArmService{
        versionService:   versionService,
        contentService:   contentService,
        trackingService:  trackingService,
        installerService: installerService,
        syncService:      syncService,
    }
}

// For backward compatibility
func NewArmServiceWithDefaults() *ArmService {
    return NewArmService(ServiceDependencies{
        ConfigManager:   config.NewFileManager(),
        ManifestManager: manifest.NewFileManager(),
        LockfileManager: lockfile.NewFileManager(),
        RegistryFactory: registry.NewFactory(),
    })
}
```

### Implementation Steps

#### Step 2.3.1: Define Interfaces
```go
// Move to separate interface files
type ConfigManager interface {
    GetSinks(ctx context.Context) (map[string]config.SinkConfig, error)
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
