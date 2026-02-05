# Decide: Build

## D4: Pick the Most Important Task from Work Tracking

Based on ready work from work tracking system:
- What is the highest priority task?
- What are the dependencies? (is it actually ready?)
- What is the scope? (can it be completed in one iteration?)
- If multiple tasks have same priority, which has most impact?
- Select exactly one task for this iteration

## D5: Determine Implementation Approach Using Parallel Subagents

Plan how to implement the selected task:
- What is the overall implementation strategy?
- Can work be parallelized across subagents?
- How many subagents are needed?
- **Critical: Only 1 subagent for build/tests** (avoid parallel test conflicts)
- What is each subagent responsible for?
- What is the execution order?

## D6: Identify Which Files to Modify

Based on the implementation approach:
- What files need to be created?
- What files need to be modified?
- What functions/classes/methods need to be added?
- What existing code needs to be refactored?
- What tests need to be written or updated?
- What documentation needs to be updated?
