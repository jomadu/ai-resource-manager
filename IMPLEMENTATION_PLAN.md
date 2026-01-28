# ARM Implementation Plan

## Status: Nearly Complete ✅

ARM is a fully functional dependency manager for AI packages with comprehensive test coverage (75 test files, 120 total Go files, 100% E2E test pass rate). All core functionality is implemented and tested.

## Summary of Outstanding Work

**Total Items: 4** (1 breaking change, 1 bug, 1 documentation gap, 1 test gap)

**Priority Breakdown:**
- Priority 1 (Breaking Change): 1 item - Archive extraction to subdirectories (v5.0)
- Priority 2 (Bug Fix): 1 item - Missing `arm list dependency` command
- Priority 3 (Documentation): 1 item - Document `arm list versions` command
- Priority 4 (Test Coverage): 1 item - E2E test for `arm list dependency`

## Outstanding Items (Priority Order)

### Priority 1: BREAKING CHANGE - Archive Extraction (v5.0)

- [ ] **Extract archives to subdirectories** (pattern-filtering.md)
  - Current: Archives merge with loose files, causing collisions
  - Required: Extract archives to subdirectories named after archive (minus extension)
  - Example: `rules.tar.gz` containing `file.yml` → extracts to `rules/file.yml`
  - Impact: Breaking change - prevents collisions, enables skillset path resolution
  - Files to update:
    - `internal/arm/core/archive.go` - Rename ExtractAndMerge → Extract, add subdirectory logic
    - `internal/arm/registry/git.go` - Change ExtractAndMerge → Extract (line 168)
    - `internal/arm/registry/gitlab.go` - Change ExtractAndMerge → Extract (line 214)
    - `internal/arm/registry/cloudsmith.go` - Change ExtractAndMerge → Extract (line 255)
    - `test/e2e/archive_test.go` - Update expectations for subdirectory structure
    - `specs/e2e-testing.md` - Update acceptance criteria checkboxes
  - Spec: `specs/pattern-filtering.md` (see BREAKING CHANGE v5.0 section)
  - Status: NOT STARTED - Required for v5.0 release

### Priority 2: Bug Fix

- [ ] **Implement `arm list dependency` command** (query-operations.md)
  - Bug: Documented command not implemented in handleList() switch
  - Current: `arm list` (no args) shows all dependencies grouped under "Dependencies:" header
  - Required: `arm list dependency` should be a dedicated subcommand with different output format
  - Expected output: Dash-prefixed list of `registry/package@version` (see docs/commands.md:527-550)
  - Files to update:
    - `cmd/arm/main.go` - Add "dependency" case to handleList() switch (line 954-971)
    - Create `handleListDependency()` function (similar to handleListRegistry/handleListSink)
  - Implementation notes:
    - Read manifest and lock file
    - Combine rulesets and promptsets
    - Format as `- registry/package@version`
    - Sort alphabetically
  - Status: CONFIRMED BUG - Documented in docs/commands.md but not implemented in code

### Priority 3: Documentation Improvements

- [ ] **Add `arm list versions` to docs/commands.md**
  - Command exists and works (cmd/arm/main.go:966)
  - Add new section under "Dependency Management" commands (after `arm info dependency`)
  - Show usage: `arm list versions REGISTRY/PACKAGE`
  - Document output format (semver descending, branches labeled)
  - Provide examples with expected output
  - Files: `docs/commands.md`
  - Status: Command implemented, documentation missing

### Priority 4: Test Coverage

- [ ] **Add E2E test for `arm list dependency` command**
  - Test after implementing the command
  - Verify output format (dash-prefixed list)
  - Verify sorting (alphabetical)
  - Files: `test/e2e/manifest_test.go` or similar
  - Status: Blocked by Priority 2 (implement command first)




## Completed Features ✅

### Recently Completed (Verified 2026-01-28)
- ✅ Pattern matching in standalone compilation - matchesPatterns uses core.MatchPattern for full path matching (service.go:1833)
- ✅ UpdateAll error handling - continues on error with partial success pattern
- ✅ UpgradeAll error handling - continues on error with partial success pattern
- ✅ Default pattern behavior in registries - all three registries apply `**/*.yml` and `**/*.yaml` defaults (git.go:199, gitlab.go:374, cloudsmith.go:337)
- ✅ CLI command for listing package versions - `arm list versions REGISTRY/PACKAGE` implemented and functional (main.go:966)
- ✅ E2E test for `arm list versions` - exists in test/e2e/version_test.go
- ✅ Help text for `arm list` command - shows all subcommands (main.go:150-158)
- ✅ Prerelease version comparison - fully implemented with comprehensive tests (version.go:33-120, version_test.go:683-750)

### Core Functionality
- ✅ Package installation (install, update, upgrade, uninstall)
- ✅ Version resolution (semver, constraints, branches, latest)
- ✅ Registry management (Git, GitLab, Cloudsmith)
- ✅ Sink management and compilation (Cursor, Amazon Q, Copilot, Markdown)
- ✅ Priority-based rule conflict resolution
- ✅ Pattern filtering (include/exclude with glob patterns)
- ✅ Default pattern behavior in registries (defaults to `**/*.yml` and `**/*.yaml` when no patterns specified)
- ✅ Archive extraction (zip, tar.gz) - NOTE: v5.0 will change to subdirectory extraction
- ✅ Cache management (storage, cleanup, file locking)
- ✅ Authentication (token-based via .armrc)
- ✅ Integrity verification (SHA256 hashing)
- ✅ Query operations (list dependencies, check outdated, info, list versions)
- ✅ Standalone compilation (local files without registry)
- ✅ CLI command for listing package versions (`arm list versions REGISTRY/PACKAGE`)

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
- ✅ Complete command reference
- ✅ Registry type documentation
- ✅ Publishing guide
- ✅ Migration guide (v2 to v3)
- ✅ 18 specification documents

## Implementation Notes

### Why So Little Left?
The project is feature-complete for v3.x. All major features are implemented and tested. The remaining items are:
1. **v5.0 Breaking Change**: Archive extraction to subdirectories (prevents collisions)
2. **Bug Fix**: 1 confirmed bug - missing `arm list dependency` command in handleList switch
3. **Documentation**: 1 gap where implemented feature isn't documented
4. **Test Coverage**: 1 E2E test missing for new functionality

### Code Quality
- All linting passes (13 linters enabled)
- All tests pass (go test ./... succeeds)
- No TODO/FIXME/HACK comments in production code (only in docs/examples and git history)
- Only 2 skipped tests (both documented with reasons in test files)

### Architecture
- Clean separation: cmd/ (CLI) → internal/arm/service/ (business logic) → internal/arm/* (components)
- Constructor injection for test isolation
- Environment variable support (ARM_HOME, ARM_CONFIG_PATH, ARM_MANIFEST_PATH)
- Registry factory pattern for extensibility

## Next Steps

**Recommended Order:**

1. **Fix bug** (Priority 2) - Confirmed issue affecting current functionality:
   - Implement `arm list dependency` command (help text exists but command missing from handleList switch)

2. **Add documentation** (Priority 3) - Quick win to improve user experience:
   - Document `arm list versions` command in docs/commands.md

3. **Add test coverage** (Priority 4) - Ensure new functionality is tested:
   - E2E test for `arm list dependency` (after implementing command)

4. **Plan v5.0 breaking change** (Priority 1) - Requires careful planning and migration guide:
   - Extract archives to subdirectories (prevents collisions)
   - Update all registry implementations
   - Update E2E tests
   - Create migration guide for users

## Maintenance Items

These are not missing features but ongoing maintenance:
- Keep dependencies updated (Dependabot handles this)
- Monitor security advisories (CodeQL handles this)
- Update documentation as needed
- Respond to user feedback and bug reports
