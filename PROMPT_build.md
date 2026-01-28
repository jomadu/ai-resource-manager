# ORIENT

Study `specs/*` with up to 500 parallel subagents to learn the application specifications.
Study @IMPLEMENTATION_PLAN.md to understand current work.
Study @AGENTS.md to learn how to build and test.

For reference, the application source code is in `internal/*`, `cmd/*`, `test/*`, and `docs/*`.

Log: [ORIENT] Starting orientation...
Log when done: [ORIENT] ✓ Complete

# TASK

Follow @IMPLEMENTATION_PLAN.md and choose the most important item to address.

Before making changes:
- Search codebase first (don't assume not implemented)
- Use up to 500 parallel subagents for searches/reads
- Use only 1 subagent for build/tests
- Use subagents when complex reasoning needed (debugging, architectural decisions)

Implement functionality per the specifications using parallel subagents.

Log: [TASK] Starting task: <task name>
Log when done: [TASK] ✓ Complete

# VALIDATE

Run tests for the code you changed: `go test ./...`

If tests fail:
- Fix the issues
- Run tests again
- Repeat until tests pass

If tests pass:
- Update @IMPLEMENTATION_PLAN.md with findings (use subagent)
- Mark task as complete

Log: [VALIDATE] Running tests...
Log when done: [VALIDATE] ✓ Tests passed

# COMMIT

When tests pass:
- `git add -A`
- `git commit -m "type: description"` (conventional commit format: feat/fix/docs/refactor/test/chore; use `!` for breaking changes like `feat!:`)
- `git push`

Log: [COMMIT] Committing changes...
Log when done: [COMMIT] ✓ Complete

# GUARDRAILS

Priority 1 (Must do):
- Keep @IMPLEMENTATION_PLAN.md current - future work depends on this to avoid duplicating efforts
- Commit only when tests pass
- No placeholders or stubs - implement completely
- If unrelated tests fail, resolve them as part of the increment
- Single sources of truth, no migrations/adapters

Priority 2 (Important):
- Update @AGENTS.md when learning operational things (keep brief and operational only)
- Capture the why in documentation - tests and implementation importance
- Clean completed items from plan when it gets large
- For any bugs noticed, resolve them or document in @IMPLEMENTATION_PLAN.md even if unrelated
- If inconsistencies found in docs/* or specs/*, use subagent with ultrathink to update them

Priority 3 (Nice to have):
- Create git tag with prefix 'ralph-' when no errors (start at ralph-0.0.0, increment patch)
- Add logging if needed for debugging

# LOGGING

At start of each section above, output the log message shown.
When done with section, output the completion message.
This helps humans see where you are and where you get stuck.
