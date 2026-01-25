# Version Resolution

## Job to be Done
Resolve package versions from registries to determine which version to install based on user constraints and available versions.

## Activities
1. Parse version strings into structured Version objects
2. Parse constraint strings into Constraint objects
3. Match available versions against constraints
4. Select the best matching version (highest semver)

## Acceptance Criteria
- [ ] Semver versions (1.0.0, v2.1.3) are parsed correctly with major, minor, patch components
- [ ] Non-semver versions (branch names like "main", "develop") are accepted but marked as non-semver
- [ ] Version comparison works for semver versions (2.0.0 > 1.9.9)
- [ ] Version comparison rejects non-semver versions with error
- [ ] Exact constraints (1.0.0) match only that specific version
- [ ] Major constraints (^1.0.0 or 1) match any version >= 1.0.0 with same major
- [ ] Minor constraints (~1.2.0 or 1.2) match any version >= 1.2.0 with same major.minor
- [ ] Latest constraint matches any version (or specific branch name)
- [ ] Version resolution returns highest matching semver version
- [ ] Version resolution returns error when no versions satisfy constraint
- [ ] Empty version string returns error
- [ ] Malformed semver strings are treated as non-semver (branch names)

## Data Structures

### Version
```json
{
  "Major": 1,
  "Minor": 2,
  "Patch": 3,
  "Prerelease": "alpha.1",
  "Build": "build.123",
  "Version": "v1.2.3-alpha.1+build.123",
  "IsSemver": true
}
```

**Fields:**
- `Major` - Major version number (breaking changes)
- `Minor` - Minor version number (new features)
- `Patch` - Patch version number (bug fixes)
- `Prerelease` - Optional prerelease identifier (e.g., "alpha.1", "beta.2")
- `Build` - Optional build metadata (e.g., "build.123")
- `Version` - Original version string
- `IsSemver` - True if version matches semver pattern, false for branch names

### Constraint
```json
{
  "Type": "major",
  "Version": {
    "Major": 1,
    "Minor": 0,
    "Patch": 0,
    "Version": "1.0.0",
    "IsSemver": true
  }
}
```

**Fields:**
- `Type` - Constraint type: "exact", "major", "minor", "latest"
- `Version` - Base version for constraint (null for "latest" without branch)

**Constraint Types:**
- `exact` - Match specific version (1.0.0)
- `major` - Match major version range (^1.0.0 or 1)
- `minor` - Match minor version range (~1.2.0 or 1.2)
- `latest` - Match any version or specific branch name

## Algorithm

### Parse Version
1. Check if version string is empty → error
2. Try to match semver regex: `^(v)?(\d+)\.(\d+)\.(\d+)(?:-([.\w-]+))?(?:\+([\w.-]+))?$`
3. If matches:
   - Extract major, minor, patch, prerelease, build
   - Set IsSemver = true
4. If no match:
   - Store as plain version string (branch name)
   - Set IsSemver = false
5. Return Version object

**Pseudocode:**
```
function ParseVersion(versionString):
    if versionString is empty:
        return error "version string cannot be empty"
    
    matches = semverRegex.match(versionString)
    if matches:
        return Version{
            Major: matches[2],
            Minor: matches[3],
            Patch: matches[4],
            Prerelease: matches[5],
            Build: matches[6],
            Version: versionString,
            IsSemver: true
        }
    else:
        return Version{
            Version: versionString,
            IsSemver: false
        }
```

### Parse Constraint
1. Check if constraint is "latest" or empty → return Latest constraint
2. Extract prefix (^, ~) if present
3. Try to match constraint regex: `^(v)?(\d+)(?:\.(\d+))?(?:\.(\d+))?$`
4. If no match and has prefix → error (prefix requires version)
5. If no match and no prefix → treat as branch name (Latest with version)
6. If matches:
   - Parse major, minor (default 0), patch (default 0)
   - Build version string
   - Determine constraint type based on prefix and input format
7. Return Constraint object

**Constraint Type Rules:**
- Prefix `^` → Major constraint
- Prefix `~` → Minor constraint
- Input has patch (1.0.0) → Exact constraint
- Input has minor (1.0) → Minor constraint
- Input has only major (1) → Major constraint

**Pseudocode:**
```
function ParseConstraint(constraintString):
    if constraintString is "latest" or empty:
        return Constraint{Type: Latest}
    
    prefix = extractPrefix(constraintString)  // ^, ~, or empty
    rest = constraintString without prefix
    
    matches = constraintRegex.match(rest)
    if not matches:
        if prefix exists:
            return error "prefix requires version"
        // Branch name
        return Constraint{
            Type: Latest,
            Version: Version{Version: rest, IsSemver: false}
        }
    
    major = matches[2]
    minor = matches[3] or 0
    patch = matches[4] or 0
    
    version = Version{
        Major: major,
        Minor: minor,
        Patch: patch,
        Version: buildVersionString(major, minor, patch),
        IsSemver: true
    }
    
    if prefix is "^":
        return Constraint{Type: Major, Version: version}
    else if prefix is "~":
        return Constraint{Type: Minor, Version: version}
    else if matches[4] exists:  // Has patch
        return Constraint{Type: Exact, Version: version}
    else if matches[3] exists:  // Has minor
        return Constraint{Type: Minor, Version: version}
    else:  // Only major
        return Constraint{Type: Major, Version: version}
```

### Version Comparison
1. Check both versions are semver → error if not
2. Compare major: if different, return result
3. Compare minor: if different, return result
4. Compare patch: if different, return result
5. Return 0 (equal)

**Pseudocode:**
```
function Compare(v1, v2):
    if not v1.IsSemver or not v2.IsSemver:
        return error "cannot compare non-semver"
    
    if v1.Major < v2.Major: return -1
    if v1.Major > v2.Major: return 1
    if v1.Minor < v2.Minor: return -1
    if v1.Minor > v2.Minor: return 1
    if v1.Patch < v2.Patch: return -1
    if v1.Patch > v2.Patch: return 1
    return 0
```

### Constraint Satisfaction
1. If constraint is Latest with version → exact string match
2. If constraint is Latest without version → accept any version
3. For all other constraints, version must be semver → error if not
4. Check constraint type:
   - Exact: version == constraint.Version
   - Major: version.Major == constraint.Version.Major AND version >= constraint.Version
   - Minor: version.Major == constraint.Version.Major AND version.Minor == constraint.Version.Minor AND version >= constraint.Version
5. Return true/false

**Pseudocode:**
```
function IsSatisfiedBy(constraint, version):
    if constraint.Type is Latest:
        if constraint.Version exists:
            return version.Version == constraint.Version.Version
        return true
    
    if not version.IsSemver:
        return error "constraint requires semver"
    
    switch constraint.Type:
        case Exact:
            return version.Compare(constraint.Version) == 0
        case Major:
            return version.Major == constraint.Version.Major 
                   AND version.Compare(constraint.Version) >= 0
        case Minor:
            return version.Major == constraint.Version.Major
                   AND version.Minor == constraint.Version.Minor
                   AND version.Compare(constraint.Version) >= 0
```

### Resolve Version
1. Parse constraint string
2. Filter available versions by constraint satisfaction
3. If no candidates → error "no version satisfies constraint"
4. Sort candidates by version (highest first)
5. Return highest version

**Pseudocode:**
```
function ResolveVersion(constraintString, availableVersions):
    constraint = ParseConstraint(constraintString)
    
    candidates = []
    for version in availableVersions:
        satisfied, err = constraint.IsSatisfiedBy(version)
        if err:
            continue  // Skip non-matching versions
        if satisfied:
            candidates.append(version)
    
    if candidates is empty:
        return error "no version satisfies constraint"
    
    sort candidates by Compare (highest first)
    return candidates[0]
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Empty version string | Error: "version string cannot be empty" |
| Non-semver version (branch name) | Accepted, IsSemver = false |
| Comparing non-semver versions | Error: "cannot compare non-semver version" |
| Constraint with prefix but no version | Error: "prefix requires version" |
| No versions available | Error: "no versions available" |
| No versions satisfy constraint | Error: "no version satisfies constraint" |
| Multiple versions satisfy constraint | Return highest semver version |
| Malformed semver (missing components) | Treated as non-semver (branch name) |
| Semver with prerelease/build metadata | Parsed correctly, metadata ignored in comparison |
| Version with 'v' prefix | Accepted, prefix preserved in Version string |
| Constraint without 'v' prefix | Accepted, normalized to semver format |

## Dependencies

- Go standard library: `regexp`, `strconv`, `sort`, `fmt`
- No external dependencies

## Implementation Mapping

**Source files:**
- `internal/arm/core/version.go` - Version struct, parsing, comparison
- `internal/arm/core/constraint.go` - Constraint struct, parsing, satisfaction
- `internal/arm/core/helpers.go` - GetBestMatching helper function

**Related specs:**
- `registry-management.md` - How registries provide available versions
- `package-installation.md` - How version resolution is used during install

## Examples

### Example 1: Parse Semver Version

**Input:**
```
"v1.2.3-alpha.1+build.123"
```

**Expected Output:**
```json
{
  "Major": 1,
  "Minor": 2,
  "Patch": 3,
  "Prerelease": "alpha.1",
  "Build": "build.123",
  "Version": "v1.2.3-alpha.1+build.123",
  "IsSemver": true
}
```

**Verification:**
- Major = 1, Minor = 2, Patch = 3
- Prerelease = "alpha.1", Build = "build.123"
- IsSemver = true

### Example 2: Parse Branch Name

**Input:**
```
"main"
```

**Expected Output:**
```json
{
  "Version": "main",
  "IsSemver": false
}
```

**Verification:**
- Version = "main"
- IsSemver = false

### Example 3: Major Constraint Resolution

**Input:**
```
Constraint: "^1.0.0"
Available: ["0.9.0", "1.0.0", "1.2.3", "2.0.0"]
```

**Expected Output:**
```
"1.2.3"
```

**Verification:**
- Constraint matches 1.0.0 and 1.2.3 (same major, >= 1.0.0)
- Returns highest: 1.2.3

### Example 4: Minor Constraint Resolution

**Input:**
```
Constraint: "~1.2.0"
Available: ["1.1.9", "1.2.0", "1.2.5", "1.3.0"]
```

**Expected Output:**
```
"1.2.5"
```

**Verification:**
- Constraint matches 1.2.0 and 1.2.5 (same major.minor, >= 1.2.0)
- Returns highest: 1.2.5

### Example 5: No Matching Version

**Input:**
```
Constraint: "2.0.0"
Available: ["1.0.0", "1.2.3"]
```

**Expected Output:**
```
Error: "no version satisfies constraint: 2.0.0"
```

**Verification:**
- No versions match exact constraint 2.0.0
- Returns error

### Example 6: Abbreviated Constraint

**Input:**
```
Constraint: "1.2"
Available: ["1.1.0", "1.2.0", "1.2.5", "1.3.0"]
```

**Expected Output:**
```
"1.2.5"
```

**Verification:**
- Constraint "1.2" is interpreted as minor constraint (~1.2.0)
- Matches 1.2.0 and 1.2.5
- Returns highest: 1.2.5

## Notes

**Semver Regex**: The regex `^(v)?(\d+)\.(\d+)\.(\d+)(?:-([.\w-]+))?(?:\+([\w.-]+))?$` strictly requires all three components (major.minor.patch). Versions missing components are treated as non-semver.

**Constraint Regex**: The regex `^(v)?(\d+)(?:\.(\d+))?(?:\.(\d+))?$` allows abbreviated forms (1, 1.2, 1.2.3) which are expanded to full semver with missing components defaulting to 0.

**Prerelease/Build Metadata**: Currently parsed but not used in version comparison. Future enhancement could implement semver prerelease precedence rules.

**Branch Name Handling**: Non-semver strings are treated as branch names. The Latest constraint with a version enables exact branch name matching (e.g., "main", "develop").

**Version Priority**: In Git registries, semver tags always take precedence over branches. This is enforced at the registry level, not in version resolution.
