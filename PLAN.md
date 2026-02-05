# Draft Plan: Add Kiro CLI Tool Support

## Priority-Ordered Tasks

### 1. Update Specification - Document Kiro Tool Support
**File:** `specs/sink-compilation.md`

**Changes:**
- Add Kiro to list of supported tools in acceptance criteria
- Document Kiro default paths (`.kiro/steering/`, `.kiro/prompts/`)
- Add Kiro to priority index generation algorithm (uses `arm_index.md`)
- Add Kiro examples showing hierarchical layout
- Document that Kiro/AmazonQ/Markdown share identical format

**Acceptance:**
- [ ] Kiro listed in acceptance criteria
- [ ] Kiro paths documented
- [ ] Kiro examples added
- [ ] Format equivalence documented

**Dependencies:** None

---

### 2. Refactor Compiler - Eliminate Duplication and Add Kiro
**Files:** 
- `internal/arm/compiler/types.go` - Add Kiro constant
- `internal/arm/compiler/factory.go` - Map Kiro/AmazonQ/Markdown to same generators
- `internal/arm/compiler/amazonq.go` - DELETE (duplicate of markdown.go)
- `internal/arm/compiler/amazonq_test.go` - DELETE (covered by markdown tests)

**Changes:**
- Add `Kiro Tool = "kiro"` constant to types.go
- Update all four factory methods to map Kiro/AmazonQ/Markdown to markdown generators:
  - `NewRuleGenerator()` → `MarkdownRuleGenerator`
  - `NewPromptGenerator()` → `MarkdownPromptGenerator`
  - `NewRuleFilenameGenerator()` → `MarkdownRuleFilenameGenerator`
  - `NewPromptFilenameGenerator()` → `MarkdownPromptFilenameGenerator`
- Delete amazonq.go and amazonq_test.go files

**Acceptance:**
- [ ] Kiro constant added
- [ ] Factory maps all three tools to markdown generators
- [ ] amazonq.go deleted
- [ ] amazonq_test.go deleted
- [ ] All existing tests pass (including Amazon Q E2E tests)

**Dependencies:** Task 1 (spec defines behavior)

---

### 3. Update Sink Management - Handle Kiro Tool Type
**Files:**
- `cmd/arm/main.go` - Add "kiro" to valid tool options
- `internal/arm/sink/manager.go` - Handle kiro tool type in priority index generation

**Changes:**
- Add "kiro" to tool validation in CLI
- Ensure priority index generation handles kiro tool (should use `arm_index.md` like markdown/amazonq)

**Acceptance:**
- [ ] `arm add sink --tool kiro` command works
- [ ] `arm install ruleset ... kiro-sink` works
- [ ] `arm install promptset ... kiro-sink` works
- [ ] Priority index (`arm_index.md`) generated for kiro rulesets
- [ ] Hierarchical layout used

**Dependencies:** Task 2 (compiler must support kiro)

---

### 4. Add E2E Tests - Verify Kiro Compilation
**File:** `test/e2e/compile_test.go`

**Changes:**
- Add `TestKiroRulesetCompilation` - Verify ruleset compilation to .kiro/steering/
- Add `TestKiroPromptsetCompilation` - Verify promptset compilation to .kiro/prompts/
- Add `TestKiroPriorityIndex` - Verify arm_index.md generation
- Add `TestKiroMultiplePriorities` - Verify priority ordering

**Acceptance:**
- [ ] All four Kiro E2E tests pass
- [ ] Tests verify correct file paths
- [ ] Tests verify markdown format
- [ ] Tests verify priority index generation

**Dependencies:** Task 3 (sink management must work)

---

### 5. Update User Documentation - Add Kiro Examples
**Files:**
- `README.md` - Add Kiro to supported tools list
- `docs/concepts.md` - Document Kiro tool support
- `docs/sinks.md` - Add Kiro sink configuration examples
- `docs/commands.md` - Update examples to include Kiro

**Changes:**
- Add Kiro to "Configure sinks" examples
- Add Kiro to "Install to multiple tools" examples
- Document `.kiro/steering/` and `.kiro/prompts/` paths
- Note that Kiro uses markdown format (same as Amazon Q)

**Acceptance:**
- [ ] README.md includes Kiro
- [ ] docs/concepts.md documents Kiro
- [ ] docs/sinks.md has Kiro examples
- [ ] docs/commands.md has Kiro examples

**Dependencies:** Task 4 (implementation complete and tested)

---

## Task Dependencies

```
1. Update Spec (no dependencies)
   ↓
2. Refactor Compiler (depends on spec)
   ↓
3. Update Sink Management (depends on compiler)
   ↓
4. Add E2E Tests (depends on sink management)
   ↓
5. Update Documentation (depends on tests passing)
```

## Implementation Notes

### Key Design Decision: Eliminate Duplication
The current codebase has `amazonq.go` as a duplicate of `markdown.go`. This refactor:
- Makes format equivalence explicit (Kiro/AmazonQ/Markdown → same generators)
- Reduces code duplication (delete ~200 lines of duplicate code)
- Simplifies maintenance (one implementation for three tools)
- Preserves backward compatibility (Amazon Q still works, just mapped differently)

### Testing Strategy
- Existing Amazon Q E2E tests continue to pass (validates backward compatibility)
- New Kiro E2E tests verify identical behavior
- Unit tests for markdown generators cover all three tools

### Rollout
- No breaking changes (Amazon Q continues to work)
- Kiro is additive (new tool support)
- Documentation updated to reflect new capability
