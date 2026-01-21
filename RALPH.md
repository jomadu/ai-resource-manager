# RALPH Agent Protocol

You are an autonomous coding agent. Each iteration completes ONE user story from `prd.json`.

## CRITICAL: Execution Loop

Each iteration MUST follow this exact sequence:

### 1. SETUP
- [ ] Read `prd.json`
- [ ] Read `progress.txt` (check Codebase Patterns section first)
- [ ] Verify branch matches `prd.json` branchName (checkout/create if needed)
- [ ] Count stories where `passes: false`
- [ ] If count = 0, exit with `<promise>COMPLETE</promise>`

### 2. SELECT
- [ ] Pick highest priority story where `passes: false`
- [ ] Append to progress.txt: story ID and plan (use format below)

### 3. IMPLEMENT
- [ ] Work in logical chunks (e.g., types → functions → tests)
- [ ] After each chunk: append progress to progress.txt
- [ ] Continue until story complete

### 4. VERIFY
- [ ] Run quality checks (typecheck, lint, test)
- [ ] If fail 3 times: exit with `<promise>BLOCKED: [reason]</promise>`

### 5. COMMIT
- [ ] Commit with: `feat: [Story ID] - [Story Title]`
- [ ] Update prd.json: set `passes: true` for completed story
- [ ] Append completion report to progress.txt with learnings
- [ ] Update AGENTS.md if reusable patterns discovered

### 6. EXIT
- [ ] Re-read prd.json from disk
- [ ] Count stories where `passes: false`
- [ ] State count explicitly in response
- [ ] Output promise tag on own line at END:
  - `<promise>CONTINUE</promise>` if count > 0
  - `<promise>COMPLETE</promise>` if count = 0

---

## REQUIRED: File Formats

### progress.txt (ALWAYS APPEND, NEVER REPLACE)

**Story Start:**
```
## [Date/Time] - Starting [Story ID]
- Story: [Story Title]
- Plan: [Brief approach]
---
```

**Progress Chunk:**
```
### [Time] - Progress on [Story ID]
- Completed: [What finished]
- Next: [What's next]
```

**Story Complete:**
```
## [Date/Time] - Completed [Story ID]
- What was implemented
- Files changed
- **Learnings for future iterations:**
  - [Pattern/gotcha/context]
---
```

**Codebase Patterns (at TOP of file):**
```
## Codebase Patterns
- [General reusable pattern]
```

### Commit Message Format
```
feat: [Story ID] - [Story Title]
```

---

## GUIDANCE: Best Practices

### Why Each Step Matters

**Setup checks prevent:**
- Working on wrong branch
- Missing context from previous iterations
- Duplicate work

**Progress logging enables:**
- Debugging when things go wrong
- Understanding what was tried
- Building institutional knowledge

**Chunk-based work allows:**
- Incremental progress tracking
- Easier debugging
- Clear audit trail

**Quality gates ensure:**
- No broken code committed
- Bounded retry attempts
- Clean CI/CD pipeline

**Exit verification prevents:**
- Stale state bugs
- Infinite loops
- Premature completion

### Knowledge Management

**progress.txt Codebase Patterns:**
- Add general, reusable patterns to TOP of file
- Examples: "Use X for Y", "Always do Z when W"
- Keep it scannable for future iterations

**AGENTS.md Updates:**
- Add module-specific knowledge only
- Update when you discover non-obvious patterns
- Examples:
  - "Update Y when changing X"
  - "This module uses pattern Z for all API calls"
  - "Tests need PORT 3000"
  - "Field names must match template exactly"
- Do NOT add:
  - Story-specific implementation details
  - Temporary debugging notes
  - Information already in progress.txt

### Implementation Tips

- Follow existing code patterns
- Keep changes minimal and focused
- Commit frequently within the story
- Read Codebase Patterns before starting

---

## CONSTRAINTS

- ONE story per iteration (never continue automatically)
- ALL commits must pass quality checks
- NEVER commit broken code
- NEVER commit partial or incomplete implementations
- Story must be FULLY complete before committing
- ALWAYS append to progress.txt (never replace)
- ALWAYS re-read prd.json before exit
- Promise tag MUST be on own line at END of response
