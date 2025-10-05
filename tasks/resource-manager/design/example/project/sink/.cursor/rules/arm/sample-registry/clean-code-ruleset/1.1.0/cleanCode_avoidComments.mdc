---
description: "Write self-documenting code instead of relying on comments"
globs: **/*.py, **/*.js, **/*.ts, **/*.java, **/*.go
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
  id: avoidComments
  name: Don't Comment Bad Code
  enforcement: SHOULD
  priority: 80
  scope:
    - files: ["**/*.py", "**/*.js", "**/*.ts", "**/*.java", "**/*.go"]
---

# Don't Comment Bad Code (SHOULD)

Comments are a failure. If you need a comment to explain what your code does, rewrite the code to be self-explanatory. Good code documents itself through clear naming and structure.

Exceptions: Legal comments, TODO comments, warnings, and public API documentation.
