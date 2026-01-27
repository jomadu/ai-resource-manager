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
- [x] Default to **/*.yml and **/*.yaml if no patterns specified
- [x] Extract .zip archives automatically
- [x] Extract .tar.gz archives automatically
- [x] Merge extracted files with loose files (archives take precedence)
- [x] Apply patterns after archive extraction
- [x] Prevent directory traversal attacks in archives

## Data Structures

### Pattern Matching
```go
type PatternMatcher struct {
    Include []string  // OR'd together
    Exclude []string  // Override includes
}
```

## Algorithm

### Match Pattern
1. Normalize path separators to /
2. Split pattern and path on /
3. Match components left-to-right:
   - `**` matches zero or more path components
   - `*` matches any characters except /
   - Literal matches exact component
4. Return true if pattern matches entire path

### Filter Files
1. If no include patterns, default to ["**/*.yml", "**/*.yaml"]
2. For each file:
   - Check if matches any include pattern (OR logic)
   - If matches, check if matches any exclude pattern
   - If excluded, skip file
   - Otherwise, include file
3. Return filtered file list

### Extract Archive
1. Detect archive by extension (.zip, .tar.gz)
2. Open archive for reading
3. For each entry:
   - Sanitize path (prevent ../ traversal)
   - Extract to temporary directory
4. Return extracted files

### Merge Files
1. Extract all archives to temporary directory
2. Collect loose files from package
3. Merge with archive files taking precedence
4. Apply include/exclude patterns to merged list
5. Return final file list

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| No patterns specified | Default to **/*.yml and **/*.yaml |
| Empty include list | Include all files |
| File matches include and exclude | Exclude wins (skip file) |
| Multiple include patterns | OR logic (match any) |
| Archive with ../ paths | Sanitize to prevent traversal |
| Archive with absolute paths | Convert to relative |
| Nested archives | Extract outer only (don't recurse) |
| Corrupted archive | Return error, don't install |
| Archive + loose file collision | Archive file takes precedence |

## Dependencies

- File system operations
- Archive libraries (zip, tar, gzip)

## Implementation Mapping

**Source files:**
- `internal/arm/core/pattern.go` - MatchPattern, matchDoublestar, matchSimpleWildcard
- `internal/arm/core/archive.go` - ExtractAndMerge, extractZip, extractTarGz
- `internal/arm/registry/git.go` - matchesPatterns (applies patterns)
- `internal/arm/registry/gitlab.go` - matchesPatterns
- `internal/arm/registry/cloudsmith.go` - matchesPatterns
- `test/e2e/archive_test.go` - E2E archive extraction tests
- `test/e2e/install_test.go` - E2E pattern filtering tests

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
  - rules/clean-code.yml ✓
  - security/auth/rules.yml ✓
  - build/cursor/rule.yml ✓
Does not match:
  - rules/clean-code.md ✗
  - README.yml (no directory) ✗

Pattern: security/**/*.yml
Matches:
  - security/auth/rules.yml ✓
  - security/crypto/aes.yml ✓
Does not match:
  - rules/security.yml ✗
  - security.yml ✗

Pattern: **/typescript-*
Matches:
  - rules/typescript-strict.yml ✓
  - lang/typescript-eslint.yml ✓
Does not match:
  - rules/javascript-strict.yml ✗
  - typescript.yml ✗
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
