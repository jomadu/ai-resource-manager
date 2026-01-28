# ORIENT

Study specifications in `specs/*` with up to 250 parallel subagents.
Study @IMPLEMENTATION_PLAN.md (if present) to understand the plan so far.
Study implementation in `internal/*`, `cmd/*`, `test/*`, and `docs/*` with up to 250 parallel subagents.

Specifications drive the entire repo: README, CI/CD, configs, code, tests, and docs.

Implementation is in: `internal/*`, `cmd/*`, `test/*`, `docs/*`, root files, and infrastructure (`.github/*`, etc.)
Specification is in: `specs/*`

Log: [ORIENT] Starting orientation...
Log when done: [ORIENT] ✓ Complete

# TASK

Study @IMPLEMENTATION_PLAN.md (if present; assume it is inaccurate and incomplete).

Use up to 500 parallel subagents to:
- Study existing implementation in `internal/*`, `cmd/*`, `test/*`, and `docs/*`
- Compare it against specifications in `specs/*`
- Search for TODO, minimal implementations, placeholders, skipped/flaky tests, inconsistent patterns

Use a subagent to:
- Analyze findings
- Prioritize tasks
- Create/update @IMPLEMENTATION_PLAN.md as a bullet point list sorted by priority of items yet to be implemented

Ultrathink. Study @IMPLEMENTATION_PLAN.md to determine starting point for research and keep it up to date with items considered complete/incomplete using subagents.

Log: [TASK] Creating/updating implementation plan...
Log when done: [TASK] ✓ Complete

# VALIDATE

Verify the plan makes sense:
- Does it cover missing functionality?
- Are priorities reasonable?
- Are tasks clear and actionable?

Log: [VALIDATE] Reviewing plan...
Log when done: [VALIDATE] ✓ Plan ready

# COMMIT

When plan is ready:
- `git add -A`
- `git commit -m "docs: update implementation plan"` (conventional commit format)
- `git push`

Log: [COMMIT] Committing plan...
Log when done: [COMMIT] ✓ Complete

# GUARDRAILS

Priority 1 (Must do):
- Plan only - do NOT implement anything
- Do NOT assume functionality is missing - confirm with code search in implementation first
- Implementation is in: internal/cmd/test/docs, root files, and infrastructure (.github/*, etc.)
- Specification is in: specs

Priority 2 (Important):
- If element is missing from implementation, search first to confirm it doesn't exist
- If needed, author specification at specs/FILENAME.md
- If you create new element, document plan to implement it in @IMPLEMENTATION_PLAN.md using subagent

Priority 3 (Context):
- Ultimate goal: Robust AI resource manager with semantic versioning, reproducible installs, and flexible registry support
- Consider missing elements and plan accordingly

# LOGGING

At start of each section above, output the log message shown.
When done with section, output the completion message.
This helps humans see where you are and where you get stuck.
