Convert a PRD (markdown) to Ralph's prd.json format for autonomous execution.

## Input

Markdown PRD from `specs/[feature-name]/spec.md`

## Output Format

```json
{
  "project": "AI Resource Manager",
  "branchName": "ralph/[feature-name-kebab-case]",
  "description": "[Feature description]",
  "userStories": [
    {
      "id": "US-001",
      "title": "[Story title]",
      "description": "As a [user], I want [feature] so that [benefit]",
      "acceptanceCriteria": [
        "Specific criterion",
        "Typecheck passes",
        "Tests pass"
      ],
      "priority": 1,
      "passes": false,
      "notes": ""
    }
  ]
}
```

Save to: `specs/[feature-name]/prd.json`

## Critical Rules

### Story Size
Each story must complete in ONE context window.

**Right-sized:**
- Add database field + migration
- Add CLI command
- Update error handling in module

**Too big (split):**
- "Build entire feature"
- "Refactor codebase"
- "Add comprehensive testing"

### Story Ordering
Order by dependencies:
1. Data structures/types
2. Core logic/functions
3. CLI commands/interfaces
4. Integration features

### Acceptance Criteria
Must be verifiable, not vague.

**Good:**
- "Add `status` field to Version struct"
- "Function returns error for invalid input"
- "Typecheck passes"

**Bad:**
- "Works correctly"
- "Good error handling"

### Always Include
- "Typecheck passes" in every story
- "Tests pass" for testable logic

## Conversion Steps

1. Extract user stories from PRD
2. Split large stories into small ones
3. Order by dependencies
4. Assign sequential IDs (US-001, US-002...)
5. Set all `passes: false`
6. Add required criteria (typecheck, tests)
7. Generate branch name: `ralph/[feature-kebab-case]`

## Example

**Input PRD:**
```markdown
# Add Status Field

Add status tracking to versions.

## Requirements
- Store status in struct
- Validate transitions
- Return errors for invalid states
```

**Output prd.json:**
```json
{
  "project": "AI Resource Manager",
  "branchName": "ralph/add-status-field",
  "description": "Add status tracking to versions",
  "userStories": [
    {
      "id": "US-001",
      "title": "Add status field to Version struct",
      "description": "As a developer, I need to store version status.",
      "acceptanceCriteria": [
        "Add Status field to Version struct",
        "Update NewVersion constructor",
        "Typecheck passes"
      ],
      "priority": 1,
      "passes": false,
      "notes": ""
    },
    {
      "id": "US-002",
      "title": "Add status validation",
      "description": "As a user, I want invalid status transitions prevented.",
      "acceptanceCriteria": [
        "Add ValidateStatusTransition function",
        "Return error for invalid transitions",
        "Typecheck passes",
        "Tests pass"
      ],
      "priority": 2,
      "passes": false,
      "notes": ""
    }
  ]
}
```

## Checklist

- [ ] Each story completable in one iteration
- [ ] Stories ordered by dependencies
- [ ] All have "Typecheck passes"
- [ ] Testable stories have "Tests pass"
- [ ] Criteria are verifiable
- [ ] No story depends on later story
