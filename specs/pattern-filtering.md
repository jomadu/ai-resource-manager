# Pattern Filtering

## Job to be Done
Selectively install files from packages using glob patterns with include/exclude rules and automatic archive extraction.

## Activities
1. Filter files using glob patterns with ** wildcard support
2. Extract .zip and .tar.gz archives automatically
3. Apply include/exclude patterns to merged content
4. Handle pattern precedence (exclude overrides include)

## Acceptance Criteria
- [x] Support ** wildcard for recursive directory matching
- [x] Support * wildcard for single path component
- [x] Support --include patterns (OR logic - match any)
- [x] Support --exclude patterns (override includes)
- [ ] Default to **/*.yml and **/*.yaml if no patterns specified (BUG: registries return all files instead)
- [x] Extract .zip archives automatically
- [x] Extract .tar.gz archives automatically
- [x] Merge extracted files with loose files (archives take precedence)
- [x] Apply patterns after archive extraction
- [x] Prevent directory traversal attacks in archives
- [ ] Use consistent pattern matching across install and compile (BUG: compile uses filepath.Match on basename only)

## Data Structures

Pattern matching uses simple string slices passed to functions. No dedicated struct exists.

## Algorithm

### Match Pattern
1. Normalize path separators to /
2. If no wildcards, return exact match
3. If pattern contains `**`, use recursive matching
4. Otherwise use simple wildcard matching

### Filter Files (Registry Operations)
1. Apply defaults: if `len(include) == 0`, set `include = ["**/*.yml", "**/*.yaml"]`
2. Check exclude patterns first (if any match, skip file)
3. If no include patterns after defaults, include file
4. Check include patterns (if any match, include file)
5. Otherwise skip file

**Note:** Current implementation skips step 1 (BUG: git.go:199, gitlab.go:374, cloudsmith.go:337)

### Filter Files (Standalone Compilation)
1. Apply defaults: if `len(include) == 0`, set `include = ["*.yml", "*.yaml"]`
2. Get relative path from input directory
3. Check exclude patterns using `core.MatchPattern()` (if any match, skip file)
4. If no include patterns after defaults, include file
5. Check include patterns using `core.MatchPattern()` (if any match, include file)
6. Otherwise skip file

**Note:** Current implementation uses `filepath.Match(pattern, filepath.Base(filePath))` instead of `core.MatchPattern(pattern, filePath)` (BUG: service.go:1763)

### Extract Archive
1. Detect archive by extension (.zip, .tar.gz)
2. Write to temp file for streaming
3. For each entry:
   - Skip directories
   - Sanitize path: skip if contains `..`, is absolute, or equals `.`
   - Read content
4. Return extracted files

### Merge Files
1. Add all loose files to map
2. Extract archives and add to map (overwrites loose files with same path)
3. Convert map to slice
4. Return merged files

**Note:** Pattern filtering happens AFTER merge in registry GetPackage methods

## Edge Cases

| Condition | Expected Behavior | Current Status |
|-----------|-------------------|----------------|
| No patterns specified | Default to **/*.yml and **/*.yaml | ❌ Registries return all files |
| Empty include list | Include all files | ✅ Works |
| File matches include and exclude | Exclude wins (skip file) | ✅ Works |
| Multiple include patterns | OR logic (match any) | ✅ Works |
| Archive with ../ paths | Sanitize to prevent traversal | ✅ Works |
| Archive with absolute paths | Skip (security) | ✅ Works |
| Nested archives | Extract outer only (don't recurse) | ✅ Works |
| Corrupted archive | Return error, don't install | ✅ Works |
| Archive + loose file collision | Archive file takes precedence | ✅ Works |
| Compile with ** patterns | Should work like install | ❌ Uses filepath.Match on basename |

## Dependencies

- File system operations
- Archive libraries (zip, tar, gzip)

## Implementation Mapping

**Source files:**
- `internal/arm/core/pattern.go` - MatchPattern, matchDoublestar, matchSimpleWildcard
- `internal/arm/core/archive.go` - ExtractAndMerge, extractZip, extractTarGz
- `internal/arm/registry/git.go` - matchesPatterns (line 199) ⚠️ Missing default patterns
- `internal/arm/registry/gitlab.go` - matchesPatterns (line 374) ⚠️ Missing default patterns
- `internal/arm/registry/cloudsmith.go` - matchesPatterns (line 337) ⚠️ Missing default patterns
- `internal/arm/service/service.go` - matchesPatterns (line 1763) ⚠️ Wrong implementation, discoverFiles (line 1671) ✅ Correct defaults
- `test/e2e/archive_test.go` - E2E archive extraction tests
- `test/e2e/install_test.go` - E2E pattern filtering tests

## Known Bugs

### Bug 1: Registries Don't Apply Default Patterns
**Location:** `internal/arm/registry/{git,gitlab,cloudsmith}.go`  
**Issue:** When `len(include) == 0`, registries return ALL files instead of defaulting to `["**/*.yml", "**/*.yaml"]`  
**Fix:** Add default pattern logic before calling `matchesPatterns()`

### Bug 2: Compile Uses Wrong Pattern Matcher
**Location:** `internal/arm/service/service.go:1763`  
**Issue:** Uses `filepath.Match(pattern, filepath.Base(filePath))` instead of `core.MatchPattern(pattern, filePath)`  
**Impact:** Patterns like `security/**/*.yml` don't work in `arm compile`  
**Fix:** Replace with `core.MatchPattern()` call

## Examples

### Include Patterns
```bash
# Install only TypeScript rules
arm install ruleset --include "**/typescript-*.yml" ai-rules/language-rules cursor-rules

# Install multiple patterns (OR logic)
arm install promptset --include "review/**/*.yml" --include "refactor/**/*.yml" ai-rules/prompts cursor-commands
```

### Exclude Patterns
```bash
# Exclude experimental files
arm install ruleset --exclude "**/experimental/**" ai-rules/security cursor-rules

# Exclude multiple patterns
arm install ruleset --exclude "**/test/**" --exclude "**/draft/**" ai-rules/rules cursor-rules
```

### Combined Patterns
```bash
# Include security rules but exclude experimental
arm install ruleset --include "security/**/*.yml" --exclude "**/experimental/**" ai-rules/rules cursor-rules
```

### Archive Extraction
```bash
# Package contains rules.tar.gz
# ARM automatically extracts and applies patterns to contents
arm install ruleset --include "**/*.yml" ai-rules/archived-rules cursor-rules
```

### Pattern Matching Examples
```
Pattern: **/*.yml
Matches:
  - file.yml (root level - ** matches zero directories)
  - rules/clean-code.yml
  - security/auth/rules.yml
  - build/cursor/rule.yml
Does not match:
  - rules/clean-code.md (wrong extension)

Pattern: security/**/*.yml
Matches:
  - security/auth/rules.yml
  - security/crypto/aes.yml
Does not match:
  - rules/security.yml (wrong prefix)
  - security.yml (no directory after security/)

Pattern: **/typescript-*
Matches:
  - typescript-strict.yml (root level)
  - rules/typescript-strict.yml
  - lang/typescript-eslint.yml
Does not match:
  - rules/javascript-strict.yml (wrong prefix)
```

### Archive Precedence
```
Package contents:
  - rules/clean-code.yml (loose file)
  - rules.tar.gz containing:
    - rules/clean-code.yml (different content)
    - rules/security.yml

Result after extraction and merge:
  - rules/clean-code.yml (from archive - takes precedence)
  - rules/security.yml (from archive)
```
