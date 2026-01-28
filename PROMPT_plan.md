# ORIENT

Study @AGENTS.md for spec vs implementation definitions. If definitions don't exist, investigate and create them.
Study specifications (see @AGENTS.md for what constitutes specs) with up to 250 parallel subagents.
Study @IMPLEMENTATION_PLAN.md (if present) to understand the plan so far.
Study implementation (see @AGENTS.md for what constitutes implementation) with up to 250 parallel subagents.

Specifications drive the entire repo: README, CI/CD, configs, code, tests, and docs.

Log: [ORIENT] Starting orientation...
Log when done: [ORIENT] ✓ Complete

# TASK

Study @IMPLEMENTATION_PLAN.md (if present; assume it is inaccurate and incomplete).

Use up to 500 parallel subagents to:
- Study existing implementation (see @AGENTS.md for what constitutes implementation)
- Compare it against specifications (see @AGENTS.md for what constitutes specs)
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
- See @AGENTS.md for spec vs implementation definitions

Priority 2 (Important):
- If element is missing from implementation, search first to confirm it doesn't exist
- If needed, author specification (see @AGENTS.md for what constitutes specs)
- If you create new element, document plan to implement it in @IMPLEMENTATION_PLAN.md using subagent
- Update @AGENTS.md when learning operational things or if spec vs implementation definitions evolve (keep brief and operational only)

Priority 3 (Context):
- Ultimate goal: Robust AI resource manager with semantic versioning, reproducible installs, and flexible registry support
- Consider missing elements and plan accordingly

# LOGGING

At start of each section above, output the log message shown.
When done with section, output the completion message.
This helps humans see where you are and where you get stuck.
