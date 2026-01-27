# Sink Compilation

## Job to be Done
Compile ARM resource files (YAML) to tool-specific formats and manage output directories (sinks) with proper cleanup.

## Activities
1. Compile rulesets to tool-specific rule formats
2. Compile promptsets to tool-specific prompt formats
3. Generate priority index files for conflict resolution
4. Manage sink directory structure (hierarchical or flat)
5. Clean up empty directories on uninstall

## Acceptance Criteria
- [x] Compile to Cursor format (.mdc with frontmatter for rules, .md for prompts)
- [x] Compile to Amazon Q format (.md for both rules and prompts)
- [x] Compile to Copilot format (.instructions.md for rules)
- [x] Compile to Markdown format (.md for both)
- [x] Generate arm_index.* priority files for rulesets
- [x] Support hierarchical layout (preserves directory structure)
- [x] Support flat layout (single directory with hash prefixes)
- [x] Clean up empty directories recursively on uninstall
- [x] Remove arm-index.json when all packages uninstalled
- [x] Remove arm_index.* when all rulesets uninstalled

## Data Structures

### Sink Configuration
```json
{
  "sinks": {
    "cursor-rules": {
      "directory": ".cursor/rules",
      "tool": "cursor"
    }
  }
}
```

### ARM Index (arm-index.json)
```json
{
  "version": 1,
  "rulesets": {
    "registry/package@version": {
      "priority": 100,
      "files": [
        "arm/registry/package/version/rules/rule.mdc"
      ]
    }
  },
  "promptsets": {
    "registry/package@version": {
      "files": [
        "arm/registry/package/version/prompts/prompt.md"
      ]
    }
  }
}
```

**Note:** Tracks installed packages per sink with file paths and priority (for rulesets)

## Algorithm

### Compile Ruleset
1. Parse ARM resource YAML
2. For each rule in ruleset:
   - Generate filename using tool-specific generator
   - Generate content using tool-specific compiler
   - Embed metadata (priority, enforcement, scope)
3. Write files to sink directory
4. Update arm-index.json
5. Generate arm_index.* priority file

### Compile Promptset
1. Parse ARM resource YAML
2. For each prompt in promptset:
   - Generate filename using tool-specific generator
   - Generate content using tool-specific compiler
3. Write files to sink directory
4. Update arm-index.json

### Generate Priority Index
1. Collect all rulesets in sink
2. Sort by priority (highest first)
3. Generate tool-specific index file:
   - Cursor: arm_index.mdc with frontmatter
   - Amazon Q: arm_index.md
   - Copilot: arm_index.instructions.md
   - Markdown: arm_index.md
4. Include instructions for AI to respect priority order

### Cleanup on Uninstall
1. Remove package files from sink
2. Remove package from arm-index.json
3. If no packages remain, remove arm-index.json
4. If no rulesets remain, remove arm_index.*
5. Recursively remove empty directories:
   - Start from deepest directories
   - Remove if empty
   - Never remove sink root directory
   - Continue until no more empty directories

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Duplicate rule IDs | Later rules override earlier ones |
| Missing tool compiler | Error with supported tools list |
| Filename collision | Use hash prefix to prevent collision |
| Filename > 100 chars | Truncate progressively (path → filename → truncated) |
| Empty directory after uninstall | Remove directory and empty ancestors |
| Nested empty directories | Remove all empty ancestors recursively |
| Sink root empty | Keep sink root, only remove subdirs |
| Non-empty directory | Keep directory and all ancestors |

## Dependencies

- ARM resource parsing (parser)
- Sink configuration (manifest)
- Package tracking (arm-index.json)

## Implementation Mapping

**Source files:**
- `internal/arm/compiler/compiler.go` - CompileRuleset, CompilePromptset
- `internal/arm/compiler/cursor.go` - Cursor-specific compilation
- `internal/arm/compiler/amazonq.go` - Amazon Q compilation
- `internal/arm/compiler/copilot.go` - Copilot compilation
- `internal/arm/compiler/markdown.go` - Markdown compilation
- `internal/arm/compiler/generators.go` - GenerateRuleMetadata
- `internal/arm/sink/manager.go` - InstallRuleset, InstallPromptset, Uninstall, CleanupEmptyDirectories
- `test/e2e/compile_test.go` - E2E compilation tests
- `test/e2e/cleanup_test.go` - E2E cleanup tests

## Examples

### Hierarchical Layout (Cursor)
```
.cursor/rules/
└── arm/
    └── ai-rules/
        └── clean-code-ruleset/
            └── 1.0.0/
                └── rules/
                    ├── cleanCode_ruleOne.mdc
                    ├── cleanCode_ruleTwo.mdc
                    └── cleanCode_ruleThree.mdc
```

### Flat Layout (Copilot)
```
.github/instructions/
├── arm_1a2b_3c4d_rules_cleanCode_ruleOne.instructions.md
├── arm_1a2b_5e6f_rules_cleanCode_ruleTwo.instructions.md
└── arm_index.instructions.md
```

### Cursor Rule with Frontmatter
```markdown
---
priority: 100
enforcement: required
scope: all
---

# Rule Title

Rule body content here.
```

### Amazon Q Rule (Pure Markdown)
```markdown
# Rule Title

Rule body content here.
```

### Priority Index (arm_index.mdc)
```markdown
---
priority: 999
enforcement: required
scope: all
---

# ARM Priority Index

When multiple rules conflict, apply them in this priority order:

1. team-standards (priority: 200)
2. clean-code-ruleset (priority: 100)
3. security-ruleset (priority: 100)
```

### Cleanup Example
```bash
# Before uninstall
.cursor/rules/arm/ai-rules/clean-code/1.0.0/rules/rule.mdc

# After uninstall
# All empty directories removed:
# - .cursor/rules/arm/ai-rules/clean-code/1.0.0/rules/
# - .cursor/rules/arm/ai-rules/clean-code/1.0.0/
# - .cursor/rules/arm/ai-rules/clean-code/
# - .cursor/rules/arm/ai-rules/
# - .cursor/rules/arm/
# But .cursor/rules/ (sink root) is kept
```
