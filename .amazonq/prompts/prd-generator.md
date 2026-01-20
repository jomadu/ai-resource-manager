Generate a Product Requirements Document (PRD) for a new feature.

## Process

1. Ask 3-5 clarifying questions with lettered options (A/B/C/D)
2. Generate structured PRD based on answers
3. Save to `specs/[feature-name]/spec.md`

## Questions Format

```
1. What is the primary goal?
   A. Option 1
   B. Option 2
   C. Option 3
   D. Other: [specify]
```

User responds: "1A, 2C, 3B"

## PRD Structure

### 1. Introduction
Brief description and problem statement.

### 2. Goals
Specific, measurable objectives.

### 3. User Stories
```
### US-001: [Title]
**Description:** As a [user], I want [feature] so that [benefit].

**Acceptance Criteria:**
- [ ] Specific verifiable criterion
- [ ] Typecheck passes
- [ ] Tests pass
```

Stories must be small (completable in one session).

### 4. Functional Requirements
- FR-1: System must...
- FR-2: When user does X...

### 5. Non-Goals
What this will NOT include.

### 6. Technical Considerations (Optional)
Constraints, dependencies, performance.

### 7. Success Metrics
How to measure success.

## Guidelines

- Be explicit and unambiguous
- Verifiable acceptance criteria (not "works correctly")
- Small stories (2-3 sentences each)
- Number all requirements
- Use concrete examples

## Output

Save to: `specs/[feature-name]/spec.md` (kebab-case)
