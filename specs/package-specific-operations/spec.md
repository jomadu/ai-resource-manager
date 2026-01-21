# Package-Specific Operations

## 1. Introduction

Currently, ARM only supports operating on all packages at once (`arm uninstall`, `arm update`, `arm upgrade`). Users need the ability to target specific packages for these operations without having to specify resource types (ruleset/promptset), since they're all just dependencies.

## 2. Goals

- Support uninstalling specific packages by name
- Support updating specific packages by name
- Support upgrading specific packages by name
- Maintain resource-type agnostic interface (no need to specify ruleset vs promptset)
- Support multiple packages in a single command

## 3. User Stories

### US-001: Uninstall Specific Package
**Description:** As a developer, I want to uninstall a specific package without affecting other installed packages.

**Acceptance Criteria:**
- [ ] `arm uninstall registry/package` removes only that package
- [ ] `arm uninstall pkg1 pkg2 pkg3` removes multiple packages
- [ ] Warns if package not found, continues with others
- [ ] Removes from manifest and all sinks
- [ ] Typechecks pass
- [ ] Tests pass

### US-002: Update Specific Package
**Description:** As a developer, I want to update specific packages to their latest compatible versions.

**Acceptance Criteria:**
- [ ] `arm update registry/package` updates only that package
- [ ] `arm update pkg1 pkg2` updates multiple packages
- [ ] Respects version constraints
- [ ] Warns if package not found, continues with others
- [ ] Typechecks pass
- [ ] Tests pass

### US-003: Upgrade Specific Package
**Description:** As a developer, I want to upgrade specific packages to their absolute latest versions.

**Acceptance Criteria:**
- [ ] `arm upgrade registry/package` upgrades only that package
- [ ] `arm upgrade pkg1 pkg2` upgrades multiple packages
- [ ] Ignores version constraints
- [ ] Updates constraint to ^X.0.0 based on new version
- [ ] Warns if package not found, continues with others
- [ ] Typechecks pass
- [ ] Tests pass

## 4. Functional Requirements

**FR-1:** Commands accept zero or more package names as arguments
- Zero arguments: operate on all packages (current behavior)
- One or more: operate only on specified packages

**FR-2:** Package names use format `registry/package` (no resource type needed)

**FR-3:** System automatically detects if package is ruleset or promptset

**FR-4:** If package not found, warn and continue with remaining packages

**FR-5:** Exit code 0 if at least one package succeeds, 1 if all fail

## 5. Command Signatures

```bash
# Uninstall
arm uninstall [REGISTRY/PACKAGE...]

# Update (respects constraints)
arm update [REGISTRY/PACKAGE...]

# Upgrade (ignores constraints)
arm upgrade [REGISTRY/PACKAGE...]
```

## 6. Examples

```bash
# Uninstall specific packages
arm uninstall my-org/clean-code-ruleset
arm uninstall my-org/pkg1 my-org/pkg2

# Update specific packages
arm update my-org/clean-code-ruleset
arm update my-org/pkg1 my-org/pkg2

# Upgrade specific packages
arm upgrade my-org/clean-code-ruleset
arm upgrade my-org/pkg1 my-org/pkg2

# Operate on all (existing behavior)
arm uninstall
arm update
arm upgrade
```

## 7. Non-Goals

- Resource type filtering (e.g., `arm update ruleset`)
- Pattern matching (e.g., `arm update my-org/*`)
- Interactive selection
- Dry-run mode

## 8. Technical Considerations

- Reuse existing uninstall/update/upgrade logic
- Lookup package in manifest to determine type
- Filter operations to specified packages only
- Maintain backward compatibility with zero-argument form

## 9. Success Metrics

- All workflow scripts work without modification
- Users can selectively manage packages
- No breaking changes to existing commands
