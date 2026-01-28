# ARM Implementation Plan

## Status: Nearly Complete ✅

ARM is a fully functional dependency manager for AI packages with comprehensive test coverage (75 test files, 120 total Go files, 100% E2E test pass rate). All core functionality is implemented and tested.

## Outstanding Items (Priority Order)

### Priority 1: Missing CLI Command

- [ ] **Implement `arm list versions` command** (query-operations.md)
  - Spec: List available versions for a package from its registry
  - Backend: `ListPackageVersions()` already implemented in all registries (git.go, gitlab.go, cloudsmith.go)
  - Missing: CLI handler in `cmd/arm/main.go`
  - Implementation:
    - Add case "versions" in `handleList()` switch statement (line 951-965)
    - Create `handleListVersions()` function
    - Parse package key (registry/package) from args
    - Call service layer to get registry and list versions
    - Format output: package name header, then indented list of versions
    - Sort: semver descending, branches in config order
    - Label branches with "(branch)" suffix
  - Example: `arm list versions test-registry/clean-code-ruleset`
  - Files to modify: `cmd/arm/main.go`

### Priority 2: Pattern Filtering Bugs

- [ ] **Fix default pattern behavior in registries** (pattern-filtering.md)
  - Bug: When no patterns specified, registries return ALL files instead of defaulting to `["**/*.yml", "**/*.yaml"]`
  - Root cause: `matchesPatterns()` returns true when both include and exclude are empty
  - Files affected:
    - `internal/arm/registry/git.go:199` (GitRegistry.matchesPatterns)
    - `internal/arm/registry/gitlab.go:374` (matchesPatterns helper)
    - `internal/arm/registry/cloudsmith.go:337` (CloudsmithRegistry.matchesPatterns)
  - Fix: Add default pattern logic at start of function:
    ```go
    if len(include) == 0 && len(exclude) == 0 {
        include = []string{"**/*.yml", "**/*.yaml"}
    }
    ```
  - Impact: Users currently must explicitly specify `--include "**/*.yml"` to avoid getting non-YAML files
  - Test: Verify install without patterns only gets YAML files

- [ ] **Fix pattern matching in standalone compilation** (standalone-compilation.md, pattern-filtering.md)
  - Bug: `internal/arm/service/service.go:1763,1778` uses `filepath.Match(pattern, filepath.Base(filePath))` instead of `core.MatchPattern(pattern, filePath)`
  - Root cause: Pattern matching on basename only, not full path
  - Impact: Patterns like `security/**/*.yml` don't work in `arm compile` command
  - Fix: Replace both occurrences:
    ```go
    // Before
    if matched, _ := filepath.Match(pattern, filepath.Base(filePath)); matched {
    
    // After
    if core.MatchPattern(pattern, filePath) {
    ```
  - Files: `internal/arm/service/service.go` (matchesPatterns function)
  - Test: Verify `arm compile` with `--include "security/**/*.yml"` works correctly

### Priority 3: Update/Upgrade Error Handling

- [ ] **Fix UpdateAll to continue on error** (package-installation.md)
  - Bug: `UpdateAll()` returns on first error instead of continuing for partial success
  - Files: `internal/arm/service/service.go:731-780` (UpdateAll function)
  - Expected behavior: Continue processing remaining packages, collect errors, return combined error
  - Reference: `UpdatePackages()` (line 600-729) correctly implements partial success pattern
  - Implementation:
    - Collect errors in slice instead of returning immediately
    - Continue loop on error
    - Return combined error at end if any errors occurred
  - Test: Verify update continues when one package fails

- [ ] **Fix UpgradeAll to continue on error** (package-installation.md)
  - Bug: `UpgradeAll()` returns on first error instead of continuing for partial success
  - Files: `internal/arm/service/service.go:887-950` (UpgradeAll function)
  - Expected behavior: Continue processing remaining packages, collect errors, return combined error
  - Reference: `UpgradePackages()` correctly implements partial success pattern
  - Implementation: Same pattern as UpdateAll fix
  - Test: Verify upgrade continues when one package fails

### Priority 4: Documentation Improvements

- [ ] **Update help text for `arm list` command**
  - Current: Only shows `arm list registry` (line 149-156 in cmd/arm/main.go)
  - Should show: All subcommands (registry, sink, dependency, versions)
  - Files: `cmd/arm/main.go` (showHelp function, case "list")
  - Add lines:
    ```
    fmt.Println("  arm list sink")
    fmt.Println("  arm list dependency")
    fmt.Println("  arm list versions REGISTRY/PACKAGE")
    ```

- [ ] **Add `arm list versions` to docs/commands.md**
  - Add new section under "Core" commands
  - Show usage: `arm list versions REGISTRY/PACKAGE`
  - Document output format (semver descending, branches labeled)
  - Provide examples with expected output
  - Files: `docs/commands.md`

### Priority 5: Test Coverage

- [ ] **Add E2E test for `arm list versions` command**
  - Create new file: `test/e2e/query_test.go`
  - Test scenarios:
    - List versions from Git registry (semver tags + branches)
    - List versions from GitLab registry (test pagination)
    - List versions from Cloudsmith registry
    - Verify output format and sorting
    - Test error cases (package not found, registry not configured)
  - Use existing test helpers from `test/e2e/helpers/`

- [ ] **Add tests for pattern filtering bugs**
  - Test default pattern behavior in registries
    - Install without patterns, verify only YAML files installed
    - Files: `test/e2e/install_test.go`
  - Test ** patterns in standalone compilation
    - Compile with `--include "security/**/*.yml"`, verify correct files
    - Files: `test/e2e/compile_test.go`

### Priority 6: Version Resolution Edge Cases (Low Priority - May Not Be Bugs)

- [ ] **Verify prerelease version comparison** (version-resolution.md)
  - Spec mentions: 1.0.0-alpha.1 < 1.0.0-alpha.2 < 1.0.0-beta.1 < 1.0.0-rc.1 < 1.0.0
  - Current implementation: `internal/arm/core/version.go` (comparePrerelease function)
  - Status: Implementation looks correct, but spec notes "may not handle all edge cases"
  - Action: Review existing tests in `internal/arm/core/version_test.go`
  - If gaps found, add tests for:
    - Alpha < beta < rc < release precedence
    - Numeric vs alphanumeric prerelease identifiers
    - Multiple prerelease components (1.0.0-alpha.1.2)

- [ ] **Verify "latest" resolution with no semantic versions** (version-resolution.md)
  - Spec: When no semantic versions exist, "latest" should use first configured branch
  - Current: May use lexicographic sort instead
  - Files: `internal/arm/core/helpers.go` (ResolveVersion or GetBestMatching)
  - Status: One test is skipped for this scenario (test/e2e/version_test.go:321)
  - Action: Review implementation and determine if this is actually a bug
  - If bug confirmed, fix to use first configured branch from registry config


## Completed Features ✅

### Core Functionality
- ✅ Package installation (install, update, upgrade, uninstall)
- ✅ Version resolution (semver, constraints, branches, latest)
- ✅ Registry management (Git, GitLab, Cloudsmith)
- ✅ Sink management and compilation (Cursor, Amazon Q, Copilot, Markdown)
- ✅ Priority-based rule conflict resolution
- ✅ Pattern filtering (include/exclude with glob patterns) - has bugs but functional
- ✅ Archive extraction (zip, tar.gz)
- ✅ Cache management (storage, cleanup, file locking)
- ✅ Authentication (token-based via .armrc)
- ✅ Integrity verification (SHA256 hashing)
- ✅ Query operations (list dependencies, check outdated, info)
- ✅ Standalone compilation (local files without registry) - has pattern bug but functional
- ✅ Backend for listing package versions (ListPackageVersions in all registries)

### Infrastructure
- ✅ Cross-platform builds (Linux, macOS, Windows - amd64/arm64)
- ✅ Installation scripts (install.sh, uninstall.sh)
- ✅ CI/CD workflows (build, test, lint, security, release)
- ✅ Semantic release automation
- ✅ Code quality (13 linters, pre-commit hooks, conventional commits)
- ✅ Security scanning (CodeQL, dependency review)
- ✅ Dependency management (Dependabot)

### Testing
- ✅ 75 test files with comprehensive coverage
- ✅ 14 E2E test suites covering all workflows
- ✅ Test isolation via environment variables
- ✅ 100% E2E test pass rate
- ✅ Only 2 skipped tests (both documented with reasons)

### Documentation
- ✅ README.md (comprehensive overview)
- ✅ 12 docs files (2686 lines total)
- ✅ Complete command reference (except `arm list versions`)
- ✅ Registry type documentation
- ✅ Publishing guide
- ✅ Migration guide (v2 to v3)
- ✅ 18 specification documents

## Implementation Notes

### Why So Little Left?
The project is essentially feature-complete. The main missing piece is exposing an already-implemented backend feature (`ListPackageVersions`) through the CLI. The bugs are edge cases that don't prevent normal usage.

### Code Quality
- All linting passes (13 linters enabled)
- All tests pass (go test ./... succeeds)
- No TODO/FIXME/HACK comments in production code
- Only 2 skipped tests (both documented with reasons)

### Architecture
- Clean separation: cmd/ (CLI) → internal/arm/service/ (business logic) → internal/arm/* (components)
- Constructor injection for test isolation
- Environment variable support (ARM_HOME, ARM_CONFIG_PATH, ARM_MANIFEST_PATH)
- Registry factory pattern for extensibility

## Next Steps

1. Implement `arm list versions` CLI command (Priority 1)
2. Fix pattern filtering bugs (Priority 2)
3. Fix UpdateAll/UpgradeAll error handling (Priority 3)
4. Update documentation (Priority 4)
5. Add test coverage for new features and bug fixes (Priority 5)
6. Investigate version resolution edge cases if time permits (Priority 6)

## Maintenance Items

These are not missing features but ongoing maintenance:
- Keep dependencies updated (Dependabot handles this)
- Monitor security advisories (CodeQL handles this)
- Update documentation as needed
- Respond to user feedback and bug reports
