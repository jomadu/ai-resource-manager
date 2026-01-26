# End-to-End Testing

## Purpose

Validate that ARM respects all specifications through comprehensive end-to-end tests using local Git repositories. Tests must be runnable both locally and in CI environments.

## Test Infrastructure

### Local Git Repository Setup

Tests use local Git repositories to simulate real-world scenarios without external dependencies:

```bash
# Create test registry repository
mkdir -p /tmp/arm-test-registry
cd /tmp/arm-test-registry
git init
git config user.email "test@example.com"
git config user.name "Test User"

# Add test resources
# Commit and tag versions
git add .
git commit -m "Initial commit"
git tag v1.0.0
```

### Test Isolation

- Each test creates isolated temporary directories
- Tests clean up after themselves
- No shared state between tests
- Tests can run in parallel
- **Direct path injection** - Tests pass t.TempDir() to component constructors to avoid polluting user's ~/.arm/ and ~/.armrc
- Components accept directory paths as constructor parameters
- No direct os.UserHomeDir() calls in testable components
- **Lock file colocation** - Lock file always in same directory as manifest file
- Lock path derived from manifest path (arm.json → arm-lock.json)

### CI Compatibility

- No external network dependencies
- Deterministic test data
- Fast execution (< 5 minutes total)
- Clear pass/fail criteria

## Test Scenarios

### Registry Management

**Git Registry:**
- Add Git registry with local file:// URL
- Add Git registry with SSH-style path
- Remove registry
- Set registry configuration
- List registries
- Info for specific registry

**GitLab Registry:**
- Add GitLab registry with project ID
- Add GitLab registry with group ID
- Authentication via .armrc
- Remove registry
- Set registry configuration

**Cloudsmith Registry:**
- Add Cloudsmith registry
- Authentication via .armrc
- Remove registry
- Set registry configuration

### Sink Management

**Sink Operations:**
- Add sink for each tool type (cursor, amazonq, copilot, markdown)
- Remove sink
- Set sink configuration (tool, directory)
- List sinks
- Info for specific sink

**Layout Modes:**
- Hierarchical layout (default)
- Flat layout
- Verify file paths match spec

### Resource Installation

**Ruleset Installation:**
- Install ruleset from Git registry (semver tag)
- Install ruleset from Git registry (branch)
- Install ruleset with version constraint (@1, @1.1, @1.0.0)
- Install ruleset with @latest
- Install ruleset to multiple sinks
- Install ruleset with custom priority
- Install ruleset with --include patterns
- Install ruleset with --exclude patterns
- Install ruleset with both --include and --exclude
- Reinstall to different sinks (verify old sinks cleaned)

**Promptset Installation:**
- Install promptset from Git registry (semver tag)
- Install promptset from Git registry (branch)
- Install promptset with version constraint
- Install promptset to multiple sinks
- Install promptset with --include/--exclude patterns
- Reinstall to different sinks

**Archive Support:**
- Install from .tar.gz archive
- Install from .zip archive
- Install with mixed archives and loose files
- Verify archive precedence over loose files
- Apply --include/--exclude to extracted content

### Version Resolution

**Git Registry:**
- Resolve @latest to highest semver tag
- Resolve @latest to branch HEAD when no tags
- Resolve @1 to highest 1.x.x version
- Resolve @1.1 to highest 1.1.x version
- Resolve @1.0.0 to exact version
- Resolve branch name to commit hash
- Verify arm-lock.json contains resolved version

**Version Constraints:**
- Install with major constraint (^1.0.0 → >= 1.0.0, < 2.0.0)
- Install with minor constraint (^1.1.0 → >= 1.1.0, < 1.2.0)
- Install with exact constraint (1.0.0 → exactly 1.0.0)
- Update respects constraints
- Upgrade ignores constraints

### Dependency Management

**Install Operations:**
- `arm install` - install all configured dependencies
- `arm install ruleset` - install specific ruleset
- `arm install promptset` - install specific promptset
- Verify arm.json updated
- Verify arm-lock.json created/updated
- Verify arm-index.json created/updated

**Uninstall Operations:**
- `arm uninstall` - remove all dependencies
- Verify sink directories cleaned
- Verify arm-index.json updated
- Verify arm.json preserved

**Update Operations:**
- `arm update` - update within constraints
- Verify only compatible versions installed
- Verify arm-lock.json updated

**Upgrade Operations:**
- `arm upgrade` - upgrade to latest ignoring constraints
- Verify arm.json constraint updated to ^X.0.0
- Verify arm-lock.json updated

**Outdated Check:**
- `arm outdated` - detect outdated packages
- Verify table output format
- Verify JSON output format
- Verify list output format
- Show constraint, current, wanted, latest

### Compilation

**Tool-Specific Compilation:**
- Compile ruleset to cursor format (.mdc with frontmatter)
- Compile ruleset to amazonq format (.md)
- Compile ruleset to copilot format (.instructions.md)
- Compile ruleset to markdown format (.md)
- Compile promptset to cursor format (.md)
- Compile promptset to amazonq format (.md)
- Compile promptset to copilot format (.instructions.md)
- Compile promptset to markdown format (.md)

**Compilation Options:**
- Compile single file
- Compile multiple files
- Compile directory (non-recursive)
- Compile directory (--recursive)
- Compile with --namespace
- Compile with --include/--exclude patterns
- Compile with --force overwrite
- Compile with --validate-only (no output)
- Compile with --fail-fast

**Validation:**
- Validate valid ruleset
- Validate valid promptset
- Detect invalid YAML
- Detect missing required fields
- Detect invalid field types

### Priority Resolution

**Rule Priority:**
- Install multiple rulesets with different priorities
- Verify arm_index.* contains priority metadata
- Verify higher priority rules listed first
- Install with default priority (100)
- Install with custom priority (--priority 200)
- Update priority via `arm set ruleset`

### File Patterns

**Include Patterns:**
- Default includes (*.yml, *.yaml)
- Single include pattern
- Multiple include patterns (OR logic)
- Glob patterns (**/*.yml, security/**/*.md)

**Exclude Patterns:**
- Single exclude pattern
- Multiple exclude patterns
- Exclude overrides include
- Combined include and exclude

### Storage System

**Cache Operations:**
- Verify package cached after first install
- Verify cache reused on second install
- Verify cache key generation (registry + package + patterns)
- Clean cache with default (7 days)
- Clean cache with --max-age
- Clean cache with --nuke

**Cache Structure:**
- Verify registry metadata.json
- Verify package metadata.json
- Verify version metadata.json
- Verify Git repo/ directory for Git registries
- Verify packages/<package-key>/<version>/files/

### Manifest Files

**arm.json:**
- Created on first registry/sink/install
- Updated on configuration changes
- Preserves existing configuration
- Valid JSON format
- Contains registries, sinks, dependencies

**arm-lock.json:**
- Created on first install
- Updated on install/update/upgrade
- Contains resolved versions
- Git branches resolve to commit hash
- Semver tags remain as semver

**arm-index.json:**
- Created on first install
- Updated on install/uninstall
- Tracks installed files and metadata
- Used by `arm clean sinks`

**arm_index.*:**
- Generated for each sink
- Contains priority-ordered rules
- Tool-specific format
- Includes metadata for conflict resolution

### Authentication

**.armrc:**
- GitLab registry authentication
- Cloudsmith registry authentication
- Environment variable substitution
- Local .armrc takes precedence over ~/.armrc
- Section matching by full URL

### Error Handling

**Invalid Operations:**
- Install non-existent package
- Install non-existent version
- Install to non-existent sink
- Install from non-existent registry
- Add duplicate registry (without --force)
- Add duplicate sink (without --force)
- Remove non-existent registry
- Remove non-existent sink
- Invalid version constraint
- Invalid glob pattern

**Validation Errors:**
- Invalid YAML syntax
- Missing required fields
- Invalid field types
- Invalid resource kind

### Multi-Sink Scenarios

**Cross-Tool Installation:**
- Install same ruleset to cursor and amazonq
- Verify different compilation formats
- Verify both sinks updated
- Verify arm-index.json tracks both

**Sink Switching:**
- Install to sink A
- Reinstall to sink B
- Verify sink A cleaned
- Verify sink B populated
- Verify arm-index.json updated

### Update Workflows

**Version Updates:**
- Install v1.0.0
- Publish v1.1.0 to registry
- Run `arm update`
- Verify v1.1.0 installed
- Verify arm-lock.json updated

**Breaking Changes:**
- Install v1.0.0 with constraint ^1.0.0
- Publish v2.0.0 to registry
- Run `arm update`
- Verify v1.0.0 remains (constraint respected)
- Run `arm upgrade`
- Verify v2.0.0 installed
- Verify constraint updated to ^2.0.0

## Test Data Structure

### Minimal Test Registry

```
test-registry/
├── clean-code-ruleset.yml       # Simple ruleset
├── security-ruleset.yml         # Ruleset with priority
├── code-review-promptset.yml    # Simple promptset
├── archived-rules.tar.gz        # Archive test
└── archived-rules.zip           # Archive test
```

### Test Resource Examples

**Minimal Ruleset:**
```yaml
apiVersion: v1
kind: Ruleset
metadata:
  id: "testRuleset"
  name: "Test Ruleset"
spec:
  rules:
    ruleOne:
      body: "This is rule one."
    ruleTwo:
      priority: 150
      body: "This is rule two with priority."
```

**Minimal Promptset:**
```yaml
apiVersion: v1
kind: Promptset
metadata:
  id: "testPromptset"
  name: "Test Promptset"
spec:
  prompts:
    promptOne:
      body: "This is prompt one."
```

## Test Execution

### Local Execution

```bash
# Run all e2e tests
go test ./test/e2e/... -v

# Run specific test
go test ./test/e2e/... -v -run TestGitRegistry

# Run with coverage
go test ./test/e2e/... -v -coverprofile=coverage.out
```

### CI Execution

```yaml
# GitHub Actions example
- name: Run E2E Tests
  run: |
    go test ./test/e2e/... -v -race -timeout 5m
```

## Success Criteria

Each test must:
- ✅ Create isolated test environment
- ✅ Execute ARM commands
- ✅ Verify expected files created
- ✅ Verify file contents match spec
- ✅ Verify manifest files updated correctly
- ✅ Clean up test environment
- ✅ Complete in < 5 seconds (individual test)
- ✅ Provide clear failure messages

## Test Organization

```
test/e2e/
├── registry_test.go           # Registry management tests
├── sink_test.go              # Sink management tests
├── install_test.go           # Installation tests
├── version_test.go           # Version resolution tests
├── compile_test.go           # Compilation tests
├── priority_test.go          # Priority resolution tests
├── patterns_test.go          # File pattern tests
├── storage_test.go           # Cache/storage tests
├── manifest_test.go          # Manifest file tests
├── auth_test.go              # Authentication tests
├── errors_test.go            # Error handling tests
├── multisink_test.go         # Multi-sink scenarios
├── update_test.go            # Update workflow tests
└── helpers/
    ├── git.go                # Git repo helpers
    ├── fixtures.go           # Test data generators
    └── assertions.go         # Custom assertions
```

## Implementation Notes

- Use `testing.T` for test framework
- Use `t.TempDir()` for isolated directories
- Use `exec.Command()` to run ARM binary
- Parse JSON output for verification
- Use table-driven tests for variations
- Mock time for cache age tests
- Use golden files for complex output verification
- **Pass t.TempDir() to constructors** - Inject test directories directly as string parameters
- Components should accept directory paths via constructor for testability
- Avoid direct os.UserHomeDir() calls in components
