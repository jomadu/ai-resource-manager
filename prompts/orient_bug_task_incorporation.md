# Orient: Bug Task Incorporation

## R10: Analyze Bug from Task File

Study the bug description from the task file:
- What are the symptoms? (observable incorrect behavior)
- What is the root cause? (underlying issue)
- What functionality is affected?
- What are the reproduction steps?
- What is the expected vs actual behavior?
- What is the impact and severity?

## R11: Understand Existing Spec Structure and Patterns

Study the specifications to understand:
- How are specs currently organized? (by JTBD, by component, by feature)
- What patterns do existing specs follow?
- What level of detail is typical?
- How are acceptance criteria expressed?
- How are edge cases documented?
- How are error conditions described?

## R13: Determine How Spec Should Be Adjusted to Drive Bug Fix

Based on the bug analysis and spec patterns, determine:
- Which spec files are affected by this bug?
- What acceptance criteria are missing or incorrect?
- What edge cases need to be added?
- What clarifications are needed in existing specs?
- What error handling needs to be specified?
- Should this expose gaps in the spec structure itself?

## R14: If Draft Plan Exists - Critique It

If a draft plan file exists from a previous iteration:
- Is it complete? (covers all spec adjustments needed)
- Is it accurate? (correctly identifies the spec gaps)
- Are priorities correct? (most critical adjustments first)
- Is it clear? (tasks are well-defined and actionable)
- What needs to be adjusted?

## R15: Identify Tasks Needed

Based on the spec adjustment strategy, identify:
- What acceptance criteria need to be added or corrected?
- What edge cases need to be documented?
- What error conditions need to be specified?
- What clarifications need to be made?
- What examples need to be added to prevent regression?
- What data structure constraints need to be tightened?

Break down into discrete, implementable tasks.
