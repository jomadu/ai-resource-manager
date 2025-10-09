---
description: "A class should have only one reason to change"
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
  id: singleResponsibility
  name: Single Responsibility Principle
  enforcement: SHOULD
  priority: 85
  scope:
    - files: ["**/*.py", "**/*.js", "**/*.ts", "**/*.java", "**/*.go"]
---

# Single Responsibility Principle (SHOULD)

Every class should have a single responsibility. If a class has multiple reasons to change, it violates the Single Responsibility Principle and should be refactored into smaller, focused classes.

Signs of violation:
- Multiple responsibilities in one class
- Classes that are hard to test
- Classes that are hard to understand
