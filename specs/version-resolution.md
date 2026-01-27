# Version Resolution

## Job to be Done
Resolve package versions from semantic version constraints, tags, branches, and "latest" specifiers to concrete version identifiers.

## Activities
1. Parse semantic versions (1.0.0, v2.1.3)
2. Parse version constraints (^1.0.0, ~1.2.0, 1.x, 1.2.x)
3. Resolve "latest" to highest semantic version
4. Resolve branches to commit hashes
5. Compare versions including prerelease identifiers

## Acceptance Criteria
- [x] Parse semantic versions with optional 'v' prefix
- [x] Parse major (1), minor (1.2), and exact (1.2.3) constraints
- [x] Parse caret (^1.0.0) and tilde (~1.2.0) constraints
- [x] Resolve "latest" to highest semantic version
- [x] Resolve branch names to commit hashes
- [ ] Compare versions with prerelease precedence (1.0.0-alpha.1 < 1.0.0-alpha.2 < 1.0.0-beta.1 < 1.0.0-rc.1 < 1.0.0)
- [ ] Resolve "latest" with no semantic versions to first configured branch
- [x] Prioritize semantic versions over branches
- [x] Handle mixed version formats (v1.0.0, 1.0.0, tags without semver)

## Data Structures

### Version
```go
type Version struct {
    Major      int
    Minor      int
    Patch      int
    Prerelease string
    Build      string
}
```

### Constraint
```go
type Constraint struct {
    Operator string // "^", "~", "=", ""
    Version  Version
}
```

## Algorithm

### Parse Version
1. Strip 'v' prefix if present
2. Split on '.' to get major.minor.patch
3. Split on '-' to extract prerelease
4. Split on '+' to extract build metadata
5. Parse numeric components
6. Return Version struct

### Parse Constraint
1. Detect operator (^, ~, or none)
2. Parse version component
3. Determine constraint type:
   - `1` → Major constraint (1.x.x)
   - `1.2` → Minor constraint (1.2.x)
   - `1.2.3` → Exact version
   - `^1.2.3` → Caret (>=1.2.3 <2.0.0)
   - `~1.2.3` → Tilde (>=1.2.3 <1.3.0)

### Resolve Version
1. Parse constraint from version string
2. Filter available versions that satisfy constraint
3. If no matches, return error
4. Sort matches by version (highest first)
5. Return first match

**Special cases:**
- "latest" constraint: matches all versions, returns highest
- "latest" with no semantic versions: should return first configured branch (currently returns lexicographically highest - BUG)
- Branch name constraint: matches exact branch name only

### Compare Versions
1. Compare major, minor, patch numerically
2. If equal, compare prerelease:
   - No prerelease > any prerelease
   - Split prerelease on '.'
   - Compare identifiers left-to-right
   - Numeric identifiers compared numerically
   - Alphanumeric compared lexically
   - Shorter prerelease < longer if all equal

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Version without 'v' prefix | Parse as semantic version |
| Version with 'v' prefix | Strip 'v' and parse |
| Prerelease versions | 1.0.0-alpha.1 < 1.0.0-alpha.2 < 1.0.0-beta.1 < 1.0.0-rc.1 < 1.0.0 |
| Build metadata | Ignored in comparison |
| No semantic versions | Use first configured branch (currently: lexicographic sort - BUG) |
| Mixed versions and branches | Semantic versions take precedence |
| Invalid version string | Return error |
| Constraint not satisfied | Return error with available versions |

## Dependencies

- Registry access (registry-management.md)
- Git operations (for branch resolution)

## Implementation Mapping

**Source files:**
- `internal/arm/core/version.go` - ParseVersion, CompareTo, comparePrerelease
- `internal/arm/core/constraint.go` - ParseConstraint, IsSatisfiedBy
- `internal/arm/core/helpers.go` - GetBestMatching, ResolveVersion
- `internal/arm/core/version_test.go` - Comprehensive version comparison tests
- `internal/arm/core/constraint_test.go` - Constraint parsing and matching tests
- `test/e2e/version_test.go` - E2E version resolution tests

## Examples

### Semantic Version Parsing
```go
ParseVersion("v1.2.3") → Version{1, 2, 3, "", ""}
ParseVersion("1.2.3-alpha.1") → Version{1, 2, 3, "alpha.1", ""}
ParseVersion("1.2.3-rc.1") → Version{1, 2, 3, "rc.1", ""}
ParseVersion("2.0.0+build.123") → Version{2, 0, 0, "", "build.123"}
```

### Constraint Matching
```go
constraint := ParseConstraint("^1.2.0")
constraint.IsSatisfiedBy(ParseVersion("1.2.5")) → true
constraint.IsSatisfiedBy(ParseVersion("1.3.0")) → true
constraint.IsSatisfiedBy(ParseVersion("2.0.0")) → false
```

### Version Comparison
```go
v1 := ParseVersion("1.0.0-alpha.1")
v2 := ParseVersion("1.0.0-alpha.2")
v3 := ParseVersion("1.0.0-beta.1")
v4 := ParseVersion("1.0.0-rc.1")
v5 := ParseVersion("1.0.0")

v1.CompareTo(v2) → -1 (v1 < v2)
v2.CompareTo(v3) → -1 (v2 < v3)
v3.CompareTo(v4) → -1 (v3 < v4)
v4.CompareTo(v5) → -1 (v4 < v5)
v5.CompareTo(v1) → 1  (v5 > v1)
```

### Resolution Examples
```bash
# Latest semantic version
arm install ruleset ai-rules/clean-code@latest cursor-rules
# Resolves to highest semver tag (e.g., v2.1.0)

# Major constraint
arm install ruleset ai-rules/clean-code@1 cursor-rules
# Resolves to highest 1.x.x version (e.g., v1.9.5)

# Minor constraint
arm install ruleset ai-rules/clean-code@1.2 cursor-rules
# Resolves to highest 1.2.x version (e.g., v1.2.8)

# Branch
arm install ruleset ai-rules/clean-code@develop cursor-rules
# Resolves to commit hash of develop branch HEAD
```
