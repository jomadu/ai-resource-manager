# Orient: Build

## R5: Understand Task Requirements

Study the task from work tracking:
- What is the specific deliverable?
- What are the acceptance criteria?
- What specs define this functionality?
- What are the constraints and dependencies?
- What is the scope? (what's in, what's out)
- What does "done" look like?

## R6: Search Codebase - Don't Assume Not Implemented

**Critical: This is the Achilles' heel.**

Before implementing anything, search thoroughly:
- Search for function/class/method names related to the task
- Search for keywords from the task description
- Search for related functionality that might already exist
- Check file trees and symbols for relevant modules
- Look for partial implementations or stubs
- Verify assumptions empirically

Don't assume not implemented. Search first.

## R7: Identify What Needs to Be Built/Modified

Based on task requirements and codebase search:
- What files need to be created?
- What files need to be modified?
- What functions/classes/methods need to be added?
- What existing code needs to be refactored?
- What integration points exist?
- What dependencies need to be added?
- What can be reused vs built new?

## R8: Determine Test Strategy

Plan how to verify the implementation:
- What tests need to be written?
- What existing tests need to be updated?
- What test data is needed?
- What edge cases must be covered?
- What integration tests are required?
- How will acceptance criteria be verified?
- What manual verification steps are needed?
