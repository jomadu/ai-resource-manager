---
applyTo: "**/*.py,**/*.js,**/*.ts,**/*.java,**/*.go"
---

---
namespace: sample-registry
ruleset:
  id: cleanCode
  name: Clean Code
  description: Essential clean code principles for maintainable software
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
  description: Choose names that reveal intent and are easy to understand
  priority: 100
  enforcement: must
  scope:
    - files: ["**/*.py", "**/*.js", "**/*.ts", "**/*.java", "**/*.go"]
---

# Use Meaningful Names (MUST)

Use intention-revealing names. Avoid disinformation, make meaningful distinctions, use pronounceable names, and use searchable names. A name should tell you why it exists, what it does, and how it is used.

Good: `userAccountBalance`, `isValidEmail`, `calculateTotalPrice`
Bad: `data`, `temp`, `x`, `flag`