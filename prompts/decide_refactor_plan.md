# Decide: Refactor Plan

## D12: If Criteria Fail Threshold - Propose Refactoring

Based on quality assessment:
- Which criteria failed?
- What specific issues were identified?
- What refactoring is needed to address each failure?
- What is the scope of refactoring required?
- If all criteria passed: no refactoring plan needed

## D9: Structure Plan by Priority (Most Important First)

Order the refactoring tasks:
- What are the critical quality issues? (correctness, security)
- What are the high-impact issues? (maintainability, clarity)
- What are the low-priority issues? (style, minor improvements)
- What depends on other tasks?
- What can be parallelized?

Most important first.

## D13: Prioritize by Impact

For each refactoring task:
- What is the impact if not addressed? (high/medium/low)
- What is the effort required? (high/medium/low)
- What is the risk of the refactoring? (breaking changes, regressions)
- Prioritize high-impact, low-effort tasks first
- Consider risk vs reward
