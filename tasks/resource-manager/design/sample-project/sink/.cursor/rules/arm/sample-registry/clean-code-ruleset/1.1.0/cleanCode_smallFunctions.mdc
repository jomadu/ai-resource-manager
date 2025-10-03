---
description: "Functions should do one thing and do it well"
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
  id: smallFunctions
  name: Keep Functions Small
  enforcement: MUST
  priority: 90
  scope:
    - files: ["**/*.py", "**/*.js", "**/*.ts", "**/*.java", "**/*.go"]
---

# Keep Functions Small (MUST)

Functions should be small and do only one thing. They should fit on a screen and have a single level of abstraction. If you need to add comments to explain what a function does, it's probably too long.

- Aim for functions under 20 lines
- One level of abstraction per function
- Use descriptive names that eliminate the need for comments
