# Orient: Gap Analysis

## R16: Gap Analysis - Compare Specs vs Implementation

Perform systematic comparison:
- What is specified but not implemented? (missing features)
- What is implemented but not specified? (undocumented code)
- What is specified differently than implemented? (drift)
- What specs have no corresponding implementation files?
- What implementation files have no corresponding specs?

## R6: Search Codebase - Don't Assume Not Implemented

**Critical: This is the Achilles' heel.**

Before concluding something is missing, search thoroughly:
- Search for function/class/method names mentioned in specs
- Search for keywords from spec descriptions
- Search for related functionality that might implement the spec differently
- Check file trees and symbols for relevant modules
- Verify assumptions empirically

Don't assume not implemented. Search first.

## R17: Assess Completeness and Accuracy

Evaluate the gaps identified:
- How complete is the implementation relative to specs?
- How accurate are the specs relative to implementation?
- What are the highest priority gaps? (critical features, security, correctness)
- What are the lowest priority gaps? (nice-to-have, documentation)
- What gaps block other work?
- What gaps are independent and can be parallelized?
