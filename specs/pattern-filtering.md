# Pattern Filtering

## ⚠️ BREAKING CHANGE v5.0 - HIGH PRIORITY

**Archive extraction behavior is changing to prevent collisions and enable skillset path resolution.**

**Current (v3.x):** Archives merge with loose files, causing collisions  
**New (v5.0):** Archives extract to subdirectories named after the archive

See `BREAKING_CHANGE_ARCHIVE_EXTRACTION.md` for full details and migration guide.

## Job to be Done
Selectively install files from packages using glob patterns with include/exclude rules and automatic archive extraction.

## Activities
1. Filter files using glob patterns with ** wildcard support
2. Extract .zip and .tar.gz archives to subdirectories
3. Apply include/exclude patterns to extracted content
4. Handle pattern precedence (exclude overrides include)

## Acceptance Criteria
- [x] Support ** wildcard for recursive directory matching
- [x] Support * wildcard for single path component
- [x] Support --include patterns (OR logic - match any)
- [x] Support --exclude patterns (override includes)
- [x] Default to **/*.yml and **/*.yaml if no patterns specified
- [x] Extract .zip archives automatically
- [x] Extract .tar.gz archives automatically
- [ ] **BREAKING CHANGE v5.0**: Extract archives to subdirectories named after archive (prevents collisions, enables skillset path resolution)
- [x] Apply patterns after archive extraction
- [x] Prevent directory traversal attacks in archives
- [x] Use consistent pattern matching across install and compile

## Data Structures

Pattern matching uses simple string slices passed to functions. No dedicated struct exists.

## Algorithm

### Match Pattern (core.MatchPattern)
1. Normalize path separators to /
2. If no wildcards, return exact match
3. If pattern contains `**`, use recursive matching
4. Otherwise use simple wildcard matching

### Filter Files (Registry Operations)
1. Apply defaults: if `len(include) == 0`, set `include = ["**/*.yml", "**/*.yaml"]` (BUG: not implemented)
2. Check exclude patterns using `core.MatchPattern()` (if any match, skip file)
3. Check include patterns using `core.MatchPattern()` (if any match, include file)
4. Otherwise skip file

### Filter Files (Standalone Compilation)
1. Apply defaults: if `len(include) == 0`, set `include = ["*.yml", "*.yaml"]`
2. Get relative path from input directory
3. Check exclude patterns using `core.MatchPattern()` on full path
4. Check include patterns using `core.MatchPattern()` on full path
5. Otherwise skip file

### Extract Archive
1. Detect archive by extension (.zip, .tar.gz)
2. Determine subdirectory name: strip extension(s) from filename
   - `archive.tar.gz` → `archive/`
   - `rules.zip` → `rules/`
3. Write to temp file for streaming
4. For each entry:
   - Skip directories
   - Sanitize path: skip if contains `..`, is absolute, or equals `.`
   - Prepend subdirectory name to path
   - Read content
5. Return extracted files with subdirectory prefix

**Example:**
- Input: `my-package.tar.gz` containing `rules/rule1.yml`
- Output: File with path `my-package/rules/rule1.yml`

### Process Files (No Merge)
1. For each file:
   - If archive: extract to subdirectory (see Extract Archive)
   - If loose file: keep as-is
2. Apply include/exclude pattern filtering to all files
3. Return filtered files (archives and loose files coexist)

## Edge Cases

| Condition | Expected Behavior | Current Status |
|-----------|-------------------|----------------|
| No patterns specified | Default to **/*.yml and **/*.yaml | ✅ Works |
| Empty include list | Include all files | ✅ Works |
| File matches include and exclude | Exclude wins (skip file) | ✅ Works |
| Multiple include patterns | OR logic (match any) | ✅ Works |
| Archive with ../ paths | Sanitize to prevent traversal | ✅ Works |
| Archive with absolute paths | Skip (security) | ✅ Works |
| Nested archives | Extract outer only (don't recurse) | ✅ Works |
| Corrupted archive | Return error, don't install | ✅ Works |
| Multiple archives with same structure | Each extracts to own subdirectory (no collision) | ❌ BREAKING: Currently merges and collides |
| Archive + loose file same name | Both preserved in different paths | ❌ BREAKING: Currently archive overwrites loose |
| Compile with ** patterns | Should work like install | ✅ Works |

## Dependencies

- File system operations
- Archive libraries (zip, tar, gzip)

## Implementation Mapping

**Source files:**
- `internal/arm/core/pattern.go` - MatchPattern, matchDoublestar, matchSimpleWildcard
- `internal/arm/core/archive.go` - **NEEDS UPDATE**: Rename ExtractAndMerge → Extract, add subdirectory logic
- `internal/arm/registry/git.go` - **NEEDS UPDATE**: Change ExtractAndMerge → Extract (line 168)
- `internal/arm/registry/gitlab.go` - **NEEDS UPDATE**: Change ExtractAndMerge → Extract (line 214), missing default patterns
- `internal/arm/registry/cloudsmith.go` - **NEEDS UPDATE**: Change ExtractAndMerge → Extract (line 255), missing default patterns
- `internal/arm/service/service.go` - matchesPatterns (BUG: uses filepath.Match on basename), discoverFiles
- `test/e2e/archive_test.go` - **NEEDS UPDATE**: Update expectations for subdirectory structure
- `test/e2e/install_test.go` - E2E pattern filtering tests

## Known Bugs

### Bug: Archive Merge Causes Collisions (BREAKING CHANGE v5.0)
**Files:** `internal/arm/core/archive.go`, all registry implementations  
**Issue:** Archives are merged with loose files, causing naming collisions and breaking skillset path resolution  
**Impact:** 
- Multiple archives with same internal structure collide
- Loose files are overwritten by archive files with same path
- Skillset `source` references can't reliably point to files in archives
**Fix:** Extract archives to subdirectories named after the archive (minus extension)  
**Priority:** HIGH - Required for v5.0 release

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

### Archive Extraction (v5.0+ Behavior)
```bash
# Package contains rules.tar.gz
# ARM extracts to subdirectory named after archive
arm install ruleset --include "**/*.yml" ai-rules/archived-rules cursor-rules
```

**Package contents:**
```
- file.txt (loose file)
- archive.tar.gz containing:
  - file.txt
  - other.txt
```

**After extraction (v5.0+):**
```
- file.txt (loose file)
- archive/
  - file.txt (from archive)
  - other.txt (from archive)
```

**Before extraction (v3.x - DEPRECATED):**
```
- file.txt (from archive - overwrote loose file)
- other.txt (from archive)
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

### Archive Precedence (DEPRECATED in v5.0)

**v3.x behavior (DEPRECATED):**
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

**v5.0+ behavior (NO MERGE):**
```
Package contents:
  - rules/clean-code.yml (loose file)
  - rules.tar.gz containing:
    - rules/clean-code.yml (different content)
    - rules/security.yml

Result after extraction (no merge):
  - rules/clean-code.yml (loose file - preserved)
  - rules/
    - rules/clean-code.yml (from archive - preserved)
    - rules/security.yml (from archive)
```

Both files are preserved in different paths. No collision.
