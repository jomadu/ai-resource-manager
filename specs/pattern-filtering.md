# Pattern Filtering

## Job to be Done
Filter package files using glob patterns to install only relevant resources, enabling selective installation from large package repositories.

## Activities
1. **Normalize patterns** - Sanitize and sort patterns for consistent cache keys
2. **Match patterns** - Apply include/exclude logic to file paths
3. **Extract archives** - Extract tar.gz and zip archives before filtering
4. **Sanitize paths** - Prevent directory traversal attacks

## Acceptance Criteria
- [ ] Include patterns are OR'd together (match any)
- [ ] Exclude patterns override includes
- [ ] No patterns returns all files
- [ ] Archives extracted before pattern matching
- [ ] Path traversal attacks prevented (no `..`, no absolute paths)
- [ ] Patterns normalized for consistent cache keys (sorted, forward slashes)
- [ ] Glob patterns support `*`, `**`, and literal paths
- [ ] Empty pattern list treated as "match all"

## Data Structures

### Pattern Matching Input
```json
{
  "filePath": "string",
  "include": ["string"],
  "exclude": ["string"]
}
```

**Fields:**
- `filePath` - File path to test against patterns (relative, forward slashes)
- `include` - List of glob patterns to include (OR'd together)
- `exclude` - List of glob patterns to exclude (overrides includes)

### Archive Extraction Input
```json
{
  "path": "string",
  "content": "[]byte",
  "size": "int64"
}
```

**Fields:**
- `path` - Archive file path (`.tar.gz` or `.zip`)
- `content` - Archive file content as bytes
- `size` - Archive file size in bytes

## Algorithm

### 1. Normalize Patterns
Sanitize patterns for consistent cache keys and matching.

**Pseudocode:**
```
function normalizePatterns(patterns []string) []string:
    if len(patterns) == 0:
        return patterns
    
    normalized = []
    for pattern in patterns:
        // Trim whitespace and normalize path separators
        clean = trim(pattern)
        clean = replaceAll(clean, "\\", "/")
        normalized.append(clean)
    
    // Sort for consistent cache keys
    sort(normalized)
    return normalized
```

**Implementation:** `internal/arm/registry/git.go:normalizePatterns()`

### 2. Match Patterns
Apply include/exclude logic to determine if file should be included.

**Pseudocode:**
```
function matchesPatterns(filePath string, include []string, exclude []string) bool:
    // No patterns = include all
    if len(include) == 0 and len(exclude) == 0:
        return true
    
    // Check excludes first (they override includes)
    for pattern in exclude:
        if matchPattern(pattern, filePath):
            return false
    
    // No includes = include all (after excludes)
    if len(include) == 0:
        return true
    
    // Check includes (OR'd together)
    for pattern in include:
        if matchPattern(pattern, filePath):
            return true
    
    return false
```

**Implementation:** 
- `internal/arm/registry/gitlab.go:matchesPatterns()`
- `internal/arm/registry/cloudsmith.go:matchesPatterns()`

### 3. Match Pattern
Match a single glob pattern against a file path.

**Pseudocode:**
```
function matchPattern(pattern string, path string) bool:
    // Normalize path separators
    pattern = replaceAll(pattern, "\\", "/")
    path = replaceAll(path, "\\", "/")
    
    // Handle **/ prefix (match anywhere in path)
    if startsWith(pattern, "**/"):
        suffix = pattern[3:]
        return endsWith(path, suffix) or contains(path, "/" + suffix)
    
    // Handle /** suffix (match directory and descendants)
    if endsWith(pattern, "/**"):
        prefix = pattern[:len(pattern)-3]
        return startsWith(path, prefix + "/") or path == prefix
    
    // Handle literal paths (no wildcards)
    if not contains(pattern, "*"):
        return pattern == path
    
    // Handle wildcards using simple matching
    parts = split(pattern, "*")
    pos = 0
    for i, part in enumerate(parts):
        if part == "":
            continue
        
        idx = indexOf(path[pos:], part)
        if idx == -1 or (i == 0 and idx != 0):
            return false
        
        pos += idx + len(part)
    
    // Check last part matches end of path
    if len(parts) > 0 and parts[len(parts)-1] != "" and not endsWith(path, parts[len(parts)-1]):
        return false
    
    return true
```

**Implementation:** `internal/arm/registry/cloudsmith.go:matchPattern()`

**Note:** GitLab registry uses `filepath.Match()` which doesn't support `**` patterns. Cloudsmith uses custom implementation with `**` support.

### 4. Extract Archive
Extract tar.gz or zip archives and return extracted files.

**Pseudocode:**
```
function extractArchive(file File) []File:
    if endsWith(file.Path, ".tar.gz"):
        return extractTarGz(file)
    else if endsWith(file.Path, ".zip"):
        return extractZip(file)
    else:
        return error("unsupported archive format")

function extractTarGz(file File) []File:
    // Create temp file for streaming
    tempFile = createTempFile("arm-extract-*.tar.gz")
    defer remove(tempFile)
    
    write(tempFile, file.Content)
    seek(tempFile, 0)
    
    // Open gzip and tar readers
    gzReader = gzip.NewReader(tempFile)
    defer close(gzReader)
    
    tarReader = tar.NewReader(gzReader)
    
    extractedFiles = []
    while true:
        header = tarReader.Next()
        if header == EOF:
            break
        
        // Skip directories
        if header.Typeflag == TypeDir:
            continue
        
        // Sanitize path (prevent directory traversal)
        cleanName = filepath.Clean(header.Name)
        if cleanName == "." or isAbsolute(header.Name) or contains(cleanName, ".."):
            continue
        
        // Read file content
        content = readAll(tarReader)
        
        extractedFiles.append(File{
            Path: cleanName,
            Content: content,
            Size: header.Size
        })
    
    return extractedFiles

function extractZip(file File) []File:
    // Create temp file
    tempFile = createTempFile("arm-extract-*.zip")
    defer remove(tempFile)
    
    write(tempFile, file.Content)
    close(tempFile)
    
    // Open zip reader
    zipReader = zip.OpenReader(tempFile.Name())
    defer close(zipReader)
    
    extractedFiles = []
    for zipFile in zipReader.Files:
        // Skip directories
        if zipFile.IsDir():
            continue
        
        // Sanitize path (prevent directory traversal)
        cleanName = filepath.Clean(zipFile.Name)
        if cleanName == "." or isAbsolute(zipFile.Name) or contains(cleanName, ".."):
            continue
        
        // Read file content
        rc = zipFile.Open()
        content = readAll(rc)
        close(rc)
        
        extractedFiles.append(File{
            Path: cleanName,
            Content: content,
            Size: zipFile.UncompressedSize64
        })
    
    return extractedFiles
```

**Implementation:** `internal/arm/core/archive.go`

### 5. Extract and Merge
Extract archives and merge with loose files (archives win on conflicts).

**Pseudocode:**
```
function extractAndMerge(files []File) []File:
    fileMap = {}
    
    // Add loose files first
    for file in files:
        if not isArchive(file.Path):
            fileMap[file.Path] = file
    
    // Extract archives (overwrite loose files)
    for file in files:
        if isArchive(file.Path):
            extractedFiles = extractArchive(file)
            for extracted in extractedFiles:
                fileMap[extracted.Path] = extracted
    
    // Convert map to slice
    mergedFiles = []
    for path, file in fileMap:
        mergedFiles.append(file)
    
    return mergedFiles
```

**Implementation:** `internal/arm/core/archive.go:ExtractAndMerge()`

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| No patterns provided | Return all files |
| Empty include list | Include all files (after excludes) |
| Empty exclude list | Apply includes only |
| Include and exclude both match | Exclude wins |
| Invalid glob pattern | Pattern treated as literal string |
| Path traversal (`../`, absolute paths) | File skipped during extraction |
| Nested archives | Not supported (only top-level archives extracted) |
| Empty archive | Return empty file list |
| Corrupted archive | Return error |
| Pattern with backslashes | Normalized to forward slashes |
| Pattern with leading/trailing whitespace | Trimmed during normalization |
| Multiple includes match same file | File included once |
| Archive and loose file same path | Archive file wins |

## Dependencies

- `filepath` package for path operations
- `archive/tar` and `compress/gzip` for tar.gz extraction
- `archive/zip` for zip extraction
- `sort` package for pattern normalization

## Implementation Mapping

**Source files:**
- `internal/arm/core/archive.go` - Archive extraction and merging
- `internal/arm/registry/git.go` - Pattern normalization
- `internal/arm/registry/gitlab.go` - Pattern matching (simple glob)
- `internal/arm/registry/cloudsmith.go` - Pattern matching (with `**` support)

**Related specs:**
- `registry-management.md` - Registry types and GetPackage interface
- `package-installation.md` - Install workflow using pattern filtering
- `cache-management.md` - Cache keys include normalized patterns

## Examples

### Example 1: Include Only TypeScript Files

**Input:**
```bash
arm install ruleset --include "**/typescript-*.yml" ai-rules/language-rules cursor-rules
```

**Pattern Matching:**
- `include`: `["**/typescript-*.yml"]`
- `exclude`: `[]`

**Files in Package:**
- `typescript-basics.yml` ✅ (matches `**/typescript-*.yml`)
- `typescript-advanced.yml` ✅ (matches `**/typescript-*.yml`)
- `javascript-basics.yml` ❌ (no match)
- `python-basics.yml` ❌ (no match)

**Verification:**
- Only TypeScript files installed
- JavaScript and Python files excluded

### Example 2: Exclude Experimental Files

**Input:**
```bash
arm install ruleset --exclude "**/experimental/**" ai-rules/security-ruleset cursor-rules
```

**Pattern Matching:**
- `include`: `[]` (all files)
- `exclude`: `["**/experimental/**"]`

**Files in Package:**
- `security/auth.yml` ✅ (not excluded)
- `security/crypto.yml` ✅ (not excluded)
- `experimental/new-feature.yml` ❌ (matches exclude)
- `security/experimental/test.yml` ❌ (matches exclude)

**Verification:**
- All files except experimental included
- Experimental files excluded regardless of depth

### Example 3: Include and Exclude Together

**Input:**
```bash
arm install ruleset --include "*.yml" --exclude "test*" ai-rules/clean-code cursor-rules
```

**Pattern Matching:**
- `include`: `["*.yml"]`
- `exclude`: `["test*"]`

**Files in Package:**
- `rule.yml` ✅ (matches include, not excluded)
- `test.yml` ❌ (matches include but excluded)
- `test-rule.yml` ❌ (matches include but excluded)
- `doc.md` ❌ (no match)

**Verification:**
- Only `.yml` files included
- Files starting with `test` excluded
- Exclude overrides include

### Example 4: Archive Extraction with Patterns

**Input:**
```bash
arm install ruleset --include "security/**/*.yml" --exclude "**/experimental/**" ai-rules/archive-rules cursor-rules
```

**Archive Contents (`rules.tar.gz`):**
- `security/rule1.yml`
- `security/rule2.yml`
- `general/rule3.yml`
- `experimental/rule4.yml`

**Extraction Process:**
1. Extract archive to memory
2. Apply patterns to extracted paths
3. Filter results

**Pattern Matching:**
- `security/rule1.yml` ✅ (matches `security/**/*.yml`)
- `security/rule2.yml` ✅ (matches `security/**/*.yml`)
- `general/rule3.yml` ❌ (no match)
- `experimental/rule4.yml` ❌ (matches exclude)

**Verification:**
- Archive extracted before filtering
- Only security files included
- Experimental files excluded

### Example 5: Path Traversal Prevention

**Archive Contents:**
- `../../../etc/passwd` (malicious)
- `rule.yml` (valid)
- `/etc/shadow` (malicious)
- `subdir/../../../secret` (malicious)

**Extraction Process:**
1. Extract each file
2. Clean path with `filepath.Clean()`
3. Check for:
   - Absolute paths (`/etc/shadow`)
   - Parent directory references (`..`)
   - Current directory (`.`)

**Results:**
- `../../../etc/passwd` ❌ (contains `..`)
- `rule.yml` ✅ (valid relative path)
- `/etc/shadow` ❌ (absolute path)
- `subdir/../../../secret` ❌ (contains `..` after cleaning)

**Verification:**
- Only safe relative paths extracted
- Malicious paths skipped silently

## Notes

### Design Decisions

1. **Exclude overrides include** - Simplifies logic and matches user expectations (blacklist > whitelist)
2. **OR logic for includes** - Allows multiple patterns without complex boolean logic
3. **Archives extracted before filtering** - Patterns match extracted paths, not archive names
4. **Path sanitization during extraction** - Security check happens early, before pattern matching
5. **Pattern normalization for cache keys** - Ensures consistent cache hits regardless of pattern order or path separators

### Implementation Differences

- **GitLab registry** uses `filepath.Match()` which doesn't support `**` patterns
- **Cloudsmith registry** uses custom `matchPattern()` with `**` support
- **Git registry** applies patterns after extraction (archives extracted by storage layer)

### Testing Considerations

- Test include-only, exclude-only, and combined scenarios
- Test `**` patterns for deep directory matching
- Test path traversal attacks with various malicious paths
- Test archive extraction with nested directories
- Test pattern normalization with different path separators
- Test empty pattern lists (should include all files)
- Test pattern caching (same patterns = same cache key)
