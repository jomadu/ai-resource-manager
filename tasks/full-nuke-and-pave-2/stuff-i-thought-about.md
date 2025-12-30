Tasks

- [x] concurrency protection in storage at registry, repo, and package level.
- [] sink implementation with index files, ruleset index file
- [] update arm.json format:
  - Change `packages` → `dependencies` with separate `rulesets` and `promptsets` sections
  - Remove `resourceType` field (structure indicates type)
  - Simplify sinks: replace `compileTarget` + `layout` with just `tool` field
  - Add optional `include`/`exclude` patterns for selective installs
  - Keep registry `type` field as-is
