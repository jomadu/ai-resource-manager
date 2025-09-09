# ARM Code Cleanup Plan

## Overview
Refactor ARM codebase for better readability, organization, and test coverage while following clean code principles.

## Phase 1: CLI Layer Cleanup ✅ COMPLETED

### Task 1.1: Split CLI Commands ✅ DONE
**Status:** COMPLETED - All commands extracted into separate files
**Files:** CLI commands properly separated

✅ Commands extracted into separate files:
  - `cmd/arm/install.go` - install command logic
  - `cmd/arm/uninstall.go` - uninstall command logic
  - `cmd/arm/update.go` - update command logic
  - `cmd/arm/config.go` - config command logic
  - `cmd/arm/cache.go` - cache command logic
  - `cmd/arm/info.go` - info command
  - `cmd/arm/list.go` - list command
  - `cmd/arm/outdated.go` - outdated command
  - `cmd/arm/version.go` - version command
✅ `cmd/arm/parser.go` - argument parsing utilities created
✅ `cmd/arm/output.go` - formatting utilities created
✅ `main.go` minimal with just root command setup (30 lines)

### Task 1.2: Extract CLI Utilities ✅ DONE
**Status:** COMPLETED - Utilities properly extracted
**Files:** `cmd/arm/parser.go`, `cmd/arm/output.go`

✅ Parsing utilities moved to `parser.go`
✅ Output formatting moved to `output.go`
✅ Proper parser functions implemented
✅ Input validation at CLI layer

### Task 1.3: Improve Error Handling ✅ DONE
**Status:** COMPLETED - Error handling standardized
**Files:** All CLI command files

✅ Standardized error message format
✅ Context added to error messages
✅ Error logging handled by service layer

## Phase 2: Service Layer Refactoring (IN PROGRESS)

### Task 2.1: Break Down Large Methods ⚠️ NEEDS REFACTORING
**Status:** IDENTIFIED - Large methods still exist and need breaking down
**Files:** `internal/arm/service.go` (600+ lines)
**Current Issues:**

⚠️ `InstallRuleset` method is 120+ lines doing multiple responsibilities:
  - Input validation
  - Registry client creation
  - Version resolution
  - Content download
  - Manifest/lockfile updates
  - Sink installation

⚠️ `GetOutdatedRulesets` method is 60+ lines with nested loops
⚠️ `SyncSink` method is 50+ lines with complex logic

**Next Steps:**
- Extract `validateInstallRequest`
- Extract `resolveRulesetVersion`
- Extract `downloadRulesetContent`
- Extract `updateManifestAndLock`
- Extract `installToSinks`

### Task 2.2: Extract Business Logic ❌ NOT STARTED
**Status:** PENDING - Service layer still monolithic
**Files:** Need to create new service files
**Current State:** All business logic in single `service.go` file

**Required Services:**
- `installer_service.go` - Installation/uninstallation logic
- `version_service.go` - Version resolution and comparison
- `content_service.go` - Content download and validation
- `tracking_service.go` - Manifest/lockfile management
- `sync_service.go` - Sink synchronization

### Task 2.3: Improve Dependency Injection ⚠️ PARTIALLY DONE
**Status:** BASIC STRUCTURE - Still has hard-coded dependencies
**Files:** `internal/arm/service.go`

✅ Constructor exists (`NewArmService()`)
⚠️ Still has hard-coded dependencies:
  - `config.NewFileManager()`
  - `manifest.NewFileManager()`
  - `lockfile.NewFileManager()`

**Needs:**
- Interface-based dependency injection
- Testable constructor with explicit dependencies

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

## Implementation Status & Next Steps

### ✅ COMPLETED: Foundation (Phase 1)
- ✅ Split CLI commands into focused files
- ✅ Extract parsing and output utilities
- ✅ Standardize error handling
- **Outcome:** Clean, organized CLI layer

### ⚠️ CURRENT PRIORITY: Service Layer Refactoring (Phase 2)
**Status:** Ready to start - CLI foundation is solid

**Immediate Next Steps:**
1. **Break down `InstallRuleset` method** (120 lines → 5 focused methods)
2. **Extract service components** (installer, version, content, tracking, sync)
3. **Implement dependency injection** (interfaces + testable constructors)

**Estimated Effort:** 8-10 hours

### ❌ PENDING: Testing & Quality (Phase 3)
- Service layer unit tests
- CLI command tests
- Integration tests
- **Outcome:** Comprehensive test coverage

### ❌ PENDING: Polish (Phase 4-5)
- Constants extraction
- Logging improvements
- Documentation
- Domain model cleanup

## Success Metrics

### ✅ Achieved:
- **CLI Organization:** Clean command separation, focused files
- **CLI Complexity:** All command files <100 lines, main.go <50 lines

### ⚠️ Current Issues:
- **Service Complexity:** `service.go` is 600+ lines with methods >100 lines
- **Test Coverage:** Minimal unit tests, mostly integration tests
- **Maintainability:** Monolithic service layer

### 🎯 Targets:
- **Test Coverage:** >80% on service layer, >60% overall
- **Code Complexity:** No functions >50 lines, no files >300 lines
- **Maintainability:** Clear separation of CLI, service, and domain layers
- **Documentation:** All public APIs documented with examples

## Risk Mitigation

- **Regression Risk:** Comprehensive integration tests before refactoring
- **Breaking Changes:** Maintain existing CLI interface during refactoring
- **Time Overrun:** Prioritize Phase 1-2, defer Phase 5 if needed
