# ARM Implementation Plan

## Status: Verified Analysis Complete (2026-01-24)

This document tracks implementation status and prioritizes remaining work for the AI Resource Manager (ARM) project.

## Executive Summary

**Overall Status**: ~95% Complete - Core functionality fully implemented and tested

**Key Findings**:
- All core commands implemented and tested (registry, sink, dependency management)
- All registry types functional (Git, GitLab, Cloudsmith) with archive support
- All compilers fully implemented (Cursor, Amazon Q, Copilot, Markdown) - verified
- All parsers implemented (Ruleset, Promptset) - verified
- File type detection implemented - verified
- Archive extraction implemented in core with comprehensive tests
- 304 tests across 58 test files - excellent coverage
- Primary gap: `arm compile` command service layer implementation (CLI exists, service stub only)

---

## Priority 1: Critical Missing Functionality

### 1.1 Implement `arm compile` Command Service Layer ⚠️ VERIFIED MISSING
**Status**: ❌ Not Implemented  
**Location**: `internal/arm/service/service.go:1595-1598`  
**CLI Location**: `cmd/arm/main.go:2141-2250` (✅ fully implemented)  
**Spec Reference**: `specs/commands.md` - "arm compile" section

**Current State**:
```go
func (s *ArmService) CompileFiles(ctx context.Context, req *CompileRequest) error {
    // TODO: implement
    return nil
}
```

**Verified Available Components**:
- ✅ CLI handler fully implemented with all flags
- ✅ Compiler functions: `compiler.CompileRuleset()`, `compiler.CompilePromptset()`
- ✅ All tool compilers: Cursor, AmazonQ, Copilot, Markdown (verified in factory.go)
- ✅ Parser functions: `parser.ParseRuleset()`, `parser.ParsePromptset()`
- ✅ File type detection: `filetype.IsRulesetFile()`, `filetype.IsPromptsetFile()`
- ✅ CompileRequest struct with all required fields

**Requirements** (from spec):
- Accept input paths (files and/or directories)
- Support `--tool` flag (markdown, cursor, amazonq, copilot)
- Support `--namespace` flag (defaults to resource metadata ID)
- Support `--force` flag (overwrite existing files)
- Support `--recursive` flag (process directories recursively)
- Support `--validate-only` flag (validate without output)
- Support `--include` patterns (default: `**/*.yml`, `**/*.yaml`)
- Support `--exclude` patterns
- Support `--fail-fast` flag (stop on first error)
- Handle both files and directories as input
- Shell glob expansion handled by shell before ARM processes
- Optional OUTPUT_PATH (required unless `--validate-only`)

**Implementation Approach**:
1. Validate and normalize input paths
2. Discover files (handle directories with include/exclude patterns)
3. For each file:
   - Read file content into `core.File`
   - Detect type using `filetype.IsRulesetFile()` / `filetype.IsPromptsetFile()`
   - Parse YAML using `parser.ParseRuleset()` / `parser.ParsePromptset()`
   - Validate schema (parser does this automatically)
   - If validate-only: continue to next file
   - Determine namespace (use flag or resource metadata ID)
   - Compile using `compiler.CompileRuleset()` / `compiler.CompilePromptset()`
   - Write compiled files to output directory
4. Handle errors per fail-fast setting
5. Return appropriate errors with context

**Acceptance Criteria**:
- [ ] Parse and validate input paths
- [ ] Discover files using include/exclude patterns (use existing pattern matching utilities)
- [ ] Support recursive directory traversal
- [ ] Detect file types (ruleset vs promptset)
- [ ] Parse YAML resources
- [ ] Validate resources (when `--validate-only`)
- [ ] Compile resources to target tool format
- [ ] Write output files to OUTPUT_PATH
- [ ] Respect `--force` flag for overwrites
- [ ] Stop on first error when `--fail-fast`
- [ ] Return appropriate errors with context
- [ ] Update test expectations in `cmd/arm/compile_test.go`

---

## Priority 2: Documentation Gaps

### 2.1 Create Migration Guide (v2 to v3) ⚠️ VERIFIED MISSING
**Status**: ❌ Missing  
**Location**: Should be `specs/migration-v2-to-v3.md`  
**Referenced In**: `README.md:59`

**Current State**: Referenced in README but file does not exist (verified with glob search)

**Requirements**:
- Document breaking changes from v2 to v3
- Provide migration steps
- Include examples of old vs new patterns
- Explain "nuke and pave" recommendation from README
- Update README link to point to actual file

**Context from README**:
> "Upgrading from v2? See the migration guide for breaking changes and upgrade steps. TL;DR: Sorry, nuke and pave. We made some poor design choices in v1 and v2."

**Acceptance Criteria**:
- [ ] Create `specs/migration-v2-to-v3.md`
- [ ] Document all breaking changes between v2 and v3
- [ ] Provide step-by-step migration guide
- [ ] Include before/after examples
- [ ] Explain why "nuke and pave" is recommended
- [ ] Verify README link points to actual file

---

## Priority 3: Test Coverage Improvements

### 3.1 Implement Archive Support Test ⚠️ VERIFIED SKIPPED
**Status**: ⚠️ Skipped  
**Location**: `internal/arm/registry/git_test.go:958-961`  
**Current State**: Test exists but skipped with `t.Skip("TODO: implement")`

**Important Context**: 
- ✅ Archive extraction IS fully implemented in `internal/arm/core/archive.go`
- ✅ Comprehensive tests exist in `internal/arm/core/archive_test.go` (10 tests)
- ✅ Git registry uses it via `core.NewExtractor()`
- ❌ Missing: End-to-end test in Git registry context

**Test Purpose**: Verify end-to-end archive handling in Git registry, not the extraction itself

**Acceptance Criteria**:
- [ ] Remove `t.Skip()` from test
- [ ] Test `.zip` file extraction in Git registry context
- [ ] Test `.tar.gz` file extraction in Git registry context
- [ ] Verify extracted files are properly merged with loose files
- [ ] Test archive precedence over loose files (archives win conflicts)
- [ ] Ensure test passes

### 3.2 Review GetTags Test Assertion ⚠️ VERIFIED TODO
**Status**: ⚠️ Needs Review  
**Location**: `internal/arm/storage/repo_test.go:40`  
**Current State**: Comment says "TODO: implement and change assertion"

**Test Code**:
```go
// Test GetTags clones and returns tags
ctx := context.Background()
tags, err := repo.GetTags(ctx, sourceDir)

// TODO: implement and change assertion
assert.NoError(t, err)
assert.Equal(t, []string{"v1.0.0", "v1.1.0"}, tags)
```

**Note**: The test appears to be passing, but the comment suggests the assertion may not be correct or complete.

**Acceptance Criteria**:
- [ ] Review GetTags implementation behavior
- [ ] Verify test expectations match actual behavior
- [ ] Update assertion if needed
- [ ] Remove TODO comment
- [ ] Verify test accurately reflects GetTags behavior

---

## Priority 4: Code Quality & Maintenance

### 4.1 Review Hash Pattern Comments ✅ VERIFIED NOT ISSUES
**Status**: ✅ Informational (Not TODOs)  
**Locations**:
- `internal/arm/sink/manager.go:547-548` - "arm_xxxx_xxxx_" pattern comment
- `internal/arm/sink/manager_test.go:120` - Dual hash pattern comment

**Analysis**: These are helpful explanatory comments about the hash pattern format, not TODO items. They document the filename pattern used in flat layout mode.

**Recommendation**: Keep as-is. These comments aid code comprehension.

**No Action Required**: These are documentation comments, not implementation gaps.

---

## Completed Features ✅

### Core Commands (All Implemented & Tested)
- ✅ `arm version` - Display version information
- ✅ `arm help` - Display help
- ✅ `arm list` - List all entities
- ✅ `arm info` - Show detailed information

### Registry Management (All Implemented & Tested)
- ✅ `arm add registry git` - Add Git registry
- ✅ `arm add registry gitlab` - Add GitLab registry
- ✅ `arm add registry cloudsmith` - Add Cloudsmith registry
- ✅ `arm remove registry` - Remove registry
- ✅ `arm set registry` - Configure registry
- ✅ `arm list registry` - List registries
- ✅ `arm info registry` - Show registry details

### Sink Management (All Implemented & Tested)
- ✅ `arm add sink` - Add sink
- ✅ `arm remove sink` - Remove sink
- ✅ `arm set sink` - Configure sink
- ✅ `arm list sink` - List sinks
- ✅ `arm info sink` - Show sink details

### Dependency Management (All Implemented & Tested)
- ✅ `arm install` - Install all dependencies
- ✅ `arm install ruleset` - Install specific ruleset
- ✅ `arm install promptset` - Install specific promptset
- ✅ `arm uninstall` - Uninstall packages
- ✅ `arm update` - Update within constraints
- ✅ `arm upgrade` - Upgrade to latest
- ✅ `arm list dependency` - List dependencies
- ✅ `arm info dependency` - Show dependency details
- ✅ `arm outdated` - Check for outdated packages
- ✅ `arm set ruleset` - Configure ruleset
- ✅ `arm set promptset` - Configure promptset

### Utilities (Mostly Implemented & Tested)
- ✅ `arm clean cache` - Clean cache (with `--nuke` and `--max-age`)
- ✅ `arm clean sinks` - Clean sinks (with `--nuke`)
- ⚠️ `arm compile` - CLI implemented, service layer TODO

### Registry Types (All Implemented & Tested)
- ✅ Git Registry - Full implementation with branch/tag support
- ✅ GitLab Registry - Full implementation with project/group support
- ✅ Cloudsmith Registry - Full implementation
- ✅ Archive Support - `.zip` and `.tar.gz` extraction for all registries

### Compilers (All Implemented & Tested)
- ✅ Cursor Compiler - Rulesets (`.mdc`) and Promptsets (`.md`)
- ✅ Amazon Q Compiler - Rulesets and Promptsets (`.md`)
- ✅ Copilot Compiler - Instructions format (`.instructions.md`)
- ✅ Markdown Compiler - Generic markdown output

### Core Infrastructure (All Implemented & Tested)
- ✅ Version Resolution - Semantic versioning with constraints
- ✅ Package Storage - Local cache with metadata
- ✅ Manifest Management - `arm.json` handling
- ✅ Lock File Management - `arm-lock.json` handling
- ✅ Index Management - `arm-index.json` tracking
- ✅ Parser - YAML resource validation
- ✅ File Type Detection - Ruleset vs Promptset detection
- ✅ Sink Layouts - Hierarchical and flat layouts
- ✅ Priority Resolution - Conflict resolution for rulesets
- ✅ Pattern Matching - Include/exclude glob patterns
- ✅ Authentication - `.armrc` support for GitLab/Cloudsmith

---

## Implementation Approach

### For Priority 1 (arm compile):

1. **Study existing patterns**:
   - Review `InstallRuleset` and `InstallPromptset` for file handling patterns
   - Review compiler tests for expected behavior
   - Review parser usage in existing code

2. **Implementation steps**:
   ```
   a. Validate and normalize input paths
   b. Discover files (handle directories with include/exclude)
   c. For each file:
      - Detect type (ruleset/promptset)
      - Parse YAML
      - Validate schema
      - If validate-only: continue
      - Compile to target format
      - Write to output directory
   d. Handle errors per fail-fast setting
   ```

3. **Testing**:
   - Update `cmd/arm/compile_test.go` expectations
   - Add service layer tests
   - Test all tool formats
   - Test validation-only mode
   - Test error handling

### For Priority 2 (migration guide):

1. Review git history for v2 to v3 changes
2. Document breaking changes
3. Create migration examples
4. Update README reference

### For Priority 3 (test improvements):

1. Implement skipped archive test
2. Review and fix GetTags assertion
3. Ensure all tests pass

---

## Notes

- **Test Status**: 304 tests across 58 test files - excellent coverage
- **Code Quality**: Well-structured, idiomatic Go with clean separation of concerns
- **Architecture**: Clean separation (cmd, service, internal packages)
- **Documentation**: Comprehensive specs with good inline comments
- **Patterns**: Consistent use of managers, factories, and interfaces
- **Verification**: All components verified present and functional except noted gaps

**Verified Components**:
- ✅ All 4 compilers implemented (Cursor, AmazonQ, Copilot, Markdown)
- ✅ All parsers implemented (Ruleset, Promptset)
- ✅ File type detection implemented
- ✅ All 3 registry types implemented (Git, GitLab, Cloudsmith)
- ✅ Archive extraction fully implemented with 10 comprehensive tests
- ✅ All CLI commands implemented
- ✅ All service methods implemented except CompileFiles

---

## Estimated Effort

- **Priority 1** (arm compile service layer): 4-6 hours
  - Service layer implementation: 2-3 hours
  - Testing and validation: 2-3 hours

- **Priority 2** (migration guide): 1-2 hours
  - Research git history for v2→v3 changes: 30-60 minutes
  - Documentation writing: 30-60 minutes

- **Priority 3** (test improvements): 1-2 hours
  - Archive end-to-end test: 30-60 minutes
  - GetTags assertion review: 30-60 minutes

**Total Estimated Effort**: 6-10 hours to 100% completion

---

## Success Criteria

Project is considered complete when:
- [ ] All Priority 1 items implemented and tested
- [ ] All Priority 2 documentation created
- [ ] All Priority 3 tests passing without skips
- [ ] All tests pass: `go test ./...`
- [ ] No TODO comments in production code (only 1 exists: CompileFiles)
- [ ] All spec features implemented
- [ ] README accurately reflects capabilities

---

## Last Updated

2026-01-24 - Comprehensive verification analysis complete with 304 tests across 58 files confirmed
