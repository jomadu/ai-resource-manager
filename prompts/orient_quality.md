# Orient: Quality Assessment

## R17: Validate Quality Criteria Are Boolean

Before applying criteria, validate AGENTS.md defines boolean criteria:

**Check each criterion:**
- Is it PASS/FAIL? (boolean)
- Is it yes/no? (boolean)
- Is it true/false? (boolean)
- Does it have a clear threshold? (boolean)

**Warn about non-boolean criteria:**
- Subjective scores (1-5, percentages)
- Ambiguous language ("good", "adequate", "reasonable")
- Relative comparisons ("better than", "worse than")
- Vague thresholds ("most", "some", "few")

**Examples of boolean criteria:**
- ✅ "All functions have docstrings" (PASS/FAIL)
- ✅ "Test coverage > 80%" (PASS/FAIL with clear threshold)
- ✅ "No functions exceed 50 lines" (PASS/FAIL)
- ❌ "Code quality is good" (subjective, no threshold)
- ❌ "Documentation is adequate" (ambiguous)
- ❌ "Performance is acceptable" (vague)

**If non-boolean criteria found:**
- Document in Operational Learnings section of AGENTS.md
- Suggest boolean alternatives
- Proceed with available boolean criteria only

## R18: Apply Boolean Criteria per AGENTS.md

Apply the quality criteria defined in AGENTS.md:
- What boolean criteria apply? (specs or implementation)
- What are the specific thresholds?
- What metrics or checks are defined?
- Evaluate each criterion systematically
- Document findings for each criterion

## R19: Identify Human Markers

Search for indicators of quality issues:
- TODOs, FIXMEs, HACKs in code or specs
- Code smells (duplication, long functions, complex conditionals)
- Unclear language in specs (ambiguous, vague, contradictory)
- Missing documentation
- Incomplete implementations
- Inconsistent patterns
- Dead code or unused specs

## R20: Score Each Criterion PASS/FAIL

For each quality criterion:
- PASS: criterion is met, no action needed
- FAIL: criterion is not met, refactoring required

Use boolean scoring only - no subjective grades or percentages.

Determine overall assessment:
- If all criteria PASS: no refactoring needed
- If any criteria FAIL: refactoring plan required
- Identify which specific criteria failed
- Prioritize by impact and severity
