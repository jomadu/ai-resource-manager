# Draft Plan: Add Kiro CLI Tool Support

## Priority-Ordered Tasks

### 1. Update Specifications (TDD Approach)
**Files:**
- `specs/sink-compilation.md` - Add Kiro to acceptance criteria and examples
- `specs/standalone-compilation.md` - Add Kiro to acceptance criteria and examples

**Changes:**
- Add `[x] Compile to Kiro format (.md for both rules and prompts)` to acceptance criteria
- Add Kiro examples showing `.kiro/steering/` and `.kiro/prompts/` paths
- Document that Kiro/AmazonQ/Markdown share identical output format
- Add note about factory mapping three tools to same generators

**Rationale:** Define expected behavior before implementation (TDD)

### 2. Add Kiro Tool Constant
**File:** `internal/arm/compiler/types.go`
- Add `Kiro Tool = "kiro"` constant after AmazonQ

### 3. Refactor Factory to Eliminate Duplication
**Files:** 
- `internal/arm/compiler/factory.go` - Map AmazonQ/Kiro/Markdown to same generators
- `internal/arm/compiler/amazonq.go` - DELETE (duplicate of markdown.go)
- `internal/arm/compiler/amazonq_test.go` - DELETE (covered by markdown tests)

**Changes:**
- `NewRuleGenerator()`: `case Markdown, AmazonQ, Kiro: return &MarkdownRuleGenerator{}, nil`
- `NewPromptGenerator()`: `case Markdown, AmazonQ, Kiro: return &MarkdownPromptGenerator{}, nil`
- `NewRuleFilenameGenerator()`: `case Markdown, AmazonQ, Kiro: return &MarkdownRuleFilenameGenerator{}, nil`
- `NewPromptFilenameGenerator()`: `case Markdown, AmazonQ, Kiro: return &MarkdownPromptFilenameGenerator{}, nil`

**Rationale:** Eliminates ~200 lines of duplicate code, makes format equivalence explicit

**Dependencies:** Requires step 2 (Kiro constant)

### 4. Update CLI Tool Validation
**File:** `cmd/arm/main.go`
- Add "kiro" to valid tool list in `handleAddSink()`
- Add "kiro" to valid tool list in `handleCompile()`

**Dependencies:** Requires step 2 (Kiro constant)

### 5. Add E2E Tests
**File:** `test/e2e/compile_test.go`
- `TestKiroRulesetCompilation` - Verify ruleset to .kiro/steering/
- `TestKiroPromptsetCompilation` - Verify promptset to .kiro/prompts/
- `TestKiroPriorityIndex` - Verify arm_index.md generation
- Follow existing AmazonQ test patterns

**Dependencies:** Requires steps 2-4 (implementation complete)

### 6. Update User Documentation
**Files:**
- `README.md` - Add Kiro to supported tools list
- `docs/concepts.md` - Document Kiro tool support
- `docs/sinks.md` - Add Kiro sink examples
- `docs/commands.md` - Add Kiro command examples

**Changes:**
- Add Kiro examples: `arm add sink --tool kiro kiro-steering .kiro/steering`
- Add Kiro examples: `arm install ruleset my-registry/clean-code kiro-steering`
- Document default paths: `.kiro/steering/` (rulesets), `.kiro/prompts/` (promptsets)
- Note: Kiro uses markdown format (same as AmazonQ)

**Dependencies:** Can happen anytime after step 1 (specs define behavior)

## Acceptance Criteria

- [ ] Specs updated with Kiro acceptance criteria and examples
- [ ] Kiro tool constant added to types.go
- [ ] Factory maps Kiro/AmazonQ/Markdown to same generators
- [ ] amazonq.go and amazonq_test.go deleted
- [ ] `arm add sink --tool kiro` works
- [ ] `arm add sink --tool amazonq` still works (backward compatible)
- [ ] `arm install ruleset ... kiro-steering` installs to .kiro/steering/
- [ ] `arm install promptset ... kiro-prompts` installs to .kiro/prompts/
- [ ] Priority index (arm_index.md) generated for rulesets
- [ ] Hierarchical layout: arm/{registry}/{package}/{version}/
- [ ] All existing tests pass (including AmazonQ E2E tests)
- [ ] E2E tests added for Kiro compilation
- [ ] User docs updated with Kiro examples

## Task Dependencies

```
Step 1 (Specs) → Independent (can start immediately)
Step 2 (Constant) → Independent
Step 3 (Factory) → Requires Step 2
Step 4 (CLI) → Requires Step 2
Step 5 (Tests) → Requires Steps 2, 3, 4
Step 6 (Docs) → Requires Step 1 (specs define behavior)
```

## Implementation Notes

**Key Architectural Change:** Refactoring eliminates duplication by making format equivalence explicit in factory. Three tool constants map to one implementation.

**Backward Compatibility:** AmazonQ tool constant remains valid, mapped to markdown generators.

**Output Format:** All three tools (Kiro/AmazonQ/Markdown) produce:
- `.md` file extension
- YAML frontmatter with metadata
- Same filename pattern: `{id}_{id}.md`
- Hierarchical layout: `{sink}/arm/{registry}/{package}/{version}/{file}`

**Priority Index:** Standard markdown format (arm_index.md) works for all three tools.
