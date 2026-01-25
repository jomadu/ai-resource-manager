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
  id: avoidDuplication
  name: Don't Repeat Yourself (DRY)
  description: Every piece of knowledge should have a single representation
  priority: 75
  enforcement: may
  scope:
    - files: ["**/*.py", "**/*.js", "**/*.ts", "**/*.java", "**/*.go"]
---

# Don't Repeat Yourself (DRY) (MAY)

Duplication is the root of all evil in software design. When you find yourself copying and pasting code, extract it into a function, class, or module that can be reused.

Types of duplication to avoid:
- Copy-paste code
- Similar logic in different places
- Repeated configuration
- Duplicate data structures