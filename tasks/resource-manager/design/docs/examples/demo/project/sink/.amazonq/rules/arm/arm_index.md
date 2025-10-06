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
  - arm/sample-registry/clean-code-ruleset/1.1.0/cleanCode_meaningfulNames.md
  - arm/sample-registry/clean-code-ruleset/1.1.0/cleanCode_smallFunctions.md
  - arm/sample-registry/clean-code-ruleset/1.1.0/cleanCode_avoidComments.md
  - arm/sample-registry/clean-code-ruleset/1.1.0/cleanCode_singleResponsibility.md
  - arm/sample-registry/clean-code-ruleset/1.1.0/cleanCode_avoidDuplication.md
  - arm/sample-registry/clean-code-ruleset/1.1.0/cleanCode_consistentFormatting.md
  - arm/sample-registry/clean-code-ruleset/1.1.0/cleanCode_errorHandling.md
