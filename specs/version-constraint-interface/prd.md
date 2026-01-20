# Product Requirements Document: Version and Constraint Interface Refactor

## 1. Introduction

Refactor Version and Constraint types from struct-based interpretation to interface-based behavior encapsulation. Current implementation exposes internal fields requiring callers to interpret meaning. New design uses methods to encapsulate behavior and explicit error handling for invalid operations.

**Problem:** Non-semver versions (branch names, commit hashes) cannot be compared without git history, but current API allows invalid comparisons with undefined behavior.

## 2. Goals

- Encapsulate Version and Constraint behavior behind methods
- Return errors for operations that cannot logically succeed (comparing non-semver versions)
- Remove BranchHead constraint type (constraints are semver-only)
- Preserve original string format in ToString methods
- Maintain backwards compatibility where possible

## 3. User Stories

### US-001: Version Constructor
**Description:** As a developer, I want to create Version instances from strings so that I can represent both semver and non-semver versions.

**Acceptance Criteria:**
- [ ] NewVersion accepts semver strings: `1.2.3`, `v1.2.3`
- [ ] NewVersion accepts non-semver strings: `main`, `abc123def`
- [ ] NewVersion returns error for empty string
- [ ] Typecheck passes
- [ ] Tests pass

### US-002: Version Comparison
**Description:** As a developer, I want to compare two semver versions so that I can determine ordering.

**Acceptance Criteria:**
- [ ] CompareTo returns -1 when self is older than other
- [ ] CompareTo returns 0 when self equals other
- [ ] CompareTo returns 1 when self is newer than other
- [ ] CompareTo returns error when either version is non-semver
- [ ] Comparison uses major → minor → patch ordering
- [ ] Typecheck passes
- [ ] Tests pass

### US-003: Version Convenience Methods
**Description:** As a developer, I want IsNewerThan and IsOlderThan methods so that I can write readable comparison code.

**Acceptance Criteria:**
- [ ] IsNewerThan returns true when CompareTo > 0
- [ ] IsOlderThan returns true when CompareTo < 0
- [ ] Both methods return error when either version is non-semver
- [ ] Typecheck passes
- [ ] Tests pass

### US-004: Version String Representation
**Description:** As a developer, I want ToString to preserve the original format so that version prefixes are maintained.

**Acceptance Criteria:**
- [ ] ToString returns original string with `v` prefix if present
- [ ] ToString works for both semver and non-semver versions
- [ ] ToString never returns error
- [ ] Typecheck passes
- [ ] Tests pass

### US-005: Constraint Constructor
**Description:** As a developer, I want to create Constraint instances from semver strings so that I can specify version requirements.

**Acceptance Criteria:**
- [ ] NewConstraint accepts `latest` keyword
- [ ] NewConstraint accepts exact semver: `1.2.3`, `v1.2.3`
- [ ] NewConstraint accepts abbreviated: `1`, `1.2`, `v1`, `v1.2`
- [ ] NewConstraint accepts caret: `^1.2.3`, `^v1.2.3`
- [ ] NewConstraint accepts tilde: `~1.2.3`, `~v1.2.3`
- [ ] NewConstraint returns error for non-semver strings
- [ ] Typecheck passes
- [ ] Tests pass

### US-006: Constraint Satisfaction Check
**Description:** As a developer, I want to check if a version satisfies a constraint so that I can validate version requirements.

**Acceptance Criteria:**
- [ ] IsSatisfiedBy returns true for Latest constraint with any semver
- [ ] IsSatisfiedBy checks exact match for Exact constraint
- [ ] IsSatisfiedBy checks major version range for Major constraint
- [ ] IsSatisfiedBy checks minor version range for Minor constraint
- [ ] IsSatisfiedBy returns error when version is non-semver
- [ ] Typecheck passes
- [ ] Tests pass

### US-007: Constraint String Representation
**Description:** As a developer, I want ToString to preserve the original constraint format so that prefixes are maintained.

**Acceptance Criteria:**
- [ ] ToString returns original string with `v`, `^`, `~` prefix if present
- [ ] ToString never returns error
- [ ] Typecheck passes
- [ ] Tests pass

### US-008: Update ResolveVersion Function
**Description:** As a developer, I want ResolveVersion to use new Version methods so that version resolution handles errors properly.

**Acceptance Criteria:**
- [ ] ResolveVersion uses Version.CompareTo method
- [ ] ResolveVersion handles errors from IsSatisfiedBy
- [ ] Typecheck passes
- [ ] Tests pass

### US-009: Update Manifest Manager
**Description:** As a developer, I want manifest manager to use new constructors so that version parsing is consistent.

**Acceptance Criteria:**
- [ ] Uses NewVersion and NewConstraint constructors
- [ ] Handles error returns from new methods
- [ ] Typecheck passes
- [ ] Tests pass

### US-010: Update Package Lock Manager
**Description:** As a developer, I want package lock manager to use new constructors so that version parsing is consistent.

**Acceptance Criteria:**
- [ ] Uses NewVersion constructor
- [ ] Handles error returns from new methods
- [ ] Typecheck passes
- [ ] Tests pass

### US-011: Update Sink Manager
**Description:** As a developer, I want sink manager to use new constructors so that version parsing is consistent.

**Acceptance Criteria:**
- [ ] Uses NewVersion constructor
- [ ] Handles error returns from new methods
- [ ] Typecheck passes
- [ ] Tests pass

### US-012: Update Storage Package
**Description:** As a developer, I want storage package to use new constructors so that version parsing is consistent.

**Acceptance Criteria:**
- [ ] Uses NewVersion constructor in package.go and storage.go
- [ ] Handles error returns from new methods
- [ ] Typecheck passes
- [ ] Tests pass

## 4. Functional Requirements

**FR-1:** Version.CompareTo must return error when comparing semver to non-semver versions

**FR-2:** Constraint.NewConstraint must reject non-semver input strings

**FR-3:** Constraint.IsSatisfiedBy must return error when checking non-semver versions

**FR-4:** Version.ToString and Constraint.ToString must preserve original string format including prefixes

**FR-5:** Caret constraint `^1.2.3` must match versions >=1.2.3 and <2.0.0

**FR-6:** Tilde constraint `~1.2.3` must match versions >=1.2.3 and <1.3.0

**FR-7:** Abbreviated constraint `1.2` must match versions >=1.2.0 and <1.3.0

**FR-8:** Abbreviated constraint `1` must match versions >=1.0.0 and <2.0.0

**FR-9:** Latest constraint must match any semver version

## 5. Non-Goals

- Git history integration for non-semver version comparison
- Support for non-semver constraints (branch names, commit hashes)
- Automatic version normalization (removing/adding `v` prefix)
- Range constraints with multiple operators (e.g., `>=1.2.0 <2.0.0`)

## 6. Technical Considerations

**Breaking Changes:**
- Remove BranchHead constraint type
- Replace ParseVersion/ParseConstraint with NewVersion/NewConstraint
- Add error returns to comparison methods

**Dependencies:**
- All packages using Version/Constraint types must be updated
- Test suites must cover error cases

**Migration Path:**
1. Update core types (version.go, constraint.go)
2. Update tests
3. Update dependent packages in order: helpers → storage → manifest → packagelockfile → sink

## 7. Success Metrics

- All existing tests pass with updated API
- New error cases have test coverage
- No undefined behavior for non-semver comparisons
- Zero regressions in version resolution logic
