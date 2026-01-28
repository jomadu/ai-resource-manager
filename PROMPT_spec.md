# ORIENT

ULTRASTUDY the code in `internal/*`, `cmd/*`, `test/*`, and project files in root directory to understand architecture, features, components, functions, and edge cases.

Study user documentation in `./docs` as part of the implementation.

Log: [ORIENT] Starting deep code study...
Log when done: [ORIENT] ✓ Complete

# TASK

Update specs/README.md with current Jobs to be Done (JTBDs) and Topics of Concern:
- Assume existing README.md is inaccurate, incomplete, and outdated
- Verify all claims by inspection and hard study of the implementation
- Include JTBDs and topics for both code implementation and user documentation (./docs)
- README.md may not exist yet - create it
- Keep the existing format (JTBDs section, Topics section, Spec docs links)
- Update if core JTBDs or topics change
- Allow large refactorings when structure needs to change

Update or create spec documents in specs/ using TEMPLATE.md:
- Assume existing specs are inaccurate, incomplete, and outdated
- Verify all claims (especially acceptance criteria) by inspection and hard study of the implementation
- Specs may not exist yet - create them
- Implementation is the source of truth
- Deep-dive verification: Trace through actual implementation code to find exact discrepancies between spec and real behavior, then fix specs directly
- When bugs found in code: uncheck relevant acceptance criteria and add to "Known Issues" section of spec
- When opportunities for improvement found: add to "Areas for Improvement" section of spec
- Reorganize/refactor specs if JTBDs or topics have evolved
- Keep specs minimal but complete and accurate

Log: [TASK] Updating specs from implementation...
Log when done: [TASK] ✓ Complete

# VALIDATE

Verify specs match reality:
- Do acceptance criteria reflect actual behavior?
- Are known issues documented?
- Are improvement opportunities captured?

Log: [VALIDATE] Reviewing specs...
Log when done: [VALIDATE] ✓ Specs accurate

# COMMIT

When specs are accurate:
- `git add -A`
- `git commit -m "docs: update specs [specs...] from implementation"` (conventional commit format: feat/fix/docs/refactor/test/chore)
- `git push`

Log: [COMMIT] Committing spec updates...
Log when done: [COMMIT] ✓ Complete

# GUARDRAILS

Priority 1 (Must do):
- Implementation is the source of truth - specs document what exists
- Verify all claims by inspection - don't assume specs are correct, or that the JTBDs or topics of concern are complete or accurate
- Deep-dive verification - trace through actual code to find discrepancies

Priority 2 (Important):
- When bugs found: uncheck acceptance criteria, add to "Known Issues"
- When improvements found: add to "Areas for Improvement"
- Update @AGENTS.md when learning how to contextualize with codebase or run application (keep brief and operational only)

Priority 3 (Context):
- Focus: What does the code actually do? What are the acceptance criteria from tests?

# LOGGING

At start of each section above, output the log message shown.
When done with section, output the completion message.
This helps humans see where you are and where you get stuck.
