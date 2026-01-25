# Implementation Plan: Documentation Restructuring

## Goal
Separate user documentation from builder specifications. Align specs with Ralph methodology (JTBD → activities → acceptance criteria → tasks).

See [SPECIFICATION_PHILOSOPHY.md](./SPECIFICATION_PHILOSOPHY.md) for detailed guidance on writing effective builder-oriented specifications.

## Current State Analysis

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

### Phase 0: Pre-Migration Setup (Priority: CRITICAL)
Ensure loop.sh won't break during migration

- [ ] Update `PROMPT_build.md` line 1
  - Change: `Study 'specs/*' with up to 500 parallel subagents`
  - To: `Study 'docs/*' and 'specs/*' with up to 500 parallel subagents to learn the application specifications`
  
- [ ] Update `PROMPT_plan.md` line 1
  - Change: `Study 'specs/*' with up to 250 parallel subagents`
  - To: `Study 'docs/*' and 'specs/*' with up to 250 parallel subagents to learn the application specifications`

- [ ] Update `PROMPT_plan.md` line 6
  - Change: `compare it against 'specs/*'`
  - To: `compare it against 'docs/*' and 'specs/*'`

- [ ] Update `PROMPT_build.md` line 14
  - Change: `If you find inconsistencies in the specs/*`
  - To: `If you find inconsistencies in the docs/* or specs/*`

**Verification:**
- Run `./loop.sh plan 1` to verify prompts work
- Confirm no errors about missing specs

### Phase 1: File Migration (Priority: HIGH)
Move user documentation from `specs/` to `docs/`, preserving `specs/e2e-testing.md`

- [ ] Move `specs/concepts.md` → `docs/concepts.md`
- [ ] Move `specs/commands.md` → `docs/commands.md`
- [ ] Move `specs/registries.md` → `docs/registries.md`
- [ ] Move `specs/git-registry.md` → `docs/git-registry.md`
- [ ] Move `specs/gitlab-registry.md` → `docs/gitlab-registry.md`
- [ ] Move `specs/cloudsmith-registry.md` → `docs/cloudsmith-registry.md`
- [ ] Move `specs/sinks.md` → `docs/sinks.md`
- [ ] Move `specs/storage.md` → `docs/storage.md`
- [ ] Move `specs/resource-schemas.md` → `docs/resource-schemas.md`
- [ ] Move `specs/armrc.md` → `docs/armrc.md`
- [ ] Move `specs/migration-v2-to-v3.md` → `docs/migration-v2-to-v3.md`
- [ ] Merge `specs/examples/` → `docs/examples/` (docs/examples/ already exists)

**Verification:**
- Confirm `specs/` contains only `e2e-testing.md` after migration
- Confirm all 12 files exist in `docs/`

### Phase 2: Reference Updates (Priority: HIGH)
Update all references from `specs/` to `docs/` for migrated files

- [ ] Update `README.md` (10 references)
  - Line 24: `specs/git-registry.md` → `docs/git-registry.md`
  - Line 59: `specs/migration-v2-to-v3.md` → `docs/migration-v2-to-v3.md`
  - Lines 151-155: Update documentation section links
  - Lines 159-161: Update registry types section links
- [ ] Update `AGENTS.md` (1 reference)
  - Line 71: Update specs reference to clarify docs vs specs
- [ ] Update `specs/migration-v2-to-v3.md` → `docs/migration-v2-to-v3.md` (1 self-reference)
  - Line 468: Update to reference `docs/` directory
- [ ] Review `.amazonq/prompts/` references (informational, may not need updates)
- [ ] Review `SPECIFICATION_PHILOSOPHY.md` references (meta-documentation, intentional)

**Verification:**
- Run `grep -r "specs/" *.md` to find remaining references
- Confirm only intentional references remain (e2e-testing.md, meta-docs)

### Phase 3: Builder Spec Creation (Priority: MEDIUM)
Create builder-oriented specifications following Ralph methodology. See [SPEC_COVERAGE_ANALYSIS.md](./SPEC_COVERAGE_ANALYSIS.md) for complete traceability matrix.

- [ ] Create `specs/TEMPLATE.md` (foundation for all other specs)
  - Include: JTBD, Activities, Acceptance Criteria, Data Structures, Algorithm, Edge Cases, Dependencies, Examples
  - Reference: SPECIFICATION_PHILOSOPHY.md template section

- [ ] Create `specs/version-resolution.md`
  - **JTBD**: Resolve package versions from registries
  - **Maps to**: internal/arm/core/version.go, constraint.go, helpers.go
  - **User docs**: concepts.md, git-registry.md
  - **Coverage**: Semver parsing, constraint matching, tag/branch priority, version comparison
  - **Edge cases**: No versions, malformed tags, network failures, branch not found

- [ ] Create `specs/package-installation.md`
  - **JTBD**: Install, update, upgrade, uninstall packages
  - **Maps to**: internal/arm/service/service.go, manifest/, packagelockfile/
  - **User docs**: commands.md, concepts.md
  - **Coverage**: Install workflow, reinstall behavior (remove from old sinks), lock file updates, manifest updates
  - **Edge cases**: Missing sinks, version conflicts, partial failures, concurrent installs

- [ ] Create `specs/registry-management.md`
  - **JTBD**: Configure and manage registries
  - **Maps to**: internal/arm/registry/, internal/arm/manifest/
  - **User docs**: registries.md, git-registry.md, gitlab-registry.md, cloudsmith-registry.md
  - **Coverage**: Registry types (Git, GitLab, Cloudsmith), configuration storage, authentication integration, key generation
  - **Edge cases**: Invalid URLs, duplicate names, missing auth, network failures

- [ ] Create `specs/sink-compilation.md`
  - **JTBD**: Compile resources to tool-specific formats
  - **Maps to**: internal/arm/compiler/, internal/arm/sink/manager.go
  - **User docs**: sinks.md, commands.md
  - **Coverage**: Compilation algorithms per tool (Cursor, AmazonQ, Copilot, Markdown), layout modes (hierarchical vs flat), filename generation, truncation rules
  - **Edge cases**: Long filenames, special characters, empty resources, invalid YAML

- [ ] Create `specs/priority-resolution.md`
  - **JTBD**: Resolve conflicts between overlapping rules
  - **Maps to**: internal/arm/sink/manager.go, internal/arm/compiler/generators.go
  - **User docs**: sinks.md, concepts.md
  - **Coverage**: Priority merging algorithm, conflict resolution rules, index generation, metadata embedding
  - **Edge cases**: Same priority, no rules, circular dependencies, missing metadata

- [ ] Create `specs/cache-management.md`
  - **JTBD**: Cache packages locally to avoid redundant downloads
  - **Maps to**: internal/arm/storage/
  - **User docs**: storage.md, commands.md
  - **Coverage**: Storage structure (~/.arm/storage), cache key generation, metadata schemas (registry, package, version), cleanup strategies (max-age, nuke)
  - **Edge cases**: Corrupted cache, disk full, concurrent access, stale metadata

- [ ] Create `specs/pattern-filtering.md`
  - **JTBD**: Filter package files using glob patterns
  - **Maps to**: internal/arm/registry/, internal/arm/core/archive.go
  - **User docs**: concepts.md, registries.md
  - **Coverage**: Glob matching, include/exclude logic (OR for includes, exclude overrides), archive extraction (zip, tar.gz), path sanitization
  - **Edge cases**: Invalid patterns, path traversal attacks, nested archives, empty results

- [ ] Create `specs/authentication.md`
  - **JTBD**: Authenticate with registries requiring tokens
  - **Maps to**: internal/arm/config/manager.go, internal/arm/registry/gitlab.go, cloudsmith.go
  - **User docs**: armrc.md
  - **Coverage**: .armrc parsing (INI format), token resolution (local vs global), environment variable substitution, security (file permissions)
  - **Edge cases**: Missing .armrc, invalid format, expired tokens, permission errors

**Verification:**
- Each spec follows TEMPLATE.md structure
- Each spec has testable acceptance criteria
- Each spec maps to implementation (see SPEC_COVERAGE_ANALYSIS.md)
- Each spec focuses on single concern
- Traceability: Command → User Doc → Builder Spec → Implementation

### Phase 4: Validation (Priority: LOW)
Ensure migration is complete and consistent

- [ ] Verify no broken links in documentation
- [ ] Verify all user docs in `docs/`
- [ ] Verify all builder specs in `specs/`
- [ ] Verify `specs/e2e-testing.md` remains in place
- [ ] Run tests to ensure no functionality broken

## Success Criteria

- [ ] All user documentation moved to `docs/`
- [ ] All references updated to point to `docs/`
- [ ] `specs/` contains only builder-oriented specifications
- [ ] Each builder spec follows Ralph methodology template
- [ ] No broken links in documentation
- [ ] All tests pass

## Notes

- `docs/examples/` already exists, merge with `specs/examples/`
- `specs/e2e-testing.md` is already builder-oriented, keep in specs
- Builder specs should enable test-driven development
- Specs are disposable - iterate based on implementation learnings
