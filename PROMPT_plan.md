# ORIENT

Study `docs/*` and `specs/*` with up to 250 parallel subagents to learn the application specifications.
Study @IMPLEMENTATION_PLAN.md (if present) to understand the plan so far.
Study `internal/*` and `cmd/*` with up to 250 parallel subagents to understand shared utilities and components.

For reference, the application source code is in `internal/*` and `cmd/*`.

Log: [ORIENT] Starting orientation...
Log when done: [ORIENT] ✓ Complete

# TASK

Study @IMPLEMENTATION_PLAN.md (if present; it may be incorrect).

Use up to 500 parallel subagents to:
- Study existing source code in `internal/*` and `cmd/*`
- Compare it against `docs/*` and `specs/*`
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
- Do NOT assume functionality is missing - confirm with code search first
- Treat `internal/*` and `cmd/*` as the project's standard libraries
- Prefer consolidated, idiomatic implementations over ad-hoc copies

Priority 2 (Important):
- If element is missing, search first to confirm it doesn't exist
- If needed, author specification at specs/FILENAME.md
- If you create new element, document plan to implement it in @IMPLEMENTATION_PLAN.md using subagent

Priority 3 (Context):
- Ultimate goal: Robust AI resource manager with semantic versioning, reproducible installs, and flexible registry support
- Consider missing elements and plan accordingly

# LOGGING

At start of each section above, output the log message shown.
When done with section, output the completion message.
This helps humans see where you are and where you get stuck.
