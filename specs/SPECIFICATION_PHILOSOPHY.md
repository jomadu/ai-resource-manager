# Specification Philosophy

This document captures our understanding of what makes effective builder-oriented specifications, synthesized from studying the CURSED language project and the Ralph methodology.

## What Are Specifications?

**Specifications define WHAT to build, not HOW to use it.**

- **User Documentation** = How to use the system (commands, workflows, examples)
- **Specifications** = How to build the system (algorithms, data structures, acceptance criteria)

## Key Principles from CURSED

### 1. Technical Depth Over User Guidance

CURSED specs provide implementation-ready technical details:

```
❌ User Doc: "CURSED supports imports using the yeet keyword"
✅ Spec: "ImportDecl = 'yeet' ( ImportSpec | '(' { ImportSpec ';' } ')' | ImportList )"
```

**What this means:**
- Grammar rules in BNF notation
- Algorithm pseudocode with step-by-step logic
- Memory layouts with byte-level precision
- State machines with transitions
- Data structure schemas

### 2. Explicit Requirements Markers

CURSED uses clear markers to indicate implementation requirements:

- **CANONICAL** - Required behavior, must implement exactly as specified
- **IMPLEMENTATION REQUIREMENT** - Parser/compiler must support this
- **DEPRECATED** - Do not implement, will be removed
- **PARSING RULES** - Specific parsing behavior required

**Example:**
```markdown
**CANONICAL SYNTAX:** Imports declare dependencies using the `yeet` keyword.

**IMPLEMENTATION REQUIREMENT:** All parsers MUST support all four import forms.

**DEPRECATED:** Go-style `Result<T, Error>` patterns are deprecated.
```

### 3. One Spec Per Concern

Each specification file focuses on a single technical concern:

- `lexical.md` - Token structure, keywords, operators
- `grammar.md` - Language grammar rules
- `types.md` - Type system, conversions, zero values
- `memory_management.md` - GC algorithms, heap structure
- `error_handling.md` - Error types, propagation, recovery
- `compiler_stages.md` - Bootstrap process, deliverables

**Not:** "language_overview.md" covering everything

### 4. Testable Specifications

Every spec should enable test derivation:

```markdown
## Zero Values

| Type | Zero Value |
|------|------------|
| lit  | `cringe` (false) |
| numeric types | `0` |
| tea  | `""` (empty string) |
| pointers | `nah` (nil) |
```

**Derived tests:**
- Assert uninitialized `lit` equals `cringe`
- Assert uninitialized `normie` equals `0`
- Assert uninitialized `tea` equals `""`
- Assert uninitialized pointer equals `nah`

### 5. Edge Cases and Failure Modes

Specs document boundary conditions and error handling:

```markdown
**Edge cases:**
- No matching versions → error "no versions satisfy constraint"
- Branch doesn't exist → error "branch not found"
- Malformed semver tag → skip, log warning
- Network failure during resolution → retry 3x with backoff
```

### 6. Visual Diagrams for Complex Systems

Use ASCII diagrams for architecture and data flow:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Runtime       │    │ Memory Manager  │    │ Garbage         │
│   Allocator     │◄──►│                 │◄──►│ Collector       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Key Principles from Ralph Methodology

### 1. JTBD → Activities → Acceptance Criteria

Specs should follow this structure:

```markdown
# [Activity Name]

## Job to be Done
[User outcome this enables]

## Activities
[Discrete steps/operations]

## Acceptance Criteria
- [ ] Observable outcome 1
- [ ] Observable outcome 2
- [ ] Edge case handled

## Data Structures
[Schemas, types, state]

## Algorithm
[Step-by-step logic]

## Edge Cases
[Boundary conditions, errors]

## Dependencies
[What must exist first]

## Examples
[Concrete scenarios]
```

### 2. Observable, Verifiable Outcomes

Acceptance criteria must be:
- **Observable** - Can be seen/measured
- **Verifiable** - Can be tested programmatically
- **Binary** - Pass or fail, no ambiguity

```
✅ "Version resolution returns highest semver tag when multiple exist"
✅ "Installing to new sinks removes files from old sinks"
✅ "Cache key includes registry URL + package patterns"

❌ "Version resolution works correctly"
❌ "Installation is clean"
❌ "Cache is efficient"
```

### 3. Specs Define WHAT, Ralph Decides HOW

**Spec responsibility:**
- What outcomes must be achieved
- What constraints must be satisfied
- What edge cases must be handled
- What data structures are required

**Ralph's responsibility:**
- How to implement the algorithm
- What libraries to use
- How to structure the code
- What patterns to follow

### 4. Disposable and Regenerable

Specs should be:
- **Disposable** - If wrong, throw out and regenerate
- **Iterative** - Evolve through observed failures
- **Focused** - One activity per spec
- **Testable** - Enable TDD workflow

## Anti-Patterns to Avoid

### ❌ User Documentation Disguised as Specs

```markdown
# Commands

## arm install

Install a ruleset from a registry to a sink.

**Example:**
```bash
arm install ruleset my-org/clean-code-ruleset cursor-rules
```
```

**Why it's wrong:** Tells users how to use, not builders how to implement.

### ❌ Vague Acceptance Criteria

```markdown
- [ ] Installation works correctly
- [ ] Version resolution is accurate
- [ ] Cache is performant
```

**Why it's wrong:** Not observable, not verifiable, not testable.

### ❌ Implementation Details in Specs

```markdown
Use the `semver` crate to parse version strings.
Store cache in `~/.arm/storage` using SHA256 hashes.
Implement version resolution with a binary search tree.
```

**Why it's wrong:** Prescribes HOW, not WHAT. Ralph should decide implementation.

### ❌ Multiple Concerns in One Spec

```markdown
# ARM System

This spec covers registries, sinks, packages, version resolution,
compilation, caching, and authentication.
```

**Why it's wrong:** Too broad, not focused, hard to test, hard to regenerate.

## Specification Template

```markdown
# [Activity Name]

## Job to be Done
[What user outcome does this enable?]

## Activities
[What discrete operations accomplish this JTBD?]

## Acceptance Criteria
- [ ] [Observable outcome 1]
- [ ] [Observable outcome 2]
- [ ] [Edge case 1 handled]
- [ ] [Edge case 2 handled]

## Data Structures

### [Structure Name]
```
{
  "field": "type",
  "description": "purpose"
}
```

## Algorithm

1. [Step 1 with clear input/output]
2. [Step 2 with decision points]
3. [Step 3 with error handling]

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| [Edge case 1] | [How system responds] |
| [Edge case 2] | [How system responds] |

## Dependencies

- [Prerequisite 1]
- [Prerequisite 2]

## Examples

### Example 1: [Scenario Name]

**Input:**
```
[Input data]
```

**Expected Output:**
```
[Output data]
```

**Verification:**
- [How to verify outcome 1]
- [How to verify outcome 2]
```

## Application to ARM

### Current State
- `specs/*.md` are user documentation
- Missing: algorithms, data structures, acceptance criteria
- Missing: edge cases, failure modes, invariants
- Missing: testable specifications

### Target State
- `docs/*.md` for user documentation
- `specs/*.md` for builder specifications
- One spec per activity (version resolution, package installation, etc.)
- Each spec follows template with JTBD → acceptance criteria → algorithm
- Specs enable test-driven development
- Specs are disposable and regenerable

### Migration Path
1. Move current `specs/*.md` → `docs/*.md`
2. Create builder specs in `specs/` following template
3. Derive tests from acceptance criteria
4. Use specs to guide Ralph implementation
5. Iterate specs based on implementation learnings

## References

- **CURSED Project**: `/tmp/cursed/specs/` - Example of builder-oriented specs
- **Ralph Playbook**: `/tmp/how-to-ralph-wiggum/` - Methodology for spec-driven development
- **ARM Implementation Plan**: `./IMPLEMENTATION_PLAN.md` - Migration plan for ARM specs
