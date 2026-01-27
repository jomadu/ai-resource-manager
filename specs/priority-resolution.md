# Priority Resolution

## Job to be Done
Resolve conflicts when multiple rulesets define overlapping rules by applying priority-based ordering with clear precedence rules.

## Activities
1. Assign priority to rulesets (default: 100)
2. Generate priority index files for AI tools
3. Resolve conflicts by priority (higher priority wins)
4. Embed priority metadata in compiled rules

## Acceptance Criteria
- [x] Support --priority flag on install (default: 100)
- [x] Store priority in manifest
- [x] Generate arm_index.* file listing rulesets by priority
- [x] Embed priority in compiled rule metadata
- [x] Higher priority rules override lower priority rules
- [x] Same priority uses installation order (later wins)
- [x] Priority index updated on install/uninstall
- [x] Priority index removed when no rulesets remain

## Data Structures

### Ruleset with Priority
```json
{
  "dependencies": {
    "ai-rules/team-standards": {
      "type": "ruleset",
      "version": "^1.0.0",
      "priority": 200,
      "sinks": ["cursor-rules"]
    },
    "ai-rules/clean-code": {
      "type": "ruleset",
      "version": "^1.0.0",
      "priority": 100,
      "sinks": ["cursor-rules"]
    }
  }
}
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
2. clean-code (priority: 100)
3. security (priority: 100)

Rules from higher priority rulesets override rules from lower priority rulesets.
```

## Algorithm

### Assign Priority
1. Parse --priority flag (default: 100)
2. Store in manifest dependency config
3. Pass to sink manager on install

### Generate Priority Index
1. Collect all rulesets in sink
2. Extract priority from manifest
3. Sort by priority (descending)
4. For same priority, sort by installation order
5. Generate tool-specific index file:
   - Cursor: arm_index.mdc with frontmatter
   - Amazon Q: arm_index.md
   - Copilot: arm_index.instructions.md
   - Markdown: arm_index.md
6. Set index priority to 999 (highest)
7. Write to sink directory

### Embed Priority Metadata
1. Extract priority from manifest
2. Add to rule frontmatter/metadata:
   - Cursor: YAML frontmatter
   - Amazon Q: Markdown comment
   - Copilot: Frontmatter
   - Markdown: Frontmatter
3. Include enforcement and scope

### Resolve Conflicts
1. AI tool reads arm_index.* first (priority 999)
2. AI tool applies rules in priority order
3. Higher priority rules override lower priority
4. Same priority uses later installation

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| No priority specified | Default to 100 |
| Negative priority | Allowed (lower than default) |
| Priority > 999 | Allowed (higher than index) |
| Same priority | Use installation order (later wins) |
| Priority index missing | Generate on next install |
| No rulesets in sink | Remove priority index |
| Promptsets in sink | Don't affect priority index (rules only) |

## Dependencies

- Manifest management
- Sink compilation (sink-compilation.md)
- Rule metadata generation

## Implementation Mapping

**Source files:**
- `internal/arm/compiler/generators.go` - GenerateRuleMetadata (embeds priority)
- `internal/arm/sink/manager.go` - generateRulesetIndexRuleFile, InstallRuleset
- `internal/arm/service/service.go` - InstallRuleset (passes priority)
- `cmd/arm/main.go` - handleInstallRuleset (parses --priority flag)
- `test/e2e/install_test.go` - TestRulesetInstallationWithPriority
- `test/e2e/compile_test.go` - TestCompilationMultiplePriorities

## Examples

### Install with Priority
```bash
# High priority team standards
arm install ruleset --priority 200 ai-rules/team-standards cursor-rules

# Default priority (100)
arm install ruleset ai-rules/clean-code cursor-rules

# Low priority experimental rules
arm install ruleset --priority 50 ai-rules/experimental cursor-rules
```

### Priority Index (Cursor)
```markdown
---
priority: 999
enforcement: required
scope: all
---

# ARM Priority Index

When multiple rules conflict, apply them in this priority order:

1. team-standards (priority: 200)
2. clean-code (priority: 100)
3. security (priority: 100)
4. experimental (priority: 50)

Rules from higher priority rulesets override rules from lower priority rulesets.
```

### Rule with Priority Metadata (Cursor)
```markdown
---
priority: 200
enforcement: required
scope: all
---

# Team Coding Standards

Always use TypeScript strict mode.
```

### Conflict Resolution Example
```
Scenario:
- team-standards (priority: 200) says: "Use 2 spaces for indentation"
- clean-code (priority: 100) says: "Use 4 spaces for indentation"

Resolution:
- AI applies team-standards rule (higher priority)
- Result: Use 2 spaces for indentation
```

### Same Priority Resolution
```
Scenario:
- clean-code (priority: 100, installed first) says: "Use camelCase"
- security (priority: 100, installed second) says: "Use snake_case"

Resolution:
- AI applies security rule (later installation)
- Result: Use snake_case
```

### Priority Index Update
```bash
# Initial state: team-standards (200), clean-code (100)
arm install ruleset --priority 150 ai-rules/security cursor-rules

# Priority index updated:
# 1. team-standards (priority: 200)
# 2. security (priority: 150)
# 3. clean-code (priority: 100)

arm uninstall ai-rules/security

# Priority index updated:
# 1. team-standards (priority: 200)
# 2. clean-code (priority: 100)

arm uninstall ai-rules/team-standards
arm uninstall ai-rules/clean-code

# Priority index removed (no rulesets remain)
```
