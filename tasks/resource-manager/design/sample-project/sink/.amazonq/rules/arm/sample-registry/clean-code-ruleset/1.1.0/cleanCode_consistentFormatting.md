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
  id: consistentFormatting
  name: Consistent Code Formatting
  enforcement: SHOULD
  priority: 70
  scope:
    - files: ["**/*.py", "**/*.js", "**/*.ts", "**/*.java", "**/*.go"]
---

# Consistent Code Formatting (SHOULD)

Consistent formatting makes code easier to read and maintain. Use automated formatting tools and establish team-wide style guidelines.

- Use linters and formatters (ESLint, Prettier, Black, gofmt)
- Configure your IDE for consistent formatting
- Establish coding standards for your team
- Use pre-commit hooks to enforce formatting
