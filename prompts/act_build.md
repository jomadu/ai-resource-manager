# Act: Build

## A3: Implement Using Parallel Subagents

Execute the implementation plan:
- Spawn subagents for parallel work as determined in decide phase
- **Critical: Only 1 subagent for build/tests** (avoid parallel test conflicts)
- Each subagent handles its assigned files/tasks
- Coordinate file modifications to avoid conflicts
- Complete all file modifications before proceeding

## A3.5: Validate Spec Structure (If Specs Modified)

If any spec files in `specs/*.md` were created or modified:
- Check filename matches naming convention (lowercase-with-hyphens.md):
  - Must be lowercase letters, hyphens, and .md extension
  - No underscores, spaces, or uppercase letters
  - Warn if non-conforming, suggest correct naming
- Check for required sections from TEMPLATE.md:
  - "Job to be Done" or "Jobs to be Done"
  - "Activities"
  - "Acceptance Criteria"
- Warn if required sections are missing
- Suggest using `specs/TEMPLATE.md` for new specs
- This is informational only - does not block commit

## A3.6: Regenerate Spec Index (If Specs Modified)

If any spec files in `specs/*.md` were created, modified, or deleted:
- Read `specs/specification-system.md` for README structure requirements
- Scan all `specs/*.md` files (excluding README.md, TEMPLATE.md, specification-system.md)
- Extract "## Job to be Done" or "## Jobs to be Done" section from each spec
- Generate `specs/README.md` following the documented structure:
  - Header and intro text
  - Links to TEMPLATE.md and specification-system.md
  - List of specs with extracted JTBDs
- This keeps the index in sync with actual specs

## A4: Run Tests per AGENTS.md (Backpressure)

Execute test commands from AGENTS.md:
- Run all relevant tests
- Verify tests pass
- If tests fail: fix issues, don't proceed
- This is backpressure - quality gate before commit

## A5: Update Work Tracking per AGENTS.md

Update the work tracking system:
- Mark task as complete if fully done
- Update status if partially done or blocked
- Use commands specified in AGENTS.md
- Document what was accomplished

## A6: Incorporate Learnings into AGENTS.md

If operational learnings occurred:
- **Incorporate into existing sections** - Update commands, paths, or criteria where they live
- **Add inline rationale** - Brief comment explaining why (e.g., "# Using X instead of Y - reason")
- **Don't append diary entries** - No dated logs like "YYYY-MM-DD: discovered X"
- Examples: command failed → update command with working version; path incorrect → fix path; pattern discovered → update definition

## A7: Commit When Tests Pass

Commit all changes with descriptive message:
- What was implemented
- What task was completed
- Reference task ID from work tracking
- Only commit when tests pass (backpressure enforced)
