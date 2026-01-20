# Version and Constraint Interface Specification

## Overview

This document specifies the interface-based design for Version and Constraint types in ARM. The refactor moves from struct-based interpretation to behavior encapsulation through methods.

## Version

### Constructor

```go
NewVersion(s string) (Version, error)
```

**Accepts:**
- Semver strings: `1.2.3`, `v1.2.3`
- Non-semver strings: `main`, `abc123def` (branch names, commit hashes)

**Returns error if:**
- Empty string

### Methods

#### CompareTo

```go
CompareTo(other Version) (int, error)
```

**Behavior:**
- Returns `-1` if self is older than other
- Returns `0` if self equals other
- Returns `1` if self is newer than other
- For semver: compares major → minor → patch

**Returns error if:**
- Either version is non-semver (cannot determine order without git history)

#### IsNewerThan

```go
IsNewerThan(other Version) (bool, error)
```

**Behavior:**
- Returns `true` if `CompareTo(other) > 0`
- Returns `false` otherwise

**Returns error if:**
- Either version is non-semver

#### IsOlderThan

```go
IsOlderThan(other Version) (bool, error)
```

**Behavior:**
- Returns `true` if `CompareTo(other) < 0`
- Returns `false` otherwise

**Returns error if:**
- Either version is non-semver

#### ToString

```go
ToString() string
```

**Behavior:**
- Returns original string (preserves `v` prefix if present)
- Never returns error

---

## Constraint

### Constructor

```go
NewConstraint(s string) (Constraint, error)
```

**Accepts:**
- `latest` - special keyword
- Semver: `1.2.3`, `v1.2.3`
- Abbreviated semver: `1`, `1.2`, `v1`, `v1.2`
- Caret constraints: `^1.2.3`, `^v1.2.3`
- Tilde constraints: `~1.2.3`, `~v1.2.3`

**Returns error if:**
- Non-semver strings (branch names, commit hashes, invalid formats)

### Constraint Types

| Input | Type | Matches |
|-------|------|---------|
| `latest` | Latest | Any semver version |
| `1.2.3` | Exact | Exactly 1.2.3 |
| `1.2` | Minor | >=1.2.0, <1.3.0 |
| `1` | Major | >=1.0.0, <2.0.0 |
| `^1.2.3` | Major | >=1.2.3, <2.0.0 |
| `~1.2.3` | Minor | >=1.2.3, <1.3.0 |

### Methods

#### IsSatisfiedBy

```go
IsSatisfiedBy(version Version) (bool, error)
```

**Behavior:**
- Returns `true` if version satisfies the constraint
- Returns `false` otherwise

**Constraint logic:**
- `Latest`: always returns `true`
- `Exact`: version equals constraint version
- `Major`: version.major == constraint.major AND version >= constraint
- `Minor`: version.major == constraint.major AND version.minor == constraint.minor AND version >= constraint

**Returns error if:**
- Version is non-semver (cannot compare against semver constraint)

#### ToString

```go
ToString() string
```

**Behavior:**
- Returns original constraint string (preserves `v`, `^`, `~` prefix)
- Never returns error

---

## Examples

### Version Usage

```go
// Semver comparison
v1, _ := NewVersion("v1.2.3")
v2, _ := NewVersion("v1.3.0")
cmp, _ := v1.CompareTo(v2)  // -1
newer, _ := v2.IsNewerThan(v1)  // true
older, _ := v1.IsOlderThan(v2)  // true

// Non-semver version
v3, _ := NewVersion("main")
_, err := v1.CompareTo(v3)  // error: cannot compare semver to non-semver

// ToString preserves format
v4, _ := NewVersion("v1.2.3")
v4.ToString()  // "v1.2.3"
```

### Constraint Usage

```go
v1, _ := NewVersion("v1.2.3")
v2, _ := NewVersion("v1.3.0")
v3, _ := NewVersion("v2.0.0")

// Caret constraint
c1, _ := NewConstraint("^1.2.0")
ok, _ := c1.IsSatisfiedBy(v1)  // true (1.2.3 >= 1.2.0, same major)
ok, _ := c1.IsSatisfiedBy(v2)  // true (1.3.0 >= 1.2.0, same major)
ok, _ := c1.IsSatisfiedBy(v3)  // false (2.0.0 has different major)

// Tilde constraint
c2, _ := NewConstraint("~1.2.0")
ok, _ := c2.IsSatisfiedBy(v1)  // true (1.2.3 >= 1.2.0, same minor)
ok, _ := c2.IsSatisfiedBy(v2)  // false (1.3.0 has different minor)

// Latest constraint
c3, _ := NewConstraint("latest")
ok, _ := c3.IsSatisfiedBy(v1)  // true
ok, _ := c3.IsSatisfiedBy(v2)  // true
ok, _ := c3.IsSatisfiedBy(v3)  // true

// Abbreviated constraints
c4, _ := NewConstraint("1.2")  // Minor constraint
ok, _ := c4.IsSatisfiedBy(v1)  // true (1.2.3 >= 1.2.0, <1.3.0)

c5, _ := NewConstraint("1")  // Major constraint
ok, _ := c5.IsSatisfiedBy(v1)  // true (1.2.3 >= 1.0.0, <2.0.0)
ok, _ := c5.IsSatisfiedBy(v3)  // false (2.0.0 >= 2.0.0)

// Error cases
_, err := NewConstraint("main")  // error: not semver

v4, _ := NewVersion("abc123")  // non-semver version
_, err = c1.IsSatisfiedBy(v4)  // error: cannot check non-semver against semver constraint
```

---

## Key Changes from Current Implementation

1. **Remove BranchHead constraint type** - constraints are semver-only
2. **Add error returns** - comparison methods return errors for non-semver operations
3. **Constructor functions** - replace `ParseConstraint`/`ParseVersion` with `NewConstraint`/`NewVersion`
4. **Interface methods** - encapsulate behavior instead of exposing struct fields for interpretation
5. **Explicit error handling** - operations that cannot logically succeed return errors instead of undefined behavior
