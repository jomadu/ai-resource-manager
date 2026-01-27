# Priority Resolution

## Job to be Done
Resolve conflicts when multiple rulesets define overlapping rules by applying priority-based ordering with clear precedence rules.

## Activities
1. Assign priority to rulesets (default: 100)
2. Generate priority index files for AI tools
3. Resolve conflicts by priority (higher priority wins)

## Acceptance Criteria
- [x] Support --priority flag on install (default: 100)
- [x] Store priority in manifest and arm-index.json
- [x] Generate arm_index.* file listing rulesets by priority
- [x] Higher priority rulesets override lower priority rulesets
- [x] Same priority uses installation order (later wins)
- [x] Priority index updated on install/uninstall
- [x] Priority index removed when no rulesets remain

## Data Structures

### Ruleset with Priority (arm.json)
```json
{
  "dependencies": {
    "ai-rules/team-standards": {
      "type": "ruleset",
      "version": "^1.0.0",
      "priority": 200,
      "sinks": ["cursor-rules"]
    }
  }
}
```

### Ruleset Index Entry (arm-index.json)
```json
{
  "version": 1,
  "rulesets": {
    "ai-rules/team-standards@1.0.0": {
      "priority": 200,
      "files": ["arm/ai-rules/team-standards/1.0.0/rules/rule.mdc"]
    }
  }
}
```

### Priority Index File (arm_index.*)
```markdown
# ARM Rulesets

This file defines the installation priorities for rulesets managed by ARM.

## Priority Rules

**This index is the authoritative source of truth for ruleset priorities.** When conflicts arise between rulesets, follow this priority order:

1. **Higher priority numbers take precedence** over lower priority numbers
2. **Rules from higher priority rulesets override** conflicting rules from lower priority rulesets
3. **Always consult this index** to resolve any ambiguity about which rules to follow

## Installed Rulesets

### ai-rules/team-standards@1.0.0
- **Priority:** 200
- **Rules:**
  - arm/ai-rules/team-standards/1.0.0/rules/rule.mdc

### ai-rules/clean-code@1.0.0
- **Priority:** 100
- **Rules:**
  - arm/ai-rules/clean-code/1.0.0/rules/rule.mdc
```

**Note:** RULESET priority is a deployment concern, not a compilation concern. It appears in manifest files and the priority index, but NOT in individual compiled rule files.

## Algorithm

### Assign Priority
1. Parse --priority flag (default: 100)
2. Store in manifest dependency config
3. Store in arm-index.json per sink on install

### Generate Priority Index
1. Load arm-index.json from sink
2. Collect all rulesets with priorities
3. Sort by priority (descending)
4. For same priority, sort by installation order
5. Generate markdown file listing rulesets
6. Write to sink directory as arm_index.* (tool-specific extension)

### Resolve Conflicts
1. AI tool reads arm_index.* to understand ruleset priorities
2. When rules conflict, apply rule from higher priority ruleset
3. Same priority uses later installation (last wins)

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

- Manifest management (arm.json)
- Sink index management (arm-index.json)
- Sink compilation (sink-compilation.md)

## Implementation Mapping

**Source files:**
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
# ARM Rulesets

This file defines the installation priorities for rulesets managed by ARM.

## Priority Rules

**This index is the authoritative source of truth for ruleset priorities.** When conflicts arise between rulesets, follow this priority order:

1. **Higher priority numbers take precedence** over lower priority numbers
2. **Rules from higher priority rulesets override** conflicting rules from lower priority rulesets
3. **Always consult this index** to resolve any ambiguity about which rules to follow

## Installed Rulesets

### ai-rules/team-standards@1.0.0
- **Priority:** 200
- **Rules:**
  - arm/ai-rules/team-standards/1.0.0/rules/teamStandards_indentation.mdc

### ai-rules/clean-code@1.0.0
- **Priority:** 100
- **Rules:**
  - arm/ai-rules/clean-code/1.0.0/rules/cleanCode_formatting.mdc

### ai-rules/security@1.0.0
- **Priority:** 100
- **Rules:**
  - arm/ai-rules/security/1.0.0/rules/security_auth.mdc

### ai-rules/experimental@1.0.0
- **Priority:** 50
- **Rules:**
  - arm/ai-rules/experimental/1.0.0/rules/experimental_feature.mdc
```

### Individual Rule File (Cursor)
```markdown
---
description: "Enforce consistent indentation"
globs: **/*.ts, **/*.js
alwaysApply: true
---

---
namespace: ai-rules/team-standards@1.0.0
ruleset:
  id: team-standards
  name: Team Standards
  rules:
    - indentation
    - naming
rule:
  id: indentation
  name: Indentation Rule
  enforcement: MUST
---

# Team Coding Standards

Always use 2 spaces for indentation.
```

**Note:** Individual rules do NOT contain RULESET priority. Priority is only in arm_index.* and arm-index.json.

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
