# ADR 0001: Manifest File Loading Pattern

**Status:** Accepted  
**Date:** 2024-12-19

## Context

We need to decide how `FileManager.loadManifest()` should behave when `arm.json` doesn't exist:
- Return error (like lockfile manager does)
- Return empty manifest (like old implementation did)

## Decision

Return empty manifest when file doesn't exist (no error).

## Rationale

1. **Most operations are "Add" operations** - they should work even if manifest doesn't exist yet
2. **Config file semantics** - manifest is a config file, it's okay if it doesn't exist initially (new project)
3. **Simpler code** - no boilerplate error handling in every operation
4. **Different from lockfile** - lockfile is a state file that only exists when packages are installed, so error makes sense there

## Consequences

- Operations don't need to check `os.IsNotExist(err)` and create empty manifest
- Simpler, more readable code
- Consistent with old implementation behavior
- Different pattern from lockfile manager (which is fine - different use cases)


