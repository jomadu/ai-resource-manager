---
description: "Don't create abstractions until you have multiple concrete examples"
globs: **/*
---

---
namespace: grug-brained-dev
ruleset:
  id: grug-brained-dev
  name: Grug Brained Developer Rules
  rules:
    - simple-tools
    - avoid-premature-optimization
    - debuggable-code
    - fear-of-looking-dumb
    - avoid-abstractions
    - test-simply
    - complexity-enemy
    - say-no
    - eighty-twenty
    - readable-over-clever
rule:
  id: avoid-abstractions
  name: Avoid Unnecessary Abstractions
  enforcement: SHOULD
  priority: 80
  scope:
    - files: ["**/*"]
---

# Avoid Unnecessary Abstractions (SHOULD)

grug see developer make abstraction too early

Rule of three:
- First time: write code
- Second time: copy and modify
- Third time: maybe abstract

Abstraction without real need make code complex and hard change
