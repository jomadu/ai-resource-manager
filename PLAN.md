# Draft Plan: Add Kiro CLI Tool Support

## Priority-Ordered Tasks

### 1. Add Kiro Tool Constant
**File:** `internal/arm/compiler/types.go`
- Add `Kiro Tool = "kiro"` constant after AmazonQ

### 2. Refactor Factory to Eliminate Duplication
**Files:** 
- `internal/arm/compiler/factory.go` - Map AmazonQ/Kiro/Markdown to same generators
- `internal/arm/compiler/amazonq.go` - DELETE (duplicate of markdown.go)
- `internal/arm/compiler/amazonq_test.go` - DELETE (covered by markdown tests)

**Changes:**
- `NewRuleGenerator()`: Map `Markdown, AmazonQ, Kiro` to `MarkdownRuleGenerator`
- `NewPromptGenerator()`: Map `Markdown, AmazonQ, Kiro` to `MarkdownPromptGenerator`
- `NewRuleFilenameGenerator()`: Map `Markdown, AmazonQ, Kiro` to `MarkdownRuleFilenameGenerator`
- `NewPromptFilenameGenerator()`: Map `Markdown, AmazonQ, Kiro` to `MarkdownPromptFilenameGenerator`

**Rationale:** AmazonQ, Kiro, and Markdown produce identical output (.md files with same structure)

### 3. Update CLI Tool Validation
**File:** `cmd/arm/main.go`
- Add "kiro" to valid tool list in `handleAddSink()`
- Add "kiro" to valid tool list in `handleCompile()`

### 4. Update Sink Manager
**File:** `internal/arm/sink/manager.go`
- Verify kiro tool type is handled (should work automatically via factory)
- Priority index generation should work for kiro (uses markdown format)

### 5. Add E2E Tests
**File:** `test/e2e/compile_test.go`
- `TestKiroRulesetCompilation` - Verify ruleset compilation to .kiro/steering/
- `TestKiroPromptsetCompilation` - Verify promptset compilation to .kiro/prompts/
- `TestKiroPriorityIndex` - Verify arm_index.md generation
- Follow existing test patterns (copy from AmazonQ tests)

### 6. Update Specifications
**Files:**
- `specs/sink-compilation.md` - Add Kiro examples, update tool list
- `specs/standalone-compilation.md` - Add Kiro to tool list and examples

**Changes:**
- Add Kiro to "Acceptance Criteria" tool lists
- Add Kiro examples showing .kiro/steering/ and .kiro/prompts/
- Document that Kiro uses markdown format (same as AmazonQ)

### 7. Update User Documentation
**Files:**
- `README.md` - Add Kiro to supported tools list
- `docs/concepts.md` - Document Kiro tool support
- `docs/sinks.md` - Add Kiro sink configuration examples
- `docs/commands.md` - Update examples to include Kiro

**Changes:**
- Add Kiro examples for `arm add sink --tool kiro`
- Add Kiro examples for `arm install ruleset/promptset`
- Document default paths: `.kiro/steering/` and `.kiro/prompts/`

## Acceptance Criteria

- [ ] Kiro tool constant added to types.go
- [ ] Factory maps Kiro/AmazonQ/Markdown to same generators
- [ ] amazonq.go and amazonq_test.go deleted
- [ ] `arm add sink --tool kiro` works
- [ ] `arm add sink --tool amazonq` still works (mapped to markdown)
- [ ] `arm install ruleset ... kiro-steering` installs to .kiro/steering/
- [ ] `arm install promptset ... kiro-prompts` installs to .kiro/prompts/
- [ ] Priority index (arm_index.md) generated for rulesets
- [ ] Hierarchical layout used (arm/{registry}/{package}/{version}/)
- [ ] All existing tests pass (including AmazonQ E2E tests)
- [ ] E2E tests added for Kiro compilation
- [ ] Specs updated with Kiro examples
- [ ] User docs updated with Kiro examples

## Dependencies

- Steps 1-2 must complete before 3-4 (need tool constant and factory)
- Step 5 depends on 1-4 (need implementation to test)
- Steps 6-7 can happen after understanding implementation

## Notes

- Kiro, AmazonQ, and Markdown all produce identical output
- Refactoring eliminates ~200 lines of duplicate code
- All three tools use .md extension, YAML frontmatter, and same filename pattern
- Hierarchical layout: `{sink}/arm/{registry}/{package}/{version}/{file}`
- Priority index uses standard markdown format (arm_index.md)
