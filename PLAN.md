# Draft Plan: Add Kiro CLI Tool Support

## Priority-Ordered Tasks

### 1. Update Specifications
- [ ] Update `specs/sink-compilation.md` - Add Kiro tool to acceptance criteria, examples, and implementation mapping
- [ ] Update `specs/standalone-compilation.md` - Add Kiro to tool list and examples

**Rationale**: Specs define "what should exist" before implementation

**Acceptance**: 
- Kiro documented in both specs with examples
- Priority index generation documented for Kiro
- Default paths documented (`.kiro/steering/`, `.kiro/prompts/`)

---

### 2. Add Kiro Tool Constant
- [ ] Add `Kiro Tool = "kiro"` constant to `internal/arm/compiler/types.go`

**Rationale**: Foundation for all other changes

**Acceptance**: 
- Kiro constant added after AmazonQ in types.go

---

### 3. Refactor Factory to Eliminate Duplication
- [ ] Update `internal/arm/compiler/factory.go` - Map Kiro, AmazonQ, and Markdown to same generators
- [ ] Delete `internal/arm/compiler/amazonq.go` (duplicate of markdown.go)
- [ ] Delete `internal/arm/compiler/amazonq_test.go` (covered by markdown tests)

**Rationale**: Eliminates code duplication - all three tools produce identical output

**Changes**:
```go
// In all four factory methods (NewRuleGenerator, NewPromptGenerator, etc.)
case Markdown, AmazonQ, Kiro:
    return &MarkdownRuleGenerator{}, nil  // or appropriate generator
```

**Acceptance**:
- Factory maps all three tools to markdown generators
- amazonq.go and amazonq_test.go deleted
- All existing tests still pass (including Amazon Q E2E tests)

**Dependencies**: Task 2 (Kiro constant must exist)

---

### 4. Update CLI to Accept Kiro Tool
- [ ] Update `cmd/arm/main.go` - Add "kiro" to valid tool options in sink commands
- [ ] Update `internal/arm/sink/manager.go` - Handle kiro tool type (if needed)

**Rationale**: Enable users to create Kiro sinks

**Acceptance**:
- `arm add sink --tool kiro` command works
- `arm install ruleset ... kiro-sink` works
- `arm install promptset ... kiro-sink` works

**Dependencies**: Task 3 (factory must support Kiro)

---

### 5. Add E2E Tests for Kiro
- [ ] Add tests to `test/e2e/compile_test.go`:
  - `TestKiroRulesetCompilation` - Verify ruleset to .kiro/steering/
  - `TestKiroPromptsetCompilation` - Verify promptset to .kiro/prompts/
  - `TestKiroPriorityIndex` - Verify arm_index.md generation
  - `TestKiroMultiplePriorities` - Verify priority ordering

**Rationale**: Verify Kiro compilation works end-to-end

**Acceptance**:
- All Kiro E2E tests pass
- Tests verify hierarchical layout
- Tests verify priority index generation
- All existing tests still pass

**Dependencies**: Task 4 (CLI must support Kiro)

---

### 6. Update Documentation
- [ ] Update `README.md` - Add Kiro to supported tools list and quick start examples
- [ ] Update `docs/concepts.md` - Document Kiro tool support
- [ ] Update `docs/sinks.md` - Add Kiro sink configuration examples
- [ ] Update `docs/commands.md` - Add Kiro examples to command reference

**Rationale**: Users need to know about Kiro support

**Acceptance**:
- Kiro mentioned in all relevant docs
- Examples show `.kiro/steering/` and `.kiro/prompts/` paths
- Documentation consistent with other tools

**Dependencies**: Task 5 (implementation complete and tested)

---

## Implementation Notes

### Key Insight: Eliminate Duplication
Amazon Q, Markdown, and Kiro all produce **identical output**:
- Same `.md` extension
- Same metadata format (YAML frontmatter)
- Same filename pattern: `{rulesetID}_{ruleID}.md`

**Current state**: `amazonq.go` is a copy-paste duplicate of `markdown.go`

**Refactored state**: Map all three tool constants to the same generators in factory

### Reuse Existing Patterns
- Hierarchical layout (same as Cursor and Amazon Q)
- Priority index generation (same markdown format)
- Filename generation (same as markdown)
- Sink management (existing infrastructure)

### Testing Strategy
- Reuse existing E2E test patterns from other tools
- Verify priority index generation
- Verify hierarchical layout structure
- Ensure Amazon Q tests still pass (validates refactor)

## Related Specifications
- `specs/sink-compilation.md` - Sink management and compilation
- `specs/standalone-compilation.md` - Local file compilation
- `specs/priority-resolution.md` - Priority-based conflict resolution
