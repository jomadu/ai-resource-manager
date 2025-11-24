---
description: "Don't optimize until you have real performance problems"
globs: **/*
alwaysApply: true
---

---
namespace: grug-brained-dev
ruleset:
  id: grug-brained-dev
  name: Grug Brained Developer Rules
  rules:
    - complexity-enemy
    - say-no
    - eighty-twenty
    - readable-over-clever
    - test-simply
    - avoid-premature-optimization
    - debuggable-code
    - fear-of-looking-dumb
    - avoid-abstractions
    - simple-tools
rule:
  id: avoid-premature-optimization
  name: No Premature Optimization
  enforcement: MUST
  priority: 70
  scope:
    - files: ["**/*"]
---

# No Premature Optimization (MUST)

grug see many developer optimize too early

This bad because:
- Make code complex for no reason
- Waste time on problems that not exist
- Hard to change later

First make it work, then make it fast if needed
