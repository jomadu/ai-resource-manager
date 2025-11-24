---
description: "Write code that is easy to debug and understand"
globs: **/*.js, **/*.ts, **/*.py, **/*.java, **/*.go
alwaysApply: true
---

---
namespace: grug-brained-dev
ruleset:
  id: grug-brained-dev
  name: Grug Brained Developer Rules
  rules:
    - test-simply
    - complexity-enemy
    - say-no
    - eighty-twenty
    - readable-over-clever
    - simple-tools
    - avoid-premature-optimization
    - debuggable-code
    - fear-of-looking-dumb
    - avoid-abstractions
rule:
  id: debuggable-code
  name: Keep Code Debuggable
  enforcement: MUST
  priority: 85
  scope:
    - files: ["**/*.js", "**/*.ts", "**/*.py", "**/*.java", "**/*.go"]
---

# Keep Code Debuggable (MUST)

grug spend much time debugging, so code must be debuggable

Good practices:
- Use clear variable names
- Add logging at important points
- Avoid deep nesting
- One thing per function
- Comments explain why, not what
