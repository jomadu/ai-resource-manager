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
  id: errorHandling
  name: Proper Error Handling
  enforcement: MUST
  priority: 85
  scope:
    - files: ["**/*.py", "**/*.js", "**/*.ts", "**/*.java", "**/*.go"]
---

# Proper Error Handling (MUST)

Error handling is crucial for robust applications. Always handle errors explicitly and provide meaningful error messages to help with debugging.

- Use try-catch blocks appropriately
- Don't ignore errors silently
- Provide context in error messages
- Log errors with sufficient detail
- Use specific exception types when possible
