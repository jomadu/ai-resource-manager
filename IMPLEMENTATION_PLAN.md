# ARM Implementation Plan

## Status: Nearly Complete ✅

ARM is a fully functional dependency manager for AI packages with comprehensive test coverage (75 test files, 120 total Go files, 100% E2E test pass rate). All core functionality is implemented and tested.

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
  - Spec: `specs/pattern-filtering.md` (see BREAKING CHANGE v5.0 section)
  - Status: NOT STARTED - Required for v5.0 release

### Priority 2: Documentation Improvements

- [ ] **Add `arm list versions` to docs/commands.md**
  - Command exists and works (cmd/arm/main.go:965)
  - Add new section under "Core" commands
  - Show usage: `arm list versions REGISTRY/PACKAGE`
  - Document output format (semver descending, branches labeled)
  - Provide examples with expected output
  - Files: `docs/commands.md`
  - Status: Command implemented, documentation missing

### Priority 3: Test Coverage

- [ ] **Add E2E test for `arm list versions` command**
  - Test listing versions from different registry types
  - Verify semver sorting (descending)
  - Verify branch labeling
  - Files: `test/e2e/version_test.go` or new `test/e2e/list_versions_test.go`
  - Status: Command works, E2E test missing

### Priority 4: Version Resolution Edge Cases (Low Priority - May Not Be Bugs)

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

### Recently Completed (Verified 2026-01-28)
- ✅ Pattern matching in standalone compilation - matchesPatterns uses core.MatchPattern for full path matching
- ✅ UpdateAll error handling - continues on error with partial success pattern
- ✅ UpgradeAll error handling - continues on error with partial success pattern
- ✅ Help text for `arm list` command - shows all subcommands (registry, sink, dependency, versions)
- ✅ Default pattern behavior in registries - all three registries (Git, GitLab, Cloudsmith) apply `**/*.yml` and `**/*.yaml` defaults when no patterns specified
- ✅ CLI command for listing package versions - `arm list versions REGISTRY/PACKAGE` implemented and functional

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
2. **Documentation**: Add `arm list versions` to docs/commands.md
3. **Test Coverage**: Add E2E test for `arm list versions`
4. **Edge Cases**: Version resolution edge cases (may not be actual bugs)

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

1. **v5.0 Breaking Change**: Implement archive extraction to subdirectories (Priority 1)
2. Document `arm list versions` command (Priority 2)
3. Add E2E test for `arm list versions` (Priority 3)
4. Investigate version resolution edge cases if time permits (Priority 4)

## Maintenance Items

These are not missing features but ongoing maintenance:
- Keep dependencies updated (Dependabot handles this)
- Monitor security advisories (CodeQL handles this)
- Update documentation as needed
- Respond to user feedback and bug reports
