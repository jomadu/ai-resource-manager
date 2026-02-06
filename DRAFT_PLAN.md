# Draft Plan: Add Kiro CLI Tool Support

## Priority-Ordered Tasks

### 1. Update Specification (HIGHEST PRIORITY)
**File**: `specs/sink-compilation.md`

**Changes**:
- Add Kiro to acceptance criteria
- Document Kiro/AmazonQ/Markdown consolidation rationale
- Add Kiro examples (hierarchical layout, priority index)
- Update tool list in "Compile to X format" criteria
- Add note about code deduplication (delete amazonq.go)

**Acceptance**:
- [ ] Kiro tool documented in acceptance criteria
- [ ] Consolidation rationale explained (why delete amazonq.go)
- [ ] Kiro examples added (directory structure, priority index)
- [ ] Spec follows TEMPLATE.md structure

---

### 2. Add Kiro Tool Constant
**File**: `internal/arm/compiler/types.go`

**Changes**:
- Add `Kiro Tool = "kiro"` constant

**Acceptance**:
- [ ] Kiro constant added to Tool type

---

### 3. Refactor Factory to Consolidate Tools
**File**: `internal/arm/compiler/factory.go`

**Changes**:
- Map `Markdown`, `AmazonQ`, `Kiro` to same generators in all four factory methods:
  - `NewRuleGenerator()` → `MarkdownRuleGenerator`
  - `NewPromptGenerator()` → `MarkdownPromptGenerator`
  - `NewRuleFilenameGenerator()` → `MarkdownRuleFilenameGenerator`
  - `NewPromptFilenameGenerator()` → `MarkdownPromptFilenameGenerator`

**Acceptance**:
- [ ] All three tools map to markdown generators
- [ ] Factory tests pass

---

### 4. Delete Duplicate Amazon Q Files
**Files**: 
- `internal/arm/compiler/amazonq.go`
- `internal/arm/compiler/amazonq_test.go`

**Rationale**: These files are identical to markdown.go/markdown_test.go. Factory now maps amazonq tool to markdown generators.

**Acceptance**:
- [ ] amazonq.go deleted
- [ ] amazonq_test.go deleted
- [ ] All tests still pass (Amazon Q E2E tests use factory, not direct imports)

---

### 5. Update Sink Manager for Kiro
**File**: `internal/arm/sink/manager.go`

**Changes**:
- Handle "kiro" tool type in sink operations
- Generate `arm_index.md` for Kiro rulesets (same as AmazonQ/Markdown)

**Acceptance**:
- [ ] Kiro tool recognized in sink operations
- [ ] Priority index generated correctly

---

### 6. Update CLI Command Validation
**File**: `cmd/arm/main.go`

**Changes**:
- Add "kiro" to valid tool options in `add sink` command

**Acceptance**:
- [ ] `arm add sink --tool kiro` command works
- [ ] Help text includes kiro

---

### 7. Add E2E Tests for Kiro
**File**: `test/e2e/compile_test.go`

**Tests**:
- `TestKiroRulesetCompilation` - Verify ruleset compilation to .kiro/steering/
- `TestKiroPromptsetCompilation` - Verify promptset compilation to .kiro/prompts/
- `TestKiroPriorityIndex` - Verify arm_index.md generation
- `TestKiroMultiplePriorities` - Verify priority ordering

**Acceptance**:
- [ ] All Kiro E2E tests pass
- [ ] Tests verify hierarchical layout
- [ ] Tests verify priority index format

---

### 8. Update Documentation
**Files**:
- `README.md` - Add Kiro to supported tools list
- `docs/concepts.md` - Document Kiro tool support
- `docs/sinks.md` - Add Kiro sink configuration examples
- `docs/commands.md` - Update examples to include Kiro

**Acceptance**:
- [ ] Kiro mentioned in all relevant docs
- [ ] Examples show Kiro sink configuration
- [ ] Examples show Kiro installation commands

---

## Dependencies

- Task 1 (spec) can be done independently
- Tasks 2-6 (implementation) depend on spec completion
- Task 7 (tests) depends on tasks 2-6
- Task 8 (docs) can be done after task 1 (spec)

## Notes

**Refactor Rationale**: Amazon Q, Markdown, and Kiro all produce identical output:
- Same `.md` extension
- Same YAML frontmatter for rules
- Same plain body for prompts
- Same filename pattern: `{id}_{id}.md`

Current state has `amazonq.go` as a duplicate of `markdown.go`. Refactored state eliminates this duplication by mapping all three tools to the same generators in the factory.

**Testing Strategy**: Existing Amazon Q E2E tests will continue to work because they use the factory pattern, not direct imports of amazonq.go. The factory now maps amazonq → markdown generators.
