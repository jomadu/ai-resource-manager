# Act: Bootstrap

## A1: Create AGENTS.md with Operational Guide

Write AGENTS.md with the following sections:

**Work Tracking System:**
- What system is used
- How to query ready work
- How to update status
- How to mark complete
- Specific commands

**Story/Bug Input:**
- Where agents read story/bug descriptions
- File location or command to retrieve

**Planning System:**
- Draft plan location
- Publishing mechanism

**Build/Test/Lint Commands:**
- Specific commands to run tests
- Specific commands to run builds
- Specific commands to run linters
- Any setup required

**Specification Definition:**
- What file paths/patterns constitute specs
- What format are specs

**Implementation Definition:**
- What file paths/patterns constitute implementation
- What should be excluded

**Quality Criteria:**
- Boolean criteria for specs (if applicable)
- Boolean criteria for implementation (if applicable)

Capture the why - document rationale for decisions made.

## A2: Validate AGENTS.md Structure

Check AGENTS.md for required sections per agents-md-format.md:
- Work Tracking System
- Story/Bug Input
- Planning System
- Build/Test/Lint Commands
- Specification Definition
- Implementation Definition
- Quality Criteria

If sections are missing or incomplete:
- Add warning to AGENTS.md under "Operational Learnings"
- Document which sections need attention
- Provide guidance on what to add

## A3: Commit Changes

Commit AGENTS.md with descriptive message:
- What was created
- Why these definitions were chosen
- What was learned during bootstrap
- What validation warnings were added (if any)
