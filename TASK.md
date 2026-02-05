# TASK: Add Kiro CLI Tool Support

## Status: ✅ COMPLETE

## Goal

Enable ARM to install promptsets to `.kiro/prompts/` and rulesets to `.kiro/steering/` directories, supporting the Kiro CLI tool format.

## Background

Kiro CLI is the next evolution of Amazon Q Developer CLI. It uses:
- **`.kiro/steering/`** - For persistent project knowledge (rules/conventions) in markdown format
- **`.kiro/prompts/`** - For reusable prompt templates
- **Scope levels**: Workspace (`.kiro/`) and Global (`~/.kiro/`)

Kiro steering files are markdown documents that provide persistent context to the AI assistant, similar to Amazon Q rules but with a different directory structure.

## Requirements

### 1. Add Kiro Tool Support + Refactor Amazon Q

**Location**: `internal/arm/compiler/factory.go`

**Behavior**:
- Kiro, Amazon Q, and Markdown all produce **identical output**
- Map all three tools to the same markdown generators in the factory
- **Delete** `internal/arm/compiler/amazonq.go` (no longer needed)
- **Delete** `internal/arm/compiler/amazonq_test.go` (no longer needed)

**Factory Changes**:
```go
func (f *DefaultRuleGeneratorFactory) NewRuleGenerator(tool Tool) (RuleGenerator, error) {
	switch tool {
	case Cursor:
		return &CursorRuleGenerator{}, nil
	case Markdown, AmazonQ, Kiro:  // All three use markdown format
		return &MarkdownRuleGenerator{}, nil
	case Copilot:
		return &CopilotRuleGenerator{}, nil
	default:
		return nil, fmt.Errorf("unsupported tool: %s", tool)
	}
}
```

Apply same pattern to:
- `NewPromptGenerator()` 
- `NewRuleFilenameGenerator()`
- `NewPromptFilenameGenerator()`

**Benefits**:
- Eliminates duplicate code (amazonq.go is identical to markdown.go)
- Simplifies maintenance (one implementation for three tools)
- Makes it clear that these tools share the same format

### 2. Update Sink Configuration

**Files**: 
- `cmd/arm/main.go` - Add "kiro" as valid tool option
- `internal/arm/sink/manager.go` - Handle kiro tool type

**Default Sink Paths**:
- Rulesets: `.kiro/steering/`
- Promptsets: `.kiro/prompts/`

**Example Commands**:
```bash
# Add sinks for Kiro
arm add sink --tool kiro kiro-steering .kiro/steering
arm add sink --tool kiro kiro-prompts .kiro/prompts

# Install ruleset to steering
arm install ruleset my-registry/clean-code kiro-steering

# Install promptset to prompts
arm install promptset my-registry/code-review kiro-prompts
```

### 3. Priority Index Generation

**File**: `internal/arm/sink/manager.go`

**Behavior**:
- Generate `arm_index.md` for rulesets in `.kiro/steering/`
- Use same priority resolution logic as other tools
- Format: Standard markdown with priority ordering

**Example `arm_index.md`**:
```markdown
# ARM Rulesets

This file defines the installation priorities for rulesets managed by ARM.

## Priority Rules

**This index is the authoritative source of truth for ruleset priorities.** When conflicts arise between rulesets, follow this priority order:

1. **Higher priority numbers take precedence** over lower priority numbers
2. **Rules from higher priority rulesets override** conflicting rules from lower priority rulesets
3. **Always consult this index** to resolve any ambiguity about which rules to follow

## Installed Rulesets

### my-registry/team-standards@1.0.0
- **Priority:** 200
- **Rules:**
  - arm/my-registry/team-standards/1.0.0/rules/teamStandards_rule.md

### my-registry/clean-code@1.0.0
- **Priority:** 100
- **Rules:**
  - arm/my-registry/clean-code/1.0.0/rules/cleanCode_ruleOne.md
  - arm/my-registry/clean-code/1.0.0/rules/cleanCode_ruleTwo.md
```

### 4. Layout Mode

**Behavior**:
- Use **hierarchical layout** (same as Cursor and Amazon Q)
- Path structure: `{sink}/arm/{registry}/{package}/{version}/{file}`

**Example Structure**:
```
.kiro/steering/
└── arm/
    └── my-registry/
        └── clean-code/
            └── 1.0.0/
                └── rules/
                    ├── cleanCode_ruleOne.md
                    ├── cleanCode_ruleTwo.md
                    └── cleanCode_ruleThree.md

.kiro/prompts/
└── arm/
    └── my-registry/
        └── code-review/
            └── 1.0.0/
                └── prompts/
                    ├── codeReview_security.md
                    └── codeReview_performance.md
```

### 5. Documentation Updates

**Files to Update**:
- `README.md` - Add Kiro to list of supported tools
- `docs/concepts.md` - Document Kiro tool support
- `docs/sinks.md` - Add Kiro sink configuration examples
- `docs/commands.md` - Update examples to include Kiro

**Example Documentation Additions**:

```markdown
### Kiro CLI

```bash
# Configure sinks
arm add sink --tool kiro kiro-steering .kiro/steering
arm add sink --tool kiro kiro-prompts .kiro/prompts

# Install resources
arm install ruleset my-registry/clean-code kiro-steering
arm install promptset my-registry/code-review kiro-prompts
```

Kiro steering files provide persistent project knowledge in markdown format.
```

## Acceptance Criteria

- [ ] Kiro tool constant added to `internal/arm/compiler/types.go`
- [ ] Factory maps Kiro, AmazonQ, and Markdown to same generators
- [ ] `internal/arm/compiler/amazonq.go` deleted (duplicate of markdown.go)
- [ ] `internal/arm/compiler/amazonq_test.go` deleted (covered by markdown tests)
- [ ] `arm add sink --tool kiro` command works
- [ ] `arm add sink --tool amazonq` still works (mapped to markdown)
- [ ] `arm install ruleset ... kiro-steering` installs to `.kiro/steering/`
- [ ] `arm install promptset ... kiro-prompts` installs to `.kiro/prompts/`
- [ ] Priority index (`arm_index.md`) generated for rulesets
- [ ] Hierarchical layout used (arm/{registry}/{package}/{version}/)
- [ ] Uninstall removes files and cleans up empty directories
- [ ] arm-index.json tracks installed packages per sink
- [ ] All existing tests pass (including Amazon Q E2E tests)
- [ ] E2E tests added for Kiro compilation
- [ ] Documentation updated with Kiro examples

## Implementation Notes

### Refactor to Eliminate Duplication

Amazon Q, Markdown, and Kiro all produce **identical output**:
- All three use `.md` extension
- All three use `GenerateRuleMetadata()` for rules (YAML frontmatter)
- All three use plain body for prompts
- All three use same filename pattern: `{id}_{id}.md`

**Current state**: `amazonq.go` is a duplicate of `markdown.go` (copy-paste code)

**Refactored state**: 
- Delete `amazonq.go` and `amazonq_test.go`
- Map `AmazonQ`, `Kiro`, and `Markdown` tool constants to the same generators in factory
- Reduces code duplication and maintenance burden

### Reuse Existing Code
- Reuse filename generation logic from existing compilers
- Reuse priority index generation logic (same markdown format)
- Follow existing patterns for tool registration in factory

### Testing Strategy

Add E2E tests in `test/e2e/compile_test.go`:
- `TestKiroRulesetCompilation` - Verify ruleset compilation to .kiro/steering/
- `TestKiroPromptsetCompilation` - Verify promptset compilation to .kiro/prompts/
- `TestKiroPriorityIndex` - Verify arm_index.md generation
- `TestKiroMultiplePriorities` - Verify priority ordering

## Related Specifications

- `specs/sink-compilation.md` - Sink management and compilation
- `specs/standalone-compilation.md` - Local file compilation
- `specs/priority-resolution.md` - Priority-based conflict resolution

## References

- Kiro CLI Steering Documentation: `KIRO-CLI/steering.md`
- Kiro CLI Migration Guide: `KIRO-CLI/upgrading-from-amazon-q.md`
- Kiro uses `.kiro/steering/` for workspace scope, `~/.kiro/steering/` for global scope
- Kiro uses `.kiro/prompts/` for workspace scope, `~/.kiro/prompts/` for global scope

---

## Implementation Summary

**Completed in commits:**
- `1d2eccc` - feat: add Kiro tool support and refactor Amazon Q compiler
- `4a15543` - test: add E2E tests for Kiro tool compilation
- `be94c0a` - docs: add Kiro CLI support to documentation

**What was implemented:**
1. ✅ Added `Kiro` tool constant to `internal/arm/compiler/types.go`
2. ✅ Refactored factory to map Kiro, AmazonQ, and Markdown to same generators
3. ✅ Deleted `internal/arm/compiler/amazonq.go` (duplicate of markdown.go)
4. ✅ Deleted `internal/arm/compiler/amazonq_test.go` (covered by markdown tests)
5. ✅ Updated `cmd/arm/main.go` to handle "kiro" tool in sink commands
6. ✅ Updated `internal/arm/service/service.go` to map "kiro" string to compiler.Kiro
7. ✅ Added E2E tests: `TestCompilationToolFormats/KiroFormat` and `TestCompilationPromptsets/KiroPrompts`
8. ✅ Updated documentation: README.md, docs/commands.md, docs/sinks.md, docs/concepts.md

**Verification:**
- All tests pass: `go test ./...`
- E2E tests verify:
  - Ruleset installation to `.kiro/steering/`
  - Promptset installation to `.kiro/prompts/`
  - Hierarchical layout (arm/{registry}/{package}/{version}/)
  - Priority index generation (`arm_index.md`)
  - Markdown format (`.md` extension)

**Branch:** `kiro-cli-support` (ready to merge to main)
