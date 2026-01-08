# ARM Rulesets

This file defines the installation priorities for rulesets managed by ARM.

## Priority Rules

**This index is the authoritative source of truth for ruleset priorities.** When conflicts arise between rulesets, follow this priority order:

1. **Higher priority numbers take precedence** over lower priority numbers
2. **Rules from higher priority rulesets override** conflicting rules from lower priority rulesets
3. **Always consult this index** to resolve any ambiguity about which rules to follow

## Installed Rulesets

### sample-registry/clean-code-ruleset@1.1.0
- **Priority:** 100
- **Rules:**
  - arm_1a2b_3c4d_cleanCode_meaningfulNames.instructions.md
  - arm_1a2b_5e6f_cleanCode_smallFunctions.instructions.md
  - arm_1a2b_7g8h_cleanCode_avoidComments.instructions.md
  - arm_1a2b_9i0j_cleanCode_singleResponsibility.instructions.md
  - arm_1a2b_k1l2_cleanCode_avoidDuplication.instructions.md