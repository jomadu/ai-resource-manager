# ORIENT

Study @AGENTS.md for spec vs implementation definitions. If definitions don't exist, investigate and create them.
ULTRASTUDY the implementation (see @AGENTS.md for what constitutes implementation) to understand architecture, features, components, functions, and edge cases.

Specifications should drive the entire repo: README, CI/CD, configs, code, tests, and docs.

Log: [ORIENT] Starting deep implementation study...
Log when done: [ORIENT] ✓ Complete

# TASK

Update spec README (see @AGENTS.md for what constitutes specs) with current Jobs to be Done (JTBDs) and Topics of Concern:
- Assume existing README.md is inaccurate, incomplete, and outdated
- Verify all claims by inspection and hard study of the implementation
- Include JTBDs and topics for code, documentation, root files, and infrastructure (workflows, CI/CD, etc.)
- README.md may not exist yet - create it
- Keep the existing format (JTBDs section, Topics section, Spec docs links)
- Update if core JTBDs or topics change
- Allow large refactorings when structure needs to change

Update or create spec documents (see @AGENTS.md for what constitutes specs) using TEMPLATE.md:
- Assume existing specs are inaccurate, incomplete, and outdated
- Verify all claims (especially acceptance criteria) by inspection and hard study of the implementation
- Specs may not exist yet - create them
- Implementation is the source of truth
- Deep-dive verification: Trace through actual implementation code to find exact discrepancies between spec and real behavior, then fix specs directly
- When bugs found in implementation: uncheck relevant acceptance criteria and add to "Known Issues" section of spec
- When opportunities for improvement found: add to "Areas for Improvement" section of spec
- Reorganize/refactor specs if JTBDs or topics have evolved
- Keep specs minimal but complete and accurate
- If @AGENTS.md spec vs implementation definitions don't exist, create them based on investigation

Log: [TASK] Updating specifications from implementation...
Log when done: [TASK] ✓ Complete

# VALIDATE

Verify specs match implementation reality:
- Do acceptance criteria reflect actual behavior?
- Are known issues documented?
- Are improvement opportunities captured?

Log: [VALIDATE] Reviewing specs...
Log when done: [VALIDATE] ✓ Specs accurate

# COMMIT

When specs are accurate:
- `git add -A`
- `git commit -m "docs: update [specs...] from implementation"` (conventional commit format: feat/fix/docs/refactor/test/chore) (fill in specs... with relevant specifications)
- `git push`

Log: [COMMIT] Committing spec updates...
Log when done: [COMMIT] ✓ Complete

# GUARDRAILS

Priority 1 (Must do):
- Implementation is the source of truth - specs document what exists
- Verify all claims by inspection - don't assume specs are correct, or that the JTBDs or topics of concern are complete or accurate
- Deep-dive verification - trace through actual implementation code to find discrepancies
- See @AGENTS.md for spec vs implementation definitions

Priority 2 (Important):
- When bugs found in implementation: uncheck acceptance criteria, add to "Known Issues"
- When improvements found: add to "Areas for Improvement"
- Update @AGENTS.md when learning how to contextualize with implementation or run application, or if spec vs implementation definitions evolve (keep brief and operational only)
- Root files and infrastructure need specs too (CI/CD, build configs, install scripts, etc.)

Priority 3 (Context):
- Focus: What does the implementation actually do? What are the acceptance criteria from tests?

# LOGGING

At start of each section above, output the log message shown.
When done with section, output the completion message.
This helps humans see where you are and where you get stuck.
