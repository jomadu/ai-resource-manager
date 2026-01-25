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
  id: avoidComments
  name: Don't Comment Bad Code
  description: Write self-documenting code instead of relying on comments
  priority: 80
  enforcement: should
  scope:
    - files: ["**/*.py", "**/*.js", "**/*.ts", "**/*.java", "**/*.go"]
---

# Don't Comment Bad Code (SHOULD)

Comments are a failure. If you need a comment to explain what your code does, rewrite the code to be self-explanatory. Good code documents itself through clear naming and structure.

Exceptions: Legal comments, TODO comments, warnings, and public API documentation.