# Branch Constraint Resolution Issue Analysis

## Problem Statement

Users cannot install rulesets from Git branches using `@branch` syntax (e.g., `arm install ai-rules/amazonq-rules@main`) due to incorrect branch constraint detection logic.

## Root Cause Analysis

### Issue 1: Overly Restrictive Branch Detection
Branch constraint detection only allows branches explicitly listed in the registry's `branches` configuration array. This violates the principle that users should be able to install from any existing branch.

**Problem:** If `main` isn't in the configured branches list, it gets treated as a semantic version constraint, leading to "no matching version found" errors.

### Issue 2: Semantic Purpose Confusion
The `branches` configuration field serves two conflicting purposes:
1. **Intended:** Priority ordering for "latest" version resolution
2. **Actual:** Access control restricting which branches can be installed

### Issue 3: Default Configuration Gap
When registries are created from manifest files, the `branches` field often isn't populated, defaulting to an empty array and blocking all branch installations.

## Architectural Design Flaw

The current design conflates **version resolution strategy** with **access permissions**. The `branches` config should only influence resolution priority, not restrict user access to existing repository branches.

## Holistic Solution

### 1. Separate Concerns
- **Branch Detection:** Allow any non-semver identifier as a potential branch
- **Priority Resolution:** Use `branches` config only for "latest" resolution ordering
- **Access Validation:** Check actual repository branch existence, not configuration

### 2. Updated Branch Detection Strategy
Change branch detection to be permissive:
- Exclude semantic version patterns (contains dots, follows semver format)
- Exclude special keywords ("latest")
- Treat everything else as potential branch names

### 3. Runtime Validation
Move branch existence validation to resolution time with helpful error messages that list available branches when a requested branch doesn't exist.

### 4. Default Configuration
Provide sensible defaults for the `branches` configuration to ensure "latest" resolution works out of the box with common branch names like "main" and "master".

## Implementation Areas

1. **Git Registry** - Update branch constraint detection and version resolution
2. **Registry Factory** - Add default branches configuration
3. **Constraint Resolver** - Ensure branch constraints are properly parsed

## Testing Strategy

Create integration tests that verify:
1. Installation from `main` branch succeeds
2. Installation from any existing branch succeeds
3. Installation from non-existent branch fails with helpful error
4. "latest" resolution respects `branches` priority order
5. Backward compatibility with existing configurations

## Key Principles

1. **Permissive by Default:** Allow access to any existing branch
2. **Fail Fast with Context:** Provide clear error messages with available options
3. **Configuration as Preference:** Use config for optimization, not restriction
4. **Separation of Concerns:** Keep resolution strategy separate from access control

This solution addresses the root architectural issue rather than patching symptoms, ensuring robust branch handling across all use cases.
