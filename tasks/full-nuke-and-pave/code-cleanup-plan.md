# ARM Code Cleanup Plan

## Overview
Refactor ARM codebase for better readability, organization, and test coverage while following clean code principles.

## Phase 1: CLI Layer Cleanup (High Priority)

### Task 1.1: Split CLI Commands
**Goal:** Break down monolithic main.go into focused command files
**Files:** `cmd/arm/main.go` → multiple command files
**Effort:** 2-3 hours

- Extract commands into separate files:
  - `cmd/arm/install.go` - install command logic
  - `cmd/arm/config.go` - config command logic (already exists, needs cleanup)
  - `cmd/arm/cache.go` - cache command logic (already exists)
  - `cmd/arm/info.go` - info/list/outdated commands
  - `cmd/arm/version.go` - version command
- Create `cmd/arm/parser.go` for argument parsing utilities
- Create `cmd/arm/output.go` for formatting utilities
- Keep `main.go` minimal with just root command setup

### Task 1.2: Extract CLI Utilities
**Goal:** Remove business logic from CLI handlers
**Files:** New utility files in `cmd/arm/`
**Effort:** 1-2 hours

- Move `parseRulesetArg`, `splitOnFirst`, `findFirst` to `parser.go`
- Move `printRulesetInfo` and table formatting to `output.go`
- Replace manual string parsing with proper parser functions
- Add input validation at CLI layer

### Task 1.3: Improve Error Handling
**Goal:** Consistent error messages and handling
**Files:** All CLI command files
**Effort:** 1 hour

- Standardize error message format
- Add context to error messages
- Remove error logging from CLI (let service layer handle)

## Phase 2: Service Layer Refactoring (Critical Priority)

### Task 2.1: Break Down Large Methods
**Goal:** Split 100+ line methods into focused functions
**Files:** `internal/arm/service.go`
**Effort:** 4-5 hours

- `InstallRuleset` (120 lines) → split into:
  - `validateInstallRequest`
  - `resolveRulesetVersion`
  - `downloadRulesetContent`
  - `updateManifestAndLock`
  - `installToSinks`
- `Install` (40 lines) → split into:
  - `determineInstallStrategy`
  - `installFromManifest`
  - `installFromLockfile`
- `Outdated` (60 lines) → split into:
  - `collectOutdatedRulesets`
  - `checkRulesetVersions`

### Task 2.2: Extract Business Logic
**Goal:** Create focused service components
**Files:** New files in `internal/arm/`
**Effort:** 3-4 hours

- Create `internal/arm/installer_service.go` for installation logic
- Create `internal/arm/version_service.go` for version resolution
- Create `internal/arm/sync_service.go` for sink synchronization
- Update main service to orchestrate these components

### Task 2.3: Improve Dependency Injection
**Goal:** Make dependencies explicit and testable
**Files:** `internal/arm/service.go`
**Effort:** 2 hours

- Add constructor with explicit dependencies
- Create interfaces for all dependencies
- Remove hard-coded `NewFileManager()` calls

## Phase 3: Add Comprehensive Tests (High Priority)

### Task 3.1: Service Layer Unit Tests
**Goal:** 80%+ coverage on core business logic
**Files:** New test files in `internal/arm/`
**Effort:** 6-8 hours

- `service_test.go` - main service orchestration tests
- `installer_service_test.go` - installation logic tests
- `version_service_test.go` - version resolution tests
- `sync_service_test.go` - sink sync tests
- Mock all external dependencies (registry, filesystem, etc.)

### Task 3.2: CLI Command Tests
**Goal:** Test CLI argument parsing and output formatting
**Files:** New test files in `cmd/arm/`
**Effort:** 3-4 hours

- `parser_test.go` - argument parsing tests
- `output_test.go` - formatting tests
- Integration tests with mocked service layer

### Task 3.3: Integration Tests
**Goal:** End-to-end workflow validation
**Files:** `internal/arm/integration_test.go` (enhance existing)
**Effort:** 4-5 hours

- Complete install → list → update → uninstall workflows
- Error scenario testing
- Multi-sink configuration testing

## Phase 4: Code Quality Improvements (Medium Priority)

### Task 4.1: Extract Constants and Configuration
**Goal:** Remove magic strings and numbers
**Files:** New `internal/arm/constants.go`
**Effort:** 1-2 hours

- Default patterns (`**/*`)
- File paths and extensions
- Error messages
- Configuration defaults

### Task 4.2: Improve Logging and Error Context
**Goal:** Better observability and debugging
**Files:** All service files
**Effort:** 2-3 hours

- Structured logging with consistent fields
- Error wrapping with context
- Performance logging for slow operations

### Task 4.3: Documentation and Examples
**Goal:** Self-documenting code
**Files:** All refactored files
**Effort:** 2-3 hours

- Add godoc comments to all public functions
- Include usage examples in complex functions
- Document error conditions and return values

## Phase 5: Architecture Improvements (Lower Priority)

### Task 5.1: Domain Model Cleanup
**Goal:** Clear separation of concerns
**Files:** `internal/types/`, `internal/arm/types.go`
**Effort:** 2-3 hours

- Consolidate type definitions
- Add validation methods to domain objects
- Remove data structure duplication

### Task 5.2: Configuration Management
**Goal:** Unified configuration handling
**Files:** `internal/config/`
**Effort:** 2-3 hours

- Merge manifest and config management concerns
- Add configuration validation
- Improve configuration file handling

## Implementation Order

### Week 1: Foundation (Tasks 1.1-1.3, 2.3)
- Split CLI commands
- Extract utilities
- Fix dependency injection
- **Outcome:** Cleaner, testable structure

### Week 2: Core Logic (Tasks 2.1-2.2, 3.1)
- Refactor service methods
- Extract business components
- Add service tests
- **Outcome:** Well-tested business logic

### Week 3: Testing & Quality (Tasks 3.2-3.3, 4.1-4.2)
- CLI tests
- Integration tests
- Constants and logging
- **Outcome:** Comprehensive test coverage

### Week 4: Polish (Tasks 4.3, 5.1-5.2)
- Documentation
- Domain model cleanup
- Configuration improvements
- **Outcome:** Production-ready codebase

## Success Metrics

- **Test Coverage:** >80% on service layer, >60% overall
- **Code Complexity:** No functions >50 lines, no files >300 lines
- **Maintainability:** Clear separation of CLI, service, and domain layers
- **Documentation:** All public APIs documented with examples

## Risk Mitigation

- **Regression Risk:** Comprehensive integration tests before refactoring
- **Breaking Changes:** Maintain existing CLI interface during refactoring
- **Time Overrun:** Prioritize Phase 1-2, defer Phase 5 if needed
