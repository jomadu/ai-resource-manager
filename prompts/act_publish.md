# Act: Publish

## A9: Execute Work Tracking Commands per AGENTS.md to Create Issues from Draft Plan

For each task in the draft plan, in dependency order:
- Execute work tracking commands from AGENTS.md
- Create issue with title, description, type, priority
- Set dependencies (blockers, parent-child relationships)
- Set labels or tags as appropriate
- Capture issue IDs for dependency references
- Follow the order determined in decide phase

## A6: Incorporate Learnings into AGENTS.md

If operational learnings occurred during publishing:
- **Incorporate into existing sections** - Update commands, paths, or criteria where they live
- **Add inline rationale** - Brief comment explaining why (e.g., "# Using X instead of Y - reason")
- **Don't append diary entries** - No dated logs like "YYYY-MM-DD: discovered X"
- Examples: work tracking commands needed adjustment → update Work Tracking System section with corrected commands

## A2: Commit Changes

Commit any AGENTS.md updates and work tracking state:
- Descriptive commit message
- What plan was published
- How many issues were created
- Reference draft plan file
