# [Activity Name]

## Job to be Done
[What user outcome does this enable? Focus on the value delivered, not the mechanism.]

## Activities
[What discrete operations accomplish this JTBD? List the key steps or sub-activities.]

## Acceptance Criteria
- [ ] [Observable outcome 1 - must be verifiable through testing]
- [ ] [Observable outcome 2 - must be verifiable through testing]
- [ ] [Edge case 1 handled - specific boundary condition]
- [ ] [Edge case 2 handled - specific error scenario]

## Data Structures

### [Structure Name]
```json
{
  "field": "type",
  "description": "purpose"
}
```

**Fields:**
- `field` - Description of field purpose and constraints

## Algorithm

1. [Step 1 with clear input/output]
2. [Step 2 with decision points]
3. [Step 3 with error handling]

**Pseudocode:**
```
function ActivityName(input):
    if condition:
        return result
    else:
        handle_error()
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| [Edge case 1] | [How system responds] |
| [Edge case 2] | [How system responds] |
| [Edge case 3] | [How system responds] |

## Dependencies

- [Prerequisite 1 - what must exist before this can work]
- [Prerequisite 2 - what other components are required]

## Implementation Mapping

**Source files:**
- `path/to/file.go` - [Brief description of what this file implements]

**Related specs:**
- `other-spec.md` - [How this spec relates to others]

## Examples

### Example 1: [Scenario Name]

**Input:**
```
[Input data or command]
```

**Expected Output:**
```
[Output data or result]
```

**Verification:**
- [How to verify outcome 1]
- [How to verify outcome 2]

### Example 2: [Error Scenario Name]

**Input:**
```
[Input that triggers error]
```

**Expected Output:**
```
[Error message or behavior]
```

**Verification:**
- [How to verify error handling]

## Notes

[Any additional context, design decisions, or rationale that helps builders understand the specification.]

## Known Issues

[List any known bugs or limitations discovered during deep-dive analysis, with references to source files if applicable.]

## Areas for Improvement

[List any known limitations or potential enhancements for future consideration.]