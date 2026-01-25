# Implementation Plan

## Current Status

✅ **Documentation Restructuring Complete** (ralph-0.0.43)

All phases of the documentation restructuring project have been completed:
- Phase 0: Pre-Migration Setup - PROMPT files updated for dual docs/specs support
- Phase 1: File Migration - All user docs moved from specs/ to docs/
- Phase 2: Reference Updates - All links updated to point to docs/
- Phase 3: Builder Spec Creation - All 9 builder specs created following Ralph methodology
- Phase 4: Validation - Documentation structure and links verified

The project now has:
- User documentation in `docs/` (11 files + examples/)
- Builder specifications in `specs/` (10 files including TEMPLATE.md and e2e-testing.md)
- All tests passing
- No broken links

## Next Steps

The documentation restructuring is complete. Current work:

1. **Fix pattern matching for ** glob support** (IN PROGRESS - ralph-0.0.45)
   - ✅ Created shared `core.MatchPattern()` function with full `**` support
   - ✅ Updated Git, GitLab, and Cloudsmith registries to use shared function
   - ✅ All unit tests pass for pattern matching
   - ❌ `TestArchiveWithIncludeExcludePatterns` still fails - integration issue
   - **Issue**: Pattern `security/**/*.yml` should match `security/rule1.yml` from extracted archive, but files aren't being installed
   - **Next**: Add debug logging to trace file paths through extraction → filtering → parsing → compilation flow
   - **Hypothesis**: Files may be filtered/dropped at parsing or compilation stage, not pattern matching stage

2. **Merge to main**: The `docs-migration` branch is ready for merge (all tests pass, documentation complete)
3. **Feature enhancements** based on user feedback
4. **Performance optimizations**
5. **Additional registry types** if needed
6. **Enhanced error messages and diagnostics**

See [SPECIFICATION_PHILOSOPHY.md](./SPECIFICATION_PHILOSOPHY.md) for guidance on writing new specifications.

## Archive: Completed Phases

<details>
<summary>Click to expand completed work</summary>

### Confirmed Existing Files (specs/)
- ✅ `specs/concepts.md` - User documentation
- ✅ `specs/commands.md` - User documentation (27KB, comprehensive CLI reference)
- ✅ `specs/registries.md` - User documentation
- ✅ `specs/git-registry.md` - User documentation
- ✅ `specs/gitlab-registry.md` - User documentation
- ✅ `specs/cloudsmith-registry.md` - User documentation
- ✅ `specs/sinks.md` - User documentation
- ✅ `specs/storage.md` - User documentation
- ✅ `specs/resource-schemas.md` - User documentation
- ✅ `specs/armrc.md` - User documentation
- ✅ `specs/migration-v2-to-v3.md` - User documentation (11KB migration guide)
- ✅ `specs/e2e-testing.md` - Builder-oriented (KEEP IN SPECS)
- ✅ `specs/examples/` - Example files

### Confirmed Existing Files (docs/)
- ✅ `docs/examples/` - Already exists (partial migration started)

### References to Update
- ✅ `README.md` - 10 references to `specs/` files
- ✅ `AGENTS.md` - 1 reference to `specs/`
- ✅ `specs/migration-v2-to-v3.md` - 1 self-reference
- ✅ `.amazonq/prompts/` - Multiple references (workflow documentation)
- ✅ `SPECIFICATION_PHILOSOPHY.md` - 4 references (meta-documentation)

## Critical: Loop.sh Compatibility

**PROMPT_build.md references `specs/*` on line 1**. The migration must be atomic to avoid breaking the build loop.

### Migration Strategy
1. **Create all docs/** - Ensure target directory exists
2. **Copy (don't move) specs → docs** - Keep both during transition
3. **Update PROMPT_build.md and PROMPT_plan.md** - Change `specs/*` references
4. **Verify loop.sh works** - Test with `./loop.sh plan 1`
5. **Delete old specs/** - Only after verification

### Updated PROMPT References Needed
- `PROMPT_build.md` line 1: `specs/*` → `docs/* and specs/*`
- `PROMPT_plan.md` line 1: `specs/*` → `docs/* and specs/*`
- After builder specs created: `docs/* and specs/*` (both needed)

## Prioritized Tasks

### Phase 0: Pre-Migration Setup (Priority: CRITICAL) ✅ COMPLETE
Ensure loop.sh won't break during migration

- [x] Update `PROMPT_build.md` line 1
  - Change: `Study 'specs/*' with up to 500 parallel subagents`
  - To: `Study 'docs/*' and 'specs/*' with up to 500 parallel subagents to learn the application specifications`
  
- [x] Update `PROMPT_plan.md` line 1
  - Change: `Study 'specs/*' with up to 250 parallel subagents`
  - To: `Study 'docs/*' and 'specs/*' with up to 250 parallel subagents to learn the application specifications`

- [x] Update `PROMPT_plan.md` line 6
  - Change: `compare it against 'specs/*'`
  - To: `compare it against 'docs/*' and 'specs/*'`

- [x] Update `PROMPT_build.md` line 14
  - Change: `If you find inconsistencies in the specs/*`
  - To: `If you find inconsistencies in the docs/* or specs/*`

**Verification:**
- [x] All tests pass (go test ./...)
- [x] PROMPT files updated successfully

### Phase 1: File Migration (Priority: HIGH) ✅ COMPLETE
Move user documentation from `specs/` to `docs/`, preserving `specs/e2e-testing.md`

- [x] Move `specs/concepts.md` → `docs/concepts.md`
- [x] Move `specs/commands.md` → `docs/commands.md`
- [x] Move `specs/registries.md` → `docs/registries.md`
- [x] Move `specs/git-registry.md` → `docs/git-registry.md`
- [x] Move `specs/gitlab-registry.md` → `docs/gitlab-registry.md`
- [x] Move `specs/cloudsmith-registry.md` → `docs/cloudsmith-registry.md`
- [x] Move `specs/sinks.md` → `docs/sinks.md`
- [x] Move `specs/storage.md` → `docs/storage.md`
- [x] Move `specs/resource-schemas.md` → `docs/resource-schemas.md`
- [x] Move `specs/armrc.md` → `docs/armrc.md`
- [x] Move `specs/migration-v2-to-v3.md` → `docs/migration-v2-to-v3.md`
- [x] Merge `specs/examples/` → `docs/examples/` (docs/examples/ already exists)

**Verification:**
- [x] Confirm `specs/` contains only `e2e-testing.md` after migration
- [x] Confirm all 12 files exist in `docs/`
- [x] All tests pass

### Phase 2: Reference Updates (Priority: HIGH) ✅ COMPLETE
Update all references from `specs/` to `docs/` for migrated files

- [x] Update `README.md` (10 references)
  - Line 24: `specs/git-registry.md` → `docs/git-registry.md`
  - Line 59: `specs/migration-v2-to-v3.md` → `docs/migration-v2-to-v3.md`
  - Lines 151-155: Update documentation section links
  - Lines 159-161: Update registry types section links
- [x] Update `AGENTS.md` (1 reference)
  - Line 71: Update specs reference to clarify docs vs specs
- [x] Update `docs/migration-v2-to-v3.md` (1 self-reference)
  - Line 468: Update to reference `docs/` directory
- [x] Review `.amazonq/prompts/` references (informational, workflow docs - intentional)
- [x] Review `SPECIFICATION_PHILOSOPHY.md` references (meta-documentation - intentional)

**Verification:**
- [x] All tests pass
- [x] Only intentional references remain (PROMPT files, IMPLEMENTATION_PLAN, meta-docs)

### Phase 3: Builder Spec Creation (Priority: MEDIUM)
Create builder-oriented specifications following Ralph methodology. See [SPEC_COVERAGE_ANALYSIS.md](./SPEC_COVERAGE_ANALYSIS.md) for complete traceability matrix.

**IMPORTANT**: Each spec below should be completed as its own independent task. Do NOT attempt to complete multiple specs in a single session to avoid context bloat. Pick one spec, complete it fully, verify tests pass, commit, and move to the next spec in a fresh session.

- [x] Create `specs/TEMPLATE.md` (foundation for all other specs)
  - Include: JTBD, Activities, Acceptance Criteria, Data Structures, Algorithm, Edge Cases, Dependencies, Examples
  - Reference: SPECIFICATION_PHILOSOPHY.md template section
  - **Verification**: Template created successfully, all tests pass

- [x] Create `specs/version-resolution.md`
  - **JTBD**: Resolve package versions from registries
  - **Maps to**: internal/arm/core/version.go, constraint.go, helpers.go
  - **User docs**: concepts.md, git-registry.md
  - **Coverage**: Semver parsing, constraint matching, tag/branch priority, version comparison
  - **Edge cases**: No versions, malformed tags, network failures, branch not found
  - **Verification**: Spec created with complete algorithm, data structures, and acceptance criteria

- [x] Create `specs/package-installation.md`
  - **JTBD**: Install, update, upgrade, uninstall packages
  - **Maps to**: internal/arm/service/service.go, manifest/, packagelockfile/
  - **User docs**: commands.md, concepts.md
  - **Coverage**: Install workflow, reinstall behavior (remove from old sinks), lock file updates, manifest updates
  - **Edge cases**: Missing sinks, version conflicts, partial failures, concurrent installs
  - **Verification**: Comprehensive spec with algorithms for install, update, upgrade, uninstall, and install-all

- [x] Create `specs/registry-management.md`
  - **JTBD**: Configure and manage registries
  - **Maps to**: internal/arm/registry/, internal/arm/manifest/
  - **User docs**: registries.md, git-registry.md, gitlab-registry.md, cloudsmith-registry.md
  - **Coverage**: Registry types (Git, GitLab, Cloudsmith), configuration storage, authentication integration, key generation
  - **Edge cases**: Invalid URLs, duplicate names, missing auth, network failures
  - **Verification**: Comprehensive spec with algorithms for add, list, remove, and factory pattern. All tests pass.

- [x] Create `specs/sink-compilation.md`
  - **JTBD**: Compile resources to tool-specific formats
  - **Maps to**: internal/arm/compiler/, internal/arm/sink/manager.go
  - **User docs**: sinks.md, commands.md
  - **Coverage**: Compilation algorithms per tool (Cursor, AmazonQ, Copilot, Markdown), layout modes (hierarchical vs flat), filename generation, truncation rules
  - **Edge cases**: Long filenames, special characters, empty resources, invalid YAML
  - **Verification**: Comprehensive spec with algorithms for compilation, filename generation, layout modes, and priority index. All tests pass.

- [x] Create `specs/priority-resolution.md`
  - **JTBD**: Resolve conflicts between overlapping rules
  - **Maps to**: internal/arm/sink/manager.go, internal/arm/compiler/generators.go
  - **User docs**: sinks.md, concepts.md
  - **Coverage**: Priority merging algorithm, conflict resolution rules, index generation, metadata embedding
  - **Edge cases**: Same priority, no rules, circular dependencies, missing metadata
  - **Verification**: Comprehensive spec created with algorithms for install, update, generate index, and embed metadata. All tests pass.

- [x] Create `specs/cache-management.md`
  - **JTBD**: Cache packages locally to avoid redundant downloads
  - **Maps to**: internal/arm/storage/
  - **User docs**: storage.md, commands.md
  - **Coverage**: Storage structure (~/.arm/storage), cache key generation, metadata schemas (registry, package, version), cleanup strategies (max-age, nuke)
  - **Edge cases**: Corrupted cache, disk full, concurrent access, stale metadata
  - **Verification**: Comprehensive spec created with algorithms for key generation, store, retrieve, clean by age, clean by last access, and nuke. All tests pass.

- [x] Create `specs/pattern-filtering.md`
  - **JTBD**: Filter package files using glob patterns
  - **Maps to**: internal/arm/registry/, internal/arm/core/archive.go
  - **User docs**: concepts.md, registries.md
  - **Coverage**: Glob matching, include/exclude logic (OR for includes, exclude overrides), archive extraction (zip, tar.gz), path sanitization
  - **Edge cases**: Invalid patterns, path traversal attacks, nested archives, empty results
  - **Verification**: Comprehensive spec created with algorithms for normalize patterns, match patterns, extract archives, and sanitize paths. All tests pass.

- [x] Create `specs/authentication.md`
  - **JTBD**: Authenticate with registries requiring tokens
  - **Maps to**: internal/arm/config/manager.go, internal/arm/registry/gitlab.go, cloudsmith.go
  - **User docs**: armrc.md
  - **Coverage**: .armrc parsing (INI format), token resolution (local vs global), environment variable substitution, security (file permissions)
  - **Edge cases**: Missing .armrc, invalid format, expired tokens, permission errors
  - **Verification**: Comprehensive spec created with algorithms for parse .armrc, hierarchical lookup, environment variable expansion, and auth key generation. All tests pass.

**Verification:**
- Each spec follows TEMPLATE.md structure
- Each spec has testable acceptance criteria
- Each spec maps to implementation (see SPEC_COVERAGE_ANALYSIS.md)
- Each spec focuses on single concern
- Traceability: Command → User Doc → Builder Spec → Implementation

### Phase 4: Validation (Priority: LOW) ✅ COMPLETE
Ensure migration is complete and consistent

- [x] Verify no broken links in documentation
- [x] Verify all user docs in `docs/`
- [x] Verify all builder specs in `specs/`
- [x] Verify `specs/e2e-testing.md` remains in place
- [x] Run tests to ensure no functionality broken

</details>

## Success Criteria ✅ COMPLETE

- [x] All user documentation moved to `docs/`
- [x] All references updated to point to `docs/`
- [x] `specs/` contains only builder-oriented specifications
- [x] Each builder spec follows Ralph methodology template
- [x] No broken links in documentation
- [x] All tests pass

## Notes

- `docs/examples/` already exists, merge with `specs/examples/`
- `specs/e2e-testing.md` is already builder-oriented, keep in specs
- Builder specs should enable test-driven development
- Specs are disposable - iterate based on implementation learnings
