# Implementation Plan: `arm info dependency` Command

## Overview
Add `arm info dependency [REGISTRY_NAME/DEPENDENCY_NAME...]` command to display detailed information about installed dependencies (rulesets and promptsets).

**Specification:** See [specs/commands.md](specs/commands.md#arm-info-dependency)

## Prioritized Tasks

### 1. Add command routing in main.go
**File:** `cmd/arm/main.go`
- Add `case "dependency":` to `handleInfo()` switch statement (line ~1057)
- Call new `handleInfoDependency()` function

### 2. Implement handleInfoDependency() function
**File:** `cmd/arm/main.go`
- Parse dependency names from args (format: `registry/package`)
- If no args provided, get all dependencies from manifest
- For each dependency:
  - Call `svc.GetDependencyInfo(ctx, registry, name)`
  - Display formatted output with:
    - Dependency name (registry/package)
    - Type (ruleset/promptset)
    - Version constraint from manifest
    - Installed version from lockfile
    - Sinks where installed
    - Include/exclude patterns (if any)
    - Priority (for rulesets only)

### 3. Add help text
**File:** `cmd/arm/main.go`
- Update `printCommandHelp()` for "info" command (line ~155)
- Add example: `arm info dependency [REGISTRY/PACKAGE...]`

### 4. Add tests
**File:** `cmd/arm/list_info_dependency_test.go` (new file)
- Test with no dependencies configured
- Test with specific dependency names
- Test with all dependencies (no args)
- Test with invalid dependency format
- Test with non-existent dependency

## Implementation Notes

- Reuse existing `service.GetDependencyInfo()` method (already exists)
- Follow pattern from `handleInfoRegistry()` for arg parsing and display
- Use `manifest.ParseDependencyKey()` to split registry/package format
- Display both ruleset and promptset dependencies together
- Show "(none)" if no dependencies configured

## Estimated Complexity
**Low** - Follows existing patterns, service method already exists
