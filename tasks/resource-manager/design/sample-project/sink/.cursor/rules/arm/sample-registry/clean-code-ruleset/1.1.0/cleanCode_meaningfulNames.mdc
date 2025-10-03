---
description: "Choose names that reveal intent and are easy to understand"
globs: **/*.py, **/*.js, **/*.ts, **/*.java, **/*.go
alwaysApply: true
---

---
namespace: sample-registry
ruleset:
  id: cleanCode
  name: Clean Code
  rules:
    - meaningfulNames
    - smallFunctions
    - avoidComments
    - singleResponsibility
    - avoidDuplication
    - consistentFormatting
    - errorHandling
rule:
  id: meaningfulNames
  name: Use Meaningful Names
  enforcement: MUST
  priority: 100
  scope:
    - files: ["**/*.py", "**/*.js", "**/*.ts", "**/*.java", "**/*.go"]
---

# Use Meaningful Names (MUST)

Use intention-revealing names. Avoid disinformation, make meaningful distinctions, use pronounceable names, and use searchable names. A name should tell you why it exists, what it does, and how it is used.

Good: `userAccountBalance`, `isValidEmail`, `calculateTotalPrice`
Bad: `data`, `temp`, `x`, `flag`
